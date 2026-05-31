<template>
  <div class="pr-review-page animate-fade-in">
    <!-- Top Header -->
    <div class="welcome-banner">
      <div class="welcome-text">
        <h2>{{ t[currentLang].header }}</h2>
        <p>{{ t[currentLang].subHeader }}</p>
      </div>
      <div class="system-time">
        <span class="pulse-green"></span>
        <span class="clock">{{ currentLang === 'zh' ? '扫描引擎就绪' : 'SCAN ENGINE READY' }}</span>
      </div>
    </div>

    <!-- Scanner Input Panel -->
    <div class="cyber-card input-card">
      <div class="card-glow"></div>
      <div class="card-content">
        <h3>{{ t[currentLang].inputTitle }}</h3>
        <p class="subtitle">{{ t[currentLang].inputSub }}</p>
        
        <!-- Futuristic LLM Engine Selector -->
        <div class="engine-selector">
          <span class="selector-label">{{ currentLang === 'zh' ? '审计引擎大脑:' : 'AUDIT ENGINE BRAIN:' }}</span>
          <div class="engine-tabs">
            <button 
              :class="['engine-tab-btn', { active: selectedModelId === 'deepseek-v3' }]"
              @click="switchModel('deepseek-v3')"
            >
              <span class="pulse-cyan" v-if="selectedModelId === 'deepseek-v3'"></span>
              <span>DeepSeek-V3 (Standard)</span>
            </button>
            <button 
              :class="['engine-tab-btn', { active: selectedModelId === 'deepseek-r1' }]"
              @click="switchModel('deepseek-r1')"
            >
              <span class="pulse-cyan" v-if="selectedModelId === 'deepseek-r1'"></span>
              <span>DeepSeek-R1 (Reasoning)</span>
            </button>
          </div>
        </div>

        <div class="input-group">
          <div class="url-input-wrapper" @click="focusInput">
            <el-icon class="search-icon"><Link /></el-icon>
            <input 
              ref="prUrlInputRef"
              v-model="prUrl" 
              type="text" 
              placeholder="e.g. https://github.com/google/ai-pr-reviewer/pull/42" 
              class="cyber-input"
              @keyup.enter="startAnalysis"
            />
          </div>
          <button @click="startAnalysis" :disabled="loading" class="cyber-btn">
            <span class="btn-glow"></span>
            <span class="btn-text">
              <el-icon v-if="loading" class="is-loading"><Loading /></el-icon>
              <span v-else>{{ t[currentLang].btnText }}</span>
            </span>
          </button>
        </div>
      </div>
    </div>

    <!-- Audit Result Section (Visible when reviewData exists) -->
    <div v-if="reviewData" class="results-grid">
      <!-- Left side: Radar chart & Stats -->
      <div class="results-left">
        <div class="cyber-card chart-card">
          <div class="card-glow"></div>
          <div class="card-content">
            <h3>{{ t[currentLang].radarTitle }}</h3>
            <div ref="radarChartRef" class="radar-chart-container"></div>
            
            <div class="scores-grid">
              <div class="score-pill">
                <span class="label">{{ t[currentLang].s1 }}</span>
                <span class="val highlight">{{ reviewData.summary.overall_score }}</span>
              </div>
              <div class="score-pill">
                <span class="label">{{ t[currentLang].s2 }}</span>
                <span class="val security">{{ reviewData.summary.security_score }}</span>
              </div>
              <div class="score-pill">
                <span class="label">{{ t[currentLang].s3 }}</span>
                <span class="val performance">{{ reviewData.summary.performance_score }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="cyber-card summary-card">
          <div class="card-glow"></div>
          <div class="card-content">
            <h3>{{ t[currentLang].summaryTitle }}</h3>
            <div class="summary-section">
              <h4>{{ t[currentLang].sum1 }}</h4>
              <p class="summary-text">{{ currentLang === 'zh' ? (reviewData.summary.one_sentence_summary_zh || reviewData.summary.one_sentence_summary) : reviewData.summary.one_sentence_summary }}</p>
            </div>
            
            <div class="summary-section">
              <h4>{{ t[currentLang].sum2 }}</h4>
              <div class="tags-group">
                <span v-for="mod in reviewData.summary.affected_modules" :key="mod" class="cyber-tag module-tag">
                  {{ mod }}
                </span>
              </div>
            </div>

            <div class="summary-section warning-box" v-if="reviewData.summary.breaking_changes?.length">
              <h4>⚠️ {{ t[currentLang].sum3 }}</h4>
              <ul>
                <li v-for="(bc, index) in reviewData.summary.breaking_changes" :key="bc">
                  {{ currentLang === 'zh' ? (reviewData.summary.breaking_changes_zh?.[index] || bc) : bc }}
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Right side: Code Diff & Inline Cards -->
      <div class="results-right">
        <div class="cyber-card inline-cards-holder">
          <div class="card-glow"></div>
          <div class="card-content">
            <div class="reviews-header">
              <h3>{{ t[currentLang].cardsHeader }}</h3>
              <div class="header-actions">
                <span class="badge">{{ reviewData.summary.cards?.length || 0 }} {{ t[currentLang].issuesBadge }}</span>
                
                <!-- Advanced PDF Export Button and Dropdown -->
                <div class="export-pdf-dropdown">
                  <button class="export-btn" @click.stop="showExportMenu = !showExportMenu">
                    <el-icon><Download /></el-icon>
                    <span>{{ currentLang === 'zh' ? '导出 PDF 报告' : 'EXPORT PDF' }}</span>
                    <el-icon class="arrow-icon"><ArrowDown /></el-icon>
                  </button>
                  <div class="dropdown-menu animate-fade-in" v-if="showExportMenu">
                    <button class="menu-item-btn" @click="triggerPdfExport('zh')">
                      <span class="icon">🇨🇳</span>
                      <span>中文 PDF 报告</span>
                    </button>
                    <button class="menu-item-btn" @click="triggerPdfExport('en')">
                      <span class="icon">🇺🇸</span>
                      <span>English PDF Report</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div class="cards-list">
              <div 
                v-for="(card, index) in reviewData.summary.cards" 
                :key="index" 
                :class="['review-inline-card', card.level.toLowerCase()]"
              >
                <!-- Card Header -->
                <div class="card-meta">
                  <div class="card-level">
                    <span class="dot"></span>
                    <span class="level-text">{{ card.level }}</span>
                  </div>
                  <div class="card-file-line">
                    <span class="file">{{ card.file }}</span>
                    <span class="line">{{ currentLang === 'zh' ? '行' : 'Line' }} {{ card.line }}</span>
                  </div>
                </div>

                <!-- Card Body -->
                <div class="card-reason">
                  {{ currentLang === 'zh' ? (card.reason_zh || card.reason) : card.reason }}
                </div>

                <!-- Code segment comparison box -->
                <div class="suggested-code-box">
                  <div class="code-header">
                    <span>{{ t[currentLang].suggestCode }}</span>
                    <button class="copy-btn" @click="copyCode(card.suggested_code, $event)">
                      <el-icon><DocumentCopy /></el-icon>
                      <span>{{ currentLang === 'zh' ? '复制' : 'COPY' }}</span>
                    </button>
                  </div>
                  <pre class="code-content"><code>{{ card.suggested_code }}</code></pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, inject, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'

const currentLang = inject('lang', ref('zh'))
const route = useRoute()

const t = {
  zh: {
    header: "智能 PR 评审控制台",
    subHeader: "对单个代码合并请求执行深层 AST 和大模型混合上下文安全审计",
    inputTitle: "安全 PR 代码扫描器",
    inputSub: "输入 GitHub PR 链接以开启混合上下文代码深度审计",
    btnText: "开启 AI 智能审计",
    radarTitle: "漏洞与质量度量雷达",
    s1: "综合得分",
    s2: "安全性",
    s3: "性能评分",
    summaryTitle: "PR 变更结构化总结",
    sum1: "一句话改动概述",
    sum2: "受影响模块清单",
    sum3: "检测到破坏性变更 (Breaking Changes)",
    cardsHeader: "行内智能评审批注卡片",
    issuesBadge: "个缺陷/优化点已发现",
    suggestCode: "推荐优化修复代码",
    t1: "规范度",
    t2: "安全性",
    t3: "性能",
    t4: "可读性",
    successMsg: "AI 审计成功完成！",
    fallbackMsg: "Go 后端离线。运行客户端模拟审计引擎。"
  },
  en: {
    header: "INTELLIGENT PR SCANNER",
    subHeader: "Perform deep AST & LLM hybrid context audits on individual Pull Requests",
    inputTitle: "SECURE PR CODE SCANNER",
    inputSub: "Input a GitHub Pull Request URL to perform deep AST and LLM hybrid context reviews",
    btnText: "TRIGGER AI AUDIT",
    radarTitle: "VULNERABILITY & METRICS RADAR",
    s1: "OVERALL",
    s2: "SECURITY",
    s3: "PERFORMANCE",
    summaryTitle: "PR STRUCTURAL SUMMARY",
    sum1: "ONE-SENTENCE OVERVIEW",
    sum2: "AFFECTED MODULES",
    sum3: "BREAKING CHANGES DETECTED",
    cardsHeader: "INLINE SCHEMATIC AUDIT CARDS",
    issuesBadge: "ISSUES FOUND",
    suggestCode: "SUGGESTED OPTIMIZATION",
    t1: "Standards",
    t2: "Security",
    t3: "Performance",
    t4: "Readability",
    successMsg: "AI Audit completed successfully!",
    fallbackMsg: "Go backend offline. Running client mock review engine."
  }
}

const prUrl = ref('')
const prUrlInputRef = ref(null)
const loading = ref(false)
const reviewData = ref(null)
const radarChartRef = ref(null)
let myChart = null

const focusInput = () => {
  if (prUrlInputRef.value) {
    prUrlInputRef.value.focus()
  }
}

const selectedModelId = ref('deepseek-v3')
const showExportMenu = ref(false)

const switchModel = (id) => {
  selectedModelId.value = id
  let baseUrl = 'https://api.deepseek.com/v1'
  let enforcedModel = 'deepseek-chat'
  if (id === 'deepseek-r1') {
    enforcedModel = 'deepseek-reasoner'
  }
  localStorage.setItem('baseUrl', baseUrl)
  localStorage.setItem('enforcedModel', enforcedModel)
  localStorage.setItem('selectedModelId', id)
}

const currentTime = ref(new Date().toLocaleTimeString())
setInterval(() => {
  currentTime.value = new Date().toLocaleTimeString()
}, 1000)

// Watch global language to re-initialize ECharts axis labels seamlessly
watch(currentLang, () => {
  if (reviewData.value) {
    nextTick(() => {
      initRadarChart()
    })
  }
})

// Standard Mock Response matching GORM MySQL structures
const mockAuditPayload = {
  id: 1024,
  pr_url: "https://github.com/google/ai-pr-reviewer/pull/42",
  summary: {
    one_sentence_summary: "Refactored JWT signature algorithm from HS256 to RS256, deprecating weak credentials and establishing env configuration loaders.",
    one_sentence_summary_zh: "重构 JWT 签名算法，从对称 HS256 升级为非对称 RS256 算法，废弃了弱密钥凭据并建立了环境变量加载机制。",
    affected_modules: ["auth/jwt.go", "auth/service.go", "config/config.go"],
    breaking_changes: [
      "Changed JWT token decoding signing keys from symmetric to asymmetric public key structure.",
      "Removed legacy hardcoded fallback authentication hooks."
    ],
    breaking_changes_zh: [
      "将 JWT 令牌解码签名密钥从对称结构更改为非对称公钥结构。",
      "移除了遗留的硬编码回退身份验证钩子。"
    ],
    overall_score: 87.5,
    normative_score: 92.0,
    security_score: 68.0,
    performance_score: 89.0,
    readability_score: 88.0,
    cards: [
      {
        file: "auth/jwt.go",
        line: 24,
        level: "CRITICAL",
        reason: "Hardcoded symmetric key detected inside production build! This allows attackers to forge authorization tokens if they download the compiled binary.",
        reason_zh: "生产构建中检测到硬编码的对称密钥！如果攻击者下载并反编译二进制文件，这将允许他们任意伪造合规的授权 Token。",
        suggested_code: "privateKeyPEM := os.Getenv(\"JWT_PRIVATE_KEY_PEM\")\nif privateKeyPEM == \"\" {\n    log.Fatal(\"Critical vulnerability: JWT_PRIVATE_KEY_PEM is empty\")\n}"
      },
      {
        file: "auth/service.go",
        line: 45,
        level: "WARNING",
        reason: "Stripe connection handles network events directly on the router's thread pools without context limits. Unresponsive API providers will exhaust Gin request pipelines.",
        reason_zh: "Stripe 支付连接在路由器的线程池上直接处理网络事件，未设置上下文超时限制。如果第三方 API 无响应，将导致 Gin 请求管道迅速耗尽。",
        suggested_code: "ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)\ndefer cancel()\nresp, err := stripeClient.DoWithContext(ctx, payload)"
      },
      {
        file: "config/config.go",
        line: 12,
        level: "SUGGESTION",
        reason: "Structure tags missing in struct definitions. Use standard JSON and DB bindings to enforce compliance with GORM frameworks.",
        reason_zh: "结构体定义中缺失 Tag 字段。推荐使用标准的 JSON 和 DB 字段绑定以增强与 GORM 框架的合规性。",
        suggested_code: "type Config struct {\n    ServerPort string `json:\"server_port\" gorm:\"column:server_port\"`\n    DatabaseURL string `json:\"database_url\" gorm:\"column:db_url\"`\n}"
      }
    ]
  }
}

const startAnalysis = () => {
  if (!prUrl.value) {
    ElMessage.warning(currentLang.value === 'zh' ? '请输入有效的 GitHub PR 链接。' : 'Please input a valid GitHub PR URL.')
    return
  }
  loading.value = true
  
  setTimeout(async () => {
    try {
      const githubToken = localStorage.getItem('githubToken') || ''
      const apiKey = localStorage.getItem('apiKey') || ''
      const baseUrl = localStorage.getItem('baseUrl') || ''
      const enforcedModel = localStorage.getItem('enforcedModel') || ''

      const res = await fetch('http://localhost:8080/api/v1/review/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          github_url: prUrl.value,
          pr_url: prUrl.value,
          github_token: githubToken,
          api_key: apiKey,
          base_url: baseUrl,
          model: enforcedModel
        })
      })
      const payload = await res.json()
      if (res.ok) {
        reviewData.value = payload.data
        ElMessage.success(t[currentLang.value].successMsg)
      } else {
        ElMessage.error(currentLang.value === 'zh' ? `审计失败: ${payload.error || '服务器错误'}` : `Audit Failed: ${payload.error || 'Server error'}`)
      }
      
      nextTick(() => {
        initRadarChart()
      })
    } catch (e) {
      reviewData.value = mockAuditPayload
      ElMessage.info(t[currentLang.value].fallbackMsg)
      nextTick(() => {
        initRadarChart()
      })
    } finally {
      loading.value = false
    }
  }, 1200)
}

