import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { 
  Shield, 
  Zap, 
  TrendingDown, 
  Code, 
  CheckCircle,
  ArrowRight,
  Activity,
  Layers
} from 'lucide-react'
import Navbar from '@/components/Navbar'

// WeAuth Logo
function WeAuthLogo({ className = "w-6 h-6" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
    </svg>
  )
}

export default function Landing() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/20">
      {/* Navigation */}
      <Navbar />

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="text-center max-w-4xl mx-auto"
          >
            <div className="inline-flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 text-blue-700 text-sm font-medium rounded-full mb-6 shadow-sm">
              <Zap className="w-4 h-4 text-blue-500" />
              新一代人机验证服务
            </div>
            
            <h1 className="text-5xl md:text-6xl font-bold text-slate-900 mb-6 leading-tight">
              基于工作量证明的
              <br />
              <span className="relative inline-block">
                <span className="bg-clip-text text-transparent bg-gradient-to-r from-rose-500 via-orange-500 via-amber-400 via-emerald-500 via-cyan-500 via-blue-500 to-violet-600 animate-gradient bg-[length:300%_auto]">
                  智能人机验证
                </span>
              </span>
            </h1>
            
            <p className="text-xl text-slate-600 mb-10 leading-relaxed">
              无需图片识别，无需滑块拼图。利用客户端算力完成验证，
              <br />
              为您的应用提供更安全、更智能的防护方案
            </p>

            <div className="flex flex-wrap items-center gap-4 justify-center">
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.98 }}>
                <Link
                  to="/login"
                  className="inline-flex items-center gap-2 px-8 py-4 rounded-xl font-medium text-lg bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 hover:text-blue-700 transition-all shadow-sm"
                >
                  免费开始
                  <ArrowRight className="w-5 h-5 transition-transform" />
                </Link>
              </motion.div>
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.98 }}>
                <Link
                  to="/demo"
                  className="inline-flex items-center justify-center px-8 py-4 bg-white border-2 border-slate-200 text-slate-700 font-medium rounded-xl hover:border-slate-300 hover:bg-slate-50 transition-all text-lg shadow-sm"
                >
                  查看演示
                </Link>
              </motion.div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-3 gap-8 mt-16 pt-16 border-t border-slate-200">
              <div>
                <div className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-emerald-600 to-teal-600 mb-2">99.9%</div>
                <div className="text-sm text-slate-600">服务可用性</div>
              </div>
              <div>
                <div className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-cyan-600 mb-2">&lt;100ms</div>
                <div className="text-sm text-slate-600">验证响应时间</div>
              </div>
              <div>
                <div className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-violet-600 to-purple-600 mb-2">0元</div>
                <div className="text-sm text-slate-600">起步价格</div>
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 mb-4">
              为什么选择 WeAuth
            </h2>
            <p className="text-lg text-slate-600 max-w-2xl mx-auto">
              结合工作量证明（PoW）与智能风控，为您的应用提供企业级安全防护
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {[
              {
                icon: Shield,
                title: '高安全性',
                description: '基于加密学的工作量证明，杜绝自动化攻击和爬虫',
                color: 'from-blue-100 to-cyan-200',
                iconColor: 'text-blue-600'
              },
              {
                icon: Zap,
                title: '极速验证',
                description: '客户端计算，服务端快速验证，毫秒级响应',
                color: 'from-amber-100 to-yellow-200',
                iconColor: 'text-amber-600'
              },
              {
                icon: TrendingDown,
                title: '节省成本',
                description: '利用客户端算力，大幅降低服务器负载和运营成本',
                color: 'from-emerald-100 to-teal-200',
                iconColor: 'text-emerald-600'
              },
              {
                icon: Activity,
                title: '智能风控',
                description: '基于 IP 和行为的动态风控策略，自动调整验证难度',
                color: 'from-indigo-100 to-purple-200',
                iconColor: 'text-indigo-600'
              },
              {
                icon: Layers,
                title: '双重模式',
                description: '支持动态难度和批量计算两种模式，灵活应对不同场景',
                color: 'from-rose-100 to-pink-200',
                iconColor: 'text-rose-600'
              },
              {
                icon: Code,
                title: '开发友好',
                description: '类似 Cloudflare Turnstile 的集成方式，5 分钟即可接入',
                color: 'from-violet-100 to-purple-200',
                iconColor: 'text-violet-600'
              }
            ].map((feature, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                viewport={{ once: true }}
                className="group relative bg-gradient-to-br from-white to-slate-50/50 rounded-2xl p-8 hover:shadow-xl hover:shadow-slate-200/50 transition-all duration-300 border border-slate-200/60 hover:border-slate-300/60"
              >
                <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${feature.color} border border-white/80 shadow-sm flex items-center justify-center mb-6 group-hover:scale-110 group-hover:shadow-md transition-all duration-300`}>
                  <feature.icon className={`w-7 h-7 ${feature.iconColor}`} />
                </div>
                <h3 className="text-xl font-semibold text-slate-900 mb-3">
                  {feature.title}
                </h3>
                <p className="text-slate-600 leading-relaxed">
                  {feature.description}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* How it works */}
      <section className="py-20 bg-gradient-to-br from-slate-50 via-blue-50/20 to-indigo-50/30">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 mb-4">
              工作原理
            </h2>
            <p className="text-lg text-slate-600 max-w-2xl mx-auto">
              简单三步，完成人机验证
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            {[
              {
                step: '01',
                title: '用户触发验证',
                description: '用户点击验证按钮，系统根据风控策略分配验证任务',
                color: 'from-indigo-500 to-blue-500'
              },
              {
                step: '02',
                title: '客户端计算',
                description: '浏览器在后台利用 Web Worker 进行 SHA-256 工作量证明计算',
                color: 'from-blue-500 to-cyan-500'
              },
              {
                step: '03',
                title: '服务端验证',
                description: '提交计算结果，服务端快速验证并签发有效期内的 Token',
                color: 'from-emerald-500 to-teal-500'
              }
            ].map((step, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                transition={{ delay: index * 0.2 }}
                viewport={{ once: true }}
                className="relative bg-white rounded-2xl p-8 border border-slate-200 hover:shadow-lg transition-all"
              >
                <div className={`text-6xl font-bold bg-clip-text text-transparent bg-gradient-to-r ${step.color} mb-4 opacity-20`}>{step.step}</div>
                <h3 className="text-xl font-semibold text-slate-900 mb-3">
                  {step.title}
                </h3>
                <p className="text-slate-600 leading-relaxed">
                  {step.description}
                </p>
                {index < 2 && (
                  <div className={`hidden md:block absolute top-12 left-full w-8 h-0.5 bg-gradient-to-r ${step.color} opacity-30`}></div>
                )}
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Pricing */}
      <section id="pricing" className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 mb-4">
              简单透明的定价
            </h2>
            <p className="text-lg text-slate-600 max-w-2xl mx-auto">
              从免费开始，按需扩展
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {[
              {
                name: '免费版',
                price: '¥0',
                period: '永久免费',
                features: [
                  '每月 10,000 次验证',
                  '基础风控策略',
                  '社区支持',
                  '标准 SLA'
                ],
                cta: '开始使用',
                popular: false
              },
              {
                name: '专业版',
                price: '¥0',
                period: '/月',
                features: [
                  '每月 100,000 次验证',
                  '高级风控策略',
                  '邮件支持',
                  '99.9% SLA',
                  '自定义域名'
                ],
                cta: '立即升级',
                popular: true
              },
              {
                name: '企业版',
                price: '定制',
                period: '联系我们',
                features: [
                  '无限验证次数',
                  '定制风控策略',
                  '专属技术支持',
                  '99.99% SLA',
                  '私有部署'
                ],
                cta: '联系销售',
                popular: false
              }
            ].map((plan, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                viewport={{ once: true }}
                whileHover={{ y: -4 }}
                className="relative"
              >
                <div
                  className={`relative h-full rounded-2xl p-8 transition-all duration-300 ${
                    plan.popular
                      ? 'bg-white border-2 border-blue-300 shadow-xl shadow-blue-500/5'
                      : 'bg-white border border-slate-200 hover:border-slate-300 hover:shadow-lg'
                  }`}
                >
                  {plan.popular && (
                    <div className="absolute -top-3 left-1/2 -translate-x-1/2 px-4 py-1 bg-gradient-to-r from-blue-500 to-cyan-500 text-white text-xs font-medium rounded-full shadow-md">
                      最受欢迎
                    </div>
                  )}
                
                  <div className="mb-6">
                    <h3 className="text-2xl font-bold mb-2 text-slate-900">
                      {plan.name}
                    </h3>
                    <div className="flex items-baseline gap-1">
                      <span className="text-4xl font-bold text-slate-900">
                        {plan.price}
                      </span>
                      <span className="text-sm text-slate-500">
                        {plan.period}
                      </span>
                    </div>
                  </div>

                  <ul className="space-y-4 mb-8">
                    {plan.features.map((feature, i) => (
                      <li key={i} className="flex items-start gap-3">
                        <CheckCircle className={`w-5 h-5 flex-shrink-0 ${
                          plan.popular ? 'text-blue-500' : 'text-emerald-500'
                        }`} />
                        <span className="text-sm text-slate-600">
                          {feature}
                        </span>
                      </li>
                    ))}
                  </ul>

                  <motion.div whileHover={{ scale: 1.02, y: -2 }} whileTap={{ scale: 0.98 }}>
                    <Link
                      to="/login"
                      className={`block w-full py-3 text-center font-medium rounded-xl transition-all ${
                        plan.popular
                          ? 'bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 shadow-sm'
                          : 'bg-slate-50 text-slate-900 hover:bg-slate-100 border border-slate-200'
                      }`}
                    >
                      {plan.cta}
                    </Link>
                  </motion.div>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-br from-indigo-600 via-blue-600 to-cyan-600 relative overflow-hidden">
        {/* 装饰性背景 */}
        <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGRlZnM+PHBhdHRlcm4gaWQ9ImdyaWQiIHdpZHRoPSI2MCIgaGVpZ2h0PSI2MCIgcGF0dGVyblVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+PHBhdGggZD0iTSAxMCAwIEwgMCAwIDAgMTAiIGZpbGw9Im5vbmUiIHN0cm9rZT0id2hpdGUiIHN0cm9rZS13aWR0aD0iMSIgb3BhY2l0eT0iMC4xIi8+PC9wYXR0ZXJuPjwvZGVmcz48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSJ1cmwoI2dyaWQpIi8+PC9zdmc+')] opacity-30"></div>
        
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
          >
            <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
              准备好开始了吗？
            </h2>
            <p className="text-xl text-white/90 mb-10">
              5 分钟即可完成接入，立即体验新一代人机验证
            </p>
            <div className="flex flex-wrap items-center justify-center gap-4">
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.98 }}>
                <Link
                  to="/login"
                  className="inline-flex items-center gap-2 px-8 py-4 rounded-xl font-medium text-lg bg-white border-2 border-white text-blue-600 hover:bg-white/90 shadow-lg hover:shadow-xl transition-all"
                >
                  免费开始
                  <ArrowRight className="w-5 h-5 transition-transform" />
                </Link>
              </motion.div>
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.98 }}>
                <Link
                  to="/demo"
                  className="inline-flex items-center justify-center px-8 py-4 bg-white/10 backdrop-blur-sm border-2 border-white/40 text-white font-medium rounded-xl hover:bg-white/20 hover:border-white/60 transition-all text-lg"
                >
                  查看演示
                </Link>
              </motion.div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-slate-900 text-slate-400 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-50 to-cyan-50 border border-blue-200 flex items-center justify-center shadow-sm">
                  <WeAuthLogo className="w-4 h-4 text-blue-600" />
                </div>
                <span className="text-white font-bold">微验 WeAuth</span>
              </div>
              <p className="text-sm">
                新一代工作量证明人机验证服务
              </p>
            </div>

            <div>
              <h3 className="text-white font-semibold mb-4">产品</h3>
              <ul className="space-y-2 text-sm">
                <li><a href="#features" className="hover:text-white transition-colors">特性</a></li>
                <li><a href="#pricing" className="hover:text-white transition-colors">定价</a></li>
                <li><Link to="/demo" className="hover:text-white transition-colors">演示</Link></li>
              </ul>
            </div>

            <div>
              <h3 className="text-white font-semibold mb-4">开发者</h3>
              <ul className="space-y-2 text-sm">
                <li><a href="https://home.miaojilab.cn/" target="_blank" rel="noopener noreferrer" className="hover:text-white transition-colors">妙计实验室</a></li>
                <li><Link to="/demo" className="hover:text-white transition-colors">API 文档</Link></li>
              </ul>
            </div>

            <div>
              <h3 className="text-white font-semibold mb-4">关于</h3>
              <ul className="space-y-2 text-sm">
                <li><Link to="/login" className="hover:text-white transition-colors">登录</Link></li>
                <li><Link to="/privacy" className="hover:text-white transition-colors">隐私政策</Link></li>
              </ul>
            </div>
          </div>

          <div className="border-t border-slate-800 mt-12 pt-8 text-sm text-center">
            <div className="flex items-center justify-center gap-6 mb-3">
              <Link to="/privacy" className="text-slate-400 hover:text-white transition-colors">
                隐私政策
              </Link>
              <span className="text-slate-600">|</span>
              <a 
                href="mailto:support@weauth.cn" 
                className="text-slate-400 hover:text-white transition-colors"
              >
                联系我们
              </a>
            </div>
            <p>© 2026 WeAuth. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
