<template>
  <div class="dashboard-page animate-fade-in">
    <!-- Top Header -->
    <div class="welcome-banner">
      <div class="welcome-text">
        <h2>{{ t[currentLang].header }}</h2>
        <p>{{ t[currentLang].subHeader }}</p>
      </div>
      <div class="system-time">
        <span class="pulse-green"></span>
        <span class="clock">SYSTEM MONITORING</span>
      </div>
    </div>

    <!-- Statistics Grid -->
    <div class="stats-grid">
      <div v-for="stat in stats" :key="stat.label" class="cyber-card stat-card">
        <div class="card-glow"></div>
        <div class="card-content stat-box">
          <span class="stat-label">{{ currentLang === 'zh' ? stat.labelZh : stat.label }}</span>
          <h3 :class="['stat-val', stat.colorClass]">{{ stat.value }}</h3>
          <p class="stat-trend">{{ currentLang === 'zh' ? stat.trendZh : stat.trend }}</p>
        </div>
      </div>
    </div>

    <!-- Analytics Charts & Overview -->
    <div class="analytics-grid">
      <!-- Historical Radar -->
      <div class="cyber-card chart-card">
        <div class="card-glow"></div>
        <div class="card-content">
          <h3>{{ t[currentLang].radarTitle }}</h3>
          <p class="subtitle">{{ t[currentLang].radarSub }}</p>
          <div ref="radarChartRef" class="radar-chart-container"></div>
        </div>
      </div>

      <!-- History Table Log -->
      <div class="cyber-card logs-card">
        <div class="card-glow"></div>
        <div class="card-content">
          <div>
            <h3>{{ t[currentLang].tableTitle }}</h3>
            <p class="subtitle">{{ t[currentLang].tableSub }}</p>
            
            <div class="logs-table-wrapper">
              <table class="cyber-table">
                <thead>
                  <tr>
                    <th>{{ t[currentLang].col1 }}</th>
                    <th>{{ t[currentLang].col2 }}</th>
                    <th>{{ t[currentLang].col3 }}</th>
                    <th>{{ t[currentLang].col4 }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr 
                    v-for="log in displayedReviews" 
                    :key="log.id" 
                    class="clickable-row"
                    @click="goToReview(log.id)"
                  >
                    <td>
                      <div class="repo-pr-name">
                        <el-icon class="link-icon"><Link /></el-icon>
                        <span>{{ log.pr_url }}</span>
                      </div>
                    </td>
                    <td>
                      <span :class="['score-val', getScoreClass(log.score)]">{{ log.score }} / 100</span>
                    </td>
                    <td><span class="date-txt">{{ log.date }}</span></td>
                    <td>
                      <span class="status-pill"><span class="dot"></span> {{ log.status.toUpperCase() }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Symmetrical Cyber Pagination -->
          <div class="pagination-wrapper" v-if="recentReviews.length > pageSize">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              layout="prev, pager, next, total"
              :total="recentReviews.length"
              background
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, inject, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'

const currentLang = inject('lang', ref('zh'))

const currentPage = ref(1)
const pageSize = ref(4)

const displayedReviews = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return recentReviews.value.slice(start, end)
})

const t = {
  zh: {
    header: "大盘监控中心",
    subHeader: "项目平均健康指标、安全规范趋势以及扫描足迹大盘",
    radarTitle: "全局代码质量指标",
    radarSub: "所有历史代码审查得分的综合平均值向量",
    tableTitle: "近期 PR 自动审计日志",
    tableSub: "保存在 MySQL 数据库中的最近代码扫描足迹",
    t1: "规范度",
    t2: "安全性",
    t3: "性能",
    t4: "可读性",
    col1: "代码合并请求节点",
    col2: "安全规范评分",
    col3: "审计时间",
    col4: "状态"
  },
  en: {
    header: "METRICS CORE DASHBOARD",
    subHeader: "Aggregate repo health, security compliance trends, and scanning footprints",
    radarTitle: "GLOBAL CODE COMPLIANCE METRICS",
    radarSub: "Aggregated average score vectors computed from all historical audits",
    tableTitle: "RECENT PR CODE AUDIT LOGS",
    tableSub: "Latest pull requests scanned and recorded inside GORM MySQL database",
    t1: "Standards",
    t2: "Security",
    t3: "Performance",
    t4: "Readability",
    col1: "Pull Request Node",
    col2: "Compliance Score",
    col3: "Audit Date",
    col4: "Status"
  }
}

