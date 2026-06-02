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
        <button class="btn btn-warning" @click="archiveRelease" v-if="release.status === 'published'">发布完成</button>
        <button class="btn" @click="sendLanxin">发送蓝信</button>
        <button class="btn" @click="sendEmail">发送邮件</button>
        <button class="btn" @click="exportRelease('markdown')">导出 MD</button>
        <button class="btn" @click="exportRelease('html')">导出 HTML</button>
      </div>
    </div>

    <div v-if="release" class="card mb-16" style="display:flex;gap:24px;flex-wrap:wrap;align-items:center">
      <template v-if="!editingRelease">
        <div><span style="color:var(--text-tertiary)">版本：</span>{{ release.version || '-' }}</div>
        <div><span style="color:var(--text-tertiary)">基准分支：</span><span style="font-weight:500;color:var(--primary)">{{ release.parentBranch || '未设置' }}</span></div>
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
        <div class="form-group" style="flex:1;margin-bottom:0">
          <label>基准分支</label>
          <input v-model="editForm.parentBranch" placeholder="例如：main、develop" />
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

    <!-- 功能说明 -->
    <div class="card mb-16">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
        <h3 style="font-size:15px">功能说明</h3>
        <button class="btn btn-sm" @click="showFeatureForm = true; featureForm = { title: '', content: '' }">添加功能说明</button>
      </div>

      <div v-if="features.length">
        <div v-for="(f, idx) in features" :key="f.id" style="padding:16px;border:1px solid var(--border);border-radius:8px;margin-bottom:12px">
          <template v-if="editingFeatureId !== f.id">
            <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:8px">
              <h4 style="font-size:14px;margin:0">{{ idx + 1 }}. {{ f.title }}</h4>
              <div style="display:flex;gap:6px">
                <button class="btn btn-sm" @click="startEditFeature(f)">编辑</button>
                <button class="btn btn-sm btn-danger" @click="removeFeature(f.id)">删除</button>
              </div>
            </div>
            <div v-if="f.content" style="color:var(--text-secondary);line-height:1.8" v-html="f.content"></div>
          </template>
          <template v-else>
            <div class="form-group" style="margin-bottom:8px">
              <label>标题</label>
              <input v-model="editFeatureForm.title" />
            </div>
            <div class="form-group" style="margin-bottom:12px">
              <label>内容（支持 HTML）</label>
              <textarea v-model="editFeatureForm.content" rows="6" style="min-height:120px" />
            </div>
            <div style="display:flex;gap:8px">
              <button class="btn btn-primary btn-sm" @click="saveEditFeature">保存</button>
              <button class="btn btn-sm" @click="editingFeatureId = ''">取消</button>
            </div>
          </template>
        </div>
      </div>
      <div v-else style="color:var(--text-tertiary);padding:12px 0">暂无功能说明，点击「添加功能说明」</div>

      <div v-if="showFeatureForm" style="margin-top:12px;padding-top:12px;border-top:1px solid var(--border)">
        <div class="form-group">
          <label>标题 *</label>
          <input v-model="featureForm.title" placeholder="功能名称" />
        </div>
        <div class="form-group">
          <label>内容（支持 HTML 富文本）</label>
          <textarea v-model="featureForm.content" rows="8" placeholder="描述此功能的详细说明，支持 HTML 格式" style="min-height:160px" />
        </div>
        <div style="display:flex;gap:8px">
          <button class="btn btn-primary" @click="submitFeature">添加</button>
          <button class="btn" @click="showFeatureForm = false">取消</button>
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

    <!-- 分支管理 -->
    <div class="card mb-16">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
        <h3 style="font-size:15px">分支管理</h3>
        <div style="display:flex;gap:8px">
          <button class="btn btn-sm" @click="showCreateReleaseBranch = true">创建发布分支</button>
          <button class="btn btn-sm" @click="showCreateFeatureBranch = true">创建功能分支</button>
        </div>
      </div>

      <div v-if="showCreateReleaseBranch" class="card mb-16" style="background:#f8f9fa">
        <h4 style="margin-bottom:12px">创建发布分支</h4>
        <div style="display:flex;gap:8px;align-items:flex-end;flex-wrap:wrap">
          <div class="form-group" style="flex:1;min-width:200px;margin-bottom:0">
            <label style="font-size:12px">选择仓库 *</label>
            <select v-model="releaseBranchForm.repoId" @change="onReleaseRepoChange" style="width:100%">
              <option value="">请选择仓库</option>
              <option v-for="r in projectRepos" :key="r.id" :value="r.id">{{ r.repoName }}</option>
            </select>
          </div>
          <div class="form-group" style="flex:1;min-width:200px;margin-bottom:0">
            <label style="font-size:12px">基于分支（父分支） *</label>
            <select v-model="releaseBranchForm.parentBranch" style="width:100%">
              <option value="">请选择父分支</option>
              <option v-for="b in releaseRepoBranches" :key="b.name" :value="b.name">
                {{ b.name }}{{ b.isDefault ? ' (默认)' : '' }}{{ b.isProtected ? ' 🔒' : '' }}
              </option>
            </select>
            <div v-if="loadingReleaseBranches" style="font-size:12px;color:var(--text-tertiary);margin-top:4px">
              <span class="loading-spinner"></span> 加载分支中...
            </div>
          </div>
          <div class="form-group" style="flex:1;min-width:200px;margin-bottom:0">
            <label style="font-size:12px">新分支名称（可选）</label>
            <input v-model="releaseBranchForm.branchName" :placeholder="`release/${release?.version || 'v1.0.0'}`" />
          </div>
          <button class="btn btn-primary btn-sm" @click="createReleaseBranch" :disabled="!releaseBranchForm.parentBranch">创建</button>
          <button class="btn btn-sm" @click="showCreateReleaseBranch = false">取消</button>
        </div>
      </div>

      <div v-if="showCreateFeatureBranch" class="card mb-16" style="background:#f8f9fa">
        <h4 style="margin-bottom:12px">创建功能分支</h4>
        <div style="display:flex;gap:8px;align-items:flex-end;flex-wrap:wrap">
          <div class="form-group" style="flex:1;min-width:200px;margin-bottom:0">
            <label style="font-size:12px">选择仓库 *</label>
            <select v-model="featureBranchForm.repoId" @change="onFeatureRepoChange" style="width:100%">
              <option value="">请选择仓库</option>
              <option v-for="r in projectRepos" :key="r.id" :value="r.id">{{ r.repoName }}</option>
            </select>
          </div>
          <div class="form-group" style="flex:1;min-width:200px;margin-bottom:0">
            <label style="font-size:12px">基于分支（父分支） *</label>
            <select v-model="featureBranchForm.parentBranch" style="width:100%">
              <option value="">请选择父分支</option>
              <option v-for="b in featureRepoBranches" :key="b.name" :value="b.name">
                {{ b.name }}{{ b.isDefault ? ' (默认)' : '' }}{{ b.isProtected ? ' 🔒' : '' }}
              </option>
            </select>
            <div v-if="loadingFeatureBranches" style="font-size:12px;color:var(--text-tertiary);margin-top:4px">
              <span class="loading-spinner"></span> 加载分支中...
            </div>
          </div>
          <div class="form-group" style="flex:1;min-width:200px;margin-bottom:0">
            <label style="font-size:12px">新分支名称 *</label>
            <input v-model="featureBranchForm.branchName" placeholder="feature/xxx" />
          </div>
          <button class="btn btn-primary btn-sm" @click="createFeatureBranch" :disabled="!featureBranchForm.parentBranch || !featureBranchForm.branchName">创建</button>
          <button class="btn btn-sm" @click="showCreateFeatureBranch = false">取消</button>
        </div>
      </div>

      <table class="data-table" v-if="branches.length">
        <thead>
          <tr>
            <th>分支名称</th>
            <th>类型</th>
            <th>父分支</th>
            <th>描述</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in branches" :key="b.id">
            <td><a v-if="b.gitlabBranchUrl" :href="b.gitlabBranchUrl" target="_blank" style="color:var(--primary)">{{ b.branchName }}</a><span v-else>{{ b.branchName }}</span></td>
            <td><span :class="['badge', b.branchType === 'release' ? 'badge-info' : 'badge-success']">{{ b.branchType === 'release' ? '发布分支' : '功能分支' }}</span></td>
            <td>{{ b.parentBranch || '-' }}</td>
            <td style="max-width:300px">
              <template v-if="editingBranchId !== b.id">
                <span style="color:var(--text-tertiary)">{{ b.description || '点击编辑添加描述' }}</span>
                <button class="btn btn-sm" style="margin-left:8px" @click="startEditBranch(b)">编辑</button>
              </template>
              <template v-else>
                <div style="display:flex;gap:4px">
                  <input v-model="editBranchDesc" placeholder="记录这个分支的用途、需求等" style="flex:1" />
                  <button class="btn btn-primary btn-sm" @click="saveBranchDesc(b.id)">保存</button>
                  <button class="btn btn-sm" @click="editingBranchId = ''">取消</button>
                </div>
              </template>
            </td>
            <td><button class="btn btn-sm btn-danger" @click="deleteBranch(b.id)">删除</button></td>
          </tr>
        </tbody>
      </table>
      <div v-else style="color:var(--text-tertiary);padding:12px 0">暂无关联分支</div>
    </div>

    <!-- Docker 镜像 -->
    <div class="card mb-16">
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
        <h3 style="font-size:15px">Docker 镜像</h3>
        <div style="display:flex;gap:8px">
          <button class="btn btn-sm" @click="showAddImage = true">手动添加镜像</button>
          <button class="btn btn-sm" @click="showWebhookGuide = !showWebhookGuide">Webhook 配置指南</button>
        </div>
      </div>

      <div v-if="showWebhookGuide" class="card mb-16" style="background:#f0f7ff;border:1px solid #b3d4fc">
        <h4 style="margin-bottom:12px;color:#1a73e8">GitLab CI 自动关联 Docker 镜像</h4>
        <div style="font-size:13px;line-height:1.8">
          <p><strong>方式一：Pipeline Webhook（自动记录构建事件）</strong></p>
          <ol style="margin:8px 0 16px 20px">
            <li>进入 GitLab 项目 → Settings → Webhooks</li>
            <li>URL: <code style="background:#e8f0fe;padding:2px 6px;border-radius:3px">{{ webhookUrl }}</code></li>
            <li>Secret Token: 留空或填写配置的 webhook_secret</li>
            <li>Trigger: 勾选 <strong>Pipeline events</strong></li>
          </ol>
          <p><strong>方式二：CI 脚本调用（推荐，精确上报镜像信息）</strong></p>
          <p>在 <code>.gitlab-ci.yml</code> 中添加 <code>report_images</code> 阶段：</p>
          <pre style="background:#1e1e1e;color:#d4d4d4;padding:12px;border-radius:6px;margin:8px 0;overflow-x:auto"><code># 在 variables 中添加：
