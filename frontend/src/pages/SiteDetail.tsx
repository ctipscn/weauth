import { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { sitesApi } from '@/lib/api'
import { 
  ArrowLeft, 
  Copy, 
  Check, 
  Plus, 
  Trash2, 
  RefreshCw,
  BarChart3,
  Globe,
  Shield,
  TrendingUp,
  CheckCircle2,
  XCircle,
  Activity,
  Layers,
  Zap,
  Info,
  Save,
  RotateCcw
} from 'lucide-react'

// 开关组件
function Toggle({ 
  checked, 
  onChange, 
  disabled 
}: { 
  checked: boolean
  onChange: (value: boolean) => void
  disabled?: boolean 
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      onClick={() => !disabled && onChange(!checked)}
      disabled={disabled}
      className={`
        relative inline-flex h-6 w-11 flex-shrink-0 rounded-full border-2 border-transparent 
        transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
        ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
        ${checked ? 'bg-emerald-600' : 'bg-gray-200'}
      `}
    >
      <span
        className={`
          pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow-lg ring-0 
          transition-transform duration-200 ease-in-out
          ${checked ? 'translate-x-5' : 'translate-x-0'}
        `}
      />
    </button>
  )
}

// 设置项组件
function SettingRow({ 
  title, 
  description, 
  children,
  noBorder
}: { 
  title: React.ReactNode
  description?: string
  children: React.ReactNode
  noBorder?: boolean
}) {
  return (
    <div className={`flex items-center justify-between py-4 ${!noBorder ? 'border-b border-slate-100' : ''}`}>
      <div className="flex-1 pr-4">
        <div className="text-sm font-medium text-slate-900">{title}</div>
        {description && (
          <div className="text-xs text-slate-500 mt-0.5">{description}</div>
        )}
      </div>
      <div className="flex-shrink-0">{children}</div>
    </div>
  )
}

// Section 卡片组件（带保存/取消按钮）
function SettingSection({
  title,
  description,
  children,
  isDirty,
  onSave,
  onCancel,
  isSaving
}: {
  title: string
  description?: string
  children: React.ReactNode
  isDirty: boolean
  onSave: () => void
  onCancel: () => void
  isSaving?: boolean
}) {
  return (
    <section className="bg-white rounded-xl border border-slate-200 overflow-hidden">
      <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-slate-900">{title}</h2>
          {description && <p className="text-xs text-slate-500 mt-1">{description}</p>}
        </div>
        {isDirty && (
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={onCancel}
              className="h-8 px-3 text-sm text-slate-600 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition-colors flex items-center gap-1.5"
            >
              <RotateCcw className="w-3.5 h-3.5" />
              取消
            </button>
            <button
              type="button"
              onClick={onSave}
              disabled={isSaving}
              className="h-8 px-3 text-sm bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 rounded-lg transition-colors flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Save className="w-3.5 h-3.5" />
              {isSaving ? '保存中...' : '保存'}
            </button>
          </div>
        )}
      </div>
      <div className="px-6">{children}</div>
    </section>
  )
}

// 输入框组件
function SettingInput({
  value,
  onChange,
  type = 'text',
  placeholder,
  className = 'w-32',
  min,
  max,
  step
}: {
  value: string | number
  onChange: (value: string) => void
  type?: string
  placeholder?: string
  className?: string
  min?: number
  max?: number
  step?: number
}) {
  return (
    <input
      type={type}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={placeholder}
      min={min}
      max={max}
      step={step}
      className={`h-9 px-3 rounded-lg border border-slate-200 bg-white text-sm text-slate-900 
        focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 
        transition-all ${className}`}
    />
  )
}

// 下拉选择组件
function SettingSelect({
  value,
  onChange,
  options,
  className = 'w-48'
}: {
  value: string | number
  onChange: (value: string) => void
  options: { value: string | number; label: string }[]
  className?: string
}) {
  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className={`h-9 px-3 rounded-lg border border-slate-200 bg-white text-sm text-slate-900 
        focus:outline-none focus:ring-2 focus:ring-violet-500/20 focus:border-violet-500 
        transition-all ${className}`}
    >
      {options.map((opt) => (
        <option key={opt.value} value={opt.value}>{opt.label}</option>
      ))}
    </select>
  )
}

