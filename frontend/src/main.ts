import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import 'element-plus/theme-chalk/el-loading.css'
import 'element-plus/theme-chalk/el-message-box.css'
import App from './App.vue'

const pinia = createPinia()
createApp(App).use(pinia).mount('#app')