const stats = ref([
  { label: "TOTAL REVIEWS COMPLETED", labelZh: "累计审计PR完成数", value: "256", trend: "▲ +18% from last week", trendZh: "▲ 比上周增长 +18%", colorClass: "cyan" },
  { label: "AVERAGE COMPLIANCE SCORE", labelZh: "平均代码健康评分", value: "86.5", trend: "▲ +1.2% quality gains", trendZh: "▲ 比上周质量提升 +1.2%", colorClass: "green" },
  { label: "CRITICAL VULNERABILITIES DETECTED", labelZh: "累计发现的致命漏洞", value: "18", trend: "▼ -34% suppression trend", trendZh: "▼ 历史漏洞降低 -34%", colorClass: "red" },
  { label: "COGNITIVE CONTEXT RETRIEVALS", labelZh: "关联上下文拉取次数", value: "1,482", trend: "▲ +156 dependencies bound", trendZh: "▲ 新解析关联依赖 +156", colorClass: "yellow" }
])

const recentReviews = ref([
  { id: 1, pr_url: "https://github.com/gin-gonic/gin/pull/3534", score: 94.0, status: "completed", date: "2026-05-29" },
  { id: 2, pr_url: "https://github.com/google/ai-pr-reviewer/pull/42", score: 82.5, status: "completed", date: "2026-05-29" },
  { id: 3, pr_url: "https://github.com/JessieXu0329/qingjiaotongxing/pull/1", score: 89.0, status: "completed", date: "2026-05-28" },
  { id: 4, pr_url: "https://github.com/example/repo/pull/18", score: 75.0, status: "completed", date: "2026-05-27" }
])

const router = useRouter()

const goToReview = (id) => {
  router.push({ path: '/review', query: { id } })
}

const radarChartRef = ref(null)
let myChart = null

// Watch global language to re-initialize ECharts axis labels seamlessly
watch(currentLang, () => {
  initRadarChart()
})

onMounted(async () => {
  try {
    const res = await fetch('http://localhost:8080/api/v1/review/dashboard')
    if (res.ok) {
      const payload = await res.json()
      const data = payload.data
      stats.value[0].value = data.total_reviews.toString()
      stats.value[1].value = data.average_score.toFixed(1)
      stats.value[2].value = data.critical_issues.toString()
      if (data.recent_reviews && data.recent_reviews.length > 0) {
        recentReviews.value = data.recent_reviews
      }
    }
  } catch (e) {
    // Graceful fallback to offline stats
  }

  nextTick(() => {
    initRadarChart()
  })
})

