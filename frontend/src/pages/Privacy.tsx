import { Link } from 'react-router-dom'
import { Shield, ArrowLeft } from 'lucide-react'

// WeAuth Logo 组件
function WeAuthLogo({ className = "" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 6c1.4 0 2.8 1.1 2.8 2.5V11c.6.2 1 .8 1 1.5v3c0 .8-.7 1.5-1.5 1.5h-4.6c-.8 0-1.5-.7-1.5-1.5v-3c0-.7.4-1.3 1-1.5V9.5C9.2 8.1 10.6 7 12 7zm0 1.2c-.8 0-1.5.7-1.5 1.3V11h3V9.5c0-.6-.7-1.3-1.5-1.3z"/>
    </svg>
  )
}

export default function Privacy() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30">
      {/* 简洁头部 */}
      <header className="sticky top-0 z-50 backdrop-blur-sm bg-white/80 border-b border-slate-200/50">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2 hover:opacity-80 transition-opacity">
              <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-50 to-cyan-50 border border-blue-200/50 flex items-center justify-center">
                <WeAuthLogo className="w-4 h-4 text-blue-600" />
              </div>
              <span className="text-lg font-semibold text-slate-900">微验 WeAuth</span>
            </Link>

            {/* 返回首页 */}
            <Link
              to="/"
              className="flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900 transition-colors"
            >
              <ArrowLeft className="w-4 h-4" />
              返回首页
            </Link>
          </div>
        </div>
      </header>

      {/* 主要内容 */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* 标题区 */}
        <div className="text-center mb-12">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-50 to-cyan-50 border border-blue-200/50 rounded-full mb-6">
            <Shield className="w-4 h-4 text-blue-600" />
            <span className="text-sm font-medium text-blue-700">隐私政策</span>
          </div>
          <h1 className="text-4xl font-bold text-slate-900 mb-4">
            隐私政策
          </h1>
          <p className="text-lg text-slate-600">
            最后更新时间：2026 年 1 月 11 日
          </p>
        </div>

        {/* 内容区 */}
        <div className="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
          <div className="p-8 sm:p-12 prose prose-slate max-w-none">
            {/* 引言 */}
            <section className="mb-10">
              <p className="text-lg text-slate-700 leading-relaxed">
                欢迎使用微验 WeAuth！我们深知隐私安全的重要性，本隐私政策旨在帮助您了解我们如何收集、使用、存储和保护您的信息。
              </p>
            </section>

            {/* 1. 信息收集 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                1. 我们收集哪些信息
              </h2>
              
              <h3 className="text-xl font-semibold text-slate-800 mt-6 mb-3">
                1.1 账户信息
              </h3>
              <p className="text-slate-700 mb-4">
                当您注册微验 WeAuth 账户时，我们会收集：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-2 ml-4">
                <li>邮箱地址（用于账户登录和通知）</li>
                <li>密码（加密存储，我们无法查看明文密码）</li>
                <li>用户名或昵称</li>
              </ul>

              <h3 className="text-xl font-semibold text-slate-800 mt-6 mb-3">
                1.2 验证过程中的技术信息
              </h3>
              <p className="text-slate-700 mb-4">
                为了提供人机验证服务并防止滥用，我们会收集：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-2 ml-4">
                <li><strong>IP 地址</strong>：用于分析访问模式和检测异常行为</li>
                <li><strong>浏览器指纹信息</strong>：包括浏览器类型、操作系统、屏幕分辨率等（不包含个人身份信息）</li>
                <li><strong>验证结果数据</strong>：包括验证成功/失败次数、验证耗时等统计信息</li>
                <li><strong>设备信息</strong>：设备类型（移动端或桌面端）、操作系统版本</li>
              </ul>

              <h3 className="text-xl font-semibold text-slate-800 mt-6 mb-3">
                1.3 网站使用数据
              </h3>
              <p className="text-slate-700 mb-4">
                我们会收集您在使用我们服务时的基本统计数据：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-2 ml-4">
                <li>访问时间和频率</li>
                <li>功能使用情况</li>
                <li>错误日志（用于技术支持和服务改进）</li>
              </ul>
            </section>

            {/* 2. 信息使用 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                2. 我们如何使用您的信息
              </h2>
              <p className="text-slate-700 mb-4">
                我们收集的信息仅用于以下目的：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-3 ml-4">
                <li><strong>提供核心服务</strong>：实现基于工作量证明（Proof of Work）的人机验证功能</li>
                <li><strong>安全防护</strong>：检测和防止恶意攻击、机器人滥用和欺诈行为</li>
                <li><strong>服务优化</strong>：分析使用数据以改进验证算法和用户体验</li>
                <li><strong>技术支持</strong>：响应您的服务请求和技术问题</li>
                <li><strong>账户管理</strong>：维护您的账户安全，发送重要通知</li>
                <li><strong>合规要求</strong>：遵守法律法规和监管要求</li>
              </ul>
              <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
                <p className="text-sm text-blue-900">
                  <strong>重要承诺：</strong>我们绝不会将您的个人信息出售给第三方，也不会用于广告投放或营销推广。
                </p>
              </div>
            </section>

            {/* 3. 信息存储 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                3. 信息存储与保护
              </h2>
              
              <h3 className="text-xl font-semibold text-slate-800 mt-6 mb-3">
                3.1 数据存储
              </h3>
              <ul className="list-disc list-inside text-slate-700 space-y-2 ml-4">
                <li><strong>存储位置</strong>：您的数据存储在位于中国境内的安全服务器上</li>
                <li><strong>存储时长</strong>：
                  <ul className="list-circle list-inside ml-6 mt-2 space-y-1">
                    <li>账户信息：在您的账户有效期内</li>
                    <li>验证日志：保留 90 天（用于安全分析和统计）</li>
                    <li>IP 记录：保留 30 天（用于防止滥用）</li>
                  </ul>
                </li>
              </ul>

              <h3 className="text-xl font-semibold text-slate-800 mt-6 mb-3">
                3.2 安全措施
              </h3>
              <p className="text-slate-700 mb-4">
                我们采用业界标准的安全措施保护您的数据：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-2 ml-4">
                <li><strong>数据加密</strong>：使用 HTTPS/TLS 加密传输，敏感数据采用 AES-256 加密存储</li>
                <li><strong>访问控制</strong>：严格限制数据访问权限，仅授权人员可访问</li>
                <li><strong>密码安全</strong>：使用 bcrypt 等安全算法加密存储密码</li>
                <li><strong>安全审计</strong>：定期进行安全审计和渗透测试</li>
                <li><strong>备份机制</strong>：定期备份数据，确保数据可靠性</li>
              </ul>
            </section>

            {/* 4. 信息共享 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                4. 信息共享与披露
              </h2>
              <p className="text-slate-700 mb-4">
                我们承诺不会出售、出租或交易您的个人信息。仅在以下情况下可能共享您的信息：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-3 ml-4">
                <li><strong>您的明确授权</strong>：在获得您明确同意的情况下</li>
                <li><strong>服务提供商</strong>：与协助我们提供服务的第三方（如云服务提供商），但他们仅能在必要范围内访问数据，且受严格的保密协议约束</li>
                <li><strong>法律要求</strong>：根据法律法规、法律程序或政府机关的强制性要求</li>
                <li><strong>安全保护</strong>：为保护微验 WeAuth、用户或公众的权利、财产或安全</li>
                <li><strong>业务转让</strong>：在合并、收购或资产出售的情况下，但接收方需遵守本隐私政策</li>
              </ul>
            </section>

            {/* 5. Cookie 和追踪技术 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                5. Cookie 和追踪技术
              </h2>
              <p className="text-slate-700 mb-4">
                我们使用 Cookie 和类似技术来提升服务体验：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-3 ml-4">
                <li><strong>必要 Cookie</strong>：用于维持登录状态和基本功能（无法禁用）</li>
                <li><strong>功能 Cookie</strong>：记住您的偏好设置（如主题选择）</li>
                <li><strong>分析 Cookie</strong>：帮助我们了解网站使用情况以改进服务</li>
              </ul>
              <p className="text-slate-700 mt-4">
                您可以通过浏览器设置管理或删除 Cookie，但这可能影响部分功能的使用。
              </p>
            </section>

            {/* 6. 您的权利 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                6. 您的隐私权利
              </h2>
              <p className="text-slate-700 mb-4">
                根据相关法律法规，您享有以下权利：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-3 ml-4">
                <li><strong>访问权</strong>：查看我们持有的您的个人信息</li>
                <li><strong>更正权</strong>：要求更正不准确或不完整的信息</li>
                <li><strong>删除权</strong>：要求删除您的个人信息（法律要求保留的除外）</li>
                <li><strong>限制处理权</strong>：要求限制对您个人信息的处理</li>
                <li><strong>数据携带权</strong>：以结构化、常用格式获取您的数据</li>
                <li><strong>撤回同意权</strong>：随时撤回您对数据处理的同意</li>
                <li><strong>投诉权</strong>：向有关监管机构投诉</li>
              </ul>
              <div className="mt-4 p-4 bg-slate-50 border border-slate-200 rounded-lg">
                <p className="text-sm text-slate-700">
                  如需行使上述权利，请通过以下方式联系我们：
                  <br />
                  <strong>邮箱：</strong> <a href="mailto:privacy@weauth.cn" className="text-blue-600 hover:underline">privacy@weauth.cn</a>
                  <br />
                  我们将在 15 个工作日内响应您的请求。
                </p>
              </div>
            </section>

            {/* 7. 未成年人保护 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                7. 未成年人保护
              </h2>
              <p className="text-slate-700">
                微验 WeAuth 主要面向企业和开发者提供服务。如果您未满 18 周岁，请在监护人的陪同下使用我们的服务。我们不会故意收集未成年人的个人信息。如果您是家长或监护人，发现未成年人在未经许可的情况下向我们提供了个人信息，请及时联系我们，我们将尽快删除相关信息。
              </p>
            </section>

            {/* 8. 隐私政策更新 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                8. 隐私政策更新
              </h2>
              <p className="text-slate-700 mb-4">
                我们可能会不定期更新本隐私政策。更新时，我们会：
              </p>
              <ul className="list-disc list-inside text-slate-700 space-y-2 ml-4">
                <li>在本页面发布更新后的政策并更新"最后更新时间"</li>
                <li>对于重大变更，我们会通过邮件或网站公告提前通知您</li>
                <li>如果变更涉及您的重要权益，我们会征求您的同意</li>
              </ul>
              <p className="text-slate-700 mt-4">
                建议您定期查看本页面，以了解最新的隐私保护措施。
              </p>
            </section>

            {/* 9. 国际用户 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                9. 国际用户须知
              </h2>
              <p className="text-slate-700">
                微验 WeAuth 主要服务于中国大陆用户。如果您从中国大陆以外的地区访问和使用我们的服务，请注意您的信息可能会被传输到中国大陆并在此处理和存储。使用我们的服务即表示您同意此类传输。
              </p>
            </section>

            {/* 10. 联系我们 */}
            <section className="mb-10">
              <h2 className="text-2xl font-bold text-slate-900 mb-4 pb-3 border-b border-slate-200">
                10. 联系我们
              </h2>
              <p className="text-slate-700 mb-4">
                如果您对本隐私政策或我们的隐私保护措施有任何疑问、意见或建议，欢迎通过以下方式联系我们：
              </p>
              <div className="bg-gradient-to-r from-blue-50 to-cyan-50 border border-blue-200 rounded-lg p-6">
                <div className="space-y-3">
                  <p className="text-slate-800">
                    <strong>公司名称：</strong>微验科技
                  </p>
                  <p className="text-slate-800">
                    <strong>电子邮箱：</strong> <a href="mailto:privacy@weauth.cn" className="text-blue-600 hover:underline">privacy@weauth.cn</a>
                  </p>
                  <p className="text-slate-800">
                    <strong>技术支持：</strong> <a href="mailto:support@weauth.cn" className="text-blue-600 hover:underline">support@weauth.cn</a>
                  </p>
                  <p className="text-slate-800">
                    <strong>官方网站：</strong> <a href="https://www.weauth.cn" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">www.weauth.cn</a>
                  </p>
                </div>
              </div>
              <p className="text-slate-600 mt-4 text-sm">
                我们承诺在收到您的问题后 15 个工作日内给予回复。
              </p>
            </section>

            {/* 结语 */}
            <section className="mt-12 pt-8 border-t border-slate-200">
              <div className="text-center">
                <p className="text-slate-600 mb-4">
                  感谢您信任微验 WeAuth！
                </p>
                <p className="text-sm text-slate-500">
                  我们致力于保护您的隐私安全，为您提供安全可靠的人机验证服务。
                </p>
              </div>
            </section>
          </div>
        </div>

        {/* 底部返回按钮 */}
        <div className="text-center mt-12">
          <Link
            to="/"
            className="inline-flex items-center gap-2 px-6 py-3 bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 hover:text-blue-700 rounded-xl transition-all shadow-sm font-medium"
          >
            <ArrowLeft className="w-4 h-4" />
            返回首页
          </Link>
        </div>
      </main>

      {/* 简洁页脚 */}
      <footer className="border-t border-slate-200 bg-white/50 backdrop-blur-sm mt-16">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center text-sm text-slate-500">
            <p>© 2026 微验 WeAuth. 保留所有权利.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
