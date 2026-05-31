<template>
  <div class="cyber-container">
    <!-- Neon grid background accent -->
    <div class="cyber-grid-overlay"></div>
    
    <!-- Top glowing navbar -->
    <header class="cyber-navbar">
      <div class="brand">
        <span class="pulse-dot"></span>
        <h1 class="glow-text">{{ t[currentLang].title }}</h1>
        <span class="badge">DEVOPS CORE v1.0</span>
      </div>
      
      <div class="system-status">
        <!-- Premium language switch button in top-right corner -->
        <button class="lang-toggle-btn" @click="toggleLang">
          <el-icon><Compass /></el-icon>
          <span>{{ currentLang === 'zh' ? 'ENGLISH' : '中文' }}</span>
        </button>

        <div class="status-item">
          <span class="label">{{ t[currentLang].api }}:</span>
          <span :class="['value', serverStatus]">{{ serverStatus.toUpperCase() }}</span>
        </div>
        <div class="status-item">
          <span class="label">{{ t[currentLang].db }}:</span>
          <span class="value success">CONNECTED</span>
        </div>
        <div class="status-item">
          <span class="label">{{ t[currentLang].latency }}:</span>
          <span class="value warning">18ms</span>
        </div>
      </div>
    </header>

    <div class="main-layout">
      <!-- Cyber glassmorphism sidebar -->
      <aside class="cyber-sidebar">
        <div class="menu-section">
          <div class="section-title">{{ t[currentLang].sec1 }}</div>
          <nav class="menu-list">
            <router-link to="/" class="menu-item" exact-active-class="active">
              <el-icon><Odometer /></el-icon>
              <span>{{ t[currentLang].dashboard }}</span>
            </router-link>
            <router-link to="/review" class="menu-item" active-class="active">
              <el-icon><DocumentChecked /></el-icon>
              <span>{{ t[currentLang].review }}</span>
            </router-link>
            <router-link to="/context" class="menu-item" active-class="active">
              <el-icon><Connection /></el-icon>
              <span>{{ t[currentLang].context }}</span>
            </router-link>
            <router-link to="/false-positive" class="menu-item" active-class="active">
              <el-icon><Filter /></el-icon>
              <span>{{ t[currentLang].filter }}</span>
            </router-link>
          </nav>
        </div>

        <div class="menu-section">
          <div class="section-title">{{ t[currentLang].sec2 }}</div>
          <nav class="menu-list">
            <router-link to="/config" class="menu-item" active-class="active">
              <el-icon><Setting /></el-icon>
              <span>{{ t[currentLang].config }}</span>
            </router-link>
          </nav>
        </div>

        <!-- Cyber user info profile -->
        <div class="cyber-profile">
          <div class="avatar-glow">
            <div class="avatar">AI</div>
          </div>
          <div class="profile-info">
            <div class="name">DevOps_Architect</div>
            <div class="role">{{ t[currentLang].role }}</div>
          </div>
        </div>
      </aside>

      <!-- Central main viewport -->
      <main class="cyber-main">
        <router-view v-if="$route.matched.length > 0" />
        <div v-else class="fallback-wrapper">
          <!-- Fallback direct component load to ensure dashboard runs if router has not loaded yet -->
          <DashboardMockUp />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, provide } from 'vue'
import DashboardMockUp from './views/Dashboard.vue'

const serverStatus = ref('connecting')
const currentLang = ref('zh') // Default to Chinese as requested

// Translation dictionary
const t = {
  zh: {
    title: "AI 代码评审助手",
    api: "API 服务",
    db: "数据库",
    latency: "系统延迟",
    dashboard: "Dashboard 大盘",
    review: "PR Review 智能评审",
    context: "Hybrid Context 关联项",
    filter: "False Positive 过滤",
    config: "AI 引擎配置",
    tokens: "Access Tokens 令牌管理",
    role: "系统超级管理员",
    sec1: "核心控制中心",
    sec2: "预设与引擎配置"
  },
  en: {
    title: "AI PR REVIEWER",
    api: "API SERVER",
    db: "DB NODE",
    latency: "LATENCY",
    dashboard: "Dashboard Index",
    review: "PR Review Audit",
    context: "Hybrid Context Retriever",
    filter: "False Positive Engine",
    config: "AI Engine Config",
    tokens: "Access Tokens",
    role: "Admin Operator",
    sec1: "METRICS & CONTROL",
    sec2: "PRESETS & ENGINE"
  }
}

// Provide language reactive state globally so all child views translate instantly!
provide('lang', currentLang)

const toggleLang = () => {
  currentLang.value = currentLang.value === 'zh' ? 'en' : 'zh'
}

// Test connection to Go backend on startup
onMounted(async () => {
  try {
    const res = await fetch('http://localhost:8080/api/v1/ping')
    if (res.ok) {
      serverStatus.value = 'online'
    } else {
      serverStatus.value = 'offline'
    }
  } catch (e) {
    serverStatus.value = 'offline'
  }
})
</script>

