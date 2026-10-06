<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { Shield, Lock, Eye, Server, FileText, Mail, ChevronRight, Scale, Database, Globe, AlertTriangle, Cookie } from 'lucide-vue-next'
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-cyan-50/30 pt-14">

    <!-- 主要内容 -->
    <main class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 pt-24">
      <!-- 标题区 -->
      <div class="text-center mb-12">
        <div class="inline-flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-50 to-cyan-50 border border-blue-200/50 rounded-full mb-6">
          <Shield class="w-4 h-4 text-blue-600" />
          <span class="text-sm font-medium text-blue-700">隐私政策</span>
        </div>
        <h1 class="text-4xl font-bold text-slate-900 mb-4">
          隐私政策
        </h1>
        <p class="text-lg text-slate-600">
          最后更新时间：2026 年 1 月 11 日
        </p>
      </div>

      <!-- 内容区 -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-sm overflow-hidden">
        <div class="p-8 sm:p-12 prose prose-slate max-w-none">
          <!-- 引言 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Shield class="w-5 h-5 text-blue-500" />
              <p class="text-lg text-slate-700 leading-relaxed">
                欢迎使用微验 WeAuth！我们深知隐私安全的重要性，本隐私政策旨在帮助您了解我们如何收集、使用、存储和保护您的信息。
              </p>
            </div>
          </section>

          <!-- 1. 信息收集 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Database class="w-5 h-5 text-blue-500" />
              <h2 class="text-2xl font-bold text-slate-900">1. 我们收集哪些信息</h2>
            </div>

            <h3 class="text-xl font-semibold text-slate-800 mt-6 mb-3 flex items-center gap-2">
              <Lock class="w-4 h-4 text-slate-500" />
              1.1 账户信息
            </h3>
            <p class="text-slate-700 mb-4">
              当您注册微验 WeAuth 账户时，我们会收集：
            </p>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li>邮箱地址（用于账户登录和通知）</li>
              <li>密码（加密存储，我们无法查看明文密码）</li>
              <li>用户名或昵称</li>
            </ul>

            <h3 class="text-xl font-semibold text-slate-800 mt-6 mb-3 flex items-center gap-2">
              <Eye class="w-4 h-4 text-slate-500" />
              1.2 验证过程中的技术信息
            </h3>
            <p class="text-slate-700 mb-4">
              为了提供人机验证服务并防止滥用，我们会收集：
            </p>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li><strong>IP 地址</strong>：用于分析访问模式和检测异常行为</li>
              <li><strong>浏览器指纹信息</strong>：包括浏览器类型、操作系统、屏幕分辨率等（不包含个人身份信息）</li>
              <li><strong>交互行为数据</strong>：例如鼠标移动轨迹、点击模式等（仅在验证过程中临时处理，不存储）</li>
            </ul>

            <h3 class="text-xl font-semibold text-slate-800 mt-6 mb-3">
              1.3 使用信息
            </h3>
            <p class="text-slate-700">
              我们会收集与您使用服务相关的信息，包括 API 调用次数、站点配置变更记录、登录日志等，以改进服务和维护安全。
            </p>
          </section>

          <!-- 2. 信息使用 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Server class="w-5 h-5 text-blue-500" />
              <h2 class="text-2xl font-bold text-slate-900">2. 我们如何使用您的信息</h2>
            </div>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li><strong>提供和维护服务</strong>：身份验证、人机验证、站点管理</li>
              <li><strong>安全防护</strong>：检测和防止滥用行为、DDoS 攻击、自动化脚本</li>
              <li><strong>服务改进</strong>：分析使用模式和性能指标、优化验证算法</li>
              <li><strong>合规要求</strong>：满足法律法规要求</li>
            </ul>
          </section>

          <!-- 3. 信息存储 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Server class="w-5 h-5 text-green-500" />
              <h2 class="text-2xl font-bold text-slate-900">3. 信息存储和保护</h2>
            </div>
            <p class="text-slate-700 mb-4">我们采取以下措施保护您的信息安全：</p>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li>采用 TLS 1.3 加密传输所有数据</li>
              <li>密码使用 bcrypt 算法加盐哈希存储</li>
              <li>实施严格的访问控制和审计日志</li>
              <li>服务器部署在中国大陆，符合国内数据安全法规</li>
              <li>定期进行安全审计和渗透测试</li>
            </ul>
          </section>

          <!-- 4. Cookie 使用 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Cookie class="w-5 h-5 text-amber-500" />
              <h2 class="text-2xl font-bold text-slate-900">4. Cookie 使用说明</h2>
            </div>
            <p class="text-slate-700 mb-4">我们仅使用必要的 Cookie 来维持服务运行：</p>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li><strong>会话 Cookie</strong>：维持登录状态，过期后自动删除</li>
              <li><strong>安全 Cookie</strong>：用于 CSRF 防护，防止跨站请求伪造</li>
            </ul>
            <p class="text-slate-600 mt-4 text-sm">
              我们不会使用跟踪 Cookie 或第三方广告 Cookie。
            </p>
          </section>

          <!-- 5. 数据共享 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Globe class="w-5 h-5 text-blue-500" />
              <h2 class="text-2xl font-bold text-slate-900">5. 数据共享与披露</h2>
            </div>
            <p class="text-slate-700 mb-4">我们不会向第三方出售您的个人信息。在以下情况下，我们可能会披露您的信息：</p>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li>获得您的明确同意</li>
              <li>法律法规要求</li>
              <li>保护我们的合法权益（如防止欺诈和滥用）</li>
              <li>与服务提供商合作（如云服务提供商），但仅限于必要范围</li>
            </ul>
          </section>

          <!-- 6. 用户权利 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Scale class="w-5 h-5 text-purple-500" />
              <h2 class="text-2xl font-bold text-slate-900">6. 您的权利</h2>
            </div>
            <p class="text-slate-700 mb-4">根据适用的数据保护法律，您享有以下权利：</p>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li><strong>知情权</strong>：了解我们如何收集和使用您的信息</li>
              <li><strong>访问权</strong>：请求访问我们持有的您的个人信息</li>
              <li><strong>更正权</strong>：更正不准确的个人信息</li>
              <li><strong>删除权</strong>：在特定条件下请求删除您的个人信息</li>
              <li><strong>撤回同意权</strong>：撤回您对信息处理的同意</li>
            </ul>
            <p class="text-slate-600 mt-4">
              如需行使上述权利，请通过下方联系方式与我们联系。我们将在 15 个工作日内回复您的请求。
            </p>
          </section>

          <!-- 7. 未成年人 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <AlertTriangle class="w-5 h-5 text-amber-500" />
              <h2 class="text-2xl font-bold text-slate-900">7. 未成年人保护</h2>
            </div>
            <p class="text-slate-700">
              我们的服务面向 18 周岁以上的用户。如果您未满 18 周岁，请在监护人的陪同下使用我们的服务。我们不会故意收集未成年人的个人信息。如果您是家长或监护人，发现未成年人在未经许可的情况下向我们提供了个人信息，请及时联系我们，我们将尽快删除相关信息。
            </p>
          </section>

          <!-- 8. 隐私政策更新 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <FileText class="w-5 h-5 text-blue-500" />
              <h2 class="text-2xl font-bold text-slate-900">8. 隐私政策更新</h2>
            </div>
            <p class="text-slate-700 mb-4">
              我们可能会不定期更新本隐私政策。更新时，我们会：
            </p>
            <ul class="list-disc list-inside text-slate-700 space-y-2 ml-4">
              <li>在本页面发布更新后的政策并更新"最后更新时间"</li>
              <li>对于重大变更，我们会通过邮件或网站公告提前通知您</li>
              <li>如果变更涉及您的重要权益，我们会征求您的同意</li>
            </ul>
            <p class="text-slate-700 mt-4">
              建议您定期查看本页面，以了解最新的隐私保护措施。
            </p>
          </section>

          <!-- 9. 国际用户 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Globe class="w-5 h-5 text-blue-500" />
              <h2 class="text-2xl font-bold text-slate-900">9. 国际用户须知</h2>
            </div>
            <p class="text-slate-700">
              微验 WeAuth 主要服务于中国大陆用户。如果您从中国大陆以外的地区访问和使用我们的服务，请注意您的信息可能会被传输到中国大陆并在此处理和存储。使用我们的服务即表示您同意此类传输。
            </p>
          </section>

          <!-- 10. 联系我们 -->
          <section class="mb-10">
            <div class="flex items-center gap-2 mb-4">
              <Mail class="w-5 h-5 text-blue-500" />
              <h2 class="text-2xl font-bold text-slate-900">10. 联系我们</h2>
            </div>
            <p class="text-slate-700 mb-4">
              如果您对本隐私政策或我们的隐私保护措施有任何疑问、意见或建议，欢迎通过以下方式联系我们：
            </p>
            <div class="bg-gradient-to-r from-blue-50 to-cyan-50 border border-blue-200 rounded-lg p-6">
              <div class="space-y-3">
                <p class="text-slate-800">
                  <strong>公司名称：</strong>微验科技
                </p>
                <p class="text-slate-800">
                  <strong>电子邮箱：</strong> <a href="mailto:privacy@weauth.cn" class="text-blue-600 hover:underline">privacy@weauth.cn</a>
                </p>
                <p class="text-slate-800">
                  <strong>技术支持：</strong> <a href="mailto:support@weauth.cn" class="text-blue-600 hover:underline">support@weauth.cn</a>
                </p>
                <p class="text-slate-800">
                  <strong>官方网站：</strong> <a href="https://www.weauth.cn" target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:underline">www.weauth.cn</a>
                </p>
              </div>
            </div>
            <p class="text-slate-600 mt-4 text-sm">
              我们承诺在收到您的问题后 15 个工作日内给予回复。
            </p>
          </section>

          <!-- 结语 -->
          <section class="mt-12 pt-8 border-t border-slate-200">
            <div class="text-center">
              <p class="text-slate-600 mb-4">
                感谢您信任微验 WeAuth！
              </p>
              <p class="text-sm text-slate-500">
                我们致力于保护您的隐私安全，为您提供安全可靠的人机验证服务。
              </p>
            </div>
          </section>
        </div>
      </div>

      <!-- 底部返回按钮 -->
      <div class="text-center mt-12">
        <RouterLink
          to="/"
          class="inline-flex items-center gap-2 px-6 py-3 bg-white border-2 border-blue-600 text-blue-600 hover:bg-blue-50 hover:border-blue-700 hover:text-blue-700 rounded-xl transition-all shadow-sm font-medium"
        >
          <ChevronRight class="w-4 h-4" />
          返回首页
        </RouterLink>
      </div>
    </main>

    <!-- 简洁页脚 -->
    <footer class="border-t border-slate-200 bg-white/50 backdrop-blur-sm mt-16">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="text-center text-sm text-slate-500">
          <p>© 2026 微验 WeAuth. 保留所有权利.</p>
        </div>
      </div>
    </footer>
  </div>
</template>