// 基本设置表单状态
interface BasicFormState {
  name: string
  enabled: boolean
}

// 验证模式表单状态
interface ModeFormState {
  batch_mode_enabled: boolean
  batch_difficulty: number
  min_batch_count: number
  max_batch_count: number
  base_difficulty: number
  max_difficulty: number
}

// 风控策略表单状态
interface PolicyFormState {
  enabled: boolean
  time_window_seconds: number
  request_threshold: number
  batch_increment: number
  difficulty_increment: number
  failure_weight: number
}

export default function SiteDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const siteId = parseInt(id!)

  const [copiedField, setCopiedField] = useState<string | null>(null)
  const [newDomain, setNewDomain] = useState('')
  const [showSecretConfirm, setShowSecretConfirm] = useState(false)
  const [newSecret, setNewSecret] = useState<string | null>(null)
  const [showViewSecret, setShowViewSecret] = useState(false)
  const [viewingSecret, setViewingSecret] = useState<string | null>(null)

  // 表单状态
  const [basicForm, setBasicForm] = useState<BasicFormState | null>(null)
  const [modeForm, setModeForm] = useState<ModeFormState | null>(null)
  const [policyForm, setPolicyForm] = useState<PolicyFormState | null>(null)

  const { data: site, isLoading } = useQuery({
    queryKey: ['site', siteId],
    queryFn: () => sitesApi.get(siteId),
  })

  const { data: ipPolicy } = useQuery({
    queryKey: ['site-policy', siteId],
    queryFn: () => sitesApi.getIpPolicy(siteId),
  })

  const { data: ipRecords } = useQuery({
    queryKey: ['site-ip-records', siteId],
    queryFn: () => sitesApi.listIpRecords(siteId),
  })

  const { data: stats } = useQuery({
    queryKey: ['site-stats', siteId],
    queryFn: () => sitesApi.getStats(siteId),
  })

  // 初始化表单状态
  useEffect(() => {
    if (site && !basicForm) {
      setBasicForm({
        name: site.name,
        enabled: site.enabled,
      })
    }
    if (site && !modeForm) {
      setModeForm({
        batch_mode_enabled: site.batch_mode_enabled,
        batch_difficulty: site.batch_difficulty,
        min_batch_count: site.min_batch_count,
        max_batch_count: site.max_batch_count,
        base_difficulty: site.base_difficulty,
        max_difficulty: site.max_difficulty,
      })
    }
  }, [site, basicForm, modeForm])

  useEffect(() => {
    if (ipPolicy && !policyForm) {
      setPolicyForm({
        enabled: ipPolicy.enabled,
        time_window_seconds: ipPolicy.time_window_seconds,
        request_threshold: ipPolicy.request_threshold,
        batch_increment: ipPolicy.batch_increment || 1,
        difficulty_increment: ipPolicy.difficulty_increment,
        failure_weight: ipPolicy.failure_weight,
      })
    }
  }, [ipPolicy, policyForm])

  // 检测是否有未保存的更改
  const isBasicDirty = useCallback(() => {
    if (!site || !basicForm) return false
    return basicForm.name !== site.name || basicForm.enabled !== site.enabled
  }, [site, basicForm])

  const isModeDirty = useCallback(() => {
    if (!site || !modeForm) return false
    return (
      modeForm.batch_mode_enabled !== site.batch_mode_enabled ||
      modeForm.batch_difficulty !== site.batch_difficulty ||
      modeForm.min_batch_count !== site.min_batch_count ||
      modeForm.max_batch_count !== site.max_batch_count ||
      modeForm.base_difficulty !== site.base_difficulty ||
      modeForm.max_difficulty !== site.max_difficulty
    )
  }, [site, modeForm])

  const isPolicyDirty = useCallback(() => {
    if (!ipPolicy || !policyForm) return false
    return (
      policyForm.enabled !== ipPolicy.enabled ||
      policyForm.time_window_seconds !== ipPolicy.time_window_seconds ||
      policyForm.request_threshold !== ipPolicy.request_threshold ||
      policyForm.batch_increment !== (ipPolicy.batch_increment || 1) ||
      policyForm.difficulty_increment !== ipPolicy.difficulty_increment ||
      policyForm.failure_weight !== ipPolicy.failure_weight
    )
  }, [ipPolicy, policyForm])

  // 重置表单
  const resetBasicForm = useCallback(() => {
    if (site) {
      setBasicForm({
        name: site.name,
        enabled: site.enabled,
      })
    }
  }, [site])

  const resetModeForm = useCallback(() => {
    if (site) {
      setModeForm({
        batch_mode_enabled: site.batch_mode_enabled,
        batch_difficulty: site.batch_difficulty,
        min_batch_count: site.min_batch_count,
        max_batch_count: site.max_batch_count,
        base_difficulty: site.base_difficulty,
        max_difficulty: site.max_difficulty,
      })
    }
  }, [site])

  const resetPolicyForm = useCallback(() => {
    if (ipPolicy) {
      setPolicyForm({
        enabled: ipPolicy.enabled,
        time_window_seconds: ipPolicy.time_window_seconds,
        request_threshold: ipPolicy.request_threshold,
        batch_increment: ipPolicy.batch_increment || 1,
        difficulty_increment: ipPolicy.difficulty_increment,
        failure_weight: ipPolicy.failure_weight,
      })
    }
  }, [ipPolicy])

  // Mutations
  const updateMutation = useMutation({
    mutationFn: (data: Parameters<typeof sitesApi.update>[1]) =>
      sitesApi.update(siteId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['site', siteId] })
    },
  })

  const updatePolicyMutation = useMutation({
    mutationFn: (data: Parameters<typeof sitesApi.updateIpPolicy>[1]) =>
      sitesApi.updateIpPolicy(siteId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['site-policy', siteId] })
    },
  })

  const addDomainMutation = useMutation({
    mutationFn: (domain: string) => sitesApi.addDomain(siteId, domain),
    onSuccess: () => {
      setNewDomain('')
      queryClient.invalidateQueries({ queryKey: ['site', siteId] })
    },
  })

  const deleteDomainMutation = useMutation({
    mutationFn: (domainId: number) => sitesApi.deleteDomain(siteId, domainId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['site', siteId] })
    },
  })

  const regenerateSecretMutation = useMutation({
    mutationFn: () => sitesApi.regenerateSecret(siteId),
    onSuccess: (data) => {
      setNewSecret(data.secret)
      setShowSecretConfirm(false)
    },
  })

  const viewSecretMutation = useMutation({
    mutationFn: () => sitesApi.getSecret(siteId),
    onSuccess: (data) => {
      setViewingSecret(data.secret)
      setShowViewSecret(false)
    },
  })

  const deleteIpRecordMutation = useMutation({
    mutationFn: (ip: string) => sitesApi.deleteIpRecord(siteId, ip),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['site-ip-records', siteId] })
    },
  })

  // 保存处理
  const handleSaveBasic = async () => {
    if (!basicForm) return
    await updateMutation.mutateAsync(basicForm)
  }

  const handleSaveMode = async () => {
    if (!modeForm) return
    await updateMutation.mutateAsync(modeForm)
  }

  const handleSavePolicy = async () => {
    if (!policyForm) return
    await updatePolicyMutation.mutateAsync(policyForm)
  }

  const handleCopy = async (text: string, field: string) => {
    await navigator.clipboard.writeText(text)
    setCopiedField(field)
    setTimeout(() => setCopiedField(null), 2000)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-6 w-6 border-2 border-violet-600 border-t-transparent"></div>
      </div>
    )
  }

  if (!site) {
    return (
      <div className="text-center py-12">
        <p className="text-slate-500">站点不存在</p>
      </div>
    )
  }

  const successRate = stats && stats.total_challenges > 0
    ? Math.round((stats.solved_challenges / stats.total_challenges) * 100)
    : 0

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={() => navigate('/console/sites')}
          className="p-2 hover:bg-slate-100 rounded-lg transition-colors"
        >
          <ArrowLeft className="h-5 w-5 text-slate-600" />
        </button>
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <h1 className="text-xl font-semibold text-slate-900">{site.name}</h1>
            <span className={`px-2 py-0.5 text-xs rounded-full font-medium ${
              site.enabled
                ? 'bg-emerald-100 text-emerald-700'
                : 'bg-slate-100 text-slate-500'
            }`}>
              {site.enabled ? '运行中' : '已禁用'}
            </span>
          </div>
          <div className="flex items-center gap-2 mt-1">
            <code className="text-xs text-slate-500 bg-slate-100 px-2 py-0.5 rounded font-mono">
              {site.sitekey}
            </code>
            <button 
              onClick={() => handleCopy(site.sitekey, 'sitekey')}
              className="text-slate-400 hover:text-slate-600 transition-colors"
            >
              {copiedField === 'sitekey' ? (
                <Check className="h-3.5 w-3.5 text-emerald-500" />
              ) : (
                <Copy className="h-3.5 w-3.5" />
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Stats - 彩色图标 */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-white rounded-xl border border-slate-200 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="p-1.5 rounded-lg bg-blue-50">
                <Activity className="h-4 w-4 text-blue-500" />
              </div>
              <span className="text-xs text-slate-500">24h 请求</span>
            </div>
            <div className="text-2xl font-semibold text-slate-900">{stats.last_24h_challenges}</div>
          </div>
          <div className="bg-white rounded-xl border border-slate-200 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="p-1.5 rounded-lg bg-emerald-50">
                <CheckCircle2 className="h-4 w-4 text-emerald-500" />
              </div>
              <span className="text-xs text-slate-500">24h 成功</span>
            </div>
            <div className="text-2xl font-semibold text-slate-900">{stats.last_24h_solved}</div>
          </div>
          <div className="bg-white rounded-xl border border-slate-200 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="p-1.5 rounded-lg bg-purple-50">
                <BarChart3 className="h-4 w-4 text-purple-500" />
              </div>
              <span className="text-xs text-slate-500">总请求</span>
            </div>
            <div className="text-2xl font-semibold text-slate-900">{stats.total_challenges}</div>
          </div>
          <div className="bg-white rounded-xl border border-slate-200 p-4">
            <div className="flex items-center gap-2 mb-2">
              <div className="p-1.5 rounded-lg bg-amber-50">
                <TrendingUp className="h-4 w-4 text-amber-500" />
              </div>
              <span className="text-xs text-slate-500">成功率</span>
            </div>
            <div className="text-2xl font-semibold text-slate-900">{successRate}%</div>
          </div>
        </div>
      )}

      {/* 基本设置 */}
      {basicForm && (
        <SettingSection
          title="基本设置"
          isDirty={isBasicDirty()}
          onSave={handleSaveBasic}
          onCancel={resetBasicForm}
          isSaving={updateMutation.isPending}
        >
          <SettingRow title="站点名称" description="用于标识此站点">
            <SettingInput
              value={basicForm.name}
              onChange={(v) => setBasicForm({ ...basicForm, name: v })}
              className="w-48"
            />
          </SettingRow>
          <SettingRow title="启用状态" description="关闭后将拒绝所有验证请求">
            <Toggle
              checked={basicForm.enabled}
              onChange={(v) => setBasicForm({ ...basicForm, enabled: v })}
            />
          </SettingRow>
          <SettingRow title="Secret 密钥" description="用于后端验证，请妥善保管" noBorder>
            <div className="flex items-center gap-2">
                <button
                  type="button"
                  onClick={() => setShowViewSecret(true)}
                  className="h-8 px-3 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg transition-colors flex items-center gap-1.5"
                >
                <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                查看
              </button>
              <button
                type="button"
                onClick={() => setShowSecretConfirm(true)}
                className="h-8 px-3 text-sm text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors flex items-center gap-1.5"
              >
                <RefreshCw className="w-3.5 h-3.5" />
                重新生成
              </button>
            </div>
          </SettingRow>
        </SettingSection>
      )}

      {/* 验证模式 */}
      {modeForm && (
        <SettingSection
          title="验证模式"
          description="选择 PoW 验证的计算策略"
          isDirty={isModeDirty()}
          onSave={handleSaveMode}
          onCancel={resetModeForm}
          isSaving={updateMutation.isPending}
        >
          <SettingRow
            title={
              <div className="flex items-center gap-2">
                <span>批量模式</span>
                <span className="px-1.5 py-0.5 text-xs rounded bg-emerald-100 text-emerald-700 font-medium">推荐</span>
              </div>
            }
            description="固定难度，根据风控动态调整计算数量"
          >
            <Toggle
              checked={modeForm.batch_mode_enabled}
              onChange={(v) => setModeForm({ ...modeForm, batch_mode_enabled: v })}
            />
          </SettingRow>

          {/* 模式说明 */}
          <div className="py-4 border-b border-slate-100">
            <div className={`flex items-start gap-3 p-3 rounded-lg ${
              modeForm.batch_mode_enabled ? 'bg-purple-50' : 'bg-slate-50'
            }`}>
              <div className={`p-1.5 rounded-lg ${
                modeForm.batch_mode_enabled ? 'bg-purple-100' : 'bg-slate-200'
              }`}>
                {modeForm.batch_mode_enabled ? (
                  <Layers className="h-4 w-4 text-purple-600" />
                ) : (
                  <Zap className="h-4 w-4 text-amber-600" />
                )}
              </div>
              <div className="flex-1 text-sm">
                <div className="font-medium text-slate-900">
                  {modeForm.batch_mode_enabled ? '批量模式' : '动态难度模式'}
                </div>
                <p className="text-slate-600 mt-0.5">
                  {modeForm.batch_mode_enabled 
                    ? `客户端计算 ${modeForm.min_batch_count}-${modeForm.max_batch_count} 个难度 ${modeForm.batch_difficulty} 的 PoW，风险越高数量越多`
                    : `客户端计算 1 个 PoW，难度范围 ${modeForm.base_difficulty}-${modeForm.max_difficulty}，风险越高难度越大`
                  }
                </p>
              </div>
            </div>
          </div>

          {/* 批量模式配置 */}
          {modeForm.batch_mode_enabled ? (
            <>
              <SettingRow title="固定难度" description="每个 PoW 的计算难度">
                <SettingSelect
                  value={modeForm.batch_difficulty}
                  onChange={(v) => setModeForm({ ...modeForm, batch_difficulty: Number(v) })}
                  options={[
                    { value: 3, label: '难度 3 (轻量)' },
                    { value: 4, label: '难度 4 (标准)' },
                    { value: 5, label: '难度 5 (较强)' },
                  ]}
                />
              </SettingRow>
              <SettingRow title="最小数量" description="低风险时的 PoW 数量">
                <SettingInput
                  type="number"
                  value={modeForm.min_batch_count}
                  onChange={(v) => setModeForm({ ...modeForm, min_batch_count: Number(v) })}
                  min={1}
                  max={10}
                  className="w-24"
                />
              </SettingRow>
              <SettingRow title="最大数量" description="高风险时的 PoW 数量" noBorder>
                <SettingInput
                  type="number"
                  value={modeForm.max_batch_count}
                  onChange={(v) => setModeForm({ ...modeForm, max_batch_count: Number(v) })}
                  min={1}
                  max={20}
                  className="w-24"
                />
              </SettingRow>
            </>
          ) : (
            <>
              <SettingRow title="基础难度" description="低风险时的难度等级">
                <SettingInput
                  type="number"
                  value={modeForm.base_difficulty}
                  onChange={(v) => setModeForm({ ...modeForm, base_difficulty: Number(v) })}
                  min={1}
                  max={10}
                  className="w-24"
                />
              </SettingRow>
              <SettingRow title="最大难度" description="高风险时的难度等级" noBorder>
                <SettingInput
                  type="number"
                  value={modeForm.max_difficulty}
                  onChange={(v) => setModeForm({ ...modeForm, max_difficulty: Number(v) })}
                  min={1}
                  max={12}
                  className="w-24"
                />
              </SettingRow>
            </>
          )}
        </SettingSection>
      )}

      {/* 风控策略 */}
      {policyForm && (
        <SettingSection
          title="风控策略"
          description="根据 IP 请求频率动态调整验证强度"
          isDirty={isPolicyDirty()}
          onSave={handleSavePolicy}
          onCancel={resetPolicyForm}
          isSaving={updatePolicyMutation.isPending}
        >
          <SettingRow title="启用风控" description="关闭后所有请求使用基础配置">
            <Toggle
              checked={policyForm.enabled}
              onChange={(v) => setPolicyForm({ ...policyForm, enabled: v })}
            />
          </SettingRow>

          {policyForm.enabled && (
            <>
              {/* 策略说明 */}
              <div className="py-4 border-b border-slate-100">
                <div className="flex items-start gap-2 p-3 rounded-lg bg-amber-50 text-amber-800">
                  <Info className="h-4 w-4 mt-0.5 flex-shrink-0" />
                  <p className="text-sm">
                    {modeForm?.batch_mode_enabled 
                      ? `同一 IP 在 ${policyForm.time_window_seconds} 秒内每 ${policyForm.request_threshold} 次请求，增加 ${policyForm.batch_increment} 个计算任务`
                      : `同一 IP 在 ${policyForm.time_window_seconds} 秒内每 ${policyForm.request_threshold} 次请求，增加 ${policyForm.difficulty_increment} 级难度`
                    }
                  </p>
                </div>
              </div>

              <SettingRow title="时间窗口" description="统计请求的时间范围（秒）">
                <SettingInput
                  type="number"
                  value={policyForm.time_window_seconds}
                  onChange={(v) => setPolicyForm({ ...policyForm, time_window_seconds: Number(v) })}
                  className="w-24"
                />
              </SettingRow>
              <SettingRow title="请求阈值" description="触发增量的请求次数">
                <SettingInput
                  type="number"
                  value={policyForm.request_threshold}
                  onChange={(v) => setPolicyForm({ ...policyForm, request_threshold: Number(v) })}
                  className="w-24"
                />
              </SettingRow>
              <SettingRow 
                title={modeForm?.batch_mode_enabled ? "批次增量" : "难度增量"} 
                description={modeForm?.batch_mode_enabled ? "每次触发增加的 PoW 数量" : "每次触发增加的难度等级"}
              >
                <SettingInput
                  type="number"
                  value={modeForm?.batch_mode_enabled ? policyForm.batch_increment : policyForm.difficulty_increment}
                  onChange={(v) => setPolicyForm(
                    modeForm?.batch_mode_enabled 
                      ? { ...policyForm, batch_increment: Number(v) }
                      : { ...policyForm, difficulty_increment: Number(v) }
                  )}
                  min={1}
                  max={5}
                  className="w-24"
                />
              </SettingRow>
              <SettingRow title="失败权重" description="失败请求的计数权重" noBorder>
                <SettingInput
                  type="number"
                  value={policyForm.failure_weight}
                  onChange={(v) => setPolicyForm({ ...policyForm, failure_weight: Number(v) })}
                  min={1}
                  max={10}
                  step={0.1}
                  className="w-24"
                />
              </SettingRow>
            </>
          )}
        </SettingSection>
      )}

      {/* 域名白名单 */}
      <section className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="px-6 py-4 border-b border-slate-100">
          <h2 className="font-semibold text-slate-900">域名白名单</h2>
          <p className="text-xs text-slate-500 mt-1">留空允许所有域名，支持通配符 *.example.com</p>
        </div>
        <div className="px-6 py-4">
          <div className="flex gap-2 mb-4">
            <input
              type="text"
              value={newDomain}
              onChange={(e) => setNewDomain(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && newDomain && addDomainMutation.mutate(newDomain)}
              placeholder="输入域名，按回车添加"
              className="flex-1 h-9 px-3 rounded-lg border border-slate-200 bg-white text-sm 
                focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
            />
            <button
              type="button"
              onClick={() => newDomain && addDomainMutation.mutate(newDomain)}
              disabled={!newDomain || addDomainMutation.isPending}
              className="h-9 px-4 bg-slate-900 hover:bg-slate-800 text-white text-sm font-medium rounded-lg 
                transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1.5"
            >
              <Plus className="w-4 h-4" />
              添加
            </button>
          </div>
          
          {site.domains && site.domains.length > 0 ? (
            <div className="space-y-2">
              {site.domains.map((domain) => (
                <div
                  key={domain.id}
                  className="flex items-center justify-between px-3 py-2.5 bg-slate-50 rounded-lg group"
                >
                  <div className="flex items-center gap-2">
                    <Globe className="h-4 w-4 text-slate-400" />
                    <code className="text-sm text-slate-700">{domain.domain}</code>
                  </div>
                  <button
                    type="button"
                    onClick={() => deleteDomainMutation.mutate(domain.id)}
                    className="p-1.5 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded 
                      opacity-0 group-hover:opacity-100 transition-all"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-slate-400 text-sm">
              未配置域名限制，允许所有域名
            </div>
          )}
        </div>
      </section>

      {/* IP 访问记录 */}
      <section className="bg-white rounded-xl border border-slate-200 overflow-hidden">
        <div className="px-6 py-4 border-b border-slate-100">
          <h2 className="font-semibold text-slate-900">IP 访问记录</h2>
          <p className="text-xs text-slate-500 mt-1">当前活跃的 IP 风控记录</p>
        </div>
        <div className="overflow-x-auto">
          {ipRecords && ipRecords.length > 0 ? (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-100 bg-slate-50/50">
                  <th className="text-left py-3 px-6 font-medium text-slate-500">IP 地址</th>
                  <th className="text-left py-3 px-4 font-medium text-slate-500">请求数</th>
                  <th className="text-left py-3 px-4 font-medium text-slate-500">失败数</th>
                  <th className="text-left py-3 px-4 font-medium text-slate-500">当前等级</th>
                  <th className="text-left py-3 px-4 font-medium text-slate-500">最后请求</th>
                  <th className="text-right py-3 px-6 font-medium text-slate-500">操作</th>
                </tr>
              </thead>
              <tbody>
                {ipRecords.map((record) => (
                  <tr key={record.id} className="border-b border-slate-50 hover:bg-slate-50/50 transition-colors">
                    <td className="py-3 px-6 font-mono text-slate-700">{record.client_ip}</td>
                    <td className="py-3 px-4 text-slate-600">{record.request_count}</td>
                    <td className="py-3 px-4">
                      <span className={record.failure_count > 0 ? 'text-rose-600 font-medium' : 'text-slate-600'}>
                        {record.failure_count}
                      </span>
                    </td>
                    <td className="py-3 px-4">
                      <span className="px-2 py-1 bg-blue-50 text-blue-700 rounded text-xs font-medium">
                        {record.current_difficulty}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-slate-500 text-xs">
                      {new Date(record.last_request_at).toLocaleString('zh-CN')}
                    </td>
                    <td className="py-3 px-6 text-right">
                      <button
                        type="button"
                        onClick={() => deleteIpRecordMutation.mutate(record.client_ip)}
                        className="text-xs text-rose-600 hover:text-rose-700 font-medium"
                      >
                        重置
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <div className="text-center py-12 text-slate-400 text-sm">
              暂无访问记录
            </div>
          )}
        </div>
      </section>

      {/* Secret Confirm Modal */}
      {showSecretConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-100">
              <h2 className="font-semibold text-slate-900 flex items-center gap-2">
                <XCircle className="w-5 h-5 text-red-500" />
                确认重新生成 Secret
              </h2>
            </div>
            <div className="p-6">
              <p className="text-slate-600 text-sm mb-6">
                重新生成后，旧的 Secret 将立即失效。所有使用旧 Secret 的服务需要更新配置。
              </p>
              <div className="flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowSecretConfirm(false)}
                  className="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium"
                >
                  取消
                </button>
                <button
                  type="button"
                  onClick={() => regenerateSecretMutation.mutate()}
                  disabled={regenerateSecretMutation.isPending}
                  className="h-9 px-4 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors text-sm font-medium disabled:opacity-50"
                >
                  确认重新生成
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* View Secret Confirm Modal */}
      {showViewSecret && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-100">
              <h2 className="font-semibold text-slate-900 flex items-center gap-2">
                <Shield className="w-5 h-5 text-amber-500" />
                确认查看 Secret
              </h2>
            </div>
            <div className="p-6">
              <p className="text-slate-600 text-sm mb-6">
                Secret 是敏感信息，请确保周围环境安全。建议在私密场所查看。
              </p>
              <div className="flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowViewSecret(false)}
                  className="h-9 px-4 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium"
                >
                  取消
                </button>
                <button
                  type="button"
                  onClick={() => viewSecretMutation.mutate()}
                  disabled={viewSecretMutation.isPending}
                  className="h-9 px-4 bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 rounded-lg transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  确认查看
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Viewing Secret Modal */}
      {viewingSecret && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
            <div className="px-6 py-4 border-b border-amber-100 bg-amber-50">
              <h2 className="font-semibold text-amber-700 flex items-center gap-2">
                <Shield className="w-5 h-5" />
                站点 Secret
              </h2>
            </div>
            <div className="p-6 space-y-4">
              <div className="flex items-center gap-2 p-3 rounded-lg bg-amber-50 border border-amber-200 text-amber-700 text-sm">
                <svg className="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                <span>请妥善保管此密钥，不要泄露给他人</span>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">
                  Secret 密钥
                </label>
                <div className="flex items-center gap-2">
                  <code className="flex-1 px-3 py-2.5 bg-slate-100 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono break-all">
                    {viewingSecret}
                  </code>
                  <button
                    type="button"
                    onClick={() => handleCopy(viewingSecret, 'viewing-secret')}
                    className="h-10 w-10 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors flex-shrink-0"
                  >
                    {copiedField === 'viewing-secret' ? (
                      <Check className="h-5 w-5 text-emerald-500" />
                    ) : (
                      <Copy className="h-5 w-5" />
                    )}
                  </button>
                </div>
              </div>

              <div className="flex justify-end pt-2">
                <button
                  type="button"
                  onClick={() => setViewingSecret(null)}
                  className="h-9 px-4 bg-slate-900 hover:bg-slate-800 text-white rounded-lg transition-colors text-sm font-medium"
                >
                  关闭
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* New Secret Modal */}
      {newSecret && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
            <div className="px-6 py-4 border-b border-violet-100 bg-gradient-to-r from-violet-50 to-purple-50">
              <h2 className="font-semibold text-violet-700 flex items-center gap-2">
                <CheckCircle2 className="w-5 h-5" />
                Secret 已重新生成
              </h2>
            </div>
            <div className="p-6 space-y-4">
              <div className="flex items-center gap-2 p-3 rounded-lg bg-blue-50 border border-blue-200 text-blue-700 text-sm">
                <Shield className="w-4 h-4 flex-shrink-0" />
                <span>新的 Secret 已生成，可随时在此页面查看</span>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">
                  新 Secret
                </label>
                <div className="flex items-center gap-2">
                  <code className="flex-1 px-3 py-2.5 bg-slate-100 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono break-all">
                    {newSecret}
                  </code>
                  <button
                    type="button"
                    onClick={() => handleCopy(newSecret, 'new-secret')}
                    className="h-10 w-10 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors flex-shrink-0"
                  >
                    {copiedField === 'new-secret' ? (
                      <Check className="h-5 w-5 text-emerald-500" />
                    ) : (
                      <Copy className="h-5 w-5" />
                    )}
                  </button>
                </div>
              </div>

              <div className="flex justify-end pt-2">
                <button
                  type="button"
                  onClick={() => setNewSecret(null)}
                  className="h-9 px-4 bg-slate-900 hover:bg-slate-800 text-white rounded-lg transition-colors text-sm font-medium"
                >
                  我已保存
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
