package v1

type UsageDashboardMajorVersion string

const (
	UsageDashboardMajorVersion_UsageDashboardMajorVersionUnspecified UsageDashboardMajorVersion = "USAGE_DASHBOARD_MAJOR_VERSION_UNSPECIFIED"
	UsageDashboardMajorVersion_UsageDashboardMajorVersion1           UsageDashboardMajorVersion = "USAGE_DASHBOARD_MAJOR_VERSION_1"
	UsageDashboardMajorVersion_UsageDashboardMajorVersion2           UsageDashboardMajorVersion = "USAGE_DASHBOARD_MAJOR_VERSION_2"
)

// String representation for [fmt.Print].
func (f *UsageDashboardMajorVersion) String() string {
	return string(*f)
}

type UsageDashboardType string

const (
	UsageDashboardType_UsageDashboardTypeUnspecified UsageDashboardType = "USAGE_DASHBOARD_TYPE_UNSPECIFIED"
	UsageDashboardType_UsageDashboardTypeWorkspace   UsageDashboardType = "USAGE_DASHBOARD_TYPE_WORKSPACE"
	UsageDashboardType_UsageDashboardTypeGlobal      UsageDashboardType = "USAGE_DASHBOARD_TYPE_GLOBAL"
)

// String representation for [fmt.Print].
func (f *UsageDashboardType) String() string {
	return string(*f)
}

type CreateBillingUsageDashboard struct {
	WorkspaceId   *int64                      `json:"workspace_id"`
	AccountId     *string                     `json:"account_id"`
	DashboardType *UsageDashboardType         `json:"dashboard_type"`
	MajorVersion  *UsageDashboardMajorVersion `json:"major_version"`
}

type CreateBillingUsageDashboard_Response struct {
	DashboardId *string `json:"dashboard_id"`
}

type GetBillingUsageDashboard struct {
	WorkspaceId   *int64              `json:"workspace_id"`
	AccountId     *string             `json:"account_id"`
	DashboardType *UsageDashboardType `json:"dashboard_type"`
}

type GetBillingUsageDashboard_Response struct {
	DashboardId  *string `json:"dashboard_id"`
	DashboardUrl *string `json:"dashboard_url"`
}