<style>
/* CSS Reset inside App.vue for absolute design control */
* {
  box-sizing: border-box;
}

body {
  background-color: #080c14;
  font-family: 'Outfit', sans-serif;
  color: #e2e8f0;
  margin: 0;
  padding: 0;
}

/* Base cyber layout & neon elements */
.cyber-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  position: relative;
  overflow: hidden;
}

.cyber-grid-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-size: 40px 40px;
  background-image: 
    linear-gradient(to right, rgba(0, 240, 255, 0.03) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(0, 240, 255, 0.03) 1px, transparent 1px);
  pointer-events: none;
  z-index: 1;
}

/* Glowing Navbar styling */
.cyber-navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 30px;
  background: rgba(13, 20, 35, 0.85);
  border-bottom: 2px solid rgba(0, 240, 255, 0.3);
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.15);
  backdrop-filter: blur(12px);
  z-index: 10;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pulse-dot {
  width: 10px;
  height: 10px;
  background-color: #00ff66;
  border-radius: 50%;
  box-shadow: 0 0 10px #00ff66, 0 0 20px #00ff66;
  animation: pulse 1.8s infinite;
}

.glow-text {
  font-size: 22px;
  font-weight: 900;
  letter-spacing: 2px;
  margin: 0;
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  filter: drop-shadow(0 0 8px rgba(0, 240, 255, 0.5));
}

.badge {
  font-size: 10px;
  font-family: 'JetBrains Mono', monospace;
  background: rgba(0, 240, 255, 0.1);
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: #00f0ff;
  padding: 2px 8px;
  border-radius: 4px;
}

.system-status {
  display: flex;
  align-items: center;
  gap: 25px;
}

.lang-toggle-btn {
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

.lang-toggle-btn:hover {
  background: rgba(0, 240, 255, 0.15);
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.4);
  transform: translateY(-1px);
}

.status-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
}

.status-item .label {
  color: #64748b;
}

.status-item .value {
  font-weight: bold;
}

.status-item .value.online,
.status-item .value.success {
  color: #00ff66;
  text-shadow: 0 0 8px rgba(0, 255, 102, 0.4);
}

.status-item .value.offline {
  color: #ef4444;
  text-shadow: 0 0 8px rgba(239, 68, 68, 0.4);
}

.status-item .value.warning {
  color: #eab308;
  text-shadow: 0 0 8px rgba(234, 179, 8, 0.4);
}

/* Sidebar & Main view styling */
.main-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
  z-index: 5;
}

.cyber-sidebar {
  width: 260px;
  background: rgba(9, 15, 29, 0.85);
  border-right: 1px solid rgba(0, 240, 255, 0.1);
  padding: 30px 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  backdrop-filter: blur(8px);
}

.menu-section {
  margin-bottom: 30px;
}

.section-title {
  font-size: 11px;
  font-family: 'JetBrains Mono', monospace;
  letter-spacing: 1.5px;
  color: #475569;
  margin-bottom: 15px;
  padding-left: 10px;
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #94a3b8;
  text-decoration: none;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 13.5px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  border: 1px solid transparent;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.menu-item:hover {
  color: #00f0ff;
  background: rgba(0, 240, 255, 0.05);
  border-color: rgba(0, 240, 255, 0.1);
  transform: translateX(4px);
}

.menu-item.active {
  color: #ffffff;
  background: linear-gradient(90deg, rgba(0, 240, 255, 0.15) 0%, rgba(112, 0, 255, 0.05) 100%);
  border: 1px solid rgba(0, 240, 255, 0.4);
  box-shadow: 0 0 15px rgba(0, 240, 255, 0.1);
}

.menu-item .el-icon {
  font-size: 18px;
}

/* User profile inside sidebar */
.cyber-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.avatar-glow {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%);
  padding: 2px;
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.3);
}

.avatar {
  background: #0d1321;
  color: #00f0ff;
  font-weight: bold;
  font-size: 14px;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.profile-info .name {
  font-size: 13px;
  font-weight: bold;
  color: #f8fafc;
}

.profile-info .role {
  font-size: 10px;
  color: #64748b;
  font-family: 'JetBrains Mono', monospace;
}

/* Main Viewport scrollable config */
.cyber-main {
  flex: 1;
  padding: 30px;
  overflow-y: auto;
  position: relative;
}

.fallback-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 0.9; }
  50% { transform: scale(1.15); opacity: 1; box-shadow: 0 0 15px #00ff66, 0 0 25px #00ff66; }
  100% { transform: scale(1); opacity: 0.9; }
}

@media (max-width: 1200px) {
  .cyber-sidebar {
    width: 230px;
    padding: 20px 10px;
  }
  .menu-item {
    font-size: 12.5px;
    padding: 10px 12px;
    gap: 8px;
  }
}
</style>
