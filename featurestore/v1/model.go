package v1

import "time"

type OnlineStore_State string

const (
	OnlineStore_State_StateUnspecified OnlineStore_State = "STATE_UNSPECIFIED"
	OnlineStore_State_Starting         OnlineStore_State = "STARTING"
	OnlineStore_State_Available        OnlineStore_State = "AVAILABLE"
	OnlineStore_State_Deleting         OnlineStore_State = "DELETING"
	OnlineStore_State_Stopped          OnlineStore_State = "STOPPED"
	OnlineStore_State_Updating         OnlineStore_State = "UPDATING"
	OnlineStore_State_FailingOver      OnlineStore_State = "FAILING_OVER"
)

// String representation for [fmt.Print].
func (f *OnlineStore_State) String() string {
	return string(*f)
}

type PublishSpec_PublishMode string

const (
	PublishSpec_PublishMode_PublishModeUnspecified PublishSpec_PublishMode = "PUBLISH_MODE_UNSPECIFIED"
	PublishSpec_PublishMode_Continuous             PublishSpec_PublishMode = "CONTINUOUS"
	PublishSpec_PublishMode_Triggered              PublishSpec_PublishMode = "TRIGGERED"
	PublishSpec_PublishMode_Snapshot               PublishSpec_PublishMode = "SNAPSHOT"
)

// String representation for [fmt.Print].
func (f *PublishSpec_PublishMode) String() string {
	return string(*f)
}

type CreateOnlineStoreRequest struct {
	OnlineStore *OnlineStore `json:"online_store"`
}

type DeleteOnlineStoreRequest struct {
	Name *string `json:"name"`
}

type DeleteOnlineTableRequest struct {
	OnlineTableName *string `json:"online_table_name"`
}

type GetOnlineStoreRequest struct {
	Name *string `json:"name"`
}

type ListOnlineStoresRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListOnlineStoresResponse struct {
	OnlineStores  []OnlineStore `json:"online_stores"`
	NextPageToken *string       `json:"next_page_token"`
}

// An OnlineStore is a logical database instance that stores and serves features
// online..
type OnlineStore struct {
	Name             *string            `json:"name"`
	Creator          *string            `json:"creator"`
	CreationTime     *time.Time         `json:"creation_time"`
	State            *OnlineStore_State `json:"state"`
	Capacity         *string            `json:"capacity"`
	ReadReplicaCount *int               `json:"read_replica_count"`
	UsagePolicyId    *string            `json:"usage_policy_id"`
}

type PublishSpec struct {
	OnlineStore     *string                  `json:"online_store"`
	OnlineTableName *string                  `json:"online_table_name"`
	PublishMode     *PublishSpec_PublishMode `json:"publish_mode"`
}

type PublishTableRequest struct {
	SourceTableName *string      `json:"source_table_name"`
	PublishSpec     *PublishSpec `json:"publish_spec"`
}

type PublishTableResponse struct {
	OnlineTableName *string `json:"online_table_name"`
	PipelineId      *string `json:"pipeline_id"`
}

type UpdateOnlineStoreRequest struct {
	OnlineStore *OnlineStore `json:"online_store"`
	UpdateMask  *string      `json:"update_mask"`
}