const initRadarChart = (isPrint = false) => {
  if (!radarChartRef.value) return
  
  if (myChart) {
    myChart.dispose()
  }

  myChart = echarts.init(radarChartRef.value)
  const option = {
    backgroundColor: 'transparent',
    animation: !isPrint,
    tooltip: {
      show: !isPrint,
      trigger: 'item',
      backgroundColor: 'rgba(9, 15, 29, 0.95)',
      borderColor: '#00f0ff',
      borderWidth: 1,
      textStyle: {
        color: '#e2e8f0',
        fontFamily: 'JetBrains Mono',
        fontSize: 12
      },
      formatter: function (params) {
        const indicators = [
          t[currentLang.value].t1,
          t[currentLang.value].t2,
          t[currentLang.value].t3,
          t[currentLang.value].t4
        ];
        let html = `<div style="padding: 2px; font-weight: bold; border-bottom: 1px solid rgba(0, 240, 255, 0.2); margin-bottom: 6px; color: #00f0ff;">${currentLang.value === 'zh' ? '当前 PR 审计指标' : 'Current PR Metrics'}</div>`;
        for (let i = 0; i < indicators.length; i++) {
          html += `<div style="display: flex; justify-content: space-between; gap: 20px; margin: 4px 0;">
            <span style="color: #94a3b8;">${indicators[i]}:</span>
            <span style="font-weight: bold; color: #00ff66;">${params.value[i]}%</span>
          </div>`;
        }
        return html;
      }
    },
    radar: {
      indicator: [
        { name: t[currentLang.value].t1, max: 100 },
        { name: t[currentLang.value].t2, max: 100 },
        { name: t[currentLang.value].t3, max: 100 },
        { name: t[currentLang.value].t4, max: 100 }
      ],
      axisName: {
        color: isPrint ? '#0f172a' : '#94a3b8',
        fontSize: 12,
        fontFamily: 'JetBrains Mono',
        fontWeight: isPrint ? 'bold' : 'normal'
      },
      splitArea: {
        areaStyle: {
          color: isPrint ? ['rgba(0, 0, 0, 0.01)', 'rgba(0, 0, 0, 0.03)'] : ['rgba(0, 240, 255, 0.02)', 'rgba(0, 240, 255, 0.05)']
        }
      },
      splitLine: {
        lineStyle: {
          color: isPrint ? 'rgba(15, 23, 42, 0.12)' : 'rgba(0, 240, 255, 0.15)'
        }
      },
      axisLine: {
        lineStyle: {
          color: isPrint ? 'rgba(15, 23, 42, 0.12)' : 'rgba(0, 240, 255, 0.15)'
        }
      }
    },
    series: [
      {
        name: 'Vulnerability Scoring Matrix',
        type: 'radar',
        data: [
          {
            value: [
              reviewData.value.summary.normative_score,
              reviewData.value.summary.security_score,
              reviewData.value.summary.performance_score,
              reviewData.value.summary.readability_score
            ],
            name: currentLang.value === 'zh' ? '安全质量维度' : 'Vulnerabilities',
            symbol: 'circle',
            symbolSize: 6,
            label: {
              show: true,
              color: isPrint ? '#0f172a' : '#00f0ff',
              fontFamily: 'JetBrains Mono',
              fontSize: 10.5,
              fontWeight: 'bold',
              formatter: function (params) {
                return params.value + '%';
              }
            },
            itemStyle: {
              color: isPrint ? '#0369a1' : '#00f0ff'
            },
            areaStyle: {
              color: isPrint ? 'rgba(3, 105, 161, 0.15)' : 'rgba(0, 240, 255, 0.25)'
            },
            lineStyle: {
              color: isPrint ? '#0369a1' : '#00f0ff',
              width: 2,
              shadowBlur: isPrint ? 0 : 8,
              shadowColor: '#00f0ff'
            }
          }
        ]
      }
    ]
  }

  myChart.setOption(option)
  window.addEventListener('resize', () => myChart && myChart.resize())
}

