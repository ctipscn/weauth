import { createApp } from 'vue'
import WidgetApp from './Widget.vue'
import './assets/main.css'

const root = document.getElementById('widget-root')
if (root) {
  const app = createApp(WidgetApp)
  app.mount('#widget-root')
} else {
  console.error('WeAuth: widget-root element not found')
}