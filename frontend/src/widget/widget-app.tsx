import { createRoot } from 'react-dom/client'
import Widget from '@/pages/Widget'
import '@/index.css'

const root = document.getElementById('widget-root')
if (root) {
  createRoot(root).render(<Widget />)
} else {
  console.error('WeAuth: widget-root element not found')
}
