package v2

type AwsAvailability string

const (
	AwsAvailability_Spot             AwsAvailability = "SPOT"
	AwsAvailability_OnDemand         AwsAvailability = "ON_DEMAND"
	AwsAvailability_SpotWithFallback AwsAvailability = "SPOT_WITH_FALLBACK"
)

// String representation for [fmt.Print].
func (f *AwsAvailability) String() string {
	return string(*f)
}

type AzureAvailability string

const (
	AzureAvailability_SpotAzure             AzureAvailability = "SPOT_AZURE"
	AzureAvailability_OnDemandAzure         AzureAvailability = "ON_DEMAND_AZURE"
	AzureAvailability_SpotWithFallbackAzure AzureAvailability = "SPOT_WITH_FALLBACK_AZURE"
)

// String representation for [fmt.Print].
func (f *AzureAvailability) String() string {
	return string(*f)
}

type AzureDiskVolumeType string

const (
	AzureDiskVolumeType_PremiumLrs  AzureDiskVolumeType = "PREMIUM_LRS"
	AzureDiskVolumeType_StandardLrs AzureDiskVolumeType = "STANDARD_LRS"
)

// String representation for [fmt.Print].
func (f *AzureDiskVolumeType) String() string {
	return string(*f)
}

type EbsVolumeType string

const (
	EbsVolumeType_GeneralPurposeSsd      EbsVolumeType = "GENERAL_PURPOSE_SSD"
	EbsVolumeType_ThroughputOptimizedHdd EbsVolumeType = "THROUGHPUT_OPTIMIZED_HDD"
)

// String representation for [fmt.Print].
func (f *EbsVolumeType) String() string {
	return string(*f)
}

type GcpAvailability string

const (
	GcpAvailability_PreemptibleGcp             GcpAvailability = "PREEMPTIBLE_GCP"
	GcpAvailability_OnDemandGcp                GcpAvailability = "ON_DEMAND_GCP"
	GcpAvailability_PreemptibleWithFallbackGcp GcpAvailability = "PREEMPTIBLE_WITH_FALLBACK_GCP"
)

// String representation for [fmt.Print].
func (f *GcpAvailability) String() string {
	return string(*f)
}

type InstancePoolState string

const (
	InstancePoolState_Active  InstancePoolState = "ACTIVE"
	InstancePoolState_Stopped InstancePoolState = "STOPPED"
	InstancePoolState_Deleted InstancePoolState = "DELETED"
)

// String representation for [fmt.Print].
func (f *InstancePoolState) String() string {
	return string(*f)
}

