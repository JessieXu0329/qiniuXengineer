package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aipr/ai-pr-reviewer/model"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

const (
	defaultLLMBaseURL = "https://api.deepseek.com/v1"
	defaultLLMModel   = "deepseek-chat"
)

type ReviewService struct {
	db          *gorm.DB
	rdb         *redis.Client
	ctx         context.Context
	restyClient *resty.Client
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

type LLMSummaryPart struct {
	OneSentenceSummary   string   `json:"one_sentence_summary"`
	OneSentenceSummaryZh string   `json:"one_sentence_summary_zh"`
	AffectedModules      []string `json:"affected_modules"`
	BreakingChanges      []string `json:"breaking_changes"`
	BreakingChangesZh    []string `json:"breaking_changes_zh"`
	OverallScore         float64  `json:"overall_score"`
	NormativeScore       float64  `json:"normative_score"`
	SecurityScore        float64  `json:"security_score"`
	PerformanceScore     float64  `json:"performance_score"`
	ReadabilityScore     float64  `json:"readability_score"`
}

type LLMCommentPart struct {
	File          string `json:"file"`
	Line          int    `json:"line"`
	Level         string `json:"level"`
	Reason        string `json:"reason"`
	ReasonZh      string `json:"reason_zh"`
	SuggestedCode string `json:"suggested_code"`
}

type DashboardResponse struct {
	TotalReviews        int64              `json:"total_reviews"`
	AverageScore        float64            `json:"average_score"`
	AvgNormativeScore   float64            `json:"avg_normative_score"`
	AvgSecurityScore    float64            `json:"avg_security_score"`
	AvgPerformanceScore float64            `json:"avg_performance_score"`
	AvgReadabilityScore float64            `json:"avg_readability_score"`
	CriticalIssues      int64              `json:"critical_issues"`
	WarningIssues       int64              `json:"warning_issues"`
	SuggestionIssues    int64              `json:"suggestion_issues"`
	RecentReviews       []RecentReviewItem `json:"recent_reviews"`
}

type RecentReviewItem struct {
	ID     uint    `json:"id"`
	PRURL  string  `json:"pr_url"`
	Score  float64 `json:"score"`
	Status string  `json:"status"`
	Date   string  `json:"date"`
}

func NewReviewService(db *gorm.DB, rdb *redis.Client) *ReviewService {
	return &ReviewService{
		db:          db,
		rdb:         rdb,
		ctx:         context.Background(),
		restyClient: resty.New().SetTimeout(60 * time.Second),
	}
}

func (s *ReviewService) AnalyzePR(prURL, owner, repo, prNumberText, customGitHubToken, customAPIKey, customBaseURL, customModel string) (*model.PRReviewTask, error) {
	apiKey := customAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("LLM_API_KEY")
	}
	if apiKey == "" {
		return nil, &HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Backend misconfiguration: LLM_API_KEY environment variable is empty.",
		}
	}

	prNumber, err := strconv.Atoi(prNumberText)
	if err != nil {
		return nil, &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid GitHub PR link format! Must match: https://github.com/{owner}/{repo}/pull/{number}",
		}
	}

	rawDiff, err := s.getPRDiff(owner, repo, prNumber, customGitHubToken)
	if err != nil {
		return nil, err
	}

	task := &model.PRReviewTask{
		GitHubURL: prURL,
		RepoOwner: owner,
		RepoName:  repo,
		PRNumber:  prNumber,
		Status:    "processing",
		RawDiff:   rawDiff,
	}

	if s.db != nil {
		if err := s.db.Create(task).Error; err != nil {
			return nil, fmt.Errorf("failed to create database task entry: %w", err)
		}
	} else {
		task.ID = 0
		task.CreatedAt = time.Now()
		task.UpdatedAt = task.CreatedAt
	}

	log.Printf("[LLM] Sending PR diff (%d bytes) to review engine...", len(rawDiff))
	comments, err := s.generateReviewComments(rawDiff, apiKey, customBaseURL, customModel)
	if err != nil {
		task.Status = "failed"
		if s.db != nil {
			_ = s.db.Save(task).Error
		}
		return nil, fmt.Errorf("ai code audit generation failed: %w", err)
	}

	summary, err := s.generatePRSummary(rawDiff, comments, apiKey, customBaseURL, customModel)
	if err != nil {
		task.Status = "failed"
		if s.db != nil {
			_ = s.db.Save(task).Error
		}
		return nil, fmt.Errorf("ai summary generation failed: %w", err)
	}

	summaryBytes, err := json.Marshal(summary)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal summary content: %w", err)
	}

	task.Status = "success"
	task.Summary = string(summaryBytes)
	task.Comments = make([]model.ReviewComment, 0, len(comments))
	for _, comment := range comments {
		task.Comments = append(task.Comments, model.ReviewComment{
			TaskID:        task.ID,
			FilePath:      comment.File,
			LineNumber:    comment.Line,
			Severity:      comment.Level,
			Reason:        comment.Reason,
			ReasonZh:      comment.ReasonZh,
			SuggestedCode: comment.SuggestedCode,
		})
	}

	if s.db != nil {
		if err := s.db.Save(task).Error; err != nil {
			return nil, fmt.Errorf("failed to finalize task review: %w", err)
		}
	}

	return task, nil
}

