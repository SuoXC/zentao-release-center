namespace go release.center

struct BaseResp {
    1: i32 code
    2: string message
}

// ==================== 项目 ====================

struct Project {
    1: string id
    2: string name
    3: string description
    4: i32 zentaoProductId
    5: i32 zentaoProjectId
    6: string zentaoProductName
    7: string zentaoProjectName
    8: string zentaoServer
    9: string status
    10: string createdAt
    11: string updatedAt
}

struct CreateProjectReq {
    1: string name (api.body="name")
    2: optional string description (api.body="description")
    3: optional i32 zentaoProductId (api.body="zentaoProductId")
    4: optional i32 zentaoProjectId (api.body="zentaoProjectId")
    5: optional string zentaoProductName (api.body="zentaoProductName")
    6: optional string zentaoProjectName (api.body="zentaoProjectName")
}

struct UpdateProjectReq {
    1: string id (api.body="id")
    2: optional string name (api.body="name")
    3: optional string description (api.body="description")
    4: optional i32 zentaoProductId (api.body="zentaoProductId")
    5: optional i32 zentaoProjectId (api.body="zentaoProjectId")
    6: optional string zentaoProductName (api.body="zentaoProductName")
    7: optional string zentaoProjectName (api.body="zentaoProjectName")
    8: optional string status (api.body="status")
}

struct ListProjectsReq {
    1: optional string status (api.query="status")
    2: optional i32 page (api.query="page")
    3: optional i32 pageSize (api.query="pageSize")
}

struct GetProjectReq {
    1: string id (api.query="id")
}

struct DeleteProjectReq {
    1: string id (api.body="id")
}

struct ListProjectsResp {
    1: BaseResp base
    2: list<Project> list
    3: i32 total
}

struct ProjectResp {
    1: BaseResp base
    2: optional Project data
}

// ==================== 发布单 ====================

struct Release {
    1: string id
    2: string projectId
    3: string projectName
    4: string name
    5: string version
    6: string status
    7: string summary
    8: i32 publishCount
    9: string firstPublishedAt
    10: string lastPublishedAt
    11: i32 itemCount
    12: i32 bugCount
    13: i32 taskCount
    14: i32 noteCount
    15: string createdAt
    16: string updatedAt
}

struct CreateReleaseReq {
    1: string projectId (api.body="projectId")
    2: string name (api.body="name")
    3: optional string version (api.body="version")
    4: optional string summary (api.body="summary")
}

struct UpdateReleaseReq {
    1: string id (api.body="id")
    2: optional string name (api.body="name")
    3: optional string version (api.body="version")
    4: optional string summary (api.body="summary")
    5: optional string status (api.body="status")
}

struct ListReleasesReq {
    1: string projectId (api.query="projectId")
    2: optional string status (api.query="status")
    3: optional i32 page (api.query="page")
    4: optional i32 pageSize (api.query="pageSize")
}

struct GetReleaseReq {
    1: string id (api.query="id")
}

struct DeleteReleaseReq {
    1: string id (api.body="id")
}

struct ListReleasesResp {
    1: BaseResp base
    2: list<Release> list
    3: i32 total
}

struct ReleaseResp {
    1: BaseResp base
    2: optional Release data
}

// ==================== 发布单条目 ====================

struct ReleaseItem {
    1: string id
    2: string releaseId
    3: string itemType
    4: i32 sortOrder
    5: optional i32 zentaoId
    6: optional string zentaoType
    7: optional string title
    8: optional string severity
    9: optional string priority
    10: optional string status
    11: optional string assignedTo
    12: optional string resolvedBy
    13: optional string zentaoUrl
    14: optional string steps
    15: optional string noteTitle
    16: optional string noteContent
    17: string createdAt
    18: string updatedAt
}

struct AddItemReq {
    1: string releaseId (api.body="releaseId")
    2: string itemType (api.body="itemType")
    3: optional i32 zentaoId (api.body="zentaoId")
    4: optional string zentaoType (api.body="zentaoType")
    5: optional string title (api.body="title")
    6: optional string severity (api.body="severity")
    7: optional string priority (api.body="priority")
    8: optional string status (api.body="status")
    9: optional string assignedTo (api.body="assignedTo")
    10: optional string resolvedBy (api.body="resolvedBy")
    11: optional string zentaoUrl (api.body="zentaoUrl")
    12: optional string steps (api.body="steps")
    13: optional string noteTitle (api.body="noteTitle")
    14: optional string noteContent (api.body="noteContent")
}

