<template>
  <div>
    <div class="page-header">
      <div style="display:flex;align-items:center;gap:12px">
        <router-link :to="`/project?id=${release?.projectId}`" class="btn btn-sm">← 返回</router-link>
        <h1 class="page-title">{{ release?.name || '发布单' }}</h1>
        <span v-if="release" :class="['badge', release.status === 'published' ? 'badge-success' : 'badge-info']">
          {{ release.status === 'draft' ? '草稿' : release.status === 'published' ? '已发布' : '已归档' }}
        </span>
        <span v-if="loading" class="loading-spinner"></span>
      </div>
      <div class="actions" v-if="release">
        <button class="btn" @click="preview" :disabled="refreshing">预览</button>
        <button class="btn" @click="handleRefresh" :disabled="refreshing">
          <span v-if="refreshing" class="loading-spinner"></span>
          {{ refreshing ? '刷新中...' : '刷新数据' }}
        </button>
        <button class="btn btn-primary" @click="publish">发布</button>
        <button class="btn" @click="exportRelease('markdown')">导出 MD</button>
        <button class="btn" @click="exportRelease('html')">导出 HTML</button>
      </div>
    </div>

    <div v-if="release" class="card mb-16" style="display:flex;gap:24px;flex-wrap:wrap;align-items:center">
      <template v-if="!editingRelease">
        <div><span style="color:var(--text-tertiary)">版本：</span>{{ release.version || '-' }}</div>
        <div><span style="color:var(--text-tertiary)">条目：</span>{{ release.itemCount }}</div>
        <div><span style="color:var(--text-tertiary)">Bug：</span>{{ release.bugCount }}</div>
        <div><span style="color:var(--text-tertiary)">Task：</span>{{ release.taskCount }}</div>
        <div><span style="color:var(--text-tertiary)">备注：</span>{{ release.noteCount }}</div>
        <div><span style="color:var(--text-tertiary)">部署：</span>{{ deployments.length }}</div>
        <div><span style="color:var(--text-tertiary)">发布次数：</span>{{ release.publishCount }}</div>
        <div v-if="release.summary" style="width:100%"><span style="color:var(--text-tertiary)">概述：</span>{{ release.summary }}</div>
        <button class="btn btn-sm" @click="startEditRelease" style="margin-left:auto">编辑</button>
      </template>
      <template v-else>
        <div class="form-group" style="flex:1;margin-bottom:0">
          <label>名称</label>
          <input v-model="editForm.name" />
        </div>
        <div class="form-group" style="flex:0.5;margin-bottom:0">
          <label>版本</label>
          <input v-model="editForm.version" />
        </div>
        <div class="form-group" style="flex:2;margin-bottom:0">
          <label>概述</label>
          <input v-model="editForm.summary" />
        </div>
        <button class="btn btn-primary btn-sm" @click="saveEditRelease">保存</button>
        <button class="btn btn-sm" @click="editingRelease = false">取消</button>
      </template>
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
            <template v-if="editingDeploymentId !== d.id">
              <td>{{ d.moduleName }}</td>
              <td><a :href="d.address" target="_blank" style="color:var(--primary)">{{ d.address }}</a></td>
              <td>{{ d.description || '-' }}</td>
              <td>{{ formatTime(d.createdAt) }}</td>
              <td class="actions">
                <button class="btn btn-sm" @click="startEditDeployment(d)">编辑</button>
                <button class="btn btn-sm btn-danger" @click="removeDeployment(d.id)">删除</button>
              </td>
            </template>
            <template v-else>
              <td><input v-model="editDeployForm.moduleName" style="width:100%" /></td>
              <td><input v-model="editDeployForm.address" style="width:100%" /></td>
              <td><input v-model="editDeployForm.description" style="width:100%" /></td>
              <td></td>
              <td class="actions">
                <button class="btn btn-primary btn-sm" @click="saveEditDeployment">保存</button>
                <button class="btn btn-sm" @click="editingDeploymentId = ''">取消</button>
              </td>
            </template>
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
      <div style="display:flex;gap:8px;margin-bottom:16px;align-items:center">
        <button class="btn" @click="showBugSelector = true">添加 Bug</button>
        <button class="btn" @click="showTaskSelector = true">添加 Task</button>
        <button class="btn" @click="addNote">添加备注</button>
        <div style="margin-left:auto;display:flex;gap:4px">
          <button :class="['btn btn-sm', activeTab === 'all' ? 'btn-primary' : '']" @click="activeTab = 'all'">全部</button>
          <button :class="['btn btn-sm', activeTab === 'bug' ? 'btn-primary' : '']" @click="activeTab = 'bug'">Bug ({{ bugs.length }})</button>
          <button :class="['btn btn-sm', activeTab === 'task' ? 'btn-primary' : '']" @click="activeTab = 'task'">Task ({{ tasks.length }})</button>
          <button :class="['btn btn-sm', activeTab === 'note' ? 'btn-primary' : '']" @click="activeTab = 'note'">备注 ({{ notes.length }})</button>
        </div>
      </div>

      <div v-if="(activeTab === 'all' || activeTab === 'bug') && bugs.length" style="margin-bottom:20px">
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

      <div v-if="(activeTab === 'all' || activeTab === 'task') && tasks.length" style="margin-bottom:20px">
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

      <div v-if="(activeTab === 'all' || activeTab === 'note') && notes.length">
        <h3 style="font-size:15px;margin-bottom:8px;color:var(--primary)">备注（{{ notes.length }}）</h3>
        <div v-for="item in notes" :key="item.id" style="padding:12px;border:1px solid var(--border);border-radius:6px;margin-bottom:8px">
          <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px">
            <strong>{{ item.noteTitle }}</strong>
            <button class="btn btn-sm btn-danger" @click="removeItem(item.id)">删除</button>
          </div>
          <p style="color:var(--text-secondary);white-space:pre-wrap">{{ item.noteContent }}</p>
        </div>
      </div>

      <div v-if="!items.length && !listLoading" class="empty-state">
        <p>暂无条目，点击上方按钮添加</p>
      </div>
      <div v-if="listLoading" style="text-align:center;padding:20px;color:var(--text-tertiary)">
        <span class="loading-spinner"></span> 加载中...
      </div>
    </div>

    <!-- Bug 选择器 -->
    <div v-if="showBugSelector" class="card mb-16">
      <h3 style="margin-bottom:12px">选择 Bug</h3>
      <div style="display:flex;gap:8px;margin-bottom:12px;align-items:center">
        <button class="btn btn-primary btn-sm" @click="loadBugs">加载 Bug</button>
        <button class="btn btn-sm" @click="confirmBugs" :disabled="!selectedBugs.length">确认添加 ({{ selectedBugs.length }})</button>
        <button class="btn btn-sm" @click="showBugSelector = false">关闭</button>
        <span v-if="bugTotal > 50" style="color:var(--text-tertiary);font-size:12px;margin-left:auto">
          共 {{ bugTotal }} 条，仅显示前 50 条
        </span>
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
      <div style="display:flex;gap:8px;margin-bottom:12px;align-items:center">
        <button class="btn btn-primary btn-sm" @click="loadTasks">加载 Task</button>
        <button class="btn btn-sm" @click="confirmTasks" :disabled="!selectedTasks.length">确认添加 ({{ selectedTasks.length }})</button>
        <button class="btn btn-sm" @click="showTaskSelector = false">关闭</button>
        <span v-if="taskTotal > 50" style="color:var(--text-tertiary);font-size:12px;margin-left:auto">
          共 {{ taskTotal }} 条，仅显示前 50 条
        </span>
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
        <div v-if="previewHtml" v-html="previewHtml" class="markdown-body"></div>
      </div>
    </div>

    <!-- 快照查看模态框 -->
    <div v-if="showSnapshotModal" style="position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000" @click.self="showSnapshotModal=false">
      <div class="card" style="width:800px;max-height:85vh;overflow-y:auto">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
          <h3>发布快照</h3>
          <button class="btn btn-sm" @click="showSnapshotModal=false">关闭</button>
        </div>
        <div v-if="snapshotHtml" v-html="snapshotHtml" class="markdown-body"></div>
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
const snapshotHtml = ref('')
const previewHtml = ref('')
const bugList = ref<any[]>([])
const bugLoading = ref(false)
const bugTotal = ref(0)
const taskList = ref<any[]>([])
const taskLoading = ref(false)
const taskTotal = ref(0)
const selectedBugs = ref<any[]>([])
const selectedTasks = ref<any[]>([])
const noteForm = ref({ title: '', content: '' })
const deployForm = ref({ moduleName: '', address: '', description: '' })
const zentaoBaseUrl = ref('')
const loading = ref(false)
const listLoading = ref(false)
const refreshing = ref(false)
const activeTab = ref<'all' | 'bug' | 'task' | 'note'>('all')

