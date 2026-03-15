package v1

import "time"

type DashboardView string

const (
	DashboardView_DashboardViewBasic DashboardView = "DASHBOARD_VIEW_BASIC"
)

// String representation for [fmt.Print].
func (f *DashboardView) String() string {
	return string(*f)
}

type LifecycleState string

const (
	LifecycleState_Active  LifecycleState = "ACTIVE"
	LifecycleState_Trashed LifecycleState = "TRASHED"
)

// String representation for [fmt.Print].
func (f *LifecycleState) String() string {
	return string(*f)
}

type SchedulePauseStatus string

const (
	SchedulePauseStatus_Unpaused SchedulePauseStatus = "UNPAUSED"
	SchedulePauseStatus_Paused   SchedulePauseStatus = "PAUSED"
)

// String representation for [fmt.Print].
func (f *SchedulePauseStatus) String() string {
	return string(*f)
}

type AuthorizationDetails struct {
	Type                  *string                          `json:"type"`
	ResourceName          *string                          `json:"resource_name"`
	ResourceLegacyAclPath *string                          `json:"resource_legacy_acl_path"`
	GrantRules            []AuthorizationDetails_GrantRule `json:"grant_rules"`
}

type AuthorizationDetails_GrantRule struct {
	PermissionSet *string `json:"permission_set"`
}

type CreateDashboardRequest struct {
	Dashboard      *Dashboard `json:"dashboard"`
	DatasetCatalog *string    `json:"dataset_catalog"`
	DatasetSchema  *string    `json:"dataset_schema"`
}

type CreateScheduleRequest struct {
	Schedule *Schedule `json:"schedule"`
}

type CreateSubscriptionRequest struct {
	Subscription *Subscription `json:"subscription"`
}

type CronSchedule struct {
	QuartzCronExpression *string `json:"quartz_cron_expression"`
	TimezoneId           *string `json:"timezone_id"`
}

type Dashboard struct {
	DashboardId         *string         `json:"dashboard_id"`
	DisplayName         *string         `json:"display_name"`
	Path                *string         `json:"path"`
	CreateTime          *time.Time      `json:"create_time"`
	UpdateTime          *time.Time      `json:"update_time"`
	WarehouseId         *string         `json:"warehouse_id"`
	Etag                *string         `json:"etag"`
	SerializedDashboard *string         `json:"serialized_dashboard"`
	LifecycleState      *LifecycleState `json:"lifecycle_state"`
	ParentPath          *string         `json:"parent_path"`
}

type DeleteScheduleRequest struct {
	ScheduleId  *string `json:"schedule_id"`
	DashboardId *string `json:"dashboard_id"`
	Etag        *string `json:"etag"`
}

type DeleteSubscriptionRequest struct {
	SubscriptionId *string `json:"subscription_id"`
	ScheduleId     *string `json:"schedule_id"`
	DashboardId    *string `json:"dashboard_id"`
	Etag           *string `json:"etag"`
}

type GetDashboardRequest struct {
	DashboardId *string `json:"dashboard_id"`
}

type GetPublishedDashboardEmbeddedRequest struct {
	DashboardId *string `json:"dashboard_id"`
}

type GetPublishedDashboardEmbeddedResponse struct {
}

type GetPublishedDashboardRequest struct {
	DashboardId *string `json:"dashboard_id"`
}

type GetPublishedDashboardTokenInfoRequest struct {
	DashboardId      *string `json:"dashboard_id"`
	ExternalValue    *string `json:"external_value"`
	ExternalViewerId *string `json:"external_viewer_id"`
}

type GetPublishedDashboardTokenInfoResponse struct {
	CustomClaim          *string                `json:"custom_claim"`
	Scope                *string                `json:"scope"`
	AuthorizationDetails []AuthorizationDetails `json:"authorization_details"`
}

type GetScheduleRequest struct {
	ScheduleId  *string `json:"schedule_id"`
	DashboardId *string `json:"dashboard_id"`
}