type CreateInstancePool struct {
	InstancePoolName                   *string                      `json:"instance_pool_name"`
	MinIdleInstances                   *int                         `json:"min_idle_instances"`
	MaxCapacity                        *int                         `json:"max_capacity"`
	AwsAttributes                      *InstancePoolAwsAttributes   `json:"aws_attributes"`
	NodeTypeId                         *string                      `json:"node_type_id"`
	CustomTags                         map[string]string            `json:"custom_tags"`
	IdleInstanceAutoterminationMinutes *int                         `json:"idle_instance_autotermination_minutes"`
	EnableElasticDisk                  *bool                        `json:"enable_elastic_disk"`
	DiskSpec                           *DiskSpec                    `json:"disk_spec"`
	PreloadedDockerImages              []DockerImage                `json:"preloaded_docker_images"`
	PreloadedSparkVersions             []string                     `json:"preloaded_spark_versions"`
	AzureAttributes                    *InstancePoolAzureAttributes `json:"azure_attributes"`
	GcpAttributes                      *InstancePoolGcpAttributes   `json:"gcp_attributes"`
	NodeTypeFlexibility                *NodeTypeFlexibility         `json:"node_type_flexibility"`
	EnableAutoAlternateNodeTypes       *bool                        `json:"enable_auto_alternate_node_types"`
	RemoteDiskThroughput               *int                         `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize         *int                         `json:"total_initial_remote_disk_size"`
}

type CreateInstancePool_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateInstancePool_Response struct {
	InstancePoolId *string `json:"instance_pool_id"`
}

type DeleteInstancePool struct {
	InstancePoolId *string `json:"instance_pool_id"`
}

type DeleteInstancePool_Response struct {
}

// Describes the disks that are launched for each instance in the spark cluster.
// For example, if the cluster has 3 instances, each instance is configured to
// launch 2 disks, 100 GiB each, then <Databricks> will launch a total of 6
// disks, 100 GiB each, for this cluster..
type DiskSpec struct {
	DiskType       *DiskType `json:"disk_type"`
	DiskCount      *int      `json:"disk_count"`
	DiskSize       *int      `json:"disk_size"`
	DiskIops       *int      `json:"disk_iops"`
	DiskThroughput *int      `json:"disk_throughput"`
}

// Describes the disk type..
type DiskType struct {
	EbsVolumeType       *EbsVolumeType       `json:"ebs_volume_type"`
	AzureDiskVolumeType *AzureDiskVolumeType `json:"azure_disk_volume_type"`
}

type DockerBasicAuth struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type DockerImage struct {
	Url       *string          `json:"url"`
	BasicAuth *DockerBasicAuth `json:"basic_auth"`
}

type EditInstancePool struct {
	InstancePoolId                     *string                      `json:"instance_pool_id"`
	InstancePoolName                   *string                      `json:"instance_pool_name"`
	MinIdleInstances                   *int                         `json:"min_idle_instances"`
	MaxCapacity                        *int                         `json:"max_capacity"`
	AwsAttributes                      *InstancePoolAwsAttributes   `json:"aws_attributes"`
	NodeTypeId                         *string                      `json:"node_type_id"`
	CustomTags                         map[string]string            `json:"custom_tags"`
	IdleInstanceAutoterminationMinutes *int                         `json:"idle_instance_autotermination_minutes"`
	EnableElasticDisk                  *bool                        `json:"enable_elastic_disk"`
	DiskSpec                           *DiskSpec                    `json:"disk_spec"`
	PreloadedDockerImages              []DockerImage                `json:"preloaded_docker_images"`
	PreloadedSparkVersions             []string                     `json:"preloaded_spark_versions"`
	AzureAttributes                    *InstancePoolAzureAttributes `json:"azure_attributes"`
	GcpAttributes                      *InstancePoolGcpAttributes   `json:"gcp_attributes"`
	NodeTypeFlexibility                *NodeTypeFlexibility         `json:"node_type_flexibility"`
	EnableAutoAlternateNodeTypes       *bool                        `json:"enable_auto_alternate_node_types"`
	RemoteDiskThroughput               *int                         `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize         *int                         `json:"total_initial_remote_disk_size"`
}

type EditInstancePool_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type EditInstancePool_Response struct {
}

type GetInstancePool struct {
	InstancePoolId *string `json:"instance_pool_id"`
}

