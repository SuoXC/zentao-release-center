import api from './api'

export interface BaseResp {
  code: number
  message: string
}

export interface Project {
  id: string
  name: string
  description: string
  zentaoProductId: number
  zentaoProjectId: number
  zentaoProductName: string
  zentaoProjectName: string
  zentaoServer: string
  status: string
  createdAt: string
  updatedAt: string
}

export interface Release {
  id: string
  projectId: string
  projectName: string
  name: string
  version: string
  status: string
  summary: string
  publishCount: number
  firstPublishedAt: string
  lastPublishedAt: string
  itemCount: number
  bugCount: number
  taskCount: number
  noteCount: number
  createdAt: string
  updatedAt: string
}

export interface ReleaseItem {
  id: string
  releaseId: string
  itemType: 'bug' | 'task' | 'note'
  sortOrder: number
  zentaoId?: number
  zentaoType?: string
  title?: string
  severity?: string
  priority?: string
  status?: string
  assignedTo?: string
  resolvedBy?: string
  zentaoUrl?: string
  steps?: string
  noteTitle?: string
  noteContent?: string
  createdAt: string
  updatedAt: string
}

export interface ReleaseSnapshot {
  id: string
  releaseId: string
  version: string
  content: string
  itemCount: number
  bugCount: number
  taskCount: number
  noteCount: number
  publishedAt: string
}

export interface Deployment {
  id: string
  releaseId: string
  moduleName: string
  address: string
  description: string
  sortOrder: number
  createdAt: string
  updatedAt: string
}

interface ApiResp {
  base: BaseResp
  [key: string]: unknown
}

export const projectApi = {
  list: (params?: Record<string, unknown>) => api.get('/projects', { params }),
  get: (id: string) => api.get('/projects/detail', { params: { id } }),
  create: (data: Partial<Project>) => api.post('/projects', data),
  update: (data: Partial<Project> & { id: string }) => api.post('/projects/update', data),
  delete: (id: string) => api.post('/projects/delete', { id }),
}

export const releaseApi = {
  list: (params: { projectId: string; status?: string; page?: number; pageSize?: number }) =>
    api.get('/releases', { params }),
  get: (id: string) => api.get('/releases/detail', { params: { id } }),
  create: (data: Partial<Release> & { projectId: string; name: string }) =>
    api.post('/releases', data),
  update: (data: Partial<Release> & { id: string }) => api.post('/releases/update', data),
  delete: (id: string) => api.post('/releases/delete', { id }),
  publish: (releaseId: string, version?: string) => api.post('/releases/publish', { releaseId, version }),
  export: (releaseId: string, format: string = 'markdown', snapshotId?: string) =>
    api.get('/releases/export', { params: { releaseId, format, snapshotId } }),
}

export const itemApi = {
  list: (releaseId: string) => api.get('/release-items', { params: { releaseId } }),
  add: (data: Partial<ReleaseItem> & { releaseId: string; itemType: string }) =>
    api.post('/release-items', data),
  batchAdd: (releaseId: string, items: Partial<ReleaseItem>[]) =>
    api.post('/release-items/batch', { releaseId, items }),
  update: (data: Partial<ReleaseItem> & { id: string }) => api.post('/release-items/update', data),
  delete: (id: string) => api.post('/release-items/delete', { id }),
  reorder: (releaseId: string, items: { id: string; sortOrder: number }[]) =>
    api.post('/release-items/reorder', { releaseId, items }),
  refresh: (releaseId: string) => api.post('/release-items/refresh', { releaseId }),
}

export const snapshotApi = {
  list: (releaseId: string) => api.get('/release-snapshots', { params: { releaseId } }),
  get: (id: string) => api.get('/release-snapshots/detail', { params: { id } }),
}

export const zentaoApi = {
  products: () => api.get('/zentao/products'),
  projects: (productId?: number) => api.get('/zentao/projects', { params: { productId } }),
  executions: (projectId?: number) => api.get('/zentao/executions', { params: { projectId } }),
  bugs: (params: Record<string, unknown>) => api.get('/zentao/bugs', { params }),
  tasks: (params: Record<string, unknown>) => api.get('/zentao/tasks', { params }),
}

export const deploymentApi = {
  list: (releaseId: string) => api.get('/deployments', { params: { releaseId } }),
  add: (data: { releaseId: string; moduleName: string; address: string; description?: string }) =>
    api.post('/deployments', data),
  update: (data: Partial<Deployment> & { id: string }) => api.post('/deployments/update', data),
  delete: (id: string) => api.post('/deployments/delete', { id }),
}

export const healthApi = {
  get: () => api.get('/health'),
}