func (s *ReviewService) GetReviewDetail(taskID uint) (*model.PRReviewTask, error) {
	if s.db == nil {
		return nil, errors.New("database not available")
	}

	var task model.PRReviewTask
	if err := s.db.Preload("Comments").First(&task, taskID).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *ReviewService) GetDashboard() (*DashboardResponse, error) {
	if s.db == nil {
		return nil, &HTTPError{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "database not available for dashboard metrics",
		}
	}

	var tasks []model.PRReviewTask
	if err := s.db.Order("created_at desc").Limit(10).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to load recent reviews: %w", err)
	}

	var totalReviews int64
	if err := s.db.Model(&model.PRReviewTask{}).Count(&totalReviews).Error; err != nil {
		return nil, fmt.Errorf("failed to count reviews: %w", err)
	}

	var criticalIssues int64
	var warningIssues int64
	var suggestionIssues int64
	if err := s.db.Model(&model.ReviewComment{}).Where("severity = ?", "CRITICAL").Count(&criticalIssues).Error; err != nil {
		return nil, fmt.Errorf("failed to count critical issues: %w", err)
	}
	if err := s.db.Model(&model.ReviewComment{}).Where("severity = ?", "WARNING").Count(&warningIssues).Error; err != nil {
		return nil, fmt.Errorf("failed to count warning issues: %w", err)
	}
	if err := s.db.Model(&model.ReviewComment{}).Where("severity = ?", "SUGGESTION").Count(&suggestionIssues).Error; err != nil {
		return nil, fmt.Errorf("failed to count suggestion issues: %w", err)
	}

	var scoreSum, normSum, secSum, perfSum, readSum float64
	var scoredCount float64
	recentReviews := make([]RecentReviewItem, 0, len(tasks))
	for _, task := range tasks {
			score, norm, sec, perf, read := extractAllScores(task.Summary)
		if score > 0 {
			scoreSum += score
					normSum += norm
			secSum += sec
			perfSum += perf
			readSum += read
			scoredCount++
		}

		recentReviews = append(recentReviews, RecentReviewItem{
			ID:     task.ID,
			PRURL:  task.GitHubURL,
			Score:  score,
			Status: task.Status,
			Date:   task.CreatedAt.Format("2006-01-02"),
		})
	}

	var averageScore, avgNorm, avgSec, avgPerf, avgRead float64
	if scoredCount > 0 {
		averageScore = scoreSum / scoredCount
				avgNorm = normSum / scoredCount
		avgSec = secSum / scoredCount
		avgPerf = perfSum / scoredCount
		avgRead = readSum / scoredCount
	}

	return &DashboardResponse{
		TotalReviews:     totalReviews,
		AverageScore:     averageScore,
				AvgNormativeScore:   avgNorm,
		AvgSecurityScore:    avgSec,
		AvgPerformanceScore: avgPerf,
		AvgReadabilityScore: avgRead,
		CriticalIssues:   criticalIssues,
		WarningIssues:    warningIssues,
		SuggestionIssues: suggestionIssues,
		RecentReviews:    recentReviews,
	}, nil
}