type GetInstancePool_Response struct {
	Stats                              *InstancePoolStats           `json:"stats"`
	Status                             *InstancePoolStatus          `json:"status"`
	InstancePoolId                     *string                      `json:"instance_pool_id"`
	DefaultTags                        map[string]string            `json:"default_tags"`
	State                              *InstancePoolState           `json:"state"`
	InstancePoolName                   *string                      `json:"instance_pool_name"`
	MinIdleInstances                   *int                         `json:"min_idle_instances"`
	MaxCapacity                        *int                         `json:"max_capacity"`
	AwsAttributes                      *InstancePoolAwsAttributes   `json:"aws_attributes"`
	NodeTypeId                         *string                      `json:"node_type_id"`
	CustomTags                         map[string]string            `json:"custom_tags"`
	IdleInstanceAutoterminationMinutes *int                         `json:"idle_instance_autotermination_minutes"`
	EnableElasticDisk                  *bool                        `json:"enable_elastic_disk"`
	DiskSpec                           *DiskSpec                    `json:"disk_spec"`
	PreloadedDockerImages              []DockerImage                `json:"preloaded_docker_images"`
	PreloadedSparkVersions             []string                     `json:"preloaded_spark_versions"`
	AzureAttributes                    *InstancePoolAzureAttributes `json:"azure_attributes"`
	GcpAttributes                      *InstancePoolGcpAttributes   `json:"gcp_attributes"`
	NodeTypeFlexibility                *NodeTypeFlexibility         `json:"node_type_flexibility"`
	EnableAutoAlternateNodeTypes       *bool                        `json:"enable_auto_alternate_node_types"`
	RemoteDiskThroughput               *int                         `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize         *int                         `json:"total_initial_remote_disk_size"`
}

type GetInstancePool_Response_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type GetInstancePool_Response_DefaultTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type InstancePoolAndStats struct {
	Stats                              *InstancePoolStats           `json:"stats"`
	Status                             *InstancePoolStatus          `json:"status"`
	InstancePoolId                     *string                      `json:"instance_pool_id"`
	DefaultTags                        map[string]string            `json:"default_tags"`
	State                              *InstancePoolState           `json:"state"`
	InstancePoolName                   *string                      `json:"instance_pool_name"`
	MinIdleInstances                   *int                         `json:"min_idle_instances"`
	MaxCapacity                        *int                         `json:"max_capacity"`
	AwsAttributes                      *InstancePoolAwsAttributes   `json:"aws_attributes"`
	NodeTypeId                         *string                      `json:"node_type_id"`
	CustomTags                         map[string]string            `json:"custom_tags"`
	IdleInstanceAutoterminationMinutes *int                         `json:"idle_instance_autotermination_minutes"`
	EnableElasticDisk                  *bool                        `json:"enable_elastic_disk"`
	DiskSpec                           *DiskSpec                    `json:"disk_spec"`
	PreloadedDockerImages              []DockerImage                `json:"preloaded_docker_images"`
	PreloadedSparkVersions             []string                     `json:"preloaded_spark_versions"`
	AzureAttributes                    *InstancePoolAzureAttributes `json:"azure_attributes"`
	GcpAttributes                      *InstancePoolGcpAttributes   `json:"gcp_attributes"`
	NodeTypeFlexibility                *NodeTypeFlexibility         `json:"node_type_flexibility"`
	EnableAutoAlternateNodeTypes       *bool                        `json:"enable_auto_alternate_node_types"`
	RemoteDiskThroughput               *int                         `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize         *int                         `json:"total_initial_remote_disk_size"`
}

type InstancePoolAndStats_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type InstancePoolAndStats_DefaultTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Attributes set during instance pool creation which are related to Amazon Web
// Services..
type InstancePoolAwsAttributes struct {
	Availability        *AwsAvailability `json:"availability"`
	ZoneId              *string          `json:"zone_id"`
	SpotBidPricePercent *int             `json:"spot_bid_price_percent"`
	InstanceProfileArn  *string          `json:"instance_profile_arn"`
}

// Attributes set during instance pool creation which are related to Azure..
type InstancePoolAzureAttributes struct {
	Availability    *AzureAvailability `json:"availability"`
	SpotBidMaxPrice *float64           `json:"spot_bid_max_price"`
}

// Attributes set during instance pool creation which are related to GCP..
type InstancePoolGcpAttributes struct {
	GcpAvailability *GcpAvailability `json:"gcp_availability"`
	LocalSsdCount   *int             `json:"local_ssd_count"`
	ZoneId          *string          `json:"zone_id"`
}

type InstancePoolStats struct {
	UsedCount        *int `json:"used_count"`
	IdleCount        *int `json:"idle_count"`
	PendingUsedCount *int `json:"pending_used_count"`
	PendingIdleCount *int `json:"pending_idle_count"`
}

type InstancePoolStatus struct {
	PendingInstanceErrors []PendingInstanceError `json:"pending_instance_errors"`
}

type ListInstancePools struct {
}

type ListInstancePools_Response struct {
	InstancePools []InstancePoolAndStats `json:"instance_pools"`
}

// Configuration for flexible node types, allowing fallback to alternate node
// types during cluster launch and upscale..
type NodeTypeFlexibility struct {
	AlternateNodeTypeIds []string `json:"alternate_node_type_ids"`
}

// Error message of a failed pending instances.
type PendingInstanceError struct {
	InstanceId *string `json:"instance_id"`
	Message    *string `json:"message"`
}
