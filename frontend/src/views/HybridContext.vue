<template>
  <div class="context-page animate-fade-in">
    <!-- Top Header -->
    <div class="welcome-banner">
      <div class="welcome-text">
        <h2>{{ t('context.header') }}</h2>
        <p>{{ t('context.subHeader') }}</p>
      </div>
      <div class="system-time">
        <span class="pulse-green"></span>
        <span class="clock">{{ t('context.status') }}</span>
      </div>
    </div>

    <!-- Summary Active Card Banner -->
    <div class="cyber-card banner-card">
      <div class="card-glow"></div>
      <div class="card-content overlay-banner">
        <div class="overlay-details">
          <div class="tech-icon"><el-icon><Share /></el-icon></div>
          <div class="details-text">
            <h3>{{ t('context.bannerTitle') }}</h3>
            <p>{{ t('context.bannerDesc') }}</p>
          </div>
        </div>
        <div class="overlay-stats">
          <div class="stat-node">
            <span class="value">14</span>
            <span class="label">{{ t('context.statLabel1') }}</span>
          </div>
          <div class="stat-node">
            <span class="value">96.8%</span>
            <span class="label">{{ t('context.statLabel2') }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="context-grid">
      <!-- Left side: Timeline of Auto Injections -->
      <div class="context-left">
        <div class="cyber-card timeline-card">
          <div class="card-glow"></div>
          <div class="card-content">
            <h3>{{ t('context.timelineTitle') }}</h3>
            <p class="subtitle">{{ t('context.timelineSub') }}</p>
            
            <div class="timeline-wrapper">
              <div v-for="(item, index) in injections" :key="index" class="timeline-item">
                <div class="timeline-dot"></div>
                <div class="timeline-body">
                  <div class="time-meta">
                    <span class="badge">{{ item.type }}</span>
                    <span class="timestamp">{{ item.time }}</span>
                  </div>
                  <h4 class="item-title">{{ item.title }}</h4>
                  <p class="item-desc">{{ lang === 'zh' ? item.descZh : item.desc }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right side: AST Dependency Linker Visualizer -->
      <div class="context-right">
        <div class="cyber-card graph-card">
          <div class="card-glow"></div>
          <div class="card-content">
            <h3>{{ t('context.graphTitle') }}</h3>
            <p class="subtitle">{{ t('context.graphSub') }}</p>

            <!-- ECharts Interactive AST Relation Graph Container -->
            <div class="ast-board">
              <div class="grid-layer"></div>
              <div ref="chartRef" class="echarts-graph-canvas"></div>
            </div>

            <!-- Inspect Panel Details -->
            <div class="inspect-panel" v-if="activeNode">
              <div class="inspect-header">
                <h4>{{ t('context.inspectTitle') }}: {{ activeNode.name }}</h4>
                <span :class="['node-badge', activeNode.status]">{{ lang === 'zh' ? activeNode.impactZh : activeNode.impact }}</span>
              </div>
              <div class="inspect-body">
                <div class="meta-row">
                  <span class="label">{{ t('context.colPath') }}:</span>
                  <span class="value code">{{ activeNode.path }}</span>
                </div>
                <div class="meta-row">
                  <span class="label">{{ t('context.nodeTypeLabel') }}:</span>
                  <span class="value type-tag">{{ lang === 'zh' ? activeNode.typeZh : activeNode.type }}</span>
                </div>
                <div class="meta-row">
                  <span class="label">{{ t('context.colExports') }}:</span>
                  <span class="value code exports">{{ activeNode.exports }}</span>
                </div>
                <div class="meta-row">
                  <span class="label">{{ t('context.colLinks') }}:</span>
                  <span class="value warning">{{ activeNode.links }}</span>
                </div>
                <div class="meta-row desc-row">
                  <span class="label">{{ t('context.impactLabel') }}:</span>
                  <p class="desc-text">{{ lang === 'zh' ? activeNode.descZh : activeNode.desc }}</p>
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
import { ref, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useI18n } from '@/i18n'
import * as echarts from 'echarts'

const { t, lang } = useI18n()
const chartRef = ref(null)
let myChart = null

const injections = ref([
  {
    type: "AUTO-INJECTED",
    time: "20:00:41.012",
    title: "packages/cli/package.json",
    desc: "Scanned dependency bounds. Detected GORM driver conflicts. Loaded standard MySQL driver configuration boundaries.",
    descZh: "扫描了依赖范围。检测到 GORM 驱动版本冲突。加载了标准的 MySQL 驱动配置边界。"
  },
  {
    type: "COMPILATION SCOPE",
    time: "20:00:41.458",
    title: "server/go.mod",
    desc: "Analyzed module boundary paths. Discovered dependency declarations for 'github.com/golang-jwt/jwt'. Added compilation scope to the context validator.",
    descZh: "分析了模块边界路径。发现了 'github.com/golang-jwt/jwt' 依赖声明。将编译范围添加到了上下文验证器中。"
  },
  {
    type: "PARENT ARCHITECTURE",
    time: "20:00:42.102",
    title: "server/service/review.go",
    desc: "Extracted call logic references. Found downstream dependent references 'github.com/aipr/ai-pr-reviewer/model'. Staged caller AST structure for deep audit mapping.",
    descZh: "提取了调用逻辑引用。发现了下游依赖引用 'github.com/aipr/ai-pr-reviewer/model'。暂存了调用者 AST 结构以进行深度审计分析。"
  }
])

const graphNodes = ref([
  {
    id: 'auth-jwt',
    name: 'auth/jwt.go',
    status: 'critical',
    path: 'server/auth/jwt.go',
    exports: 'GenerateToken(), ParseToken()',
    links: 'router.go, config/config.go',
    type: 'Modified File (变动源)',
    typeZh: '变更文件 (变动源头)',
    callDepth: 0,
    impact: 'HIGH',
    impactZh: '高风险',
    desc: 'Core JWT token operations. Upgraded symmetric signature (HS256) to asymmetric (RS256) key validator.',
    descZh: '核心 JWT 令牌库。已将对称 HS256 签名校验重构升级为非对称 RS256 公私钥架构。'
  },
  {
    id: 'router',
    name: 'router.go',
    status: 'warning',
    path: 'server/router.go',
    exports: 'InitRoutes(), RegisterAPI()',
    links: 'auth/jwt.go, main.go, controller/review.go',
    type: 'Direct Downstream (直接下游)',
    typeZh: '直接下游 (直接关联影响)',
    callDepth: 1,
    impact: 'HIGH',
    impactZh: '高风险',
    desc: 'Registers global HTTP routing endpoints. Holds JWT middleware token decoders intercepting all API scopes.',
    descZh: '配置系统 HTTP 路由表。拦截加载 JWT 鉴权中间件以保护所有后台控制器接口。'
  },
  {
    id: 'config',
    name: 'config.go',
    status: 'normal',
    path: 'server/config/config.go',
    exports: 'LoadEnv(), DatabaseCredentials()',
    links: 'auth/jwt.go',
    type: 'Dependency Provider (依赖提供)',
    typeZh: '依赖提供 (数据流来源)',
    callDepth: 1,
    impact: 'MEDIUM',
    impactZh: '中等风险',
    desc: 'Loads environment files and fetches private signing keys. Changes here alter token generation security parameters.',
    descZh: '负责从本地或环境变量中解析读取 RSA 签名私钥。配置变更将直接改变签名安全系数。'
  },
  {
    id: 'controller',
    name: 'controller/review.go',
    status: 'normal',
    path: 'server/controller/review.go',
    exports: 'Analyze(), GetDetail()',
    links: 'router.go',
    type: 'Indirect Downstream (间接下游)',
    typeZh: '间接下游 (跨文件引用)',
    callDepth: 2,
    impact: 'MEDIUM',
    impactZh: '中等风险',
    desc: 'Binds PR review endpoints. Relying on auth routing context checks to allow client triggers.',
    descZh: '绑定 PR 评审后台控制器。其安全性依赖于 router.go 中的 JWT 路由拦截鉴权。'
  },
  {
    id: 'stripe',
    name: 'stripe.go',
    status: 'warning',
    path: 'server/payment/stripe.go',
    exports: 'WebhookCallback(), InitClient()',
    links: 'router.go',
    type: 'Indirect Downstream (间接下游)',
    typeZh: '间接下游 (路由控制)',
    callDepth: 2,
    impact: 'MEDIUM',
    impactZh: '中等风险',
    desc: 'Handles Stripe payment webhook callbacks. Securing endpoint operations using auth middleware headers.',
    descZh: '处理外部支付渠道 Webhook 回调。其支付核心接口同样受 router.go 鉴权层保护。'
  },
  {
    id: 'main',
    name: 'main.go',
    status: 'normal',
    path: 'server/main.go',
    exports: 'main()',
    links: 'router.go',
    type: 'System Entry (主入口)',
    typeZh: '系统主入口 (最上游)',
    callDepth: 2,
    impact: 'LOW',
    impactZh: '低风险',
    desc: 'Main application bootstrap. Registers MySQL and Redis connections and invokes the initial HTTP engine.',
    descZh: '系统主引导文件。创建数据库与缓存会话，并启动 HTTP 服务器引擎承载流量。'
  }
])

const activeNode = ref(graphNodes.value[0]) // Focus on auth/jwt.go by default

const initGraphChart = () => {
  if (!chartRef.value) return
  if (myChart) {
    myChart.dispose()
  }

  myChart = echarts.init(chartRef.value)

  const isZh = lang.value === 'zh'

  const seriesData = graphNodes.value.map(node => {
    let color = '#00f0ff'
    let shadowColor = 'rgba(0, 240, 255, 0.4)'
    if (node.status === 'critical') {
      color = '#ef4444'
      shadowColor = 'rgba(239, 68, 68, 0.6)'
    } else if (node.status === 'warning') {
      color = '#eab308'
      shadowColor = 'rgba(234, 179, 8, 0.5)'
    }

    return {
      id: node.id,
      name: node.name,
      symbolSize: node.id === 'auth-jwt' ? 42 : 32,
      itemStyle: {
        color: '#090f1d',
        borderColor: color,
        borderWidth: 2.5,
        shadowColor: shadowColor,
        shadowBlur: 12
      },
      label: {
        show: true,
        position: 'top',
        color: '#f8fafc',
        fontFamily: 'Outfit, sans-serif',
        fontWeight: 'bold',
        fontSize: 11,
        formatter: node.name
      },
      raw: node
    }
  })

  // Fixed coordinates for clean horizontal propagating flow
  const coords = {
    'config': [120, 200],
    'auth-jwt': [300, 200],
    'router': [480, 200],
    'main': [660, 200],
    'controller': [660, 90],
    'stripe': [660, 310]
  }

  seriesData.forEach(item => {
    if (coords[item.id]) {
      item.x = coords[item.id][0]
      item.y = coords[item.id][1]
    }
  })

  const option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(9, 15, 29, 0.95)',
      borderColor: 'rgba(0, 240, 255, 0.3)',
      borderWidth: 1,
      textStyle: { color: '#cbd5e1', fontSize: 11 },
      formatter: (params) => {
        if (params.dataType === 'edge') {
          return t('context.tooltipEdge')
        }
        const raw = params.data.raw
        return `
          <strong style="color: #fff; font-size: 12px;">${raw.name}</strong><br/>
          <span style="color: #64748b;">${isZh ? '角色类型' : 'Type'}:</span> ${isZh ? raw.typeZh : raw.type}<br/>
          <span style="color: #64748b;">${isZh ? '影响评级' : 'Impact'}:</span> <span style="color: ${raw.status === 'critical' ? '#ef4444' : (raw.status === 'warning' ? '#eab308' : '#00ff66')}">${isZh ? raw.impactZh : raw.impact}</span>
        `
      }
    },
    series: [
      {
        type: 'graph',
        layout: 'none',
        roam: false,
        focusNodeAdjacency: true,
        edgeSymbol: ['none', 'arrow'],
        edgeSymbolSize: [5, 10],
        lineStyle: {
          color: 'rgba(0, 240, 255, 0.3)',
          width: 2,
          curveness: 0
        },
        data: seriesData,
        links: [
          { source: 'config', target: 'auth-jwt', lineStyle: { curveness: 0 } },
          { source: 'auth-jwt', target: 'router', lineStyle: { color: 'rgba(239, 68, 68, 0.5)', width: 3.5, curveness: 0 } },
          { source: 'router', target: 'main', lineStyle: { curveness: 0 } },
          { source: 'router', target: 'controller', lineStyle: { curveness: -0.18 } },
          { source: 'router', target: 'stripe', lineStyle: { curveness: 0.18 } }
        ]
      }
    ]
  }

  myChart.setOption(option)

  myChart.on('click', (params) => {
    if (params.dataType === 'node') {
      const clickedId = params.data.id
      const found = graphNodes.value.find(n => n.id === clickedId)
      if (found) {
        activeNode.value = found
      }
    }
  })
}

