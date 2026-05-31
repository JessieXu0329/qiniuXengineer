<template>
  <div class="config-page animate-fade-in">
    <!-- Top Header -->
    <div class="welcome-banner">
      <div class="welcome-text">
        <h2>{{ t[currentLang].header }}</h2>
        <p>{{ t[currentLang].subHeader }}</p>
      </div>
      <div class="system-time">
        <span class="pulse-green"></span>
        <span class="clock">{{ t[currentLang].clock }}</span>
      </div>
    </div>

    <!-- Active Token Node Status Banner -->
    <div class="cyber-card status-card">
      <div class="card-glow"></div>
      <div class="card-content status-banner">
        <h3>{{ t[currentLang].statusTitle }}</h3>
        <div class="status-nodes-group">
          <div class="status-node-pill">
            <span class="pulse-dot-green"></span>
            <span class="label">{{ t[currentLang].statusLabel1 }}</span>
          </div>
          <div class="status-node-pill">
            <span class="pulse-dot-green"></span>
            <span class="label">{{ t[currentLang].statusLabel2 }}</span>
          </div>
          <div class="status-node-pill">
            <span class="pulse-dot-green"></span>
            <span class="label">{{ t[currentLang].statusLabel3 }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Model Selector Grid -->
    <div class="model-selection-wrapper">
      <div class="section-header">
        <h3>{{ t[currentLang].selectionTitle }}</h3>
        <p class="subtitle">{{ t[currentLang].selectionSub }}</p>
      </div>

      <div class="models-grid">
        <div 
          v-for="model in models" 
          :key="model.id" 
          :class="['model-card', { active: selectedModelId === model.id }]"
          @click="selectModel(model.id)"
        >
          <div class="active-banner" v-if="selectedModelId === model.id">{{ t[currentLang].modelActive }}</div>
          <div class="model-meta">
            <span class="model-provider">{{ model.provider }}</span>
            <h4 class="model-title">{{ model.name }}</h4>
            <p class="model-desc">{{ currentLang === 'zh' ? model.descZh : model.desc }}</p>
          </div>
          <div class="model-footer">
            <span class="latency-pill"><el-icon><Timer /></el-icon> {{ model.avgLatency }}</span>
            <span class="select-text">{{ selectedModelId === model.id ? t[currentLang].modelEnabled : t[currentLang].modelSelect }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Credentials Configuration Form -->
    <div class="cyber-card credentials-wrapper">
      <div class="card-glow"></div>
      <div class="card-content">
        <h3>{{ t[currentLang].configTitle }}</h3>
        <p class="subtitle">{{ t[currentLang].configSub }}</p>

        <form @submit.prevent="saveConfig" class="config-form">
          <div class="form-row">
            <div class="form-group">
              <label>{{ t[currentLang].formLabelBaseUrl }}</label>
              <input v-model="baseUrl" type="text" class="form-input" placeholder="https://api.deepseek.com/v1" />
            </div>
            <div class="form-group">
              <label>{{ t[currentLang].formLabelModel }}</label>
              <input v-model="enforcedModel" type="text" class="form-input" placeholder="deepseek-chat" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>{{ t[currentLang].formLabelApiKey }}</label>
              <div class="password-input-wrapper">
                <input v-model="apiKey" :type="showApiKey ? 'text' : 'password'" class="form-input" placeholder="sk-..." />
                <button type="button" class="eye-btn" @click="showApiKey = !showApiKey">
                  <el-icon><View v-if="showApiKey" /><Hide v-else /></el-icon>
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>{{ t[currentLang].formLabelGithubToken }}</label>
              <div class="password-input-wrapper">
                <input v-model="githubToken" :type="showGithubToken ? 'text' : 'password'" class="form-input" placeholder="ghp_..." />
                <button type="button" class="eye-btn" @click="showGithubToken = !showGithubToken">
                  <el-icon><View v-if="showGithubToken" /><Hide v-else /></el-icon>
                </button>
              </div>
            </div>
          </div>

          <div class="actions-row">
            <button type="submit" class="save-btn">
              <span class="btn-glow"></span>
              <span class="btn-text">{{ t[currentLang].formBtn }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, inject, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const currentLang = inject('lang', ref('zh'))

const t = {
  zh: {
    header: "扫描引擎与凭据控制器",
    subHeader: "AI 大模型参数选择、API 接口配置以及访问令牌授权管理控制台",
    clock: "大模型引擎已上线",
    statusTitle: "当前系统节点连接状态",
    statusLabel1: "[GitHub API 客户端: 令牌活跃中 (剩余 5000 次请求/小时)]",
    statusLabel2: "[Redis 缓存网关: 已连接 (32 个活跃缓存项)]",
    statusLabel3: "[MySQL 数据库节点: 已迁移 (Schema 版本: v2.0)]",
    selectionTitle: "基础大语言模型选择",
    selectionSub: "选择用于代码规范与漏洞审计的基础大模型认知大脑",
    modelActive: "活动节点",
    modelEnabled: "已启用",
    modelSelect: "启用该节点",
    configTitle: "接口凭据与大模型参数偏好设置",
    configSub: "配置并指定目标大模型网关 API 基础 URL 以及各端授权密钥凭据",
    formLabelBaseUrl: "大模型接口网关 BASE URL",
    formLabelGithubToken: "GITHUB 个人访问令牌 (PAT)",
    formLabelApiKey: "大模型 API KEY (BEARER TOKEN)",
    formLabelModel: "指定的 MODEL ID (大模型请求参数)",
    formBtn: "保存并激活所有首选项配置",
    switchMsg: "已成功切换活动认知模型为: ",
    saveMsg: "基础大模型接口凭据已成功写入本地运行环境。"
  },
  en: {
    header: "ENGINE & CREDENTIALS CONTROLLER",
    subHeader: "AI model parameters selection, API configurations, and token authorizations",
    clock: "LLM ENGINE ONLINE",
    statusTitle: "ACTIVE CONNECTIVITY STATUS",
    statusLabel1: "[GitHub API Client: Token Active (5000 reqs/hr remaining)]",
    statusLabel2: "[Redis Cache Gateway: Connected (32 active entries)]",
    statusLabel3: "[MySQL Database Node: Migrated (Schema Version: v2.0)]",
    selectionTitle: "FOUNDATIONAL MODEL SELECTION",
    selectionSub: "Select the foundational LLM cognitive engine for code audit generations",
    modelActive: "ACTIVE NODE",
    modelEnabled: "ENABLED",
    modelSelect: "SELECT NODE",
    configTitle: "CREDENTIALS & CREDENTIALS PREFERENCES",
    configSub: "Assign the target LLM gateway base URL and API credential keys",
    formLabelBaseUrl: "LLM BASE URL (GATEWAY ENDPOINT)",
    formLabelGithubToken: "GITHUB PERSONAL ACCESS TOKEN",
    formLabelApiKey: "LLM API KEY (BEARER TOKEN)",
    formLabelModel: "ENFORCED MODEL ID (api model parameter)",
    formBtn: "SAVE PREFERENCES CONFIG",
    switchMsg: "Switched active cognitive model profile to: ",
    saveMsg: "Foundational LLM credentials successfully written to local environment."
  }
}

const selectedModelId = ref('deepseek-v3')
const baseUrl = ref('https://api.deepseek.com/v1')
const enforcedModel = ref('deepseek-chat')
const apiKey = ref('')
const githubToken = ref('')

const showApiKey = ref(false)
const showGithubToken = ref(false)

const models = ref([
  {
    id: 'deepseek-v3',
    provider: 'DEEPSEEK',
    name: 'DeepSeek-V3',
    desc: 'Cost-efficient foundational model. Highly compliant structured JSON outputs.',
    descZh: '极具性价比的基础大模型。具有高度合规的 JSON 结构化输出能力。',
    avgLatency: '1.2s'
  },
  {
    id: 'deepseek-r1',
    provider: 'DEEPSEEK REASONING',
    name: 'DeepSeek-R1 (Reasoning)',
    desc: 'Chain-of-thought advanced node. Ideal for scanning deep architectural flaws.',
    descZh: '深度思考推理模型。最适合检索和审查深层的系统架构设计缺陷。',
    avgLatency: '4.8s'
  }
])

const selectModel = (id) => {
  selectedModelId.value = id
  
  if (id === 'deepseek-v3') {
    baseUrl.value = 'https://api.deepseek.com/v1'
    enforcedModel.value = 'deepseek-chat'
  } else if (id === 'deepseek-r1') {
    baseUrl.value = 'https://api.deepseek.com/v1'
    enforcedModel.value = 'deepseek-reasoner'
  }

  localStorage.setItem('baseUrl', baseUrl.value)
  localStorage.setItem('enforcedModel', enforcedModel.value)
  localStorage.setItem('selectedModelId', id)

  ElMessage.success(`${t[currentLang.value].switchMsg}${id.toUpperCase()}`)
}

const saveConfig = () => {
  localStorage.setItem('baseUrl', baseUrl.value)
  localStorage.setItem('enforcedModel', enforcedModel.value)
  localStorage.setItem('apiKey', apiKey.value)
  localStorage.setItem('githubToken', githubToken.value)
  ElMessage.success(t[currentLang.value].saveMsg)
}

onMounted(() => {
  if (localStorage.getItem('selectedModelId')) {
    selectedModelId.value = localStorage.getItem('selectedModelId')
  } else {
    localStorage.setItem('selectedModelId', 'deepseek-v3')
  }

  baseUrl.value = localStorage.getItem('baseUrl') || 'https://api.deepseek.com/v1'
  enforcedModel.value = localStorage.getItem('enforcedModel') || 'deepseek-chat'
  apiKey.value = localStorage.getItem('apiKey') || ''
  githubToken.value = localStorage.getItem('githubToken') || ''
})
</script>

<style scoped>
.config-page {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

/* Header Banner */
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
  font-size: 11px;
  color: #00f0ff;
  font-weight: bold;
}

/* Card defaults */
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

/* Status Banner connectivity dots */
.status-banner {
  background: rgba(9, 15, 29, 0.5);
}

.status-nodes-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 15px;
}

.status-node-pill {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
}

.pulse-dot-green {
  width: 7px;
  height: 7px;
  background-color: #00ff66;
  border-radius: 50%;
  box-shadow: 0 0 8px #00ff66, 0 0 15px #00ff66;
  animation: pulse 2s infinite;
}

.status-node-pill .label {
  color: #00ff66;
  text-shadow: 0 0 8px rgba(0, 255, 102, 0.15);
}

/* Model selection layouts */
.model-selection-wrapper {
  margin-top: 5px;
}

.section-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: bold;
  color: #cbd5e1;
}

.section-header .subtitle {
  margin: 5px 0 20px 0;
  color: #475569;
  font-size: 12.5px;
}

.models-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.model-card {
  position: relative;
  background: rgba(13, 20, 35, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 190px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.model-card:hover {
  border-color: rgba(0, 240, 255, 0.3);
  transform: translateY(-4px);
}

.model-card.active {
  border-color: #00f0ff;
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.15);
  background: linear-gradient(180deg, rgba(0, 240, 255, 0.03) 0%, transparent 100%);
}

.active-banner {
  position: absolute;
  top: -1px;
  right: 20px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 8.5px;
  background: #00f0ff;
  color: #080c14;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 0 0 4px 4px;
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.4);
}

.model-provider {
  font-family: 'JetBrains Mono', monospace;
  font-size: 9px;
  color: #475569;
  font-weight: bold;
  letter-spacing: 0.8px;
}

.model-title {
  margin: 6px 0 8px 0;
  font-size: 14.5px;
  font-weight: bold;
  color: #f1f5f9;
}

.model-desc {
  margin: 0;
  font-size: 11.5px;
  color: #64748b;
  line-height: 1.5;
}

.model-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.03);
  padding-top: 10px;
}

