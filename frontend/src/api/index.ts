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
  parentBranch: string
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

export interface ProjectRepo {
  id: string
  projectId: string
  gitlabProjectId: number
  repoUrl: string
  repoName: string
  defaultBranch: string
  createdAt: string
}

export interface ReleaseBranch {
  id: string
  releaseId: string
  repoId: string
  branchName: string
  branchType: string
  parentBranch: string
  gitlabBranchUrl: string
  description: string
  createdAt: string
}

export interface DockerImage {
  id: string
  releaseId: string
  repoId: string
  imageUrl: string
  imageDigest: string
  ciPipelineId: number
  ciPipelineUrl: string
  commitSha: string
  commitMessage: string
  source: string
  tested: boolean
  createdAt: string
}

export interface DockerImagePoolItem {
  id: string
  imageUrl: string
  imageDigest: string
  commitSha: string
  commitMessage: string
  ciPipelineId: number
  ciPipelineUrl: string
  createdAt: string
}

export interface GitlabProject {
  id: number
  name: string
  nameWithNamespace: string
  pathWithNamespace: string
  webUrl: string
  httpUrlToRepo: string
  defaultBranch: string
}

export interface GitlabBranch {
  name: string
  isDefault: boolean
  isProtected: boolean
  webUrl: string
}

export interface ReleaseFeature {
  id: string
  releaseId: string
  title: string
  content: string
  sortOrder: number
  createdAt: string
  updatedAt: string
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
  create: (data: { projectId: string; name: string; version?: string; summary?: string; parentBranch?: string; repoId: string }) =>
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

export const repoApi = {
  list: (projectId: string) => api.get('/projects/repos', { params: { projectId } }),
  add: (data: { projectId: string; gitlabProjectId: number; repoUrl: string; repoName: string; defaultBranch?: string }) =>
    api.post('/projects/repos', data),
  delete: (id: string) => api.post('/projects/repos/delete', { id }),
}

export const branchApi = {
  list: (releaseId: string) => api.get('/release-branches', { params: { releaseId } }),
  createFeature: (data: { releaseId: string; repoId: string; branchName: string; parentBranch?: string }) =>
    api.post('/release-branches/feature', data),
  update: (data: { id: string; description?: string }) => api.post('/release-branches/update', data),
  delete: (id: string) => api.post('/release-branches/delete', { id }),
}

export const dockerImageApi = {
  list: (releaseId: string) => api.get('/docker-images', { params: { releaseId } }),
  pool: (gitlabProjectId?: number) => api.get('/docker-images/pool', { params: { gitlabProjectId } }),
  add: (data: { releaseId: string; repoId?: string; imageUrl: string; imageDigest?: string; commitSha?: string; commitMessage?: string }) =>
    api.post('/docker-images', data),
  update: (data: { id: string; imageUrl?: string; tested?: boolean; commitSha?: string; commitMessage?: string }) =>
    api.post('/docker-images/update', data),
  delete: (id: string) => api.post('/docker-images/delete', { id }),
}

export const gitlabApi = {
  search: (query: string) => api.get('/gitlab/search', { params: { query } }),
  branches: (gitlabProjectId: number) => api.get('/gitlab/branches', { params: { gitlabProjectId } }),
}

export const healthApi = {
  get: () => api.get('/health'),
}

export const featureApi = {
  list: (releaseId: string) => api.get('/features', { params: { releaseId } }),
  add: (data: { releaseId: string; title: string; content: string }) =>
    api.post('/features', data),
  update: (data: Partial<ReleaseFeature> & { id: string }) => api.post('/features/update', data),
  delete: (id: string) => api.post('/features/delete', { id }),
}

export const notifyApi = {
  preview: (data: { releaseId: string; version?: string }) => api.post('/notify/preview', data),
  send: (data: { releaseId: string; version?: string; channel?: string }) => api.post('/notify/send', data),
}
