import { ref } from 'vue'

// ── Shared reactive locale ──────────────────────────────────────────
const locale = ref('zh')

// ── Full translation catalog ────────────────────────────────────────
const messages = {
  zh: {
    // ── App frame / Navbar ──────────────────────────────────────────
    app: {
      title: 'AI 代码评审助手',
      api: 'API 服务',
      db: '数据库',
      latency: '系统延迟',
      role: '系统超级管理员',
      sec1: '核心控制中心',
      dashboard: 'Dashboard 大盘',
      review: 'PR Review 智能评审',
      context: 'Hybrid Context 关联项',
      filter: 'False Positive 过滤',
      config: 'AI 引擎配置',
      tokens: 'Access Tokens 令牌管理',
      sec2: '预设与引擎配置',
      statusOnline: 'ONLINE',
      statusOffline: 'OFFLINE',
      statusConnecting: 'CONNECTING',
      toggleLang: 'ENGLISH'
    },

    // ── Dashboard ───────────────────────────────────────────────────
    dashboard: {
      header: '大盘监控中心',
      subHeader: '项目平均健康指标、安全规范趋势以及扫描足迹大盘',
      radarTitle: '全局代码质量指标',
      radarSub: '所有历史代码审查得分的综合平均值向量',
      tableTitle: '近期 PR 自动审计日志',
      tableSub: '保存在 MySQL 数据库中的最近代码扫描足迹',
      col1: '代码合并请求节点',
      col2: '安全规范评分',
      col3: '审计时间',
      col4: '状态',
      clock: '系统监控运行中',
      radarTooltipTitle: '全局平均指标',
      radarSeriesName: '历史平均数据',
      // Radar axis indicators
      t1: '规范度',
      t2: '安全性',
      t3: '性能',
      t4: '可读性',
      // Stat cards
      stats: {
        totalReviews: { label: '累计审计PR完成数', trend: '▲ 比上周增长 +18%' },
        averageScore: { label: '平均代码健康评分', trend: '▲ 比上周质量提升 +1.2%' },
        criticalVulns: { label: '累计发现的致命漏洞', trend: '▼ 历史漏洞降低 -34%' },
        contextRetrievals: { label: '关联上下文拉取次数', trend: '▲ 新解析关联依赖 +156' }
      }
    },

    // ── PR Review ───────────────────────────────────────────────────
    review: {
      header: '智能 PR 评审控制台',
      subHeader: '对单个代码合并请求执行深层 AST 和大模型混合上下文安全审计',
      inputTitle: '安全 PR 代码扫描器',
      inputSub: '输入 GitHub PR 链接以开启混合上下文代码深度审计',
      btnText: '开启 AI 智能审计',
      radarTitle: '漏洞与质量度量雷达',
      s1: '综合得分',
      s2: '安全性',
      s3: '性能评分',
      summaryTitle: 'PR 变更结构化总结',
      sum1: '一句话改动概述',
      sum2: '受影响模块清单',
      sum3: '检测到破坏性变更 (Breaking Changes)',
      cardsHeader: '行内智能评审批注卡片',
      issuesBadge: '个缺陷/优化点已发现',
      suggestCode: '推荐优化修复代码',
      t1: '规范度',
      t2: '安全性',
      t3: '性能',
      t4: '可读性',
      // Level labels
      levels: {
        critical: '致命',
        warning: '警告',
        suggestion: '建议'
      },
      // Misc
      clockReady: '扫描引擎就绪',
      engineLabel: '审计引擎大脑:',
      line: '行',
      copy: '复制',
      copySuccess: '复制成功!',
      exportPdf: '导出 PDF 报告',
      exportZh: '中文 PDF 报告',
      exportEn: 'English PDF 报告',
      radarTooltipTitle: '当前 PR 审计指标',
      radarSeriesName: '安全质量维度',
      successMsg: 'AI 审计成功完成！',
      fallbackMsg: 'Go 后端离线。运行客户端模拟审计引擎。',
      validation: { emptyUrl: '请输入有效的 GitHub PR 链接。' },
      error: { auditFailed: '审计失败: ', fetchFailed: '获取历史审计记录失败。' }
    },

    // ── Hybrid Context ──────────────────────────────────────────────
    context: {
      header: '混合上下文检索器',
      subHeader: 'AST 语法解析、包路径规范约束发现以及超出 Diff 代码的关联认知叠加图层',
      status: '认知图层已激活',
      bannerTitle: '[当前扫描上下文配置文件: 认知图层已激活]',
      bannerDesc: '分析器动态检索您的代码仓库边界。它不孤立地审查 Diff 行，而是提取父类、接口和包配置，为大模型评审流水线注入高精度的上下文背景参数。',
      statLabel1: '已注入上下文锚点数',
      statLabel2: '上下文精准度评分',
      timelineTitle: '动态上下文注入时间线',
      timelineSub: '在最近一次扫描期间执行依赖解析的按时序记录',
      graphTitle: 'AST 逻辑依赖连通网络与跨文件影响图',
      graphSub: '实时 AST 关系连通性与跨文件调用链追踪图。点击节点可审查深层依赖安全指纹签名。',
      inspectTitle: '节点详细依赖指纹签名',
      colPath: '文件绝对路径',
      colExports: '已发现导出 API',
      colLinks: '关联网结依赖',
      // Labels
      nodeTypeLabel: '角色类型',
      impactLabel: '安全影响剖析',
      tooltipEdge: '跨文件 AST 调用依赖',
      // Injection type badges
      injectionTypes: {
        autoInjected: '自动注入',
        compilationScope: '编译边界',
        parentArchitecture: '父级架构'
      },
      // Node types
      nodeTypes: {
        modified: '变更文件 (变动源头)',
        directDownstream: '直接下游 (直接关联影响)',
        dependencyProvider: '依赖提供 (数据流来源)',
        indirectDownstream: '间接下游 (跨文件引用)',
        indirectDownstream2: '间接下游 (路由控制)',
        systemEntry: '系统主入口 (最上游)'
      },
      // Impact levels
      impact: {
        high: '高风险',
        medium: '中等风险',
        low: '低风险'
      }
    },

    // ── False Positive Engine ───────────────────────────────────────
    filter: {
      header: '误报规避过滤引擎',
      subHeader: '噪声过滤器、样式忽略规则以及认知忽略范围控制面板',
      clock: '今天已节省 4.2 分钟专注时间',
      bannerTitle: '[上次评审扫描报告: 噪声规避引擎活动中]',
      bannerDesc: '在您最近一次 PR 评审扫描期间，有 <strong>3 个误报</strong>已被您的过滤引擎边界自动拦截。这成功为开发人员节省了 <strong>4.2 分钟</strong> 的代码审查专注时间。',
      miniStatLbl: '扫描噪音已净化率',
      tableTitle: '当前生效的忽略过滤规则',
      tableSub: '在保存评审意见前，与 AST 审查批注执行匹配的匹配拦截策略',
      colName: '规则标识名',
      colPattern: '匹配路径匹配器',
      colType: '忽略分类',
      colActions: '操作',
      emptyMsg: '未配置任何忽略过滤器。系统扫描处于全量批注输出模式。',
      formTitle: '添加过滤规避规则',
      formSub: '部署自定义忽略匹配参数，定向规避特定文件或包路径的非关键批注',
      formLabelName: '规则标识名称',
      formLabelPattern: '范围路径匹配器 (GLOB / 正则)',
      formLabelType: '忽略分类',
      formLabelDesc: '过滤理由与详细描述说明',
      placeholderName: '例如: 忽略单元测试环境中的日志断言',
      placeholderDesc: '阐述为什么应该在审计合规中忽略这些行内批注...',
      formBtn: '部署并激活该规则',
      successMsg: '已成功部署到过滤引擎。',
      removedMsg: '已从引擎中成功移除。',
      // Type classification
      types: {
        style: '代码风格',
        lint: 'Lint静态检查',
        comments: '仅注释',
        all: '全部类型'
      },
      selectOptions: {
        style: '代码风格 / 格式化',
        lint: 'Lint 静态检查警告',
        comments: '仅注释修改',
        all: '所有发现的缺陷'
      }
    },

    // ── Common ──────────────────────────────────────────────────────
    common: {
      connected: 'CONNECTED',
      serverError: '服务器错误'
    }
  },

  // ── English ───────────────────────────────────────────────────────
  en: {
    app: {
      title: 'AI PR REVIEWER',
      api: 'API SERVER',
      db: 'DB NODE',
      latency: 'LATENCY',
      role: 'Admin Operator',
      sec1: 'METRICS & CONTROL',
      dashboard: 'Dashboard Index',
      review: 'PR Review Audit',
      context: 'Hybrid Context Retriever',
      filter: 'False Positive Engine',
      config: 'AI Engine Config',
      tokens: 'Access Tokens',
      sec2: 'PRESETS & ENGINE',
      statusOnline: 'ONLINE',
      statusOffline: 'OFFLINE',
      statusConnecting: 'CONNECTING',
      toggleLang: '中文'
    },

    dashboard: {
      header: 'METRICS CORE DASHBOARD',
      subHeader: 'Aggregate repo health, security compliance trends, and scanning footprints',
      radarTitle: 'GLOBAL CODE COMPLIANCE METRICS',
      radarSub: 'Aggregated average score vectors computed from all historical audits',
      tableTitle: 'RECENT PR CODE AUDIT LOGS',
      tableSub: 'Latest pull requests scanned and recorded inside GORM MySQL database',
      col1: 'Pull Request Node',
      col2: 'Compliance Score',
      col3: 'Audit Date',
      col4: 'Status',
      clock: 'SYSTEM MONITORING',
      radarTooltipTitle: 'Global Avg Metrics',
      radarSeriesName: 'Historical Averages',
      t1: 'Standards',
      t2: 'Security',
      t3: 'Performance',
      t4: 'Readability',
      stats: {
        totalReviews: { label: 'TOTAL REVIEWS COMPLETED', trend: '▲ +18% from last week' },
        averageScore: { label: 'AVERAGE COMPLIANCE SCORE', trend: '▲ +1.2% quality gains' },
        criticalVulns: { label: 'CRITICAL VULNERABILITIES DETECTED', trend: '▼ -34% suppression trend' },
        contextRetrievals: { label: 'COGNITIVE CONTEXT RETRIEVALS', trend: '▲ +156 dependencies bound' }
      }
    },

    review: {
      header: 'INTELLIGENT PR SCANNER',
      subHeader: 'Perform deep AST & LLM hybrid context audits on individual Pull Requests',
      inputTitle: 'SECURE PR CODE SCANNER',
      inputSub: 'Input a GitHub Pull Request URL to perform deep AST and LLM hybrid context reviews',
      btnText: 'TRIGGER AI AUDIT',
      radarTitle: 'VULNERABILITY & METRICS RADAR',
      s1: 'OVERALL',
      s2: 'SECURITY',
      s3: 'PERFORMANCE',
      summaryTitle: 'PR STRUCTURAL SUMMARY',
      sum1: 'ONE-SENTENCE OVERVIEW',
      sum2: 'AFFECTED MODULES',
      sum3: 'BREAKING CHANGES DETECTED',
      cardsHeader: 'INLINE SCHEMATIC AUDIT CARDS',
      issuesBadge: 'ISSUES FOUND',
      suggestCode: 'SUGGESTED OPTIMIZATION',
      t1: 'Standards',
      t2: 'Security',
      t3: 'Performance',
      t4: 'Readability',
      levels: {
        critical: 'CRITICAL',
        warning: 'WARNING',
        suggestion: 'SUGGESTION'
      },
      clockReady: 'SCAN ENGINE READY',
      engineLabel: 'AUDIT ENGINE BRAIN:',
      line: 'Line',
      copy: 'COPY',
      copySuccess: 'SUCCESS!',
      exportPdf: 'EXPORT PDF',
      exportZh: '中文 PDF 报告',
      exportEn: 'English PDF Report',
      radarTooltipTitle: 'Current PR Metrics',
      radarSeriesName: 'Vulnerabilities',
      successMsg: 'AI Audit completed successfully!',
      fallbackMsg: 'Go backend offline. Running client mock review engine.',
      validation: { emptyUrl: 'Please input a valid GitHub PR URL.' },
      error: { auditFailed: 'Audit Failed: ', fetchFailed: 'Failed to fetch historical audit logs.' }
    },

    context: {
      header: 'HYBRID CONTEXT RETRIEVER',
      subHeader: 'AST parsing, package constraints discovery, and out-of-diff cognitive overlays',
      status: 'COGNITIVE OVERLAY ACTIVE',
      bannerTitle: '[CURRENT SCAN CONTEXT PROFILE: COGNITIVE OVERLAY ACTIVE]',
      bannerDesc: 'The retriever scans your repository boundaries dynamically. Instead of isolating Diff lines, it extracts parent classes, interfaces, and packages configurations to feed high-fidelity context parameters to the LLM review pipeline.',
      statLabel1: 'Context Anchors Injected',
      statLabel2: 'Context Accuracy Score',
      timelineTitle: 'DYNAMIC CONTEXT INJECTION TIMELINE',
      timelineSub: 'Chronological footprint of dependency resolutions executed on last scanning event',
      graphTitle: 'AST LOGICAL DEPENDENCY LINKER & IMPACT SCOPE',
      graphSub: 'Live AST connection tree and cross-file call chains tracker. Click nodes to inspect detailed dependency fingerprints.',
      inspectTitle: 'NODE DETAILED SIGNATURE',
      colPath: 'Absolute Path',
      colExports: 'Discovered Exports',
      colLinks: 'Dependency Linkages',
      nodeTypeLabel: 'NodeType',
      impactLabel: 'Security Impact Analysis',
      tooltipEdge: 'Cross-file AST Dependency',
      injectionTypes: {
        autoInjected: 'AUTO-INJECTED',
        compilationScope: 'COMPILATION SCOPE',
        parentArchitecture: 'PARENT ARCHITECTURE'
      },
      nodeTypes: {
        modified: 'Modified File (Change Source)',
        directDownstream: 'Direct Downstream (Direct Impact)',
        dependencyProvider: 'Dependency Provider (Data Source)',
        indirectDownstream: 'Indirect Downstream (Cross-file)',
        indirectDownstream2: 'Indirect Downstream (Route Control)',
        systemEntry: 'System Entry (Top Level)'
      },
      impact: {
        high: 'HIGH',
        medium: 'MEDIUM',
        low: 'LOW'
      }
    },

    filter: {
      header: 'FALSE POSITIVE ENGINE',
      subHeader: 'Noise filters, styling suppressions, and cognitive ignore rules',
      clock: 'SAVED 4.2 MIN FOCUS TODAY',
      bannerTitle: '[LAST REVIEW SESSION REPORT: NOISE FILTER ENGINE]',
      bannerDesc: 'During your last PR review scan, <strong>3 False Positives</strong> were automatically filtered and suppressed by your active config boundaries. This successfully saved <strong>4.2 minutes</strong> of engineer code review focus.',
      miniStatLbl: 'Audits Noise Cleansed',
      tableTitle: 'ACTIVE CONTEXT SUPPRESSION RULES',
      tableSub: 'Enforced pattern suppressions matched against AST reviews before persisting comments',
      colName: 'Rule Identifier',
      colPattern: 'Scope Matcher',
      colType: 'Type',
      colActions: 'Actions',
      emptyMsg: 'No active ignore filters configured. System scanning is at full verbosity.',
      formTitle: 'ADD CONTEXT SUPPRESS FILTER',
      formSub: 'Deploy a custom suppression parameter to target specific file boundaries or packages',
      formLabelName: 'RULE NAME',
      formLabelPattern: 'SCOPE MATCHER (GLOB / REGEX)',
      formLabelType: 'SUPPRESS TYPE',
      formLabelDesc: 'DESCRIPTION & REASON',
      placeholderName: 'e.g. Ignore testing logger assertions',
      placeholderDesc: 'Explain why these comments should be suppressed for audit compliance...',
      formBtn: 'DEPLOY SUPPRESSION FILTER',
      successMsg: 'successfully deployed to filter engine.',
      removedMsg: 'has been removed.',
      types: {
        style: 'Style',
        lint: 'Lint',
        comments: 'Comments',
        all: 'All'
      },
      selectOptions: {
        style: 'Style / Formatting',
        lint: 'Linting Warnings',
        comments: 'Comments Only',
        all: 'All Findings'
      }
    },

    common: {
      connected: 'CONNECTED',
      serverError: 'Server error'
    }
  }
}

// ── Composable ──────────────────────────────────────────────────────
export function useI18n() {
  /**
   * Look up a dotted-path key in the current locale.
   * Returns the key itself as fallback when the path is missing.
   *
   * @example t('dashboard.header')   => '大盘监控中心' | 'METRICS CORE DASHBOARD'
   * @example t('review.line')        => '行' | 'Line'
   */
  const t = (key) => {
    const keys = key.split('.')
    let val = messages[locale.value]
    for (const k of keys) {
      if (val == null || val[k] === undefined) return key
      val = val[k]
    }
    return val
  }

  return { t, lang: locale }
}

export { locale }