// Watch language state to swap axis and tooltips instantly
watch(lang, () => {
  nextTick(() => {
    initGraphChart()
  })
})

onMounted(() => {
  nextTick(() => {
    initGraphChart()
  })
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  if (myChart) {
    myChart.dispose()
  }
  window.removeEventListener('resize', handleResize)
})

const handleResize = () => {
  if (myChart) {
    myChart.resize()
  }
}
</script>

<style scoped>
.context-page {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

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

/* Banner details stats */
.overlay-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(90deg, rgba(0, 240, 255, 0.05) 0%, rgba(112, 0, 255, 0.02) 100%);
}

.overlay-details {
  display: flex;
  align-items: center;
  gap: 20px;
}

.tech-icon {
  font-size: 32px;
  color: #00f0ff;
  filter: drop-shadow(0 0 8px rgba(0, 240, 255, 0.4));
  display: flex;
  align-items: center;
}

.details-text h3 {
  margin: 0;
  font-size: 15px;
  color: #f8fafc;
}

.details-text p {
  margin: 5px 0 0 0;
  font-size: 12.5px;
  color: #64748b;
  line-height: 1.5;
  max-width: 600px;
}

.overlay-stats {
  display: flex;
  gap: 30px;
}

.stat-node {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.stat-node .value {
  font-size: 24px;
  font-weight: 800;
  color: #00f0ff;
  text-shadow: 0 0 10px rgba(0, 240, 255, 0.3);
  font-family: 'JetBrains Mono', monospace;
}

.stat-node .label {
  font-size: 10.5px;
  color: #475569;
  letter-spacing: 0.5px;
  margin-top: 2px;
}

/* Two-column layout */
.context-grid {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 25px;
}

.timeline-wrapper {
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: relative;
  padding-left: 15px;
}

.timeline-wrapper::before {
  content: '';
  position: absolute;
  top: 5px;
  left: 3px;
  bottom: 5px;
  width: 1px;
  background: rgba(0, 240, 255, 0.15);
}

.timeline-item {
  position: relative;
}

.timeline-dot {
  position: absolute;
  top: 6px;
  left: -15px;
  width: 7px;
  height: 7px;
  background-color: #00f0ff;
  border-radius: 50%;
  box-shadow: 0 0 8px #00f0ff;
}

.time-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.time-meta .badge {
  font-size: 9px;
  font-family: 'JetBrains Mono', monospace;
  background: rgba(0, 240, 255, 0.1);
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: #00f0ff;
  padding: 1px 6px;
  border-radius: 3px;
}

.time-meta .timestamp {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  color: #475569;
}

.item-title {
  margin: 0;
  font-size: 13px;
  font-weight: bold;
  color: #f1f5f9;
}

.item-desc {
  margin: 5px 0 0 0;
  font-size: 11.5px;
  color: #64748b;
  line-height: 1.5;
}

/* ECharts relation graph layout */
.ast-board {
  height: 380px;
  background: #060913;
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  position: relative;
  overflow: hidden;
  margin-bottom: 25px;
}

.grid-layer {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-size: 30px 30px;
  background-image: 
    linear-gradient(to right, rgba(255, 255, 255, 0.01) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(255, 255, 255, 0.01) 1px, transparent 1px);
  z-index: 1;
}

.echarts-graph-canvas {
  width: 100%;
  height: 100%;
  position: relative;
  z-index: 5;
}

/* Inspect panel details styling */
.inspect-panel {
  background: rgba(10, 15, 30, 0.4);
  border: 1px solid rgba(0, 240, 255, 0.1);
  border-radius: 8px;
  padding: 20px;
  box-shadow: inset 0 0 15px rgba(0, 240, 255, 0.03);
}

.inspect-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 10px;
}

