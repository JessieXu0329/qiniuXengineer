import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'

// Custom styling setup can be loaded here (integrated into App.vue's style for single-file elegance)

// Create app instance
const app = createApp(App)

// Initialize Pinia state store
const pinia = createPinia()
app.use(pinia)

// Set up light router paths so navigation runs seamlessly out-of-the-box
const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('./views/Dashboard.vue')
    },
    {
      path: '/review',
      name: 'review',
      component: () => import('./views/PRReview.vue')
    },
    {
      path: '/context',
      name: 'context',
      component: () => import('./views/HybridContext.vue')
    },
    {
      path: '/false-positive',
      name: 'false-positive',
      component: () => import('./views/FalsePositive.vue')
    },
    {
      path: '/config',
      name: 'config',
      component: () => import('./views/EngineConfig.vue')
    }
  ]
})
app.use(router)

// Initialize all Element Plus icons globally for easy accessibility
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.use(component)
}

app.use(ElementPlus)

app.mount('#app')