variables:
  RELEASE_CENTER_URL: "{{ host }}"
  RELEASE_ID: "{{ release?.id || '发布单ID' }}"

# 在 stages 中添加 report 阶段
stages:
  - build
  - report    # 新增
  - deploy

# 新增 report_images job
report_images:
  stage: report
  script: |
    if [ -z "$RELEASE_ID" ]; then
      echo "RELEASE_ID not set, skipping"
      exit 0
    fi
    # 上报 server 镜像
    if [ -n "$IMAGE_server" ]; then
      IMAGE_NAME=$(echo "$IMAGE_server" | sed 's|.*//||' | cut -d: -f1)
      IMAGE_TAG=$(echo "$IMAGE_server" | cut -d: -f2)
      REGISTRY=$(echo "$IMAGE_server" | sed 's|.*//||' | cut -d/ -f1)
      curl -s -X POST "${RELEASE_CENTER_URL}/api/ci/build" \
        -H "Content-Type: application/json" \
        -d "{
          \"releaseId\": \"${RELEASE_ID}\",
          \"imageName\": \"${IMAGE_NAME}\",
          \"imageTag\": \"${IMAGE_TAG}\",
          \"registry\": \"${REGISTRY}\",
          \"branch\": \"${CI_COMMIT_BRANCH}\",
          \"commitSha\": \"${CI_COMMIT_SHA}\",
          \"commitMessage\": \"$(echo $CI_COMMIT_MESSAGE | head -1)\",
          \"ciPipelineId\": ${CI_PIPELINE_ID},
          \"ciPipelineUrl\": \"${CI_PIPELINE_URL}\"
        }"
    fi
    # 上报 querier 镜像（同理）</code></pre>
          <p style="margin-top:8px;color:#666">完整的 CI 配置示例请参考项目中的 <code>.gitlab-ci.yml.example</code></p>
        </div>
        <button class="btn btn-sm" @click="showWebhookGuide = false">关闭</button>
      </div>

      <div v-if="showAddImage" class="card mb-16" style="background:#f8f9fa">
        <h4 style="margin-bottom:12px">添加 Docker 镜像</h4>
        <div style="margin-bottom:12px">
          <div style="display:flex;gap:8px;margin-bottom:8px">
            <button :class="['btn btn-sm', addImageMode === 'pool' ? 'btn-primary' : '']" @click="addImageMode = 'pool'; loadPoolImages()">从构建记录选择</button>
            <button :class="['btn btn-sm', addImageMode === 'manual' ? 'btn-primary' : '']" @click="addImageMode = 'manual'">手动输入</button>
          </div>
        </div>

        <div v-if="addImageMode === 'pool'" style="margin-bottom:12px">
          <div v-if="poolImages.length" style="max-height:200px;overflow-y:auto;border:1px solid var(--border);border-radius:6px">
            <div v-for="img in poolImages" :key="img.id" style="display:flex;justify-content:space-between;align-items:center;padding:8px 12px;border-bottom:1px solid var(--border);cursor:hover" @click="selectPoolImage(img)" :style="{background: selectedPoolImage?.id === img.id ? '#e8f0fe' : ''}">
              <div>
                <div style="font-weight:500">{{ img.imageName }}:{{ img.imageTag }}</div>
                <div style="font-size:12px;color:var(--text-tertiary)">
                  {{ img.registry }} | {{ img.branch }} | {{ img.commitSha?.substring(0, 8) || '-' }} | {{ img.commitMessage?.substring(0, 30) || '' }}
                </div>
              </div>
              <span v-if="selectedPoolImage?.id === img.id" style="color:var(--primary)">✓</span>
            </div>
          </div>
          <div v-else style="color:var(--text-tertiary);padding:12px;text-align:center">
            暂无构建记录，可通过 CI 上报或 Webhook 自动记录
          </div>
          <div v-if="selectedPoolImage" style="margin-top:8px">
            <button class="btn btn-primary btn-sm" @click="addFromPool">添加选中的镜像</button>
          </div>
        </div>

        <div v-if="addImageMode === 'manual'" style="display:flex;gap:8px;align-items:flex-end;flex-wrap:wrap">
          <div class="form-group" style="flex:1;min-width:200px;margin-bottom:0">
            <label style="font-size:12px">镜像名称 *</label>
            <input v-model="imageForm.imageName" placeholder="my-app" />
          </div>
          <div class="form-group" style="flex:0.5;min-width:100px;margin-bottom:0">
            <label style="font-size:12px">Tag *</label>
            <input v-model="imageForm.imageTag" placeholder="latest" />
          </div>
          <div class="form-group" style="flex:1;min-width:150px;margin-bottom:0">
            <label style="font-size:12px">Registry</label>
            <input v-model="imageForm.registry" placeholder="registry.example.com" />
          </div>
          <div class="form-group" style="flex:1;min-width:100px;margin-bottom:0">
            <label style="font-size:12px">分支</label>
            <input v-model="imageForm.branch" placeholder="main" />
          </div>
          <button class="btn btn-primary btn-sm" @click="addDockerImage">添加</button>
        </div>
        <button class="btn btn-sm" style="margin-top:8px" @click="showAddImage = false">关闭</button>
      </div>

      <table class="data-table" v-if="dockerImages.length">
        <thead>
          <tr>
            <th>镜像</th>
            <th>Tag</th>
            <th>分支</th>
            <th>Commit</th>
            <th>来源</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="img in dockerImages" :key="img.id">
            <td>{{ img.imageName }}</td>
            <td>{{ img.imageTag }}</td>
            <td>{{ img.branch || '-' }}</td>
            <td>{{ img.commitSha ? img.commitSha.substring(0, 8) : '-' }}</td>
            <td><span :class="['badge', img.source === 'manual' ? 'badge-info' : 'badge-success']">{{ img.source === 'manual' ? '手动' : img.source === 'ci' ? 'CI构建' : 'Webhook' }}</span></td>
            <td>{{ formatTime(img.createdAt) }}</td>
            <td><button class="btn btn-sm btn-danger" @click="deleteDockerImage(img.id)">删除</button></td>
          </tr>
        </tbody>
      </table>
      <div v-else style="color:var(--text-tertiary);padding:12px 0">暂无镜像，可通过 GitLab CI 自动关联或手动添加</div>
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

    <!-- 通知预览模态框 -->
    <div v-if="showNotifyModal" style="position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,0.5);display:flex;align-items:center;justify-content:center;z-index:1000" @click.self="showNotifyModal=false">
      <div class="card" style="width:700px;max-height:85vh;overflow-y:auto">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
          <h3>{{ notifyChannel === 'lanxin' ? '蓝信消息预览' : '邮件内容预览' }}</h3>
          <button class="btn btn-sm" @click="showNotifyModal=false">关闭</button>
        </div>
        <div v-if="notifyChannel === 'lanxin'">
          <div v-if="notifyPreview?.lanxinEnabled">
            <div style="background:#f0f5ff;border:1px solid #b3d4fc;border-radius:8px;padding:16px;margin-bottom:16px;white-space:pre-wrap;font-size:13px;line-height:1.8">{{ notifyPreview?.lanxinMessage }}</div>
            <div style="display:flex;gap:8px;justify-content:flex-end">
              <button class="btn" @click="showNotifyModal=false">取消</button>
              <button class="btn btn-primary" @click="confirmSendNotify">确认发送</button>
            </div>
          </div>
          <div v-else style="color:var(--danger);padding:20px;text-align:center">
            蓝信通知未启用，请在 config.yaml 中配置 lanxin.enabled=true 和 lanxin.url
          </div>
        </div>
        <div v-if="notifyChannel === 'email'">
          <div v-if="notifyPreview?.emailEnabled">
            <div style="margin-bottom:8px;color:var(--text-secondary);font-size:13px">
              收件人：<strong>{{ notifyPreview?.emailTo?.join(', ') || '未配置' }}</strong> &nbsp;|&nbsp;
              主题：<strong>{{ notifyPreview?.emailSubject }}</strong>
            </div>
            <div style="border:1px solid var(--border);border-radius:8px;overflow:hidden;margin-bottom:16px;max-height:400px;overflow-y:auto">
              <div v-if="notifyPreview?.emailHtml" v-html="notifyPreview.emailHtml"></div>
            </div>
            <div style="display:flex;gap:8px;justify-content:flex-end">
              <button class="btn" @click="showNotifyModal=false">取消</button>
              <button class="btn btn-primary" @click="confirmSendNotify">确认发送</button>
            </div>
          </div>
          <div v-else style="color:var(--danger);padding:20px;text-align:center">
            邮件通知未启用，请在 config.yaml 中配置 email.enabled=true 和 SMTP 相关参数
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { releaseApi, itemApi, snapshotApi, zentaoApi, projectApi, deploymentApi, healthApi, repoApi, branchApi, dockerImageApi, gitlabApi, featureApi, notifyApi } from '../api'
import type { Release, ReleaseItem, ReleaseSnapshot, Deployment, ProjectRepo, ReleaseBranch, DockerImage, ReleaseFeature } from '../api'

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
const editForm = ref({ name: '', version: '', summary: '', parentBranch: '' })
const editingDeploymentId = ref('')
const editDeployForm = ref({ moduleName: '', address: '', description: '' })

