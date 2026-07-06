<template>
  <div>
    <div class="page-header">
      <div style="display:flex;align-items:center;gap:12px">
        <router-link to="/projects" class="btn btn-sm">← 返回</router-link>
        <h1 class="page-title">{{ project?.name || '项目详情' }}</h1>
        <span v-if="project" :class="['badge', project.status === 'active' ? 'badge-success' : 'badge-warning']">{{ project.status }}</span>
        <span v-if="pageLoading" class="loading-spinner"></span>
      </div>
      <button class="btn btn-primary" @click="showCreate = true">新建发布单</button>
    </div>

    <div v-if="showCreate" class="card mb-16">
      <h3 style="margin-bottom:16px">新建发布单</h3>
      <div v-if="!repos.length" class="form-group" style="background:#fff7e6;border:1px solid #ffd591;border-radius:4px;padding:12px;color:#874d00">
        请先在下方“关联 GitLab 仓库”为本项目添加至少一个仓库，然后才能新建发布单（创建时将自动在所选仓库创建发布分支，二者一一强绑定）。
      </div>
      <div class="form-group">
        <label>发布单名称 *</label>
        <input v-model="form.name" placeholder="例如：v2.1.0 系统优化版本" />
      </div>
      <div class="form-group">
        <label>所属 GitLab 仓库 *</label>
        <select v-model="form.repoId" :disabled="!repos.length">
          <option value="" disabled>请选择仓库</option>
          <option v-for="r in repos" :key="r.id" :value="r.id">{{ r.repoName }}</option>
        </select>
        <div style="font-size:12px;color:var(--text-tertiary);margin-top:4px">将在该仓库创建发布分支（强绑定，提测期间唯一）</div>
      </div>
      <div class="form-group">
        <label>版本号</label>
        <input v-model="form.version" placeholder="例如：2.1.0（可选；分支名会取 release/<version>）" />
      </div>
      <div class="form-group">
        <label>基准分支</label>
        <input v-model="form.parentBranch" placeholder="例如：main、develop（默认用所选仓库的默认分支）" />
      </div>
      <div class="form-group">
        <label>概述</label>
        <textarea v-model="form.summary" placeholder="发布概述（可选）" />
      </div>
      <div class="actions">
        <button class="btn btn-primary" :disabled="!repos.length || !form.repoId" @click="createRelease">创建</button>
        <button class="btn" @click="showCreate = false">取消</button>
      </div>
    </div>

    <div class="card mb-16">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
        <h3>关联 GitLab 仓库</h3>
        <button class="btn btn-sm btn-primary" @click="openAddRepo">添加仓库</button>
      </div>

      <div v-if="showAddRepo" class="card mb-16" style="background:#f8f9fa">
        <h4 style="margin-bottom:12px">搜索 GitLab 仓库</h4>
        <div style="display:flex;gap:8px;margin-bottom:12px">
          <input v-model="repoSearchQuery" placeholder="输入仓库名称搜索" style="flex:1" @keyup.enter="searchGitlabRepos" />
          <button class="btn btn-sm btn-primary" @click="searchGitlabRepos">搜索</button>
        </div>
        <div v-if="gitlabSearchResults.length" style="max-height:200px;overflow-y:auto">
          <div v-for="p in gitlabSearchResults" :key="p.id" style="display:flex;justify-content:space-between;align-items:center;padding:8px;border-bottom:1px solid #eee">
            <div>
              <div style="font-weight:500">{{ p.nameWithNamespace }}</div>
              <div style="font-size:12px;color:#666">{{ p.pathWithNamespace }}</div>
            </div>
            <button class="btn btn-sm btn-primary" @click="selectRepo(p)">选择</button>
          </div>
        </div>
        <div class="actions" style="margin-top:12px">
          <button class="btn btn-sm" @click="cancelAddRepo">取消</button>
        </div>
      </div>

      <div v-if="showRepoConfig" class="card mb-16" style="background:#f8f9fa">
        <h4 style="margin-bottom:12px">配置仓库</h4>
        <div style="margin-bottom:12px">
          <div style="font-weight:500">{{ selectedRepo?.nameWithNamespace }}</div>
          <div style="font-size:12px;color:#666">{{ selectedRepo?.httpUrlToRepo }}</div>
        </div>
        <div class="form-group">
          <label>默认分支 *</label>
          <select v-model="repoConfig.defaultBranch" style="width:100%">
            <option value="" disabled>请选择默认分支</option>
            <option v-for="b in repoBranches" :key="b.name" :value="b.name">
              {{ b.name }}{{ b.isDefault ? ' (默认)' : '' }}{{ b.isProtected ? ' 🔒' : '' }}
            </option>
          </select>
          <div v-if="loadingBranches" style="font-size:12px;color:var(--text-tertiary);margin-top:4px">
            <span class="loading-spinner"></span> 加载分支中...
          </div>
        </div>
        <div class="actions">
          <button class="btn btn-primary" @click="confirmAddRepo" :disabled="!repoConfig.defaultBranch">确认添加</button>
          <button class="btn" @click="cancelAddRepo">取消</button>
        </div>
      </div>

      <table class="data-table" v-if="repos.length">
        <thead>
          <tr>
            <th>仓库名称</th>
            <th>GitLab 项目 ID</th>
            <th>默认分支</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in repos" :key="r.id">
            <td><a :href="r.repoUrl" target="_blank" style="color:var(--primary)">{{ r.repoName }}</a></td>
            <td>{{ r.gitlabProjectId }}</td>
            <td>{{ r.defaultBranch }}</td>
            <td>
              <button class="btn btn-sm btn-danger" @click="deleteRepo(r.id)">移除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="empty-state">
        <p>暂未关联仓库</p>
      </div>
    </div>

    <div class="card">
      <table class="data-table" v-if="releases.length">
        <thead>
          <tr>
            <th>发布单名称</th>
            <th>版本</th>
            <th>状态</th>
            <th>条目数</th>
            <th>Bug/Task/备注</th>
            <th>发布次数</th>
            <th>更新时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in releases" :key="r.id">
            <td><router-link :to="`/release?id=${r.id}`" style="color:var(--primary);text-decoration:none;font-weight:500">{{ r.name }}</router-link></td>
            <td>{{ r.version || '-' }}</td>
            <td>
              <span :class="['badge', r.status === 'published' ? 'badge-success' : r.status === 'draft' ? 'badge-info' : 'badge-warning']">
                {{ r.status === 'draft' ? '草稿' : r.status === 'published' ? '已发布' : '已归档' }}
              </span>
            </td>
            <td>{{ r.itemCount }}</td>
            <td>{{ r.bugCount }} / {{ r.taskCount }} / {{ r.noteCount }}</td>
            <td>{{ r.publishCount }}</td>
            <td>{{ formatTime(r.updatedAt) }}</td>
            <td class="actions">
              <button class="btn btn-sm btn-danger" @click="deleteRelease(r.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else-if="!pageLoading" class="empty-state">
        <h3>暂无发布单</h3>
        <p>点击「新建发布单」开始</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { projectApi, releaseApi, repoApi, gitlabApi } from '../api'