struct BatchAddItemsReq {
    1: string releaseId (api.body="releaseId")
    2: list<AddItemReq> items (api.body="items")
}

struct UpdateItemReq {
    1: string id (api.body="id")
    2: optional string noteTitle (api.body="noteTitle")
    3: optional string noteContent (api.body="noteContent")
    4: optional i32 sortOrder (api.body="sortOrder")
}

struct DeleteItemReq {
    1: string id (api.body="id")
}

struct ListItemsReq {
    1: string releaseId (api.query="releaseId")
}

struct ReorderItemsReq {
    1: string releaseId (api.body="releaseId")
    2: list<SortItem> items (api.body="items")
}

struct SortItem {
    1: string id
    2: i32 sortOrder
}

struct RefreshItemsReq {
    1: string releaseId (api.body="releaseId")
}

struct ReleaseItemListResp {
    1: BaseResp base
    2: list<ReleaseItem> list
}

struct BaseOnlyResp {
    1: BaseResp base
}

// ==================== 发布快照 ====================

struct ReleaseSnapshot {
    1: string id
    2: string releaseId
    3: string version
    4: string content
    5: i32 itemCount
    6: i32 bugCount
    7: i32 taskCount
    8: i32 noteCount
    9: string publishedAt
}

struct PublishReleaseReq {
    1: string releaseId (api.body="releaseId")
    2: optional string version (api.body="version")
}

struct ListSnapshotsReq {
    1: string releaseId (api.query="releaseId")
}

struct GetSnapshotReq {
    1: string id (api.query="id")
}

struct ExportReq {
    1: string releaseId (api.query="releaseId")
    2: optional string snapshotId (api.query="snapshotId")
    3: optional string format (api.query="format")
}

struct SnapshotResp {
    1: BaseResp base
    2: optional ReleaseSnapshot data
}

struct SnapshotListResp {
    1: BaseResp base
    2: list<ReleaseSnapshot> list
}

struct ExportResp {
    1: BaseResp base
    2: optional string content
    3: optional string filename
    4: optional string format
}

// ==================== 禅道数据代理 ====================

struct ZentaoBugsReq {
    1: optional i32 productId (api.query="productId")
    2: optional i32 projectId (api.query="projectId")
    3: optional string status (api.query="status")
    4: optional string assignedTo (api.query="assignedTo")
    5: optional i32 page (api.query="page")
    6: optional i32 pageSize (api.query="pageSize")
}

struct ZentaoTasksReq {
    1: optional i32 executionId (api.query="executionId")
    2: optional i32 productId (api.query="productId")
    3: optional string status (api.query="status")
    4: optional string assignedTo (api.query="assignedTo")
    5: optional i32 page (api.query="page")
    6: optional i32 pageSize (api.query="pageSize")
}

struct ZentaoProductsReq {}

struct ZentaoProjectsReq {
    1: optional i32 productId (api.query="productId")
}

struct ZentaoExecutionsReq {
    1: optional i32 projectId (api.query="projectId")
}

struct ZentaoDataResp {
    1: BaseResp base
    2: optional string data
}

struct ZentaoPaginatedResp {
    1: BaseResp base
    2: optional string list
    3: i32 total
    4: i32 page
    5: i32 pageSize
}

// ==================== 部署地址 ====================

struct Deployment {
    1: string id
    2: string releaseId
    3: string moduleName
    4: string address
    5: string description
    6: i32 sortOrder
    7: string createdAt
    8: string updatedAt
}

struct AddDeploymentReq {
    1: string releaseId (api.body="releaseId")
    2: string moduleName (api.body="moduleName")
    3: string address (api.body="address")
    4: optional string description (api.body="description")
}

struct UpdateDeploymentReq {
    1: string id (api.body="id")
    2: optional string moduleName (api.body="moduleName")
    3: optional string address (api.body="address")
    4: optional string description (api.body="description")
}