const projectRepos = ref<ProjectRepo[]>([])
const branches = ref<ReleaseBranch[]>([])
const dockerImages = ref<DockerImage[]>([])
const showCreateReleaseBranch = ref(false)
const showCreateFeatureBranch = ref(false)
const showAddImage = ref(false)
const addImageMode = ref<'pool' | 'manual'>('pool')
const poolImages = ref<any[]>([])
const selectedPoolImage = ref<any>(null)
const releaseBranchForm = ref({ repoId: '', branchName: '', parentBranch: '' })
const featureBranchForm = ref({ repoId: '', branchName: '', parentBranch: '' })
const imageForm = ref({ imageName: '', imageTag: '', registry: '', branch: '' })

const releaseRepoBranches = ref<any[]>([])
const featureRepoBranches = ref<any[]>([])
const loadingReleaseBranches = ref(false)
const loadingFeatureBranches = ref(false)
const editingBranchId = ref('')
const editBranchDesc = ref('')
const showWebhookGuide = ref(false)
const webhookUrl = ref(`${window.location.origin}/api/webhook/gitlab`)
const host = ref(window.location.origin)

const features = ref<ReleaseFeature[]>([])
const showFeatureForm = ref(false)
const featureForm = ref({ title: '', content: '' })
const editingFeatureId = ref('')
const editFeatureForm = ref({ title: '', content: '' })

