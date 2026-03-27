package v1

type GetWorkspaceConfRequest struct {
	Keys *string `json:"keys"`
}

type WorkspaceConf struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}
