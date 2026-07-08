import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { sitesApi, Site } from '@/lib/api'
import { 
  Plus, 
  Trash2, 
  Copy, 
  Check, 
  Globe,
  ChevronRight,
  AlertCircle,
  X
} from 'lucide-react'

// WeAuth Logo for site icon
function SiteIcon({ enabled }: { enabled: boolean }) {
  return (
    <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${
      enabled 
        ? 'bg-violet-600 shadow-sm' 
        : 'bg-slate-200'
    }`}>
      <svg className={`w-5 h-5 ${enabled ? 'text-white' : 'text-slate-400'}`} viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
      </svg>
    </div>
  )
}

export default function Sites() {
  const queryClient = useQueryClient()
  const [showCreate, setShowCreate] = useState(false)
  const [newSiteName, setNewSiteName] = useState('')
  const [newDomains, setNewDomains] = useState('')
  const [batchModeEnabled, setBatchModeEnabled] = useState(true) // 默认批量模式
  const [batchDifficulty, setBatchDifficulty] = useState(4)
  const [createdSite, setCreatedSite] = useState<(Site & { secret: string }) | null>(null)
  const [copiedField, setCopiedField] = useState<string | null>(null)

  const { data: sites, isLoading } = useQuery({
    queryKey: ['sites'],
    queryFn: sitesApi.list,
  })

  const createMutation = useMutation({
    mutationFn: () =>
      sitesApi.create(
        newSiteName,
        newDomains.split('\n').filter((d) => d.trim()),
        {
          batch_mode_enabled: batchModeEnabled,
          batch_difficulty: batchDifficulty,
        }
      ),
    onSuccess: (site) => {
      setCreatedSite(site)
      setNewSiteName('')
      setNewDomains('')
      setBatchModeEnabled(true)
      setBatchDifficulty(4)
      setShowCreate(false)
      queryClient.invalidateQueries({ queryKey: ['sites'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: sitesApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sites'] })
    },
  })

  const handleCopy = async (text: string, field: string) => {
    await navigator.clipboard.writeText(text)
    setCopiedField(field)
    setTimeout(() => setCopiedField(null), 2000)
  }

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate()
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-xl font-semibold text-slate-900">站点管理</h1>
          <p className="text-slate-500 text-sm mt-0.5">管理您的 PoW 验证站点和配置</p>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className="inline-flex items-center justify-center h-9 px-4 bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 active:bg-blue-100 text-sm rounded-lg font-medium transition-all duration-150 shadow-sm"
        >
          <Plus className="h-4 w-4 mr-1.5" />
          创建站点
        </button>
      </div>

      {/* Loading */}
      {isLoading && (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-6 w-6 border-2 border-[#0066ff] border-t-transparent"></div>
        </div>
      )}

      {/* Empty State */}
      {!isLoading && sites?.length === 0 && (
        <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm p-12 text-center">
          <div className="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mx-auto mb-4">
            <Globe className="w-6 h-6 text-slate-400" />
          </div>
          <h3 className="font-medium text-slate-900 mb-1">还没有站点</h3>
          <p className="text-slate-500 text-sm mb-4">
            创建您的第一个站点来开始使用
          </p>
          <button
            onClick={() => setShowCreate(true)}
            className="inline-flex items-center px-4 py-2 bg-[#0066ff] hover:bg-[#0052cc] text-white text-sm rounded-lg font-medium transition-colors"
          >
            <Plus className="h-4 w-4 mr-1.5" />
            创建站点
          </button>
        </div>
      )}

      {/* Site List */}
      {!isLoading && sites && sites.length > 0 && (
        <div className="bg-white/80 backdrop-blur-xl rounded-2xl border border-slate-200/60 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
          <div className="divide-y divide-slate-100">
            {sites.map((site) => (
              <Link
                key={site.id}
                to={`/console/sites/${site.id}`}
                className="flex items-center gap-4 p-4 hover:bg-slate-50 transition-colors group"
              >
                {/* Icon */}
                <SiteIcon enabled={site.enabled} />

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <h3 className="font-medium text-slate-900 truncate">{site.name}</h3>
                    <span className={`px-1.5 py-0.5 text-xs rounded font-medium ${
                      site.enabled
                        ? 'bg-emerald-50 text-emerald-600'
                        : 'bg-slate-100 text-slate-500'
                    }`}>
                      {site.enabled ? '运行中' : '已禁用'}
                    </span>
                  </div>
                  <div className="flex items-center gap-2 mt-1">
                    <code className="text-xs text-slate-500 bg-slate-50 px-1.5 py-0.5 rounded font-mono">
                      {site.sitekey}
                    </code>
                    <button
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        handleCopy(site.sitekey, `sitekey-${site.id}`)
                      }}
                      className="text-slate-400 hover:text-slate-600 transition-colors"
                    >
                      {copiedField === `sitekey-${site.id}` ? (
                        <Check className="h-3.5 w-3.5 text-emerald-500" />
                      ) : (
                        <Copy className="h-3.5 w-3.5" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Stats */}
                <div className="hidden sm:flex items-center gap-6 text-sm">
                  <div className="text-center">
                    <p className="text-slate-900 font-medium">{site.domains?.length || 0}</p>
                    <p className="text-xs text-slate-400">域名</p>
                  </div>
                  <div className="text-center">
                    <p className="text-slate-900 font-medium">{site.base_difficulty}</p>
                    <p className="text-xs text-slate-400">难度</p>
                  </div>
                </div>

                {/* Actions */}
                <div className="flex items-center gap-1">
                  <button
                    onClick={(e) => {
                      e.preventDefault()
                      e.stopPropagation()
                      if (confirm('确定要删除此站点吗？')) {
                        deleteMutation.mutate(site.id)
                      }
                    }}
                    className="p-1.5 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded transition-colors"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                  <ChevronRight className="w-4 h-4 text-slate-300 group-hover:text-slate-400 transition-colors" />
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* Create Modal */}
      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/20 backdrop-blur-sm p-4">
          <div className="bg-white rounded-xl border border-slate-200 shadow-xl w-full max-w-md overflow-hidden">
            <div className="flex items-center justify-between px-5 py-3 border-b border-slate-100">
              <h2 className="font-semibold text-slate-900">创建新站点</h2>
              <button
                onClick={() => setShowCreate(false)}
                className="p-1 text-slate-400 hover:text-slate-600 rounded transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleCreate} className="p-5 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">
                  站点名称 <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={newSiteName}
                  onChange={(e) => setNewSiteName(e.target.value)}
                  placeholder="例如: 我的网站"
                  className="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1.5">
                  域名白名单
                </label>
                <textarea
                  value={newDomains}
                  onChange={(e) => setNewDomains(e.target.value)}
                  placeholder="example.com&#10;*.example.com&#10;localhost"
                  rows={3}
                  className="w-full px-3 py-2 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 placeholder-slate-400 text-sm font-mono focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all resize-none"
                />
                <p className="text-xs text-slate-400 mt-1">
                  每行一个域名，支持通配符
                </p>
              </div>

              {/* 验证模式选择 */}
              <div className="pt-2 border-t border-slate-100">
                <label className="block text-sm font-medium text-slate-700 mb-3">
                  验证模式
                </label>
                <div className="space-y-3">
                  <label className="flex items-start gap-3 p-3 rounded-lg border border-slate-200 hover:border-purple-300 cursor-pointer transition-colors has-[:checked]:border-purple-500 has-[:checked]:bg-purple-50/50">
                    <input
                      type="radio"
                      name="mode"
                      checked={batchModeEnabled}
                      onChange={() => setBatchModeEnabled(true)}
                      className="mt-0.5"
                    />
                    <div className="flex-1">
                      <div className="font-medium text-slate-900 text-sm">批量模式（推荐）</div>
                      <p className="text-xs text-slate-500 mt-0.5">
                        固定难度，根据风控动态调整计算数量。更灵活的防护策略。
                      </p>
                    </div>
                  </label>
                  <label className="flex items-start gap-3 p-3 rounded-lg border border-slate-200 hover:border-purple-300 cursor-pointer transition-colors has-[:checked]:border-purple-500 has-[:checked]:bg-purple-50/50">
                    <input
                      type="radio"
                      name="mode"
                      checked={!batchModeEnabled}
                      onChange={() => setBatchModeEnabled(false)}
                      className="mt-0.5"
                    />
                    <div className="flex-1">
                      <div className="font-medium text-slate-900 text-sm">动态难度模式</div>
                      <p className="text-xs text-slate-500 mt-0.5">
                        单次计算，根据风控动态调整难度等级。传统验证方式。
                      </p>
                    </div>
                  </label>
                </div>

                {/* 批量模式难度选择 */}
                {batchModeEnabled && (
                  <div className="mt-3">
                    <label className="block text-xs font-medium text-slate-600 mb-1.5">
                      固定难度
                    </label>
                    <select
                      value={batchDifficulty}
                      onChange={(e) => setBatchDifficulty(Number(e.target.value))}
                      className="w-full h-9 px-3 rounded-lg bg-slate-50 border border-slate-200 text-slate-900 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all"
                    >
                      <option value={3}>难度 3 - 轻量（约 0.1 秒）</option>
                      <option value={4}>难度 4 - 标准（约 1 秒）</option>
                      <option value={5}>难度 5 - 较强（约 10 秒）</option>
                    </select>
                  </div>
                )}
              </div>
              
              {createMutation.isError && (
                <div className="flex items-center gap-2 p-2.5 rounded-lg bg-red-50 border border-red-100 text-red-600 text-sm">
                  <AlertCircle className="w-4 h-4 flex-shrink-0" />
                  <span>创建失败，请重试</span>
                </div>
              )}

              <div className="flex justify-end gap-2 pt-2">
                <button
                  type="button"
                  onClick={() => setShowCreate(false)}
                  className="h-9 px-4 text-slate-600 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition-colors text-sm font-medium"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={createMutation.isPending || !newSiteName.trim()}
                  className="h-9 px-4 bg-[#0066ff] hover:bg-[#0052cc] text-white rounded-lg text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {createMutation.isPending ? '创建中...' : '创建'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Created Site Modal */}
      {createdSite && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/20 backdrop-blur-sm p-4">
          <div className="bg-white rounded-xl border border-slate-200 shadow-xl w-full max-w-lg overflow-hidden">
            <div className="px-5 py-3 border-b border-violet-100 bg-gradient-to-r from-violet-50 to-purple-50">
              <div className="flex items-center gap-2 text-violet-600">
                <Check className="w-5 h-5" />
                <h2 className="font-semibold">站点创建成功</h2>
              </div>
            </div>
            <div className="p-5 space-y-4">
              <div className="flex items-center gap-2 p-2.5 rounded-lg bg-blue-50 border border-blue-200 text-blue-700 text-sm">
                <AlertCircle className="w-4 h-4 flex-shrink-0" />
                <span>站点创建成功！Secret 可在站点详情页查看</span>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-600 mb-1.5">
                  Sitekey（前端使用）
                </label>
                <div className="flex items-center gap-2">
                  <code className="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono">
                    {createdSite.sitekey}
                  </code>
                  <button
                    onClick={() => handleCopy(createdSite.sitekey, 'new-sitekey')}
                    className="h-9 w-9 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
                  >
                    {copiedField === 'new-sitekey' ? (
                      <Check className="h-4 w-4 text-emerald-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-600 mb-1.5">
                  Secret（后端使用，请勿泄露）
                </label>
                <div className="flex items-center gap-2">
                  <code className="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm text-slate-700 font-mono break-all">
                    {createdSite.secret}
                  </code>
                  <button
                    onClick={() => handleCopy(createdSite.secret, 'new-secret')}
                    className="h-9 w-9 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors flex-shrink-0"
                  >
                    {copiedField === 'new-secret' ? (
                      <Check className="h-4 w-4 text-emerald-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </button>
                </div>
              </div>

              <div className="flex justify-end pt-2">
                <button
                  onClick={() => setCreatedSite(null)}
                  className="h-9 px-4 bg-[#0066ff] hover:bg-[#0052cc] text-white rounded-lg text-sm font-medium transition-colors"
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