func (s *ReviewService) getPRDiff(owner, repo string, prNumber int, customGitHubToken string) (string, error) {
	redisKey := fmt.Sprintf("pr:diff:%s:%s:%d", owner, repo, prNumber)
	if s.rdb != nil {
		if cached, err := s.rdb.Get(s.ctx, redisKey).Result(); err == nil && cached != "" {
			return cached, nil
		}
	}

	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, prNumber)
	request := s.restyClient.R().
		SetHeader("Accept", "application/vnd.github.v3.diff").
		SetHeader("User-Agent", "AI-PR-Review-Assistant")

	token := customGitHubToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token != "" {
		request.SetHeader("Authorization", "Bearer "+token)
	}

	resp, err := request.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("github api request failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("GitHub PR fetch failed with status %d. Please verify the repository and pull request exist.", resp.StatusCode()),
		}
	}

	diffText := resp.String()
	if strings.TrimSpace(diffText) == "" {
		return "", &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "GitHub PR fetch returned an empty diff.",
		}
	}

	if s.rdb != nil {
		_ = s.rdb.Set(s.ctx, redisKey, diffText, time.Hour).Err()
	}

	return diffText, nil
}

func (s *ReviewService) generateReviewComments(diffContent, apiKey, customBaseURL, customModel string) ([]LLMCommentPart, error) {
	systemPrompt := `You are an expert DevOps engineer and automated security code auditor.
Analyze the provided Git Diff payload and identify vulnerabilities.

[CRITICAL RULE]: You MUST output your analysis EXACTLY as a raw JSON array matching the schema below. 
DO NOT include any markdown code blocks like ` + "`" + "```json" + "`" + `, DO NOT include any backticks, and DO NOT add any pre-amble, explanations, or conversational text. Your entire response must start with '[' and end with ']'.

[JSON SCHEMA]:
[
  {
    "file": "string (filename)",
    "line": int (line number),
    "level": "string (MUST BE EXACTLY ONE OF THESE: CRITICAL, WARNING, SUGGESTION)",
    "reason": "string (concise bug description in English — mandatory)",
    "reason_zh": "string (concise bug description in Chinese — MANDATORY, NOT empty, NOT omitted)",
    "suggested_code": "string (exact fix code snippet)"
  }
]

[BILINGUAL ENFORCEMENT]: Both "reason" AND "reason_zh" MUST be non-empty strings for every entry. Submitting an empty "reason_zh" or omitting it entirely is INVALID. Write the Chinese description even if you are less confident — a rough Chinese translation is better than an empty field.

[EXAMPLE OF VALID OUTPUT - FOLLOW THIS EXACTLY]:
[{"file":"main.go","line":12,"level":"CRITICAL","reason":"Hardcoded credentials found.","reason_zh":"发现硬编码凭据。","suggested_code":"pass := os.Getenv(\"DB_PASS\")"}]`

	rawContent, err := s.chatCompletion(systemPrompt, "Review this Git diff:\n\n"+diffContent, 0.2, apiKey, customBaseURL, customModel)
	if err != nil {
		return nil, err
	}

	cleanContent := cleanLLMJSON(rawContent)
	var comments []LLMCommentPart
	if err := json.Unmarshal([]byte(cleanContent), &comments); err != nil {
		log.Printf("[ERROR] failed to parse LLM comments. Raw: %s", rawContent)
		return nil, fmt.Errorf("failed to parse AI comment JSON array: %w", err)
	}

	// Validate and backfill missing Chinese translations
	for i := range comments {
		if strings.TrimSpace(comments[i].ReasonZh) == "" {
			log.Printf("[WARN] LLM returned empty reason_zh for %s:%d — falling back to English", comments[i].File, comments[i].Line)
			comments[i].ReasonZh = comments[i].Reason
		}
	}

	return comments, nil
}

