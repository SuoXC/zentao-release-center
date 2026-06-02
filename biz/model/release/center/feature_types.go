package center

type ReleaseFeature struct {
	ID        string `thrift:"id,1" form:"id" json:"id" query:"id"`
	ReleaseId string `thrift:"releaseId,2" form:"releaseId" json:"releaseId" query:"releaseId"`
	Title     string `thrift:"title,3" form:"title" json:"title" query:"title"`
	Content   string `thrift:"content,4" form:"content" json:"content" query:"content"`
	SortOrder int32  `thrift:"sortOrder,5" form:"sortOrder" json:"sortOrder" query:"sortOrder"`
	CreatedAt string `thrift:"createdAt,6" form:"createdAt" json:"createdAt" query:"createdAt"`
	UpdatedAt string `thrift:"updatedAt,7" form:"updatedAt" json:"updatedAt" query:"updatedAt"`
}

type AddFeatureReq struct {
	ReleaseId string `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
	Title     string `thrift:"title,2" form:"title" json:"title" query:"title"`
	Content   string `thrift:"content,3" form:"content" json:"content" query:"content"`
}

type UpdateFeatureReq struct {
	ID      string  `thrift:"id,1" form:"id" json:"id" query:"id"`
	Title   *string `thrift:"title,2,optional" form:"title" json:"title,omitempty" query:"title"`
	Content *string `thrift:"content,3,optional" form:"content" json:"content,omitempty" query:"content"`
}

type DeleteFeatureReq struct {
	ID string `thrift:"id,1" form:"id" json:"id" query:"id"`
}

type ListFeaturesReq struct {
	ReleaseId string `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
}

type ReorderFeaturesReq struct {
	ReleaseId string      `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
	Items     []SortItem  `thrift:"items,2" form:"items" json:"items" query:"items"`
}

type FeatureResp struct {
	Base *BaseResp        `thrift:"base,1" form:"base" json:"base" query:"base"`
	Data *ReleaseFeature  `thrift:"data,2,optional" form:"data" json:"data,omitempty" query:"data"`
}

type FeatureListResp struct {
	Base *BaseResp          `thrift:"base,1" form:"base" json:"base" query:"base"`
	List []*ReleaseFeature  `thrift:"list,2" form:"list" json:"list" query:"list"`
}
