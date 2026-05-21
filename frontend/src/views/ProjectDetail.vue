<template>
  <div>
    <div class="page-header">
      <div style="display:flex;align-items:center;gap:12px">
        <router-link to="/projects" class="btn btn-sm">← 返回</router-link>
        <h1 class="page-title">{{ project?.name || '项目详情' }}</h1>
        <span v-if="project" :class="['badge', project.status === 'active' ? 'badge-success' : 'badge-warning']">{{ project.status }}</span>
      </div>
      <button class="btn btn-primary" @click="showCreate = true">新建发布单</button>
    </div>

    <div v-if="showCreate" class="card mb-16">
      <h3 style="margin-bottom:16px">新建发布单</h3>
      <div class="form-group">
        <label>发布单名称 *</label>
        <input v-model="form.name" placeholder="例如：v2.1.0 系统优化版本" />
      </div>
      <div class="form-group">
        <label>版本号</label>
        <input v-model="form.version" placeholder="例如：2.1.0（可选）" />
      </div>
      <div class="form-group">
        <label>概述</label>
        <textarea v-model="form.summary" placeholder="发布概述（可选）" />
      </div>
      <div class="actions">
        <button class="btn btn-primary" @click="createRelease">创建</button>
        <button class="btn" @click="showCreate = false">取消</button>
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
            <td>{{ r.updatedAt }}</td>
            <td class="actions">
              <button class="btn btn-sm btn-danger" @click="deleteRelease(r.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="empty-state">
        <h3>暂无发布单</h3>
        <p>点击「新建发布单」开始</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { projectApi, releaseApi } from '../api'
import type { Project, Release } from '../api'

const route = useRoute()
const project = ref<Project | null>(null)
const releases = ref<Release[]>([])
const showCreate = ref(false)
const form = ref({ name: '', version: '', summary: '' })

async function load() {
  const id = route.query.id as string
  if (!id) return
  const pResp: any = await projectApi.get(id)
  project.value = pResp?.data || null
  const rResp: any = await releaseApi.list({ projectId: id })
  releases.value = rResp?.list || []
}

async function createRelease() {
  if (!form.value.name) return alert('请输入名称')
  await releaseApi.create({ projectId: route.query.id as string, ...form.value })
  form.value = { name: '', version: '', summary: '' }
  showCreate.value = false
  load()
}

async function deleteRelease(id: string) {
  if (!confirm('确定删除？')) return
  await releaseApi.delete(id)
  load()
}

onMounted(load)
</script>
