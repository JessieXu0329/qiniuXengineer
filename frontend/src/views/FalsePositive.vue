<template>
  <div class="filter-page animate-fade-in">
    <!-- Top Header -->
    <div class="welcome-banner">
      <div class="welcome-text">
        <h2>{{ t('filter.header') }}</h2>
        <p>{{ t('filter.subHeader') }}</p>
      </div>
      <div class="system-time">
        <span class="pulse-green"></span>
        <span class="clock">{{ t('filter.clock') }}</span>
      </div>
    </div>

    <!-- Active Metric Logs Panel -->
    <div class="cyber-card metric-card">
      <div class="card-glow"></div>
      <div class="card-content metric-banner">
        <div class="metric-block">
          <span class="metric-title">{{ t('filter.bannerTitle') }}</span>
          <p class="metric-desc" v-html="t('filter.bannerDesc')"></p>
        </div>
        <div class="mini-stats">
          <div class="ms-pill">
            <span class="val">84%</span>
            <span class="lbl">{{ t('filter.miniStatLbl') }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="filter-grid">
      <!-- Left: Active Rules Table -->
      <div class="filter-left">
        <div class="cyber-card table-card">
          <div class="card-glow"></div>
          <div class="card-content">
            <h3>{{ t('filter.tableTitle') }}</h3>
            <p class="subtitle">{{ t('filter.tableSub') }}</p>
            
            <div class="rules-table-container">
              <table class="cyber-table">
                <thead>
                  <tr>
                    <th>{{ t('filter.colName') }}</th>
                    <th>{{ t('filter.colPattern') }}</th>
                    <th>{{ t('filter.colType') }}</th>
                    <th class="center">{{ t('filter.colActions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(rule, index) in rules" :key="index">
                    <td>
                      <div class="rule-name">{{ lang === 'zh' ? rule.nameZh : rule.name }}</div>
                      <div class="rule-desc">{{ lang === 'zh' ? rule.descZh : rule.desc }}</div>
                    </td>
                    <td><code class="match-code">{{ rule.pattern }}</code></td>
                    <td><span :class="['severity-badge', rule.type]">{{ t('filter.types.' + rule.type) }}</span></td>
                    <td class="center">
                      <button @click="deleteRule(index)" class="action-btn delete-btn">DELETE</button>
                    </td>
                  </tr>
                  <tr v-if="rules.length === 0">
                    <td colspan="4" class="empty-row">{{ t('filter.emptyMsg') }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Add Custom Wizard Form -->
      <div class="filter-right">
        <div class="cyber-card form-card">
          <div class="card-glow"></div>
          <div class="card-content">
            <h3>{{ t('filter.formTitle') }}</h3>
            <p class="subtitle">{{ t('filter.formSub') }}</p>

            <form @submit.prevent="saveRule" class="cyber-form">
              <div class="form-group">
                <label>{{ t('filter.formLabelName') }}</label>
                <input
                  v-model="newRule.name"
                  type="text"
                  :placeholder="t('filter.placeholderName')"
                  required
                  class="form-input"
                />
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label>{{ t('filter.formLabelPattern') }}</label>
                  <input
                    v-model="newRule.pattern"
                    type="text"
                    placeholder="e.g. **/*_test.go"
                    required
                    class="form-input code"
                  />
                </div>

                <div class="form-group">
                  <label>{{ t('filter.formLabelType') }}</label>
                  <select v-model="newRule.type" class="form-select">
                    <option value="style">{{ t('filter.selectOptions.style') }}</option>
                    <option value="lint">{{ t('filter.selectOptions.lint') }}</option>
                    <option value="comments">{{ t('filter.selectOptions.comments') }}</option>
                    <option value="all">{{ t('filter.selectOptions.all') }}</option>
                  </select>
                </div>
              </div>

              <div class="form-group">
                <label>{{ t('filter.formLabelDesc') }}</label>
                <textarea
                  v-model="newRule.desc"
                  :placeholder="t('filter.placeholderDesc')"
                  rows="3"
                  required
                  class="form-textarea"
                ></textarea>
              </div>

              <button type="submit" class="form-btn">
                <span class="btn-glow"></span>
                <span class="btn-text">{{ t('filter.formBtn') }}</span>
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/i18n'
import { ElMessage } from 'element-plus'

const { t, lang } = useI18n()

const rules = ref([
  {
    name: "Ignore test assertions",
    nameZh: "忽略单元测试断言",
    desc: "Suppress GORM connection faults in mock testing and unit suites.",
    descZh: "在 Mock 测试和单测用例中忽略 GORM 连接失败缺陷。",
    pattern: "**/*_test.go",
    type: "lint"
  },
  {
    name: "Generated files bypass",
    nameZh: "跳过自动生成文件",
    desc: "Bypass style comments on files automatically generated by protobuf systems.",
    descZh: "跳过 Protobuf 自动生成的文件的代码样式 and 格式警告。",
    pattern: "**/*.pb.go",
    type: "style"
  },
  {
    name: "Comment modifications ignore",
    nameZh: "忽略纯注释变更",
    desc: "Do not trigger review comments on pure comment block edits.",
    descZh: "不针对纯注释块的编辑和变动触发任何评审意见。",
    pattern: "**/*.go",
    type: "comments"
  }
])

const newRule = ref({
  name: '',
  pattern: '',
  type: 'style',
  desc: ''
})

const saveRule = () => {
  rules.value.push({
    name: newRule.value.name,
    nameZh: newRule.value.name,
    pattern: newRule.value.pattern,
    type: newRule.value.type,
    desc: newRule.value.desc,
    descZh: newRule.value.desc
  })
  
  ElMessage.success(
    `${lang.value === 'zh' ? '过滤规则' : 'Suppression rule'} '${newRule.value.name}' ${t('filter.successMsg')}`
  )
  
  // Reset form
  newRule.value = {
    name: '',
    pattern: '',
    type: 'style',
    desc: ''
  }
}

const deleteRule = (index) => {
  const deletedName = lang.value === 'zh' ? (rules.value[index].nameZh || rules.value[index].name) : rules.value[index].name
  rules.value.splice(index, 1)

  ElMessage.info(
    `${lang.value === 'zh' ? '过滤规则' : 'Suppression rule'} '${deletedName}' ${t('filter.removedMsg')}`
  )
}
</script>

<style scoped>
.filter-page {
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
  background: rgba(0, 255, 102, 0.05);
  border: 1px solid rgba(0, 255, 102, 0.2);
  padding: 6px 14px;
  border-radius: 20px;
}

.pulse-green {
  width: 6px;
  height: 6px;
  background-color: #00ff66;
  border-radius: 50%;
  box-shadow: 0 0 8px #00ff66;
}

.clock {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  color: #00ff66;
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

/* Metric Banner Logs */
.metric-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(90deg, rgba(0, 255, 102, 0.03) 0%, transparent 100%);
  border-left: 2px solid rgba(0, 255, 102, 0.3);
}

.metric-title {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12.5px;
  color: #00ff66;
  font-weight: bold;
  text-shadow: 0 0 10px rgba(0, 255, 102, 0.2);
}

.metric-desc {
  margin: 6px 0 0 0;
  font-size: 13.5px;
  color: #94a3b8;
  line-height: 1.5;
}

.mini-stats {
  display: flex;
}

.ms-pill {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.ms-pill .val {
  font-size: 24px;
  font-weight: 900;
  color: #00ff66;
  font-family: 'JetBrains Mono', monospace;
}

.ms-pill .lbl {
  font-size: 10px;
  color: #475569;
  font-family: 'JetBrains Mono', monospace;
  margin-top: 4px;
}

/* Grids Layout */
.filter-grid {
  display: grid;
  grid-template-columns: 1fr 440px;
  gap: 25px;
}

/* Cyber Rules Table Styling */
.rules-table-container {
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

.rule-name {
  font-size: 13.5px;
  font-weight: bold;
  color: #f1f5f9;
}

.rule-desc {
  font-size: 11.5px;
  color: #64748b;
  margin-top: 4px;
}

.match-code {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11.5px;
  color: #00f0ff;
  background: rgba(0, 240, 255, 0.05);
  border: 1px solid rgba(0, 240, 255, 0.15);
  padding: 3px 8px;
  border-radius: 4px;
}

.severity-badge {
  font-family: 'JetBrains Mono', monospace;
  font-size: 9px;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 4px;
  text-transform: uppercase;
}

.severity-badge.style { background: rgba(56, 189, 248, 0.1); color: #38bdf8; }
.severity-badge.lint { background: rgba(234, 179, 8, 0.1); color: #eab308; }
.severity-badge.comments { background: rgba(168, 85, 247, 0.1); color: #c084fc; }

.action-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 6px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.delete-btn {
  font-family: 'JetBrains Mono', monospace;
  font-size: 10px;
  font-weight: bold;
  letter-spacing: 1px;
  color: #ef4444;
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.4);
  padding: 5px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: #ef4444;
  color: #fca5a5;
  box-shadow: 0 0 10px rgba(239, 68, 68, 0.25);
}

.center {
  text-align: center;
}

.empty-row {
  color: #475569;
  font-size: 12.5px;
  text-align: center;
  padding: 30px 0;
  font-style: italic;
}

/* Wizard Form Styling */
.cyber-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
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

.form-input, .form-select, .form-textarea {
  background: rgba(10, 15, 30, 0.8);
  border: 1px solid rgba(0, 240, 255, 0.15);
  border-radius: 8px;
  color: #f8fafc;
  font-size: 13.5px;
  padding: 10px 14px;
  outline: none;
  transition: all 0.3s;
}

.form-input:focus, .form-select:focus, .form-textarea:focus {
  border-color: #00f0ff;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.15);
}

.form-input.code {
  font-family: 'JetBrains Mono', monospace;
  color: #00ff66;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 150px;
  gap: 15px;
}

.form-textarea {
  resize: vertical;
}

.form-btn {
  position: relative;
  background: linear-gradient(135deg, #00f0ff 0%, #7000ff 100%);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  margin-top: 10px;
}

.form-btn:hover {
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
</style>
