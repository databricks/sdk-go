package v1

type CreateRepo struct {
	Url            *string         `json:"url"`
	Provider       *string         `json:"provider"`
	Path           *string         `json:"path"`
	SparseCheckout *SparseCheckout `json:"sparse_checkout"`
}

type CreateRepo_Response struct {
	Id             *int64          `json:"id"`
	Path           *string         `json:"path"`
	Url            *string         `json:"url"`
	Provider       *string         `json:"provider"`
	Branch         *string         `json:"branch"`
	HeadCommitId   *string         `json:"head_commit_id"`
	SparseCheckout *SparseCheckout `json:"sparse_checkout"`
}

type DeleteProject struct {
	Id *int64 `json:"id"`
}

type DeleteProject_Response struct {
}

type GetRepo struct {
	Id *int64 `json:"id"`
}

type GetRepo_Response struct {
	Id             *int64          `json:"id"`
	Path           *string         `json:"path"`
	Url            *string         `json:"url"`
	Provider       *string         `json:"provider"`
	Branch         *string         `json:"branch"`
	HeadCommitId   *string         `json:"head_commit_id"`
	SparseCheckout *SparseCheckout `json:"sparse_checkout"`
}

type ListRepos struct {
	PathPrefix    *string `json:"path_prefix"`
	NextPageToken *string `json:"next_page_token"`
}

type ListRepos_Response struct {
	Repos         []RepoInfo `json:"repos"`
	NextPageToken *string    `json:"next_page_token"`
}

// Git folder (repo) information..
type RepoInfo struct {
	Id             *int64          `json:"id"`
	Path           *string         `json:"path"`
	Url            *string         `json:"url"`
	Provider       *string         `json:"provider"`
	Branch         *string         `json:"branch"`
	HeadCommitId   *string         `json:"head_commit_id"`
	SparseCheckout *SparseCheckout `json:"sparse_checkout"`
}

// Sparse checkout configuration, it contains options like cone patterns..
type SparseCheckout struct {
	Patterns []string `json:"patterns"`
}

// Sparse checkout configuration, it contains options like cone patterns..
type SparseCheckoutUpdate struct {
	Patterns []string `json:"patterns"`
}

type UpdateRepo struct {
	Id             *int64                `json:"id"`
	Branch         *string               `json:"branch"`
	Tag            *string               `json:"tag"`
	SparseCheckout *SparseCheckoutUpdate `json:"sparse_checkout"`
}

type UpdateRepo_Response struct {
}
