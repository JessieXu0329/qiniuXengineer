<template>
  <div class="context-page animate-fade-in">
    <!-- Top Header -->
    <div class="welcome-banner">
      <div class="welcome-text">
        <h2>{{ t[currentLang].header }}</h2>
        <p>{{ t[currentLang].subHeader }}</p>
      </div>
      <div class="system-time">
        <span class="pulse-green"></span>
        <span class="clock">{{ t[currentLang].status }}</span>
      </div>
    </div>

    <!-- Summary Active Card Banner -->
    <div class="cyber-card banner-card">
      <div class="card-glow"></div>
      <div class="card-content overlay-banner">
        <div class="overlay-details">
          <div class="tech-icon"><i class="el-icon-connection"></i></div>
          <div class="details-text">
            <h3>{{ t[currentLang].bannerTitle }}</h3>
            <p>{{ t[currentLang].bannerDesc }}</p>
          </div>
        </div>
        <div class="overlay-stats">
          <div class="stat-node">
            <span class="value">14</span>
            <span class="label">{{ t[currentLang].statLabel1 }}</span>
          </div>
          <div class="stat-node">
            <span class="value">96.8%</span>
            <span class="label">{{ t[currentLang].statLabel2 }}</span>
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
            <h3>{{ t[currentLang].timelineTitle }}</h3>
            <p class="subtitle">{{ t[currentLang].timelineSub }}</p>
            
            <div class="timeline-wrapper">
              <div v-for="(item, index) in injections" :key="index" class="timeline-item">
                <div class="timeline-dot"></div>
                <div class="timeline-body">
                  <div class="time-meta">
                    <span class="badge">{{ item.type }}</span>
                    <span class="timestamp">{{ item.time }}</span>
                  </div>
                  <h4 class="item-title">{{ item.title }}</h4>
                  <p class="item-desc">{{ currentLang === 'zh' ? item.descZh : item.desc }}</p>
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
            <h3>{{ t[currentLang].graphTitle }}</h3>
            <p class="subtitle">{{ t[currentLang].graphSub }}</p>

            <div class="ast-board">
              <div class="grid-layer"></div>
              
              <!-- Simulated Graph Nodes -->
              <div 
                v-for="node in graphNodes" 
                :key="node.id" 
                :class="['ast-node-element', node.status, { active: activeNode?.id === node.id }]"
                :style="{ top: node.y + 'px', left: node.x + 'px' }"
                @click="inspectNode(node)"
              >
                <div class="node-glow"></div>
                <span class="node-tag">{{ node.label }}</span>
              </div>

              <!-- Connecting Vector Lines -->
              <svg class="graph-svg">
                <defs>
                  <marker id="arrow" viewBox="0 0 10 10" refX="28" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
                    <path d="M 0 0 L 10 5 L 0 10 z" fill="rgba(0, 240, 255, 0.4)"/>
                  </marker>
                </defs>
                <line x1="200" y1="50" x2="100" y2="150" stroke="rgba(0, 240, 255, 0.25)" stroke-width="2" marker-end="url(#arrow)" />
                <line x1="200" y1="50" x2="300" y2="150" stroke="rgba(0, 240, 255, 0.25)" stroke-width="2" marker-end="url(#arrow)" />
                <line x1="100" y1="150" x2="200" y2="250" stroke="rgba(0, 240, 255, 0.25)" stroke-width="2" marker-end="url(#arrow)" />
                <line x1="300" y1="150" x2="200" y2="250" stroke="rgba(0, 240, 255, 0.25)" stroke-width="2" marker-end="url(#arrow)" />
              </svg>
            </div>

            <!-- Inspect Panel Details -->
            <div class="inspect-panel" v-if="activeNode">
              <div class="inspect-header">
                <h4>{{ t[currentLang].inspectTitle }}: {{ activeNode.label }}</h4>
                <span :class="['node-badge', activeNode.status]">{{ activeNode.status.toUpperCase() }}</span>
              </div>
              <div class="inspect-body">
                <div class="meta-row">
                  <span class="label">{{ t[currentLang].colPath }}:</span>
                  <span class="value code">{{ activeNode.path }}</span>
                </div>
                <div class="meta-row">
                  <span class="label">{{ t[currentLang].colExports }}:</span>
                  <span class="value">{{ activeNode.exports }}</span>
                </div>
                <div class="meta-row">
                  <span class="label">{{ t[currentLang].colLinks }}:</span>
                  <span class="value warning">{{ activeNode.links }}</span>
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
import { ref, inject } from 'vue'