const showNotifyModal = ref(false)
const notifyChannel = ref<'lanxin' | 'email'>('lanxin')
const notifyPreview = ref<any>(null)
const notifySending = ref(false)

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
    await Promise.all([loadItems(), loadSnapshots(), loadDeployments(), loadBranches(), loadDockerImages(), loadFeatures()])
    if (release.value?.projectId) {
      const reposResp: any = await repoApi.list(release.value.projectId)
      projectRepos.value = reposResp?.list || []
      console.log('Loaded projectRepos:', projectRepos.value)
    }
  } catch (e) {
    console.error('Load failed:', e)
  } finally {
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

async function loadBranches() {
  const id = route.query.id as string
  try {
    const resp: any = await branchApi.list(id)
    branches.value = resp?.list || []
  } catch {}
}

async function loadDockerImages() {
  const id = route.query.id as string
  try {
    const resp: any = await dockerImageApi.list(id)
    dockerImages.value = resp?.list || []
  } catch {}
}

async function loadFeatures() {
  const id = route.query.id as string
  try {
    const resp: any = await featureApi.list(id)
    features.value = resp?.list || []
  } catch {}
}

async function submitFeature() {
  if (!featureForm.value.title) return alert('请输入标题')
  try {
    await featureApi.add({
      releaseId: route.query.id as string,
      title: featureForm.value.title,
      content: featureForm.value.content,
    })
    featureForm.value = { title: '', content: '' }
    showFeatureForm.value = false
    await loadFeatures()
  } catch {}
}

function startEditFeature(f: ReleaseFeature) {
  editingFeatureId.value = f.id
  editFeatureForm.value = { title: f.title, content: f.content }
}

async function saveEditFeature() {
  try {
    await featureApi.update({ id: editingFeatureId.value, ...editFeatureForm.value })
    editingFeatureId.value = ''
    await loadFeatures()
  } catch {}
}

async function removeFeature(id: string) {
  if (!confirm('确定删除此功能说明？')) return
  try {
    await featureApi.delete(id)
    await loadFeatures()
  } catch {}
}

async function sendLanxin() {
  notifyChannel.value = 'lanxin'
  try {
    const resp: any = await notifyApi.preview({
      releaseId: route.query.id as string,
      version: release.value?.version || undefined,
    })
    notifyPreview.value = resp?.data || null
    showNotifyModal.value = true
  } catch {}
}

async function sendEmail() {
  notifyChannel.value = 'email'
  try {
    const resp: any = await notifyApi.preview({
      releaseId: route.query.id as string,
      version: release.value?.version || undefined,
    })
    notifyPreview.value = resp?.data || null
    showNotifyModal.value = true
  } catch {}
}

async function confirmSendNotify() {
  notifySending.value = true
  try {
    const resp: any = await notifyApi.send({
      releaseId: route.query.id as string,
      version: release.value?.version || undefined,
      channel: notifyChannel.value,
    })
    const result = resp?.data
    if (notifyChannel.value === 'lanxin') {
      if (result?.lanxinSuccess) {
        alert('蓝信发送成功！')
      } else {
        alert('蓝信发送失败：' + (result?.lanxinError || '未知错误'))
      }
    } else {
      if (result?.emailSuccess) {
        alert('邮件发送成功！')
      } else {
        alert('邮件发送失败：' + (result?.emailError || '未知错误'))
      }
    }
    showNotifyModal.value = false
  } catch (e: any) {
    alert('发送失败：' + (e?.message || '未知错误'))
  } finally {
    notifySending.value = false
  }
}

async function onReleaseRepoChange() {
  releaseBranchForm.value.parentBranch = ''
  releaseBranchForm.value.branchName = ''
  if (!releaseBranchForm.value.repoId) {
    releaseRepoBranches.value = []
    return
  }
  const repo = projectRepos.value.find(r => r.id === releaseBranchForm.value.repoId)
  if (!repo) return
  loadingReleaseBranches.value = true
  try {
    const resp: any = await gitlabApi.branches(repo.gitlabProjectId)
    releaseRepoBranches.value = resp?.list || []
    // 优先使用发布单的基准分支
    if (release.value?.parentBranch && releaseRepoBranches.value.find(b => b.name === release.value!.parentBranch)) {
      releaseBranchForm.value.parentBranch = release.value.parentBranch
    } else {
      const defaultBranch = releaseRepoBranches.value.find(b => b.isDefault)
      if (defaultBranch) {
        releaseBranchForm.value.parentBranch = defaultBranch.name
      }
    }
  } catch {} finally {
    loadingReleaseBranches.value = false
  }
}

async function onFeatureRepoChange() {
  featureBranchForm.value.parentBranch = ''
  featureBranchForm.value.branchName = ''
  if (!featureBranchForm.value.repoId) {
    featureRepoBranches.value = []
    return
  }
  const repo = projectRepos.value.find(r => r.id === featureBranchForm.value.repoId)
  if (!repo) {
    console.error('Repo not found:', featureBranchForm.value.repoId)
    return
  }
  loadingFeatureBranches.value = true
  try {
    console.log('Loading branches for repo:', repo.gitlabProjectId)
    const resp: any = await gitlabApi.branches(repo.gitlabProjectId)
    console.log('Branches response:', resp)
    featureRepoBranches.value = resp?.list || []
    console.log('featureRepoBranches:', featureRepoBranches.value)
    // 优先选择已关联的发布分支作为父分支
    const releaseBranch = branches.value.find(b => b.repoId === featureBranchForm.value.repoId && b.branchType === 'release')
    if (releaseBranch) {
      featureBranchForm.value.parentBranch = releaseBranch.branchName
    } else if (release.value?.parentBranch && featureRepoBranches.value.find(b => b.name === release.value!.parentBranch)) {
      // 如果没有发布分支，使用发布单的基准分支
      featureBranchForm.value.parentBranch = release.value.parentBranch
    } else {
      const defaultBranch = featureRepoBranches.value.find(b => b.isDefault)
      if (defaultBranch) {
        featureBranchForm.value.parentBranch = defaultBranch.name
      }
    }
  } catch (e) {
    console.error('Failed to load branches:', e)
  } finally {
    loadingFeatureBranches.value = false
  }
}

async function createReleaseBranch() {
  if (!releaseBranchForm.value.repoId) return alert('请选择仓库')
  if (!releaseBranchForm.value.parentBranch) return alert('请选择基于分支')
  try {
    await branchApi.createRelease({
      releaseId: route.query.id as string,
      repoId: releaseBranchForm.value.repoId,
      branchName: releaseBranchForm.value.branchName || undefined,
      parentBranch: releaseBranchForm.value.parentBranch,
    })
    releaseBranchForm.value = { repoId: '', branchName: '', parentBranch: '' }
    releaseRepoBranches.value = []
    showCreateReleaseBranch.value = false
    await loadBranches()
  } catch {}
}

async function createFeatureBranch() {
  if (!featureBranchForm.value.repoId) return alert('请选择仓库')
  if (!featureBranchForm.value.parentBranch) return alert('请选择基于分支')
  if (!featureBranchForm.value.branchName) return alert('请输入分支名称')
  try {
    await branchApi.createFeature({
      releaseId: route.query.id as string,
      repoId: featureBranchForm.value.repoId,
      branchName: featureBranchForm.value.branchName,
      parentBranch: featureBranchForm.value.parentBranch,
    })
    featureBranchForm.value = { repoId: '', branchName: '', parentBranch: '' }
    featureRepoBranches.value = []
    showCreateFeatureBranch.value = false
    await loadBranches()
  } catch {}
}

async function deleteBranch(id: string) {
  if (!confirm('确定删除此分支？')) return
  try {
    await branchApi.delete(id)
    await loadBranches()
  } catch {}
}

function startEditBranch(b: ReleaseBranch) {
  editingBranchId.value = b.id
  editBranchDesc.value = b.description || ''
}

async function saveBranchDesc(id: string) {
  try {
    await branchApi.update({ id, description: editBranchDesc.value })
    editingBranchId.value = ''
    await loadBranches()
  } catch {}
}

async function loadPoolImages() {
  try {
    const resp: any = await dockerImageApi.pool()
    poolImages.value = resp?.list || []
  } catch {}
}

function selectPoolImage(img: any) {
  selectedPoolImage.value = selectedPoolImage.value?.id === img.id ? null : img
}

async function addFromPool() {
  if (!selectedPoolImage.value) return alert('请选择一个镜像')
  try {
    await dockerImageApi.add({
      releaseId: route.query.id as string,
      imageName: selectedPoolImage.value.imageName,
      imageTag: selectedPoolImage.value.imageTag,
      registry: selectedPoolImage.value.registry || undefined,
      branch: selectedPoolImage.value.branch || undefined,
      commitSha: selectedPoolImage.value.commitSha || undefined,
      commitMessage: selectedPoolImage.value.commitMessage || undefined,
    })
    selectedPoolImage.value = null
    showAddImage.value = false
    await loadDockerImages()
  } catch {}
}

async function archiveRelease() {
  if (!confirm('确定将此发布单标记为「发布完成」？后续 CI 构建的镜像将不会自动添加到此发布单。')) return
  try {
    await releaseApi.update({ id: route.query.id as string, status: 'archived' })
    await load()
  } catch {}
}

async function addDockerImage() {
  if (!imageForm.value.imageName || !imageForm.value.imageTag) return alert('请输入镜像名称和Tag')
  try {
    await dockerImageApi.add({
      releaseId: route.query.id as string,
      imageName: imageForm.value.imageName,
      imageTag: imageForm.value.imageTag,
      registry: imageForm.value.registry || undefined,
      branch: imageForm.value.branch || undefined,
    })
    imageForm.value = { imageName: '', imageTag: '', registry: '', branch: '' }
    showAddImage.value = false
    await loadDockerImages()
  } catch {}
}

async function deleteDockerImage(id: string) {
  if (!confirm('确定删除此镜像？')) return
  try {
    await dockerImageApi.delete(id)
    await loadDockerImages()
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
  editForm.value = {
    name: release.value.name,
    version: release.value.version,
    summary: release.value.summary,
    parentBranch: release.value.parentBranch || ''
  }
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