.inspect-header h4 {
  margin: 0;
  font-size: 13.5px;
  font-family: 'Outfit', sans-serif;
  color: #00f0ff;
  font-weight: 800;
  letter-spacing: 0.5px;
}

.node-badge {
  font-family: 'JetBrains Mono', monospace;
  font-size: 9.5px;
  font-weight: bold;
  padding: 3px 10px;
  border-radius: 20px;
}

.node-badge.normal { background: rgba(0, 255, 102, 0.1); color: #00ff66; border: 1px solid rgba(0, 255, 102, 0.2); }
.node-badge.warning { background: rgba(234, 179, 8, 0.1); color: #eab308; border: 1px solid rgba(234, 179, 8, 0.2); }
.node-badge.critical { background: rgba(239, 68, 68, 0.1); color: #f87171; border: 1px solid rgba(239, 68, 68, 0.2); }

.meta-row {
  display: flex;
  margin-bottom: 10px;
  font-size: 12.5px;
  align-items: center;
}

.meta-row .label {
  width: 150px;
  color: #64748b;
  font-weight: 500;
}

.meta-row .value {
  color: #cbd5e1;
}

.meta-row .value.code {
  font-family: 'JetBrains Mono', monospace;
  color: #00f0ff;
  background: rgba(0, 240, 255, 0.05);
  padding: 2px 8px;
  border-radius: 4px;
  border: 1px solid rgba(0, 240, 255, 0.1);
}

.meta-row .value.code.exports {
  color: #00ff66;
  background: rgba(0, 255, 102, 0.05);
  border-color: rgba(0, 255, 102, 0.1);
}

.meta-row .value.type-tag {
  color: #f8fafc;
  font-weight: bold;
}

.desc-row {
  align-items: flex-start;
  margin-top: 15px;
  border-top: 1px dashed rgba(255, 255, 255, 0.05);
  padding-top: 15px;
}

.desc-row .desc-text {
  margin: 0;
  color: #94a3b8;
  line-height: 1.6;
  font-size: 12px;
  flex: 1;
}

.animate-fade-in {
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