const currentLang = inject('lang', ref('zh'))

const t = {
  zh: {
    header: "混合上下文检索器",
    subHeader: "AST 语法解析、包路径规范约束发现以及超出 Diff 代码的关联认知叠加图层",
    status: "认知图层已激活",
    bannerTitle: "[当前扫描上下文配置文件: 认知图层已激活]",
    bannerDesc: "分析器动态检索您的代码仓库边界。它不孤立地审查 Diff 行，而是提取父类、接口和包配置，为大模型评审流水线注入高精度的上下文背景参数。",
    statLabel1: "已注入上下文锚点数",
    statLabel2: "上下文精准度评分",
    timelineTitle: "动态上下文注入时间线",
    timelineSub: "在最近一次扫描期间执行依赖解析的按时序记录",
    graphTitle: "AST 逻辑依赖连通网络",
    graphSub: "实时 Abstract Syntax Tree 类关系链追踪。点击节点可审查依赖指纹签名。",
    inspectTitle: "节点详细依赖指纹签名",
    colPath: "文件绝对路径",
    colExports: "已发现导出 API",
    colLinks: "关联网结依赖"
  },
  en: {
    header: "HYBRID CONTEXT RETRIEVER",
    subHeader: "AST parsing, package constraints discovery, and out-of-diff cognitive overlays",
    status: "COGNITIVE OVERLAY ACTIVE",
    bannerTitle: "[CURRENT SCAN CONTEXT PROFILE: COGNITIVE OVERLAY ACTIVE]",
    bannerDesc: "The retriever scans your repository boundaries dynamically. Instead of isolating Diff lines, it extracts parent classes, interfaces, and packages configurations to feed high-fidelity context parameters to the LLM review pipeline.",
    statLabel1: "Context Anchors Injected",
    statLabel2: "Context Accuracy Score",
    timelineTitle: "DYNAMIC CONTEXT INJECTION TIMELINE",
    timelineSub: "Chronological footprint of dependency resolutions executed on last scanning event",
    graphTitle: "AST LOGICAL DEPENDENCY LINKER",
    graphSub: "Live Abstract Syntax Tree class-linkage tracking. Click nodes to inspect dependency signatures.",
    inspectTitle: "NODE DETAILED SIGNATURE",
    colPath: "Absolute Path",
    colExports: "Discovered Exports",
    colLinks: "Dependency Linkages"
  }
}

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
    time: "20:00:41.085",
    title: "tsconfig.json",
    desc: "Parsed compiler compiler configurations. Enforced strictNullChecks validations across standard model definitions.",
    descZh: "解析了编译器配置。在标准模型定义中强制执行了 strictNullChecks 严格空值校验。"
  },
  {
    type: "AST IMPORT RESOLUTION",
    time: "20:00:41.204",
    title: "server/controller/review.go",
    desc: "Found imports referencing 'github.com/aipr/ai-pr-reviewer/service'. Automatically resolved interface parameters.",
    descZh: "发现引用 'github.com/aipr/ai-pr-reviewer/service' 的导入语句。自动解析了相关接口参数。"
  },
  {
    type: "GIT ANCHOR LOAD",
    time: "20:00:41.341",
    title: ".github/workflows/ci.yml",
    desc: "Discovered GitHub actions pipelines. Context loaded to evaluate deployment constraints.",
    descZh: "检测到 GitHub Actions 流水线配置。已加载上下文用以评估部署约束条件。"
  }
])

const graphNodes = ref([
  { id: '1', label: 'router.go', x: 200, y: 30, status: 'normative', path: 'server/router.go', exports: 'InitRoutes(), RegisterAPI()', links: 'payment/stripe.go, auth/jwt.go' },
  { id: '2', label: 'auth/jwt.go', x: 60, y: 130, status: 'critical', path: 'server/auth/jwt.go', exports: 'GenerateToken(), ParseToken()', links: 'router.go, config/config.go' },
  { id: '3', label: 'stripe.go', x: 260, y: 130, status: 'warning', path: 'server/payment/stripe.go', exports: 'WebhookCallback(), InitClient()', links: 'router.go, config/config.go' },
  { id: '4', label: 'config.go', x: 200, y: 230, status: 'normative', path: 'server/config/config.go', exports: 'LoadEnv(), DatabaseCredentials()', links: 'auth/jwt.go, stripe.go' }
])

