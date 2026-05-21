<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">项目管理</h1>
      <button class="btn btn-primary" @click="openCreate()">新建项目</button>
    </div>

    <div v-if="showCreate" class="card mb-16">
      <h3 style="margin-bottom:16px">{{ editId ? '编辑项目' : '新建项目' }}</h3>
      <div class="form-group">
        <label>项目名称 *</label>
        <input v-model="form.name" placeholder="输入项目名称" />
      </div>
      <div class="form-group">
        <label>描述</label>
        <textarea v-model="form.description" placeholder="项目描述（可选）" />
      </div>
      <div class="form-group">
        <label>禅道产品</label>
        <div class="select-with-search">
          <input v-model="productSearch" placeholder="搜索产品..." class="search-input" />
          <select v-model="form.zentaoProductId" @change="onProductChange" size="6" class="search-select">
            <option :value="0">不关联</option>
            <option v-for="p in filteredProducts" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
          <p v-if="loadingProducts" style="color:var(--text-muted);font-size:12px;margin-top:4px">加载中...</p>
          <p v-else-if="zentaoProducts.length > 0" style="color:var(--text-muted);font-size:12px;margin-top:4px">共 {{ zentaoProducts.length }} 个产品，已选择: {{ selectedProductName }}</p>
        </div>
      </div>
      <div class="form-group">
        <label>禅道项目</label>
        <div class="select-with-search" v-if="form.zentaoProductId">
          <input v-model="projectSearch" placeholder="搜索项目..." class="search-input" />
          <select v-model="form.zentaoProjectId" size="4" class="search-select">
            <option :value="0">不关联</option>
            <option v-for="p in filteredProjects" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
        </div>
        <p v-else style="color:var(--text-muted);font-size:13px">请先选择禅道产品</p>
      </div>
      <div class="actions">
        <button class="btn btn-primary" @click="submitForm" :disabled="submitting">{{ submitting ? '保存中...' : '保存' }}</button>
        <button class="btn" @click="cancelForm">取消</button>
      </div>
    </div>

    <div class="card">
      <table class="data-table" v-if="projects.length">
        <thead>
          <tr>
            <th>项目名称</th>
            <th>禅道产品</th>
            <th>禅道项目</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in projects" :key="p.id">
            <td><router-link :to="`/project?id=${p.id}`" style="color:var(--primary);text-decoration:none;font-weight:500">{{ p.name }}</router-link></td>
            <td>{{ p.zentaoProductName || '-' }}</td>
            <td>{{ p.zentaoProjectName || '-' }}</td>
            <td><span :class="['badge', p.status === 'active' ? 'badge-success' : 'badge-warning']">{{ p.status === 'active' ? '活跃' : '已归档' }}</span></td>
            <td>{{ formatTime(p.createdAt) }}</td>
            <td class="actions">
              <button class="btn btn-sm" @click="editProject(p)">编辑</button>
              <button class="btn btn-sm btn-danger" @click="deleteProject(p.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="empty-state">
        <h3>暂无项目</h3>
        <p>点击上方「新建项目」开始</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { projectApi, zentaoApi } from '../api'
import type { Project } from '../api'

const projects = ref<Project[]>([])
const showCreate = ref(false)
const editId = ref('')
const submitting = ref(false)
const form = ref({ name: '', description: '', zentaoProductId: 0, zentaoProjectId: 0, zentaoProductName: '', zentaoProjectName: '' })
const zentaoProducts = ref<{ id: number; name: string }[]>([])
const zentaoProjects = ref<{ id: number; name: string }[]>([])
const loadingProducts = ref(false)
const productSearch = ref('')
const projectSearch = ref('')

const filteredProducts = computed(() => {
  if (!productSearch.value) return zentaoProducts.value
  const q = productSearch.value.toLowerCase()
  return zentaoProducts.value.filter(p => p.name.toLowerCase().includes(q) || String(p.id).includes(q))
})

