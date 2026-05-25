<template>
  <div>
    <div class="page-header">
      <div style="display:flex;align-items:center;gap:12px">
        <router-link :to="`/project?id=${release?.projectId}`" class="btn btn-sm">← 返回</router-link>
        <h1 class="page-title">{{ release?.name || '发布单' }}</h1>
        <span v-if="release" :class="['badge', release.status === 'published' ? 'badge-success' : 'badge-info']">
          {{ release.status === 'draft' ? '草稿' : release.status === 'published' ? '已发布' : '已归档' }}
        </span>
      </div>
      <div class="actions" v-if="release">
        <button class="btn" @click="preview">预览</button>
        <button class="btn" @click="refreshItems">刷新数据</button>
        <button class="btn btn-primary" @click="publish">发布</button>
        <button class="btn" @click="exportRelease('markdown')">导出 MD</button>
        <button class="btn" @click="exportRelease('html')">导出 HTML</button>
      </div>
    </div>

    <div v-if="release" class="card mb-16" style="display:flex;gap:24px;flex-wrap:wrap">
      <div><span style="color:var(--text-tertiary)">版本：</span>{{ release.version || '-' }}</div>
      <div><span style="color:var(--text-tertiary)">条目：</span>{{ release.itemCount }}</div>
      <div><span style="color:var(--text-tertiary)">Bug：</span>{{ release.bugCount }}</div>
      <div><span style="color:var(--text-tertiary)">Task：</span>{{ release.taskCount }}</div>
      <div><span style="color:var(--text-tertiary)">备注：</span>{{ release.noteCount }}</div>
      <div><span style="color:var(--text-tertiary)">部署：</span>{{ deployments.length }}</div>
      <div><span style="color:var(--text-tertiary)">发布次数：</span>{{ release.publishCount }}</div>
      <div v-if="release.summary" style="width:100%"><span style="color:var(--text-tertiary)">概述：</span>{{ release.summary }}</div>
    </div>

    <!-- 部署地址 -->
    <div class="card mb-16">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
        <h3 style="font-size:15px">部署地址</h3>
        <button class="btn btn-sm" @click="showDeploymentForm = true">添加部署地址</button>
      </div>
      <table class="data-table" v-if="deployments.length">
        <thead><tr><th>功能模块</th><th>地址</th><th>说明</th><th>添加时间</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="d in deployments" :key="d.id">
            <td>{{ d.moduleName }}</td>
            <td><a :href="d.address" target="_blank" style="color:var(--primary)">{{ d.address }}</a></td>
            <td>{{ d.description || '-' }}</td>
            <td>{{ formatTime(d.createdAt) }}</td>
            <td><button class="btn btn-sm btn-danger" @click="removeDeployment(d.id)">删除</button></td>
          </tr>
        </tbody>
      </table>
      <div v-else style="color:var(--text-tertiary);padding:12px 0">暂无部署地址，点击「添加部署地址」</div>

      <div v-if="showDeploymentForm" style="margin-top:12px;padding-top:12px;border-top:1px solid var(--border)">
        <div style="display:flex;gap:8px;align-items:flex-end">
          <div class="form-group" style="flex:1;margin-bottom:0">
            <label style="font-size:12px">功能模块</label>
            <input v-model="deployForm.moduleName" placeholder="前端服务" />
          </div>
          <div class="form-group" style="flex:1.5;margin-bottom:0">
            <label style="font-size:12px">地址</label>
            <input v-model="deployForm.address" placeholder="https://web.example.com" />
          </div>
          <div class="form-group" style="flex:1;margin-bottom:0">
            <label style="font-size:12px">说明</label>
            <input v-model="deployForm.description" placeholder="Nginx容器" />
          </div>
          <button class="btn btn-primary btn-sm" @click="submitDeployment">添加</button>
          <button class="btn btn-sm" @click="showDeploymentForm = false">取消</button>
        </div>
      </div>
    </div>

    <!-- 条目管理 -->
    <div class="card mb-16">
      <div style="display:flex;gap:8px;margin-bottom:16px">
        <button class="btn" @click="showBugSelector = true">添加 Bug</button>
        <button class="btn" @click="showTaskSelector = true">添加 Task</button>
        <button class="btn" @click="addNote">添加备注</button>
      </div>

      <div v-if="bugs.length" style="margin-bottom:20px">
        <h3 style="font-size:15px;margin-bottom:8px;color:var(--danger)">Bug 修复（{{ bugs.length }}）</h3>
        <table class="data-table">
          <thead><tr><th>ID</th><th>标题</th><th>严重程度</th><th>状态</th><th>指派给</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="item in bugs" :key="item.id">
              <td>
                <a v-if="item.zentaoUrl" :href="item.zentaoUrl" target="_blank" style="color:var(--primary)">#{{ item.zentaoId }}</a>
                <span v-else>#{{ item.zentaoId }}</span>
              </td>
              <td>
                <a v-if="zentaoBaseUrl && item.zentaoId" :href="zentaoBaseUrl + '/bug-view-' + item.zentaoId + '.html'" target="_blank" style="color:var(--text-primary);text-decoration:none">{{ item.title }}</a>
                <span v-else>{{ item.title }}</span>
              </td>
              <td>{{ item.severity }}</td>
              <td>{{ item.status }}</td>
              <td>{{ item.assignedTo }}</td>
              <td><button class="btn btn-sm btn-danger" @click="removeItem(item.id)">删除</button></td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="tasks.length" style="margin-bottom:20px">
        <h3 style="font-size:15px;margin-bottom:8px;color:var(--success)">任务完成（{{ tasks.length }}）</h3>
        <table class="data-table">
          <thead><tr><th>ID</th><th>标题</th><th>优先级</th><th>状态</th><th>指派给</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="item in tasks" :key="item.id">
              <td>
                <a v-if="item.zentaoUrl" :href="item.zentaoUrl" target="_blank" style="color:var(--primary)">#{{ item.zentaoId }}</a>
                <span v-else>#{{ item.zentaoId }}</span>
              </td>
              <td>
                <a v-if="zentaoBaseUrl && item.zentaoId" :href="zentaoBaseUrl + '/task-view-' + item.zentaoId + '.html'" target="_blank" style="color:var(--text-primary);text-decoration:none">{{ item.title }}</a>
                <span v-else>{{ item.title }}</span>
              </td>
              <td>{{ item.priority }}</td>
              <td>{{ item.status }}</td>
              <td>{{ item.assignedTo }}</td>
              <td><button class="btn btn-sm btn-danger" @click="removeItem(item.id)">删除</button></td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="notes.length">
        <h3 style="font-size:15px;margin-bottom:8px;color:var(--primary)">备注（{{ notes.length }}）</h3>
        <div v-for="item in notes" :key="item.id" style="padding:12px;border:1px solid var(--border);border-radius:6px;margin-bottom:8px">
          <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px">
            <strong>{{ item.noteTitle }}</strong>
            <button class="btn btn-sm btn-danger" @click="removeItem(item.id)">删除</button>
          </div>
          <p style="color:var(--text-secondary);white-space:pre-wrap">{{ item.noteContent }}</p>
        </div>
      </div>

      <div v-if="!items.length" class="empty-state">
        <p>暂无条目，点击上方按钮添加</p>
      </div>
    </div>

    <!-- Bug 选择器 -->
    <div v-if="showBugSelector" class="card mb-16">
      <h3 style="margin-bottom:12px">选择 Bug</h3>
      <div style="display:flex;gap:8px;margin-bottom:12px">
        <button class="btn btn-primary btn-sm" @click="loadBugs">加载 Bug</button>
        <button class="btn btn-sm" @click="confirmBugs" :disabled="!selectedBugs.length">确认添加 ({{ selectedBugs.length }})</button>
        <button class="btn btn-sm" @click="showBugSelector = false">关闭</button>
      </div>
      <table class="data-table" v-if="bugList.length">
        <thead><tr><th>选择</th><th>ID</th><th>标题</th><th>严重程度</th><th>状态</th></tr></thead>
        <tbody>
          <tr v-for="b in bugList" :key="b.id">
            <td><input type="checkbox" :value="b" v-model="selectedBugs" /></td>
            <td>{{ b.id }}</td>
            <td>{{ b.title }}</td>
            <td>{{ b.severity }}</td>
            <td>{{ b.status }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="!bugList.length" style="color:var(--text-tertiary);padding:20px;text-align:center">
        {{ bugLoading ? '加载中...' : '点击「加载 Bug」从禅道获取数据（需要项目已关联禅道产品）' }}
      </div>
    </div>

    <!-- Task 选择器 -->
    <div v-if="showTaskSelector" class="card mb-16">
      <h3 style="margin-bottom:12px">选择 Task</h3>
      <div style="display:flex;gap:8px;margin-bottom:12px">
        <button class="btn btn-primary btn-sm" @click="loadTasks">加载 Task</button>
        <button class="btn btn-sm" @click="confirmTasks" :disabled="!selectedTasks.length">确认添加 ({{ selectedTasks.length }})</button>
        <button class="btn btn-sm" @click="showTaskSelector = false">关闭</button>
      </div>
      <table class="data-table" v-if="taskList.length">
        <thead><tr><th>选择</th><th>ID</th><th>标题</th><th>优先级</th><th>状态</th><th>指派给</th></tr></thead>
        <tbody>
          <tr v-for="t in taskList" :key="t.id">
            <td><input type="checkbox" :value="t" v-model="selectedTasks" /></td>
            <td>{{ t.id }}</td>
            <td>{{ t.name }}</td>
            <td>{{ t.pri }}</td>
            <td>{{ t.status }}</td>
            <td>{{ t.assignedTo?.realname || t.assignedTo || '' }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="!taskList.length" style="color:var(--text-tertiary);padding:20px;text-align:center">
        {{ taskLoading ? '加载中...' : '点击「加载 Task」从禅道获取数据（需要项目已关联禅道产品）' }}
      </div>
    </div>

    <!-- 添加备注 -->
    <div v-if="showNoteForm" class="card mb-16">
      <h3 style="margin-bottom:12px">添加备注</h3>
      <div class="form-group">
        <label>标题</label>
        <input v-model="noteForm.title" placeholder="备注标题" />
      </div>
      <div class="form-group">
        <label>内容</label>
        <textarea v-model="noteForm.content" placeholder="备注内容（支持多行）" />
      </div>
      <div class="actions">
        <button class="btn btn-primary" @click="submitNote">添加</button>
        <button class="btn" @click="showNoteForm = false">取消</button>
      </div>
    </div>

    <!-- 发布历史 -->
    <div class="card" v-if="snapshots.length">
      <h3 style="font-size:15px;margin-bottom:12px">发布历史</h3>
      <table class="data-table">
        <thead><tr><th>版本</th><th>内容统计</th><th>发布时间</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="s in snapshots" :key="s.id">
            <td>{{ s.version || '-' }}</td>
            <td>{{ s.bugCount }} Bug / {{ s.taskCount }} Task / {{ s.noteCount }} 备注</td>
            <td>{{ s.publishedAt }}</td>
            <td class="actions">
              <button class="btn btn-sm" @click="viewSnapshot(s)">查看</button>
              <button class="btn btn-sm" @click="exportRelease('markdown', s.id)">导出</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 预览模态框 -->
    <div v-if="showPreviewModal" style="position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000" @click.self="showPreviewModal=false">
      <div class="card" style="width:800px;max-height:85vh;overflow-y:auto">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
          <h3>发布单预览</h3>
          <button class="btn btn-sm" @click="showPreviewModal=false">关闭</button>
        </div>
        <div v-if="previewHtml" v-html="previewHtml" style="font-size:14px;line-height:1.8"></div>
      </div>
    </div>

    <!-- 快照查看模态框 -->
    <div v-if="showSnapshotModal" style="position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000" @click.self="showSnapshotModal=false">
      <div class="card" style="width:700px;max-height:80vh;overflow-y:auto">
        <h3 style="margin-bottom:12px">发布快照</h3>
        <pre style="white-space:pre-wrap;font-size:13px;line-height:1.6;background:var(--bg);padding:16px;border-radius:6px">{{ snapshotContent }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { releaseApi, itemApi, snapshotApi, zentaoApi, projectApi, deploymentApi, healthApi } from '../api'
import type { Release, ReleaseItem, ReleaseSnapshot, Deployment } from '../api'

const route = useRoute()
const release = ref<Release | null>(null)
const items = ref<ReleaseItem[]>([])
const snapshots = ref<ReleaseSnapshot[]>([])
const deployments = ref<Deployment[]>([])
const showBugSelector = ref(false)
const showTaskSelector = ref(false)
const showNoteForm = ref(false)
const showSnapshotModal = ref(false)
const showPreviewModal = ref(false)
const showDeploymentForm = ref(false)
const snapshotContent = ref('')
const previewHtml = ref('')
const bugList = ref<any[]>([])
const bugLoading = ref(false)
const taskList = ref<any[]>([])
const taskLoading = ref(false)
const selectedBugs = ref<any[]>([])
const selectedTasks = ref<any[]>([])
const noteForm = ref({ title: '', content: '' })
const deployForm = ref({ moduleName: '', address: '', description: '' })
const zentaoBaseUrl = ref('')

function formatTime(t: string) {
  if (!t) return '-'
  return t.replace('T', ' ').substring(0, 19)
}

const bugs = computed(() => items.value.filter(i => i.itemType === 'bug'))
const tasks = computed(() => items.value.filter(i => i.itemType === 'task'))
const notes = computed(() => items.value.filter(i => i.itemType === 'note'))

async function loadZentaoBaseUrl() {
  try {
    const resp: any = await healthApi.get()
    zentaoBaseUrl.value = resp?.zentaoBaseUrl || ''
  } catch {}
}

async function load() {
  const id = route.query.id as string
  if (!id) return
  const rResp: any = await releaseApi.get(id)
  release.value = rResp?.data || null
  await Promise.all([loadItems(), loadSnapshots(), loadDeployments()])
}

async function loadItems() {
  const id = route.query.id as string
  const iResp: any = await itemApi.list(id)
  items.value = iResp?.list || []
}

async function loadSnapshots() {
  const id = route.query.id as string
  const sResp: any = await snapshotApi.list(id)
  snapshots.value = sResp?.list || []
}

async function loadDeployments() {
  const id = route.query.id as string
  const dResp: any = await deploymentApi.list(id)
  deployments.value = dResp?.list || []
}

async function refreshItems() {
  await itemApi.refresh(route.query.id as string)
  await loadItems()
}

async function removeItem(id: string) {
  if (!confirm('确定删除？')) return
  await itemApi.delete(id)
  loadItems()
}

async function removeDeployment(id: string) {
  if (!confirm('确定删除？')) return
  await deploymentApi.delete(id)
  loadDeployments()
}

function addNote() {
  noteForm.value = { title: '', content: '' }
  showNoteForm.value = true
}

async function submitNote() {
  await itemApi.add({
    releaseId: route.query.id as string,
    itemType: 'note',
    noteTitle: noteForm.value.title,
    noteContent: noteForm.value.content
  })
  showNoteForm.value = false
  loadItems()
}

async function submitDeployment() {
  if (!deployForm.value.moduleName || !deployForm.value.address) {
    alert('请填写功能模块和地址')
    return
  }
  await deploymentApi.add({
    releaseId: route.query.id as string,
    moduleName: deployForm.value.moduleName,
    address: deployForm.value.address,
    description: deployForm.value.description
  })
  deployForm.value = { moduleName: '', address: '', description: '' }
  showDeploymentForm.value = false
  loadDeployments()
}

async function loadBugs() {
  const projectId = release.value?.projectId
  let productId = 0
  if (projectId) {
    try {
      const pResp: any = await projectApi.get(projectId)
      productId = pResp?.data?.zentaoProductId || 0
    } catch {}
  }
  if (!productId) {
    alert('请先在项目设置中关联禅道产品')
    return
  }
  bugLoading.value = true
  try {
    const resp: any = await zentaoApi.bugs({ productId, page: 1, pageSize: 50 })
    const raw = typeof resp?.list === 'string' ? JSON.parse(resp.list) : resp?.list
    bugList.value = Array.isArray(raw) ? raw : []
  } finally {
    bugLoading.value = false
  }
}

async function confirmBugs() {
  const releaseId = route.query.id as string
  const base = zentaoBaseUrl.value
  const its = selectedBugs.value.map(b => ({
    releaseId,
    itemType: 'bug' as const,
    zentaoId: b.id,
    zentaoType: 'bug',
    title: b.title,
    severity: String(b.severity ?? ''),
    priority: String(b.pri ?? ''),
    status: b.status,
    assignedTo: b.assignedTo?.realname || '',
    zentaoUrl: base ? `${base}/bug-view-${b.id}.html` : ''
  }))
  await itemApi.batchAdd(releaseId, its)
  selectedBugs.value = []
  showBugSelector.value = false
  loadItems()
}

async function loadTasks() {
  const projectId = release.value?.projectId
  let productId = 0
  if (projectId) {
    try {
      const pResp: any = await projectApi.get(projectId)
      productId = pResp?.data?.zentaoProductId || 0
    } catch {}
  }
  if (!productId) {
    alert('请先在项目设置中关联禅道产品')
    return
  }
  taskLoading.value = true
  try {
    const resp: any = await zentaoApi.tasks({ productId, page: 1, pageSize: 50 })
    const raw = typeof resp?.list === 'string' ? JSON.parse(resp.list) : resp?.list
    taskList.value = Array.isArray(raw) ? raw : []
  } finally {
    taskLoading.value = false
  }
}

async function confirmTasks() {
  const releaseId = route.query.id as string
  const base = zentaoBaseUrl.value
  const its = selectedTasks.value.map(t => ({
    releaseId,
    itemType: 'task' as const,
    zentaoId: t.id,
    zentaoType: 'task',
    title: t.name || t.title || '',
    severity: '',
    priority: String(t.pri ?? ''),
    status: t.status,
    assignedTo: t.assignedTo?.realname || t.assignedTo || '',
    zentaoUrl: base ? `${base}/task-view-${t.id}.html` : ''
  }))
  await itemApi.batchAdd(releaseId, its)
  selectedTasks.value = []
  showTaskSelector.value = false
  loadItems()
}

async function publish() {
  const version = prompt('输入版本号（可选）', release.value?.version || '')
  if (version === null) return
  await releaseApi.publish(route.query.id as string, version || undefined)
  load()
}

async function preview() {
  const resp: any = await releaseApi.export(route.query.id as string, 'html')
  previewHtml.value = resp?.content || ''
  showPreviewModal.value = true
}

async function viewSnapshot(s: ReleaseSnapshot) {
  const resp: any = await snapshotApi.get(s.id)
  snapshotContent.value = resp?.data?.content || ''
  showSnapshotModal.value = true
}

async function exportRelease(format: string, snapshotId?: string) {
  const resp: any = await releaseApi.export(route.query.id as string, format, snapshotId)
  const content = resp?.content || ''
  const filename = resp?.filename || 'release.md'
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(() => {
  loadZentaoBaseUrl()
  load()
})
</script>