struct DeleteDeploymentReq {
    1: string id (api.body="id")
}

struct ListDeploymentsReq {
    1: string releaseId (api.query="releaseId")
}

struct DeploymentListResp {
    1: BaseResp base
    2: list<Deployment> list
}

struct DeploymentResp {
    1: BaseResp base
    2: optional Deployment data
}

// ==================== 健康 ====================

struct HealthResp {
    1: BaseResp base
    2: optional string status
    3: optional string zentaoMiniStatus
    4: optional string zentaoBaseUrl
}

struct EmptyReq {}

// ==================== 服务定义 ====================

service ReleaseCenterService {
    ProjectResp CreateProject(1: CreateProjectReq req) (api.post="/api/projects")
    ProjectResp UpdateProject(1: UpdateProjectReq req) (api.post="/api/projects/update")
    ListProjectsResp ListProjects(1: ListProjectsReq req) (api.get="/api/projects")
    ProjectResp GetProject(1: GetProjectReq req) (api.get="/api/projects/detail")
    BaseOnlyResp DeleteProject(1: DeleteProjectReq req) (api.post="/api/projects/delete")

    ReleaseResp CreateRelease(1: CreateReleaseReq req) (api.post="/api/releases")
    ReleaseResp UpdateRelease(1: UpdateReleaseReq req) (api.post="/api/releases/update")
    ListReleasesResp ListReleases(1: ListReleasesReq req) (api.get="/api/releases")
    ReleaseResp GetRelease(1: GetReleaseReq req) (api.get="/api/releases/detail")
    BaseOnlyResp DeleteRelease(1: DeleteReleaseReq req) (api.post="/api/releases/delete")

    BaseOnlyResp AddItem(1: AddItemReq req) (api.post="/api/release-items")
    BaseOnlyResp BatchAddItems(1: BatchAddItemsReq req) (api.post="/api/release-items/batch")
    BaseOnlyResp UpdateItem(1: UpdateItemReq req) (api.post="/api/release-items/update")
    BaseOnlyResp DeleteItem(1: DeleteItemReq req) (api.post="/api/release-items/delete")
    ReleaseItemListResp ListItems(1: ListItemsReq req) (api.get="/api/release-items")
    BaseOnlyResp ReorderItems(1: ReorderItemsReq req) (api.post="/api/release-items/reorder")
    BaseOnlyResp RefreshItems(1: RefreshItemsReq req) (api.post="/api/release-items/refresh")

    SnapshotResp PublishRelease(1: PublishReleaseReq req) (api.post="/api/releases/publish")
    SnapshotListResp ListSnapshots(1: ListSnapshotsReq req) (api.get="/api/release-snapshots")
    SnapshotResp GetSnapshot(1: GetSnapshotReq req) (api.get="/api/release-snapshots/detail")
    ExportResp ExportRelease(1: ExportReq req) (api.get="/api/releases/export")

    ZentaoPaginatedResp GetZentaoBugs(1: ZentaoBugsReq req) (api.get="/api/zentao/bugs")
    ZentaoPaginatedResp GetZentaoTasks(1: ZentaoTasksReq req) (api.get="/api/zentao/tasks")
    ZentaoDataResp GetZentaoProducts(1: ZentaoProductsReq req) (api.get="/api/zentao/products")
    ZentaoDataResp GetZentaoProjects(1: ZentaoProjectsReq req) (api.get="/api/zentao/projects")
    ZentaoDataResp GetZentaoExecutions(1: ZentaoExecutionsReq req) (api.get="/api/zentao/executions")

    DeploymentResp AddDeployment(1: AddDeploymentReq req) (api.post="/api/deployments")
    DeploymentResp UpdateDeployment(1: UpdateDeploymentReq req) (api.post="/api/deployments/update")
    BaseOnlyResp DeleteDeployment(1: DeleteDeploymentReq req) (api.post="/api/deployments/delete")
    DeploymentListResp ListDeployments(1: ListDeploymentsReq req) (api.get="/api/deployments")

    HealthResp Health(1: EmptyReq req) (api.get="/api/health")
}
