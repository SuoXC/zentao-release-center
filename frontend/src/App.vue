<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="logo">
        <span class="logo-icon">R</span>
        <span class="logo-text">发布中心</span>
      </div>
      <nav class="nav">
        <router-link to="/projects" class="nav-item" active-class="active">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
          项目管理
        </router-link>
      </nav>
      <div class="sidebar-footer">
        <span>zentao-release-center</span>
        <span class="copyright">© 2024 murphyyi</span>
      </div>
    </aside>
    <main class="content">
      <router-view />
      <div v-if="toast.visible" :class="['toast', 'toast-' + toast.type]">
        {{ toast.message }}
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted, onUnmounted } from 'vue'

const toast = reactive({ visible: false, type: 'error', message: '' })
let timer: ReturnType<typeof setTimeout> | null = null

function showToast(e: Event) {
  const d = (e as CustomEvent).detail
  toast.type = d.type || 'error'
  toast.message = d.message
  toast.visible = true
  if (timer) clearTimeout(timer)
  timer = setTimeout(() => { toast.visible = false }, 4000)
}

onMounted(() => window.addEventListener('app:toast', showToast))
onUnmounted(() => window.removeEventListener('app:toast', showToast))
</script>

<style>
:root {
  --primary: #4F6BF6;
  --primary-light: #E8ECFF;
  --success: #22C55E;
  --success-light: #DCFCE7;
  --warning: #F59E0B;
  --danger: #EF4444;
  --danger-light: #FEE2E2;
  --bg: #F8FAFC;
  --bg-card: #FFFFFF;
  --bg-hover: #F1F5F9;
  --text-primary: #1E293B;
  --text-secondary: #64748B;
  --text-tertiary: #94A3B8;
  --border: #E2E8F0;
  --radius: 8px;
  --shadow: 0 1px 3px rgba(0,0,0,0.08);
}

* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 14px; color: var(--text-primary); background: var(--bg); }
html, body, #app { height: 100%; }

.layout { display: flex; height: 100vh; }
.sidebar { width: 220px; background: #1E293B; display: flex; flex-direction: column; flex-shrink: 0; }
.logo { height: 56px; display: flex; align-items: center; gap: 10px; padding: 0 20px; border-bottom: 1px solid #334155; }
.logo-icon { width: 28px; height: 28px; border-radius: 6px; background: var(--primary); color: white; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 14px; }
.logo-text { color: #E2E8F0; font-size: 16px; font-weight: 600; }
.nav { flex: 1; padding: 12px 10px; display: flex; flex-direction: column; gap: 2px; }
.nav-item { display: flex; align-items: center; gap: 10px; padding: 10px 14px; border-radius: 6px; color: #94A3B8; text-decoration: none; font-size: 14px; transition: all 0.15s; }
.nav-item svg { width: 18px; height: 18px; }
.nav-item:hover { background: #334155; color: #E2E8F0; }
.nav-item.active { background: var(--primary); color: white; }
.sidebar-footer { padding: 12px 20px; border-top: 1px solid #334155; font-size: 11px; color: #64748B; text-align: center; }
.sidebar-footer .copyright { display: block; margin-top: 4px; font-size: 10px; }
.content { flex: 1; overflow-y: auto; padding: 24px; }

.btn { padding: 8px 16px; border-radius: 6px; border: 1px solid var(--border); background: white; color: var(--text-primary); font-size: 13px; cursor: pointer; transition: all 0.15s; display: inline-flex; align-items: center; gap: 6px; }
.btn:hover { border-color: var(--primary); color: var(--primary); }
.btn-primary { background: var(--primary); color: white; border-color: var(--primary); }
.btn-primary:hover { background: #3A51D4; }
.btn-danger { color: var(--danger); border-color: var(--danger); }
.btn-danger:hover { background: var(--danger-light); }
.btn-sm { padding: 4px 10px; font-size: 12px; }

.card { background: var(--bg-card); border-radius: var(--radius); border: 1px solid var(--border); padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-title { font-size: 20px; font-weight: 600; }
.badge { display: inline-block; padding: 2px 10px; border-radius: 100px; font-size: 12px; font-weight: 500; }
.badge-success { background: var(--success-light); color: var(--success); }
.badge-warning { background: #FEF3C7; color: #D97706; }
.badge-danger { background: var(--danger-light); color: var(--danger); }
.badge-info { background: var(--primary-light); color: var(--primary); }

table.data-table { width: 100%; border-collapse: collapse; }
table.data-table th, table.data-table td { padding: 10px 14px; text-align: left; border-bottom: 1px solid var(--border); }
table.data-table th { font-weight: 600; color: var(--text-secondary); font-size: 12px; text-transform: uppercase; letter-spacing: 0.5px; }
table.data-table tr:hover { background: var(--bg-hover); }

.empty-state { text-align: center; padding: 60px 20px; color: var(--text-tertiary); }
.empty-state h3 { font-size: 16px; margin-bottom: 8px; color: var(--text-secondary); }

.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; font-weight: 500; font-size: 13px; color: var(--text-secondary); }
.form-group input, .form-group textarea, .form-group select { width: 100%; padding: 8px 12px; border: 1px solid var(--border); border-radius: 6px; font-size: 14px; outline: none; transition: border-color 0.15s; }
.form-group input:focus, .form-group textarea:focus, .form-group select:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(79,107,246,0.12); }
.form-group textarea { min-height: 80px; resize: vertical; }

.actions { display: flex; gap: 8px; }
.mt-16 { margin-top: 16px; }
.mb-16 { margin-bottom: 16px; }
.flex-between { display: flex; justify-content: space-between; align-items: center; }
.gap-8 { gap: 8px; }

.toast { position: fixed; top: 20px; right: 20px; padding: 12px 20px; border-radius: 8px; font-size: 13px; z-index: 9999; animation: slideIn 0.3s ease; box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
.toast-error { background: #FEE2E2; color: #DC2626; border: 1px solid #FECACA; }
.toast-success { background: #DCFCE7; color: #16A34A; border: 1px solid #BBF7D0; }
@keyframes slideIn { from { transform: translateX(100%); opacity: 0; } to { transform: translateX(0); opacity: 1; } }

.loading-spinner { display: inline-block; width: 16px; height: 16px; border: 2px solid var(--border); border-top-color: var(--primary); border-radius: 50%; animation: spin 0.6s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