type GetSubscriptionRequest struct {
	SubscriptionId *string `json:"subscription_id"`
	ScheduleId     *string `json:"schedule_id"`
	DashboardId    *string `json:"dashboard_id"`
}

type ListDashboardsRequest struct {
	PageSize    *int           `json:"page_size"`
	PageToken   *string        `json:"page_token"`
	ShowTrashed *bool          `json:"show_trashed"`
	View        *DashboardView `json:"view"`
}

type ListDashboardsResponse struct {
	Dashboards    []Dashboard `json:"dashboards"`
	NextPageToken *string     `json:"next_page_token"`
}

type ListSchedulesRequest struct {
	DashboardId *string `json:"dashboard_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

type ListSchedulesResponse struct {
	Schedules     []Schedule `json:"schedules"`
	NextPageToken *string    `json:"next_page_token"`
}

type ListSubscriptionsRequest struct {
	DashboardId *string `json:"dashboard_id"`
	ScheduleId  *string `json:"schedule_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
	NextPageToken *string        `json:"next_page_token"`
}

type MigrateDashboardRequest struct {
	SourceDashboardId     *string `json:"source_dashboard_id"`
	DisplayName           *string `json:"display_name"`
	ParentPath            *string `json:"parent_path"`
	UpdateParameterSyntax *bool   `json:"update_parameter_syntax"`
}

type PublishDashboardRequest struct {
	DashboardId      *string `json:"dashboard_id"`
	EmbedCredentials *bool   `json:"embed_credentials"`
	WarehouseId      *string `json:"warehouse_id"`
}

type PublishedDashboard struct {
	DisplayName        *string    `json:"display_name"`
	WarehouseId        *string    `json:"warehouse_id"`
	EmbedCredentials   *bool      `json:"embed_credentials"`
	RevisionCreateTime *time.Time `json:"revision_create_time"`
}

type Schedule struct {
	ScheduleId   *string              `json:"schedule_id"`
	DashboardId  *string              `json:"dashboard_id"`
	CronSchedule *CronSchedule        `json:"cron_schedule"`
	PauseStatus  *SchedulePauseStatus `json:"pause_status"`
	DisplayName  *string              `json:"display_name"`
	Etag         *string              `json:"etag"`
	CreateTime   *time.Time           `json:"create_time"`
	UpdateTime   *time.Time           `json:"update_time"`
	WarehouseId  *string              `json:"warehouse_id"`
}

type Subscription struct {
	SubscriptionId  *string                  `json:"subscription_id"`
	ScheduleId      *string                  `json:"schedule_id"`
	DashboardId     *string                  `json:"dashboard_id"`
	Subscriber      *Subscription_Subscriber `json:"subscriber"`
	CreatedByUserId *int64                   `json:"created_by_user_id"`
	Etag            *string                  `json:"etag"`
	CreateTime      *time.Time               `json:"create_time"`
	UpdateTime      *time.Time               `json:"update_time"`
	SkipNotify      *bool                    `json:"skip_notify"`
}

type Subscription_Subscriber struct {
	UserSubscriber        *Subscription_Subscriber_User        `json:"user_subscriber"`
	DestinationSubscriber *Subscription_Subscriber_Destination `json:"destination_subscriber"`
}

type Subscription_Subscriber_Destination struct {
	DestinationId *string `json:"destination_id"`
}

type Subscription_Subscriber_User struct {
	UserId *int64 `json:"user_id"`
}

type TrashDashboardRequest struct {
	DashboardId *string `json:"dashboard_id"`
}

type TrashDashboardResponse struct {
}

type UnpublishDashboardRequest struct {
	DashboardId *string `json:"dashboard_id"`
}

type UnpublishDashboardResponse struct {
}

type UpdateDashboardRequest struct {
	Dashboard      *Dashboard `json:"dashboard"`
	DatasetCatalog *string    `json:"dataset_catalog"`
	DatasetSchema  *string    `json:"dataset_schema"`
}

type UpdateScheduleRequest struct {
	Schedule *Schedule `json:"schedule"`
}
