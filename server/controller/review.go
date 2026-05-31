package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aipr/ai-pr-reviewer/service"
	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	reviewService *service.ReviewService
}

func NewReviewController(rs *service.ReviewService) *ReviewController {
	return &ReviewController{
		reviewService: rs,
	}
}

type AnalyzeRequest struct {
	GitHubURL   string `json:"github_url"`
	PRURL       string `json:"pr_url"`
	GitHubToken string `json:"github_token"`
	APIKey      string `json:"api_key"`
	BaseURL     string `json:"base_url"`
	Model       string `json:"model"`
}

func (ctrl *ReviewController) Analyze(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GitHub link! Format must match a Pull Request, Commit, or Compare view."})
		return
	}

	urlStr := strings.TrimSpace(req.GitHubURL)
	if urlStr == "" {
		urlStr = strings.TrimSpace(req.PRURL)
	}
	if !strings.Contains(urlStr, "github.com/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain! Link must be from github.com"})
		return
	}

	parts := strings.Split(urlStr, "/")
	hubIndex := -1
	for i, part := range parts {
		if part == "github.com" {
			hubIndex = i
			break
		}
	}

	if hubIndex == -1 || len(parts) < hubIndex+5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GitHub link! Format must provide owner, repo, and a valid target action."})
		return
	}

	owner := parts[hubIndex+1]
	repo := parts[hubIndex+2]
	urlType := parts[hubIndex+3]
	targetID := parts[hubIndex+4]

	if urlType != "pull" && urlType != "commit" && urlType != "compare" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported action! Only pull requests, commits, and comparisons are supported."})
		return
	}

	if owner == "" || repo == "" || targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GitHub link! Format must provide owner, repo, and a valid target action."})
		return
	}

	log.Printf("[DEBUG] Analyzing PR: %s (owner: %s, repo: %s, target: %s) with token len: %d, token prefix: %s", urlStr, owner, repo, targetID, len(req.GitHubToken), func() string {
		if len(req.GitHubToken) > 6 {
			return req.GitHubToken[:6] + "..."
		}
		return req.GitHubToken
	}())

	task, err := ctrl.reviewService.AnalyzePR(urlStr, owner, repo, targetID, req.GitHubToken, req.APIKey, req.BaseURL, req.Model)
	if err != nil {
		status := http.StatusInternalServerError
		var httpErr *service.HTTPError
		if errors.As(err, &httpErr) {
			status = httpErr.StatusCode
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	var summaryMap map[string]interface{}
	if err := json.Unmarshal([]byte(task.Summary), &summaryMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpack summary meta"})
		return
	}

	type InlineReviewCard struct {
		File          string `json:"file"`
		Line          int    `json:"line"`
		Level         string `json:"level"`
		Reason        string `json:"reason"`
		SuggestedCode string `json:"suggested_code"`
	}

	cards := make([]InlineReviewCard, 0, len(task.Comments))
	for _, comment := range task.Comments {
		cards = append(cards, InlineReviewCard{
			File:          comment.FilePath,
			Line:          comment.LineNumber,
			Level:         comment.Severity,
			Reason:        comment.Reason,
			SuggestedCode: comment.SuggestedCode,
		})
	}
	summaryMap["cards"] = cards

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "PR review analyzed successfully",
		"data": gin.H{
			"id":         task.ID,
			"pr_url":     task.GitHubURL,
			"status":     task.Status,
			"created_at": task.CreatedAt,
			"summary":    summaryMap,
		},
	})
}

func (ctrl *ReviewController) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := ctrl.reviewService.GetReviewDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review task not found"})
		return
	}

	var summaryMap map[string]interface{}
	if err := json.Unmarshal([]byte(task.Summary), &summaryMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unpack summary payload"})
		return
	}

	type InlineReviewCard struct {
		File          string `json:"file"`
		Line          int    `json:"line"`
		Level         string `json:"level"`
		Reason        string `json:"reason"`
		SuggestedCode string `json:"suggested_code"`
	}

	cards := make([]InlineReviewCard, 0, len(task.Comments))
	for _, comment := range task.Comments {
		cards = append(cards, InlineReviewCard{
			File:          comment.FilePath,
			Line:          comment.LineNumber,
			Level:         comment.Severity,
			Reason:        comment.Reason,
			SuggestedCode: comment.SuggestedCode,
		})
	}
	summaryMap["cards"] = cards

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"id":         task.ID,
			"pr_url":     task.GitHubURL,
			"status":     task.Status,
			"created_at": task.CreatedAt,
			"summary":    summaryMap,
		},
	})
}

func (ctrl *ReviewController) GetDashboard(c *gin.Context) {
	dashboard, err := ctrl.reviewService.GetDashboard()
	if err != nil {
		status := http.StatusInternalServerError
		var httpErr *service.HTTPError
		if errors.As(err, &httpErr) {
			status = httpErr.StatusCode
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": dashboard,
	})
}