const activeNode = ref(graphNodes.value[1]) // Focus on auth/jwt.go by default

const inspectNode = (node) => {
  activeNode.value = node
}
</script>

<style scoped>
.context-page {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

/* Banner overlay styles */
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

/* Cyber card base styles */
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

/* Banner Overlay Details */
.overlay-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(90deg, rgba(0, 240, 255, 0.05) 0%, transparent 100%);
}

.overlay-details {
  display: flex;
  align-items: center;
  gap: 20px;
  max-width: 70%;
}

.tech-icon {
  font-size: 32px;
  color: #00f0ff;
  text-shadow: 0 0 15px rgba(0, 240, 255, 0.4);
}

.details-text h3 {
  color: #00f0ff !important;
  text-shadow: 0 0 10px rgba(0, 240, 255, 0.2);
}

.details-text p {
  margin: 5px 0 0 0;
  font-size: 13px;
  color: #94a3b8;
  line-height: 1.5;
}

.overlay-stats {
  display: flex;
  gap: 25px;
}

.stat-node {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.stat-node .value {
  font-size: 24px;
  font-weight: 900;
  color: #ffffff;
  font-family: 'JetBrains Mono', monospace;
}

.stat-node .label {
  font-size: 10px;
  color: #475569;
  margin-top: 4px;
  font-family: 'JetBrains Mono', monospace;
}

/* Context Grid Layout */
.context-grid {
  display: grid;
  grid-template-columns: 460px 1fr;
  gap: 25px;
}

/* Timeline Components */
.timeline-wrapper {
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: relative;
  padding-left: 20px;
  border-left: 1px solid rgba(255, 255, 255, 0.05);
}

.timeline-item {
  position: relative;
}

.timeline-dot {
  position: absolute;
  left: -25px;
  top: 6px;
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background-color: #00f0ff;
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
  font-size: 13.5px;
  font-weight: bold;
  color: #f1f5f9;
}

.item-desc {
  margin: 5px 0 0 0;
  font-size: 12px;
  color: #64748b;
  line-height: 1.5;
}

/* AST Board Interactive Canvas Graph */
.ast-board {
  height: 300px;
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
  background-size: 20px 20px;
  background-image: 
    linear-gradient(to right, rgba(255, 255, 255, 0.01) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(255, 255, 255, 0.01) 1px, transparent 1px);
}

.graph-svg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.ast-node-element {
  position: absolute;
  z-index: 5;
  cursor: pointer;
  background: rgba(13, 20, 35, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 8px 16px;
  border-radius: 20px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  transition: all 0.3s;
}

.ast-node-element:hover {
  transform: scale(1.05);
}

.ast-node-element.active {
  border-color: #00f0ff !important;
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.3);
}

.ast-node-element.normative {
  border-color: rgba(255, 255, 255, 0.2);
}

.ast-node-element.warning {
  border-color: #eab308;
}

.ast-node-element.critical {
  border-color: #ef4444;
}

.node-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  filter: blur(8px);
  opacity: 0.15;
  border-radius: 20px;
}

.normative .node-glow { background-color: #ffffff; }
.warning .node-glow { background-color: #eab308; }
.critical .node-glow { background-color: #ef4444; }

.node-tag {
  color: #e2e8f0;
}

/* Detailed Inspect Panel */
.inspect-panel {
  background: rgba(10, 15, 30, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 15px;
}

.inspect-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 8px;
}

.inspect-header h4 {
  margin: 0;
  font-size: 12px;
  font-family: 'JetBrains Mono', monospace;
  color: #00f0ff;
}

.node-badge {
  font-family: 'JetBrains Mono', monospace;
  font-size: 9px;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 4px;
}

.node-badge.normative { background: rgba(255, 255, 255, 0.05); color: #94a3b8; }
.node-badge.warning { background: rgba(234, 179, 8, 0.15); color: #eab308; }
.node-badge.critical { background: rgba(239, 68, 68, 0.15); color: #f87171; }

.meta-row {
  display: flex;
  margin-bottom: 8px;
  font-size: 12.5px;
}

.meta-row .label {
  width: 160px;
  color: #475569;
}

.meta-row .value {
  color: #cbd5e1;
}

.meta-row .value.code {
  font-family: 'JetBrains Mono', monospace;
  color: #00ff66;
}

.meta-row .value.warning {
  color: #fca5a5;
}

.animate-fade-in {
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