const initRadarChart = () => {
  if (!radarChartRef.value) return
  
  if (myChart) {
    myChart.dispose()
  }

  myChart = echarts.init(radarChartRef.value)
  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      show: true,
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
        let html = `<div style="padding: 2px; font-weight: bold; border-bottom: 1px solid rgba(0, 240, 255, 0.2); margin-bottom: 6px; color: #00f0ff;">${currentLang.value === 'zh' ? '全局平均指标' : 'Global Avg Metrics'}</div>`;
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
        color: '#94a3b8',
        fontSize: 12,
        fontFamily: 'JetBrains Mono'
      },
      splitArea: {
        areaStyle: {
          color: ['rgba(0, 240, 255, 0.02)', 'rgba(0, 240, 255, 0.05)']
        }
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(0, 240, 255, 0.15)'
        }
      },
      axisLine: {
        lineStyle: {
          color: 'rgba(0, 240, 255, 0.15)'
        }
      }
    },
    series: [
      {
        name: 'Aggregated Repo Metrics',
        type: 'radar',
        data: [
          {
            value: [88.5, 82.0, 85.0, 90.5], // Average quality vectors
            name: currentLang.value === 'zh' ? '历史平均数据' : 'Historical Averages',
            symbol: 'circle',
            symbolSize: 6,
            label: {
              show: true,
              color: '#00f0ff',
              fontFamily: 'JetBrains Mono',
              fontSize: 10,
              formatter: function (params) {
                return params.value + '%';
              }
            },
            itemStyle: {
              color: '#00f0ff'
            },
            areaStyle: {
              color: 'rgba(0, 240, 255, 0.2)'
            },
            lineStyle: {
              color: '#00f0ff',
              width: 2,
              shadowBlur: 8,
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

const getScoreClass = (score) => {
  if (score >= 90) return 'green'
  if (score >= 80) return 'cyan'
  if (score >= 70) return 'yellow'
  return 'red'
}
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

/* Header */
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

/* Stats Cards Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

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
  margin: 0;
}

.card-content .subtitle {
  margin: 0 0 20px 0;
  color: #64748b;
  font-size: 13px;
}

.stat-box {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 105px;
}

.stat-label {
  font-family: 'JetBrains Mono', monospace;
  font-size: 9.5px;
  color: #475569;
  letter-spacing: 0.8px;
}

.stat-val {
  font-size: 26px;
  font-weight: 900;
  margin: 6px 0;
}

.stat-val.cyan { color: #00f0ff; text-shadow: 0 0 10px rgba(0, 240, 255, 0.25); }
.stat-val.green { color: #00ff66; text-shadow: 0 0 10px rgba(0, 255, 102, 0.25); }
.stat-val.red { color: #ef4444; text-shadow: 0 0 10px rgba(239, 68, 68, 0.25); }
.stat-val.yellow { color: #eab308; text-shadow: 0 0 10px rgba(234, 179, 8, 0.25); }

.stat-trend {
  font-size: 11px;
  color: #64748b;
}

/* Analytics Section */
.analytics-grid {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 25px;
}

.radar-chart-container {
  height: 280px;
  width: 100%;
}

/* Logs Table */
.logs-table-wrapper {
  overflow-x: auto;
}

.cyber-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.cyber-table th {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  color: #475569;
  letter-spacing: 1px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.cyber-table td {
  padding: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
  vertical-align: middle;
}

.clickable-row {
  cursor: pointer;
  transition: all 0.3s;
}

.clickable-row:hover {
  background: rgba(0, 240, 255, 0.06) !important;
  box-shadow: inset 0 0 12px rgba(0, 240, 255, 0.1);
}

.repo-pr-name {
  font-size: 13.5px;
  font-weight: bold;
  color: #f1f5f9;
  max-width: 320px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 8px;
}

.link-icon {
  color: #00f0ff;
  font-size: 14px;
  opacity: 0.6;
  transition: all 0.3s;
}

.clickable-row:hover .link-icon {
  opacity: 1;
  transform: scale(1.1);
}

.score-val {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12.5px;
  font-weight: bold;
}

.score-val.green { color: #00ff66; }
.score-val.cyan { color: #00f0ff; }
.score-val.yellow { color: #eab308; }
.score-val.red { color: #ef4444; }

.date-txt {
  font-size: 12px;
  color: #64748b;
}

.status-pill {
  font-family: 'JetBrains Mono', monospace;
  font-size: 10px;
  color: #00ff66;
  background: rgba(0, 255, 102, 0.05);
  border: 1px solid rgba(0, 255, 102, 0.2);
  padding: 3px 8px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.status-pill .dot {
  width: 5px;
  height: 5px;
  background: #00ff66;
  border-radius: 50%;
  box-shadow: 0 0 5px #00ff66;
}

/* Symmetrical flex layout for logs card */
.logs-card {
  display: flex;
  flex-direction: column;
}

.logs-card :deep(.card-content) {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
  min-height: 380px; /* Aligns visually with the chart card */
}

/* Custom Cyber Pagination overrides */
.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-pagination.is-background .el-pager li:not(.is-active)) {
  background-color: rgba(13, 20, 35, 0.7) !important;
  border: 1px solid rgba(0, 240, 255, 0.15) !important;
  color: #94a3b8 !important;
  transition: all 0.3s;
}

:deep(.el-pagination.is-background .el-pager li:not(.is-active):hover) {
  border-color: #00f0ff !important;
  color: #00f0ff !important;
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.2);
}

:deep(.el-pagination.is-background .el-pager li.is-active) {
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%) !important;
  color: #ffffff !important;
  border: none !important;
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.4);
}

:deep(.el-pagination.is-background .btn-next),
:deep(.el-pagination.is-background .btn-prev) {
  background-color: rgba(13, 20, 35, 0.7) !important;
  border: 1px solid rgba(0, 240, 255, 0.15) !important;
  color: #94a3b8 !important;
  transition: all 0.3s;
}

:deep(.el-pagination.is-background .btn-next:hover),
:deep(.el-pagination.is-background .btn-prev:hover) {
  border-color: #00f0ff !important;
  color: #00f0ff !important;
}

:deep(.el-pagination .el-pagination__total) {
  color: #64748b !important;
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
}

.animate-fade-in {
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;
  }
  .analytics-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  .radar-chart-container {
    height: 300px;
  }
}

@media (max-width: 576px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