const filteredProjects = computed(() => {
  if (!projectSearch.value) return zentaoProjects.value
  const q = projectSearch.value.toLowerCase()
  return zentaoProjects.value.filter(p => p.name.toLowerCase().includes(q) || String(p.id).includes(q))
})

const selectedProductName = computed(() => {
  const p = zentaoProducts.value.find(p => p.id === form.value.zentaoProductId)
  return p ? p.name : '不关联'
})

function formatTime(t: string) {
  if (!t) return '-'
  return t.replace('T', ' ').substring(0, 19)
}

async function loadProducts() {
  loadingProducts.value = true
  try {
    const resp: any = await zentaoApi.products()
    const data = JSON.parse(resp?.data || '[]')
    zentaoProducts.value = Array.isArray(data) ? data.map((p: any) => ({ id: p.id, name: p.name })) : []
  } catch (e) {
    console.error('Failed to load products:', e)
    zentaoProducts.value = []
  }
  loadingProducts.value = false
}

async function loadZentaoProjects(productId: number) {
  if (!productId) {
    zentaoProjects.value = []
    return
  }
  try {
    const resp: any = await zentaoApi.projects(productId)
    const data = JSON.parse(resp?.data || '[]')
    zentaoProjects.value = Array.isArray(data) ? data.map((p: any) => ({ id: p.id, name: p.name })) : []
  } catch (e) {
    console.error('Failed to load projects:', e)
    zentaoProjects.value = []
  }
}

async function onProductChange() {
  form.value.zentaoProjectId = 0
  form.value.zentaoProjectName = ''
  const selected = zentaoProducts.value.find(p => p.id === form.value.zentaoProductId)
  form.value.zentaoProductName = selected ? selected.name : ''
  projectSearch.value = ''
  await loadZentaoProjects(form.value.zentaoProductId)
}

async function openCreate() {
  resetForm()
  showCreate.value = true
  await loadProducts()
}

async function load() {
  try {
    const resp: any = await projectApi.list()
    projects.value = resp?.list || []
  } catch (e) {
    console.error('Failed to load projects:', e)
  }
}

function resetForm() {
  form.value = { name: '', description: '', zentaoProductId: 0, zentaoProjectId: 0, zentaoProductName: '', zentaoProjectName: '' }
  editId.value = ''
  showCreate.value = false
  zentaoProducts.value = []
  zentaoProjects.value = []
  productSearch.value = ''
  projectSearch.value = ''
  submitting.value = false
}

function cancelForm() { resetForm() }

async function editProject(p: Project) {
  form.value = {
    name: p.name,
    description: p.description,
    zentaoProductId: p.zentaoProductId,
    zentaoProjectId: p.zentaoProjectId,
    zentaoProductName: p.zentaoProductName || '',
    zentaoProjectName: p.zentaoProjectName || '',
  }
  editId.value = p.id
  showCreate.value = true
  productSearch.value = ''
  projectSearch.value = ''
  await loadProducts()
  if (p.zentaoProductId) {
    await loadZentaoProjects(p.zentaoProductId)
  }
}

async function submitForm() {
  if (!form.value.name) return alert('请输入项目名称')
  submitting.value = true
  try {
    const selectedProject = zentaoProjects.value.find(p => p.id === form.value.zentaoProjectId)
    if (selectedProject) {
      form.value.zentaoProjectName = selectedProject.name
    }
    if (editId.value) {
      await projectApi.update({ id: editId.value, ...form.value })
    } else {
      await projectApi.create(form.value)
    }
    resetForm()
    load()
  } catch (e) {
    console.error('Failed to save:', e)
    alert('保存失败')
  } finally {
    submitting.value = false
  }
}

async function deleteProject(id: string) {
  if (!confirm('确定删除？')) return
  await projectApi.delete(id)
  load()
}

onMounted(load)
</script>

<style scoped>
.select-with-search {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.search-input {
  padding: 6px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  outline: none;
}
.search-input:focus {
  border-color: var(--primary);
}
.search-select {
  width: 100%;
  padding: 4px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg);
  font-size: 13px;
  cursor: pointer;
}
</style>