const copyCode = (text, event) => {
  navigator.clipboard.writeText(text)
  const btn = event.currentTarget
  const originalHtml = btn.innerHTML
  btn.innerHTML = `<span>${currentLang.value === 'zh' ? '复制成功!' : 'SUCCESS!'}</span>`
  btn.style.color = '#00ff66'
  setTimeout(() => {
    btn.innerHTML = originalHtml
    btn.style.color = ''
  }, 1500)
}

const triggerPdfExport = (targetLang) => {
  const originalLang = currentLang.value
  currentLang.value = targetLang
  showExportMenu.value = false
  
  // Add print mode selector class to body
  document.body.classList.add('printing-' + targetLang)
  
  nextTick(() => {
    // Redraw ECharts immediately with target language text and high contrast print theme
    initRadarChart(true)
    
    // Tiny delay to ensure charts are fully rendered before printing
    setTimeout(() => {
      window.print()
      
      // Restore states
      document.body.classList.remove('printing-' + targetLang)
      currentLang.value = originalLang
      
      nextTick(() => {
        initRadarChart()
      })
    }, 450)
  })
}

const closeExportMenu = () => {
  showExportMenu.value = false
}

onMounted(async () => {
  window.addEventListener('click', closeExportMenu)

  if (localStorage.getItem('selectedModelId')) {
    selectedModelId.value = localStorage.getItem('selectedModelId')
  } else {
    switchModel('deepseek-v3')
  }

  const taskId = route.query.id
  if (taskId) {
    loading.value = true
    try {
      const res = await fetch(`http://localhost:8080/api/v1/review/${taskId}`)
      if (res.ok) {
        const payload = await res.json()
        reviewData.value = payload.data
        prUrl.value = payload.data.pr_url
        nextTick(() => {
          initRadarChart()
        })
      } else {
        ElMessage.error(currentLang.value === 'zh' ? '获取历史审计记录失败。' : 'Failed to fetch historical audit logs.')
      }
    } catch (e) {
      // Fallback
      reviewData.value = mockAuditPayload
      prUrl.value = mockAuditPayload.pr_url
      nextTick(() => {
        initRadarChart()
      })
    } finally {
      loading.value = false
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('click', closeExportMenu)
})
</script>

<style scoped>
.pr-review-page {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

/* Banner Design */
.welcome-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-text h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 1.5px;
  color: #f8fafc;
}

.welcome-text p {
  margin: 5px 0 0 0;
  color: #64748b;
  font-size: 13px;
}

.system-time {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(13, 20, 35, 0.5);
  border: 1px solid rgba(0, 240, 255, 0.1);
  padding: 6px 14px;
  border-radius: 20px;
}

.pulse-green {
  width: 6px;
  height: 6px;
  background-color: #00f0ff;
  border-radius: 50%;
  box-shadow: 0 0 8px #00f0ff;
}

.clock {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  color: #00f0ff;
}

/* Cyber Card base styles */
.cyber-card {
  position: relative;
  background: rgba(15, 23, 42, 0.65);
  border: 1px solid rgba(0, 240, 255, 0.15);
  border-radius: 12px;
  overflow: hidden;
  backdrop-filter: blur(10px);
}

.card-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #00f0ff, transparent);
  opacity: 0.8;
}

.card-content {
  padding: 25px;
}

.card-content h3 {
  margin: 0 0 10px 0;
  font-size: 16px;
  font-weight: bold;
  letter-spacing: 1.2px;
  color: #f1f5f9;
}

.card-content .subtitle {
  margin: 0 0 20px 0;
  color: #64748b;
  font-size: 13px;
}

/* User URL Input Panel */
.input-group {
  display: flex;
  gap: 15px;
}

.url-input-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  height: 48px;
  background: rgba(10, 15, 30, 0.8);
  border: 1px solid rgba(0, 240, 255, 0.2);
  border-radius: 8px;
  padding: 0 15px;
  transition: all 0.3s;
  cursor: text;
}