.latency-pill {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.select-text {
  font-family: 'JetBrains Mono', monospace;
  font-size: 10px;
  font-weight: bold;
  color: #64748b;
}

.active .select-text {
  color: #00f0ff;
}

/* Credentials Forms */
.credentials-wrapper {
  margin-top: 10px;
}

.config-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 25px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  color: #475569;
  letter-spacing: 0.8px;
}

.form-input {
  background: rgba(10, 15, 30, 0.8);
  border: 1px solid rgba(0, 240, 255, 0.15);
  border-radius: 8px;
  color: #f8fafc;
  font-size: 13.5px;
  padding: 10px 14px;
  outline: none;
  transition: all 0.3s;
}

.form-input:focus {
  border-color: #00f0ff;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.15);
}

.password-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input-wrapper .form-input {
  width: 100%;
  padding-right: 45px;
}

.eye-btn {
  position: absolute;
  right: 12px;
  background: transparent;
  border: none;
  color: #475569;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.eye-btn:hover {
  color: #00f0ff;
}

.actions-row {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.save-btn {
  position: relative;
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  padding: 12px 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.save-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.3);
}

.btn-text {
  position: relative;
  color: #ffffff;
  font-weight: bold;
  font-size: 12px;
  letter-spacing: 1px;
  z-index: 2;
}

.btn-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%);
  filter: blur(8px);
  opacity: 0.3;
  z-index: 1;
}

.animate-fade-in {
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 0.9; }
  50% { transform: scale(1.1); opacity: 1; box-shadow: 0 0 10px #00ff66, 0 0 18px #00ff66; }
  100% { transform: scale(1); opacity: 0.9; }
}
</style>