const editingRelease = ref(false)
const editForm = ref({ name: '', version: '', summary: '' })
const editingDeploymentId = ref('')
const editDeployForm = ref({ moduleName: '', address: '', description: '' })

function formatTime(t: string) {
  if (!t) return '-'
  return t.replace('T', ' ').substring(0, 19)
}

const bugs = computed(() => items.value.filter(i => i.itemType === 'bug'))
const tasks = computed(() => items.value.filter(i => i.itemType === 'task'))
const notes = computed(() => items.value.filter(i => i.itemType === 'note'))

function renderMarkdown(md: string): string {
  if (!md) return ''
  let html = md
    .replace(/^### (.+)$/gm, '<h3>$1</h3>')
    .replace(/^## (.+)$/gm, '<h2>$1</h2>')
    .replace(/^# (.+)$/gm, '<h1>$1</h1>')
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.+?)\*/g, '<em>$1</em>')
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank">$1</a>')
    .replace(/^---$/gm, '<hr>')

  const lines = html.split('\n')
  let result = ''
  let inTable = false
  let isHeader = false

  for (const line of lines) {
    const trimmed = line.trim()
    if (trimmed.startsWith('<h') || trimmed.startsWith('<hr')) {
      if (inTable) { result += '</tbody></table>'; inTable = false }
      result += trimmed
    } else if (trimmed.startsWith('|') && !trimmed.startsWith('|---')) {
      if (!inTable) {
        result += '<table><thead>'
        inTable = true
        isHeader = true
      }
      const cells = trimmed.split('|').filter(c => c.trim()).map(c => c.trim())
      const tag = isHeader ? 'th' : 'td'
      result += '<tr>' + cells.map(c => `<${tag}>${c}</${tag}>`).join('') + '</tr>'
      if (isHeader) { result += '</thead><tbody>'; isHeader = false }
    } else if (trimmed.startsWith('|---')) {
      continue
    } else {
      if (inTable) { result += '</tbody></table>'; inTable = false }
      if (trimmed && !trimmed.startsWith('<')) {
        result += `<p>${trimmed}</p>`
      } else if (trimmed) {
        result += trimmed
      }
    }
  }
  if (inTable) result += '</tbody></table>'
  return result
}

async function loadZentaoBaseUrl() {
  try {
    const resp: any = await healthApi.get()
    zentaoBaseUrl.value = resp?.zentaoBaseUrl || ''
  } catch {}
}

async function load() {
  const id = route.query.id as string
  if (!id) return
  loading.value = true
  try {
    const rResp: any = await releaseApi.get(id)
    release.value = rResp?.data || null
    await Promise.all([loadItems(), loadSnapshots(), loadDeployments()])
  } catch {} finally {
    loading.value = false
  }
}

async function loadItems() {
  const id = route.query.id as string
  listLoading.value = true
  try {
    const iResp: any = await itemApi.list(id)
    items.value = iResp?.list || []
  } catch {} finally {
    listLoading.value = false
  }
}

async function loadSnapshots() {
  const id = route.query.id as string
  try {
    const sResp: any = await snapshotApi.list(id)
    snapshots.value = sResp?.list || []
  } catch {}
}

async function loadDeployments() {
  const id = route.query.id as string
  try {
    const dResp: any = await deploymentApi.list(id)
    deployments.value = dResp?.list || []
  } catch {}
}

async function handleRefresh() {
  refreshing.value = true
  try {
    await itemApi.refresh(route.query.id as string)
    await loadItems()
  } finally {
    refreshing.value = false
  }
}

async function removeItem(id: string) {
  if (!confirm('确定删除？')) return
  try {
    await itemApi.delete(id)
    await loadItems()
  } catch {}
}

async function removeDeployment(id: string) {
  if (!confirm('确定删除？')) return
  try {
    await deploymentApi.delete(id)
    await loadDeployments()
  } catch {}
}

function addNote() {
  noteForm.value = { title: '', content: '' }
  showNoteForm.value = true
}

async function submitNote() {
  try {
    await itemApi.add({
      releaseId: route.query.id as string,
      itemType: 'note',
      noteTitle: noteForm.value.title,
      noteContent: noteForm.value.content
    })
    showNoteForm.value = false
    await loadItems()
  } catch {}
}

async function submitDeployment() {
  if (!deployForm.value.moduleName || !deployForm.value.address) {
    alert('请填写功能模块和地址')
    return
  }
  try {
    await deploymentApi.add({
      releaseId: route.query.id as string,
      moduleName: deployForm.value.moduleName,
      address: deployForm.value.address,
      description: deployForm.value.description
    })
    deployForm.value = { moduleName: '', address: '', description: '' }
    showDeploymentForm.value = false
    await loadDeployments()
  } catch {}
}

function startEditDeployment(d: Deployment) {
  editingDeploymentId.value = d.id
  editDeployForm.value = { moduleName: d.moduleName, address: d.address, description: d.description }
}

async function saveEditDeployment() {
  try {
    await deploymentApi.update({ id: editingDeploymentId.value, ...editDeployForm.value })
    editingDeploymentId.value = ''
    await loadDeployments()
  } catch {}
}

function startEditRelease() {
  if (!release.value) return
  editingRelease.value = true
  editForm.value = { name: release.value.name, version: release.value.version, summary: release.value.summary }
}

async function saveEditRelease() {
  try {
    await releaseApi.update({ id: route.query.id as string, ...editForm.value })
    editingRelease.value = false
    await load()
  } catch {}
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
    bugTotal.value = resp?.total || bugList.value.length
  } catch {} finally {
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
  try {
    await itemApi.batchAdd(releaseId, its)
    selectedBugs.value = []
    showBugSelector.value = false
    await loadItems()
  } catch {}
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
    taskTotal.value = resp?.total || taskList.value.length
  } catch {} finally {
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
  try {
    await itemApi.batchAdd(releaseId, its)
    selectedTasks.value = []
    showTaskSelector.value = false
    await loadItems()
  } catch {}
}

async function publish() {
  const version = prompt('输入版本号（可选）', release.value?.version || '')
  if (version === null) return
  try {
    await releaseApi.publish(route.query.id as string, version || undefined)
    await load()
  } catch {}
}

async function preview() {
  try {
    const resp: any = await releaseApi.export(route.query.id as string, 'html')
    previewHtml.value = resp?.content || ''
    showPreviewModal.value = true
  } catch {}
}

async function viewSnapshot(s: ReleaseSnapshot) {
  try {
    const resp: any = await snapshotApi.get(s.id)
    const content = resp?.data?.content || ''
    snapshotHtml.value = renderMarkdown(content)
    showSnapshotModal.value = true
  } catch {}
}

async function exportRelease(format: string, snapshotId?: string) {
  try {
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
  } catch {}
}

onMounted(() => {
  loadZentaoBaseUrl()
  load()
})
</script>

<style scoped>
.markdown-body :deep(h1) { font-size: 20px; font-weight: 700; margin: 16px 0 8px; border-bottom: 2px solid var(--primary); padding-bottom: 6px; }
.markdown-body :deep(h2) { font-size: 17px; font-weight: 600; margin: 16px 0 8px; color: var(--text-primary); }
.markdown-body :deep(h3) { font-size: 15px; font-weight: 600; margin: 12px 0 6px; }
.markdown-body :deep(table) { width: 100%; border-collapse: collapse; margin: 12px 0; font-size: 13px; }
.markdown-body :deep(th), .markdown-body :deep(td) { border: 1px solid var(--border); padding: 8px 12px; text-align: left; }
.markdown-body :deep(th) { background: var(--bg); font-weight: 600; }
.markdown-body :deep(a) { color: var(--primary); text-decoration: none; }
.markdown-body :deep(a:hover) { text-decoration: underline; }
.markdown-body :deep(hr) { border: none; border-top: 1px solid var(--border); margin: 16px 0; }
.markdown-body :deep(p) { margin: 8px 0; line-height: 1.6; }
</style>