func (s *ReviewService) generatePRSummary(diffContent string, comments []LLMCommentPart, apiKey, customBaseURL, customModel string) (*LLMSummaryPart, error) {
	commentsBytes, _ := json.Marshal(comments)
	systemPrompt := `You are an expert PR summary generator.
Return a raw JSON object only, with no markdown, no code fences, and no explanatory text.
The object must match:
{
  "one_sentence_summary": "one sentence summary in English",
  "one_sentence_summary_zh": "one sentence summary in Chinese — MANDATORY, NOT empty (一句话改动概述)",
  "affected_modules": ["changed/file.ext"],
  "breaking_changes": ["breaking change description"],
  "breaking_changes_zh": ["breaking change description in Chinese - same length as breaking_changes, MANDATORY (破坏性变更说明)"],
  "overall_score": 90,
  "normative_score": 90,
  "security_score": 90,
  "performance_score": 90,
  "readability_score": 90
}
[BILINGUAL ENFORCEMENT]: one_sentence_summary_zh MUST be a non-empty Chinese sentence. The breaking_changes_zh array MUST have the SAME number of entries as breaking_changes, each entry being a non-empty Chinese translation. An empty Chinese field is INVALID — write even a rough translation.
The one_sentence_summary and one_sentence_summary_zh fields must be exactly one sentence. Scores must be numbers from 0 to 100.`

	userPrompt := fmt.Sprintf("Generate a one-sentence PR summary and score object for this diff and these review comments.\n\nReview comments JSON:\n%s\n\nGit diff:\n%s", string(commentsBytes), diffContent)
	rawContent, err := s.chatCompletion(systemPrompt, userPrompt, 0.2, apiKey, customBaseURL, customModel)
	if err != nil {
		return nil, err
	}

	cleanContent := cleanLLMJSON(rawContent)
	var summary LLMSummaryPart
	if err := json.Unmarshal([]byte(cleanContent), &summary); err != nil {
		log.Printf("[ERROR] failed to parse LLM summary. Raw: %s", rawContent)
		return nil, fmt.Errorf("failed to parse AI summary JSON object: %w", err)
	}

	// Validate and backfill missing Chinese translations in summary
	if strings.TrimSpace(summary.OneSentenceSummaryZh) == "" {
		log.Printf("[WARN] LLM returned empty one_sentence_summary_zh — falling back to English")
		summary.OneSentenceSummaryZh = summary.OneSentenceSummary
	}
	if len(summary.BreakingChangesZh) != len(summary.BreakingChanges) {
		log.Printf("[WARN] LLM returned mismatched breaking_changes_zh length (%d vs %d) — falling back to English", len(summary.BreakingChangesZh), len(summary.BreakingChanges))
		summary.BreakingChangesZh = make([]string, len(summary.BreakingChanges))
		copy(summary.BreakingChangesZh, summary.BreakingChanges)
	} else {
		for i := range summary.BreakingChangesZh {
			if strings.TrimSpace(summary.BreakingChangesZh[i]) == "" {
				summary.BreakingChangesZh[i] = summary.BreakingChanges[i]
			}
		}
	}

	return &summary, nil
}

func (s *ReviewService) chatCompletion(systemPrompt, userPrompt string, temperature float64, apiKey, customBaseURL, customModel string) (string, error) {
	if apiKey == "" {
		apiKey = os.Getenv("LLM_API_KEY")
	}
	if apiKey == "" {
		return "", &HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Backend misconfiguration: LLM_API_KEY environment variable is empty.",
		}
	}

	baseURL := customBaseURL
	if baseURL == "" {
		baseURL = os.Getenv("LLM_BASE_URL")
	}
	if baseURL == "" {
		baseURL = defaultLLMBaseURL
	}

	modelName := customModel
	if modelName == "" {
		modelName = os.Getenv("LLM_MODEL")
	}
	if modelName == "" {
		modelName = defaultLLMModel
	}

	payload := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": temperature,
	}

	endpoint := strings.TrimSuffix(baseURL, "/") + "/chat/completions"
	resp, err := s.restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(payload).
		Post(endpoint)
	if err != nil {
		return "", fmt.Errorf("llm api request failed: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("llm api returned error status %d: %s", resp.StatusCode(), resp.String())
	}

	type chatCompletionResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	var chatResp chatCompletionResponse
	if err := json.Unmarshal(resp.Body(), &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal llm response body: %w", err)
	}
	if len(chatResp.Choices) == 0 || strings.TrimSpace(chatResp.Choices[0].Message.Content) == "" {
		return "", errors.New("no completion content returned from the large language model")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func cleanLLMJSON(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```JSON")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	return strings.TrimSpace(cleaned)
}

func extractOverallScore(summary string) float64 {
	var payload struct {
		OverallScore float64 `json:"overall_score"`
	}
	if err := json.Unmarshal([]byte(summary), &payload); err != nil {
		return 0
	}
	return payload.OverallScore
}

func extractAllScores(summary string) (overall, normative, security, performance, readability float64) {
	var payload struct {
		OverallScore     float64 `json:"overall_score"`
		NormativeScore   float64 `json:"normative_score"`
		SecurityScore    float64 `json:"security_score"`
		PerformanceScore float64 `json:"performance_score"`
		ReadabilityScore float64 `json:"readability_score"`
	}
	if err := json.Unmarshal([]byte(summary), &payload); err != nil {
		return 0, 0, 0, 0, 0
	}
	return payload.OverallScore, payload.NormativeScore, payload.SecurityScore, payload.PerformanceScore, payload.ReadabilityScore
}