.url-input-wrapper:focus-within {
  border-color: #00f0ff;
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.15);
}

.search-icon {
  color: #00f0ff;
  font-size: 18px;
  margin-right: 12px;
  display: inline-flex !important;
  align-items: center;
  justify-content: center;
}

.cyber-input {
  background: transparent;
  border: none;
  color: #f8fafc;
  font-size: 14px;
  flex: 1;
  min-width: 0;
  height: 100%;
  padding: 0;
  margin: 0;
  outline: none;
  font-family: inherit;
  line-height: 48px; /* Force browser to vertically center text exactly in the middle of 48px */
  box-sizing: border-box !important;
}

.cyber-btn {
  position: relative;
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  padding: 0 25px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.cyber-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.4);
}

.cyber-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%);
  filter: blur(8px);
  opacity: 0.4;
  z-index: 1;
}

.btn-text {
  position: relative;
  color: #ffffff;
  font-weight: bold;
  font-size: 13px;
  letter-spacing: 1px;
  z-index: 2;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Results layout grid */
.results-grid {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 25px;
  align-items: start;
}

.results-left {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.radar-chart-container {
  height: 280px;
  width: 100%;
  margin-bottom: 15px;
}

.scores-grid {
  display: flex;
  gap: 15px;
}

.score-pill {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(13, 20, 35, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 10px;
}

.score-pill .label {
  font-size: 9px;
  font-family: 'JetBrains Mono', monospace;
  color: #64748b;
}

.score-pill .val {
  font-size: 18px;
  font-weight: 900;
  margin-top: 4px;
}

.score-pill .val.highlight {
  color: #00f0ff;
  text-shadow: 0 0 10px rgba(0, 240, 255, 0.25);
}

.score-pill .val.security {
  color: #ef4444; /* Low security warning color */
}

.score-pill .val.performance {
  color: #00ff66;
}

/* PR Summary layout */
.summary-card {
  flex: 1;
}

.summary-section {
  margin-bottom: 20px;
}

.summary-section h4 {
  font-size: 11px;
  font-family: 'JetBrains Mono', monospace;
  color: #475569;
  letter-spacing: 1px;
  margin: 0 0 8px 0;
}

.summary-text {
  font-size: 13.5px;
  line-height: 1.5;
  color: #cbd5e1;
  margin: 0;
}

.tags-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.cyber-tag {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  padding: 4px 10px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #94a3b8;
}

.module-tag {
  color: #00f0ff;
  background: rgba(0, 240, 255, 0.05);
  border-color: rgba(0, 240, 255, 0.15);
}

.warning-box {
  background: rgba(239, 68, 68, 0.05);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 8px;
  padding: 15px;
}

.warning-box h4 {
  color: #f87171 !important;
}

.warning-box ul {
  margin: 8px 0 0 0;
  padding-left: 20px;
  font-size: 12.5px;
  color: #fca5a5;
  line-height: 1.6;
}

/* Inline Card Reviews */
.reviews-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 10px;
}

.reviews-header .badge {
  font-size: 10px;
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.cards-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.review-inline-card {
  background: rgba(10, 15, 30, 0.7);
  border-left: 4px solid #cbd5e1;
  border-radius: 4px 8px 8px 4px;
  padding: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.03);
  border-right: 1px solid rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.review-inline-card.critical {
  border-left-color: #ef4444;
  background: linear-gradient(90deg, rgba(239, 68, 68, 0.03) 0%, transparent 100%);
}

.review-inline-card.warning {
  border-left-color: #eab308;
  background: linear-gradient(90deg, rgba(234, 179, 8, 0.03) 0%, transparent 100%);
}

.review-inline-card.suggestion {
  border-left-color: #00ff66;
  background: linear-gradient(90deg, rgba(0, 255, 102, 0.03) 0%, transparent 100%);
}

.card-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.card-level {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  font-weight: bold;
}

.critical .card-level { color: #ef4444; }
.warning .card-level { color: #eab308; }
.suggestion .card-level { color: #00ff66; }

.card-level .dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.critical .dot { background-color: #ef4444; box-shadow: 0 0 6px #ef4444; }
.warning .dot { background-color: #eab308; box-shadow: 0 0 6px #eab308; }
.suggestion .dot { background-color: #00ff66; box-shadow: 0 0 6px #00ff66; }

.card-file-line {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
}

.card-file-line .file {
  color: #64748b;
  margin-right: 8px;
}

.card-file-line .line {
  color: #00f0ff;
  background: rgba(0, 240, 255, 0.08);
  padding: 1px 6px;
  border-radius: 3px;
}

.card-reason {
  font-size: 13px;
  line-height: 1.5;
  color: #cbd5e1;
  margin-bottom: 15px;
}

/* Suggested Code Fragment Box */
.suggested-code-box {
  background: #060913;
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  overflow: hidden;
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  font-family: 'JetBrains Mono', monospace;
  font-size: 10px;
  color: #475569;
}

.copy-btn {
  background: transparent;
  border: none;
  color: #00f0ff;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 10px;
  font-weight: bold;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.2s;
}

.copy-btn:hover {
  background: rgba(0, 240, 255, 0.1);
}

.code-content {
  margin: 0;
  padding: 12px;
  overflow-x: auto;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  line-height: 1.6;
  color: #38bdf8; /* Soft blue code styling */
}

/* Animations */
.animate-fade-in {
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Engine Selector Tabs style */
.engine-selector {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px dashed rgba(255, 255, 255, 0.08);
}

.selector-label {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  color: #475569;
  letter-spacing: 1px;
  white-space: nowrap;
}

.engine-tabs {
  display: flex;
  background: rgba(9, 15, 29, 0.8);
  border: 1px solid rgba(0, 240, 255, 0.15);
  border-radius: 8px;
  padding: 3px;
  gap: 4px;
}

.engine-tab-btn {
  background: transparent;
  border: 1px solid transparent;
  color: #94a3b8;
  font-family: 'Outfit', sans-serif;
  font-size: 12.5px;
  font-weight: bold;
  padding: 6px 16px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.engine-tab-btn:hover {
  color: #00f0ff;
}

.engine-tab-btn.active {
  background: linear-gradient(135deg, rgba(0, 240, 255, 0.15) 0%, rgba(112, 0, 255, 0.05) 100%);
  border-color: rgba(0, 240, 255, 0.4);
  color: #ffffff;
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.1);
}

.pulse-cyan {
  width: 6px;
  height: 6px;
  background-color: #00f0ff;
  border-radius: 50%;
  box-shadow: 0 0 8px #00f0ff;
}

@media (max-width: 576px) {
  .engine-selector {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  .engine-tabs {
    width: 100%;
  }
  .engine-tab-btn {
    flex: 1;
    justify-content: center;
    font-size: 11.5px;
    padding: 6px 8px;
  }
}

/* Screen-only Dropdown Styles */
.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.export-pdf-dropdown {
  position: relative;
  display: inline-block;
}

.export-btn {
  background: rgba(0, 240, 255, 0.05);
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: #00f0ff;
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  font-weight: bold;
  padding: 6px 14px;
  border-radius: 20px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all 0.3s;
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.1);
}

.export-btn:hover {
  background: rgba(0, 240, 255, 0.15);
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.4);
  transform: translateY(-1px);
}

.export-btn .arrow-icon {
  font-size: 10px;
  margin-left: 2px;
  transition: transform 0.3s;
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: rgba(13, 20, 35, 0.95);
  border: 1px solid rgba(0, 240, 255, 0.3);
  border-radius: 8px;
  padding: 6px;
  width: 170px;
  box-shadow: 0 5px 25px rgba(0, 0, 0, 0.5), 0 0 15px rgba(0, 240, 255, 0.15);
  backdrop-filter: blur(12px);
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.menu-item-btn {
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #cbd5e1;
  padding: 8px 12px;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  text-align: left;
  width: 100%;
  transition: all 0.2s;
}

.menu-item-btn:hover {
  background: rgba(0, 240, 255, 0.1);
  color: #00f0ff;
}

.menu-item-btn .icon {
  font-size: 14px;
}
</style>

<style>
/* ==========================================
   HIGH-FIDELITY VECTOR PDF PRINT STYLES 
   ========================================== */
@media print {
  /* 0. Universal print reset: disable animations, transitions, and force full opacity */
  *, *::before, *::after {
    animation: none !important;
    transition: none !important;
    opacity: 1 !important;
    box-shadow: none !important;
    text-shadow: none !important;
    backdrop-filter: none !important;
    filter: none !important;
  }

  /* 1. Hide interactive/navigation elements entirely */
  .cyber-sidebar,
  .cyber-navbar,
  .welcome-banner,
  .input-card,
  .export-pdf-dropdown,
  .copy-btn,
  .card-glow,
  .cyber-grid-overlay {
    display: none !important;
  }

  /* 2. Format basic layout for paper and avoid collapsing containers */
  #app,
  body,
  html,
  .cyber-container,
  .cyber-main,
  .main-layout,
  .pr-review-page {
    background: #ffffff !important;
    color: #000000 !important;
    height: auto !important;
    min-height: 0 !important;
    overflow: visible !important;
    padding: 0 !important;
    margin: 0 !important;
    display: block !important;
  }

  /* 3. Restructure grid columns into clean vertical blocks */
  .results-grid {
    display: block !important;
    width: 100% !important;
  }

  .results-left,
  .results-right {
    width: 100% !important;
    display: block !important;
    float: none !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  /* 4. White background for cards, clear dark text */
  .cyber-card {
    background: #ffffff !important;
    border: 1px solid #cbd5e1 !important;
    border-radius: 8px !important;
    box-shadow: none !important;
    color: #000000 !important;
    margin-bottom: 25px !important;
    padding: 20px !important;
    page-break-inside: avoid !important;
    backdrop-filter: none !important;
    filter: none !important;
  }

  .card-content {
    padding: 0 !important;
  }

  /* 5. Formal Typography overrides */
  h3 {
    font-size: 18px !important;
    color: #0f172a !important;
    border-bottom: 2px solid #e2e8f0 !important;
    padding-bottom: 8px !important;
    margin-bottom: 15px !important;
  }

  h4 {
    font-size: 13px !important;
    color: #475569 !important;
    margin-bottom: 8px !important;
  }

  .summary-text {
    color: #334155 !important;
    font-size: 13px !important;
    line-height: 1.6 !important;
  }

  .cyber-tag {
    background: #f1f5f9 !important;
    border: 1px solid #cbd5e1 !important;
    color: #0f172a !important;
    box-shadow: none !important;
  }

  .warning-box {
    background: #fef2f2 !important;
    border: 1px solid #fca5a5 !important;
    color: #991b1b !important;
    padding: 15px !important;
    border-radius: 6px !important;
  }

  .warning-box h4,
  .warning-box ul,
  .warning-box li {
    color: #991b1b !important;
  }

  /* 6. Score pills design for printing */
  .score-pill {
    background: #f8fafc !important;
    border: 1px solid #cbd5e1 !important;
    box-shadow: none !important;
  }

  .score-pill .label {
    color: #64748b !important;
  }

  .score-pill .val {
    color: #0f172a !important;
    text-shadow: none !important;
  }

  .score-pill .val.highlight {
    color: #0369a1 !important; /* High contrast blue */
  }

  .score-pill .val.performance {
    color: #047857 !important; /* High contrast green */
  }

  .score-pill .val.security {
    color: #b91c1c !important; /* High contrast red */
  }

  /* 7. Timeline & Inline Cards constraints */
  .review-inline-card {
    background: #f8fafc !important;
    border: 1px solid #e2e8f0 !important;
    border-left: 5px solid #64748b !important;
    color: #0f172a !important;
    page-break-inside: avoid !important;
  }
  
  .review-inline-card.critical {
    border-left-color: #dc2626 !important; /* Formal red border */
  }
  
  .review-inline-card.warning {
    border-left-color: #d97706 !important; /* Formal amber border */
  }
  
  .review-inline-card.suggestion {
    border-left-color: #059669 !important; /* Formal green border */
  }

  .card-level {
    font-size: 11px !important;
    font-weight: bold !important;
  }

  .critical .card-level { color: #dc2626 !important; }
  .warning .card-level { color: #d97706 !important; }
  .suggestion .card-level { color: #059669 !important; }

  .card-meta .dot {
    display: none !important; /* Hide neon dots inside PDF */
  }

  .card-file-line .file {
    color: #475569 !important;
    font-weight: bold !important;
  }

  .card-file-line .line {
    color: #0f172a !important;
    background: #e2e8f0 !important;
    border: 1px solid #cbd5e1 !important;
    font-weight: bold !important;
    padding: 1px 6px !important;
  }

  .card-reason {
    color: #1e293b !important;
    font-size: 13.5px !important;
    line-height: 1.6 !important;
  }

  /* Code suggestion boxes */
  .suggested-code-box {
    background: #f8fafc !important;
    border: 1px solid #cbd5e1 !important;
    border-radius: 6px !important;
    overflow: hidden !important;
    page-break-inside: avoid !important;
  }

  .code-header {
    background: #e2e8f0 !important;
    border-bottom: 1px solid #cbd5e1 !important;
    color: #334155 !important;
  }

  pre.code-content {
    background: #f8fafc !important;
    border: none !important;
    padding: 12px !important;
    margin: 0 !important;
    white-space: pre-wrap !important;
  }

  pre.code-content code {
    color: #0f172a !important;
    font-family: "JetBrains Mono", Courier, monospace !important;
    font-size: 11.5px !important;
    line-height: 1.5 !important;
  }

  /* 8. ECharts Radar Chart Page constraints */
  .radar-chart-container {
    height: 280px !important;
    width: 100% !important;
    border: 1px solid #cbd5e1 !important;
    background: #f8fafc !important;
    page-break-inside: avoid !important;
  }

  /* 9. Insert Language-Specific Formal Headers */
  body.printing-zh::before {
    content: "AI 代码审计与软件合规性评估报告";
    display: block;
    font-size: 24px;
    font-weight: 800;
    text-align: center;
    margin-bottom: 25px;
    border-bottom: 3px double #0f172a;
    padding-bottom: 12px;
    color: #0f172a;
    font-family: "SimSun", "Microsoft YaHei", sans-serif !important;
  }

  body.printing-en::before {
    content: "AI Code Security Audit & Compliance Assessment Report";
    display: block;
    font-size: 24px;
    font-weight: 800;
    text-align: center;
    margin-bottom: 25px;
    border-bottom: 3px double #0f172a;
    padding-bottom: 12px;
    color: #0f172a;
    font-family: Arial, sans-serif !important;
  }
}
</style>