import type { Project, Release, ProjectRepo, GitlabProject, GitlabBranch } from '../api'

const route = useRoute()
const project = ref<Project | null>(null)
const releases = ref<Release[]>([])
const repos = ref<ProjectRepo[]>([])
const showCreate = ref(false)
const showAddRepo = ref(false)
const showRepoConfig = ref(false)
const form = ref({ name: '', version: '', summary: '', parentBranch: '', repoId: '' })
const pageLoading = ref(false)
const repoSearchQuery = ref('')
const gitlabSearchResults = ref<GitlabProject[]>([])

const selectedRepo = ref<GitlabProject | null>(null)
const repoBranches = ref<GitlabBranch[]>([])
const loadingBranches = ref(false)
const repoConfig = ref({ defaultBranch: '' })

function formatTime(t: string) {
  if (!t) return '-'
  return t.replace('T', ' ').substring(0, 19)
}

async function load() {
  const id = route.query.id as string
  if (!id) return
  pageLoading.value = true
  try {
    const pResp: any = await projectApi.get(id)
    project.value = pResp?.data || null
    const rResp: any = await releaseApi.list({ projectId: id })
    releases.value = rResp?.list || []
    const reposResp: any = await repoApi.list(id)
    repos.value = reposResp?.list || []
  } catch {} finally {
    pageLoading.value = false
  }
}

function openAddRepo() {
  showAddRepo.value = true
  showRepoConfig.value = false
  selectedRepo.value = null
  repoBranches.value = []
  repoConfig.value = { defaultBranch: '' }
}

function cancelAddRepo() {
  showAddRepo.value = false
  showRepoConfig.value = false
  selectedRepo.value = null
  repoBranches.value = []
  repoConfig.value = { defaultBranch: '' }
  gitlabSearchResults.value = []
  repoSearchQuery.value = ''
}

async function searchGitlabRepos() {
  if (!repoSearchQuery.value) return
  try {
    const resp: any = await gitlabApi.search(repoSearchQuery.value)
    gitlabSearchResults.value = resp?.list || []
  } catch {}
}

async function selectRepo(p: GitlabProject) {
  selectedRepo.value = p
  showAddRepo.value = false
  showRepoConfig.value = true
  repoConfig.value = { defaultBranch: p.defaultBranch || '' }
  loadingBranches.value = true
  try {
    const resp: any = await gitlabApi.branches(p.id)
    repoBranches.value = resp?.list || []
  } catch {} finally {
    loadingBranches.value = false
  }
}

async function confirmAddRepo() {
  if (!selectedRepo.value || !repoConfig.value.defaultBranch) return
  const id = route.query.id as string
  try {
    await repoApi.add({
      projectId: id,
      gitlabProjectId: selectedRepo.value.id,
      repoUrl: selectedRepo.value.httpUrlToRepo,
      repoName: selectedRepo.value.nameWithNamespace,
      defaultBranch: repoConfig.value.defaultBranch,
    })
    cancelAddRepo()
    await load()
  } catch {}
}

async function deleteRepo(repoId: string) {
  if (!confirm('确定移除此仓库关联？')) return
  try {
    await repoApi.delete(repoId)
    await load()
  } catch {}
}

async function createRelease() {
  if (!form.value.name) return alert('请输入名称')
  if (!form.value.repoId) return alert('请选择 GitLab 仓库')
  try {
    await releaseApi.create({
      projectId: route.query.id as string,
      name: form.value.name,
      version: form.value.version || undefined,
      summary: form.value.summary || undefined,
      parentBranch: form.value.parentBranch || undefined,
      repoId: form.value.repoId,
    })
    form.value = { name: '', version: '', summary: '', parentBranch: '', repoId: '' }
    showCreate.value = false
    await load()
  } catch {}
}

async function deleteRelease(id: string) {
  if (!confirm('确定删除？')) return
  try {
    await releaseApi.delete(id)
    await load()
  } catch {}
}

onMounted(load)
</script>
