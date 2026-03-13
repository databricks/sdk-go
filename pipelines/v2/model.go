package v2

type CloneMode string

const (
	CloneMode_MigrateToUc CloneMode = "MIGRATE_TO_UC"
)

// String representation for [fmt.Print].
func (f *CloneMode) String() string {
	return string(*f)
}

type ConnectorType string

const (
	ConnectorType_ConnectorTypeUnspecified ConnectorType = "CONNECTOR_TYPE_UNSPECIFIED"
	ConnectorType_Cdc                      ConnectorType = "CDC"
	ConnectorType_QueryBased               ConnectorType = "QUERY_BASED"
)

// String representation for [fmt.Print].
func (f *ConnectorType) String() string {
	return string(*f)
}

type DayOfWeek string

const (
	DayOfWeek_DayOfWeekUnspecified DayOfWeek = "DAY_OF_WEEK_UNSPECIFIED"
	DayOfWeek_Monday               DayOfWeek = "MONDAY"
	DayOfWeek_Tuesday              DayOfWeek = "TUESDAY"
	DayOfWeek_Wednesday            DayOfWeek = "WEDNESDAY"
	DayOfWeek_Thursday             DayOfWeek = "THURSDAY"
	DayOfWeek_Friday               DayOfWeek = "FRIDAY"
	DayOfWeek_Saturday             DayOfWeek = "SATURDAY"
	DayOfWeek_Sunday               DayOfWeek = "SUNDAY"
)

// String representation for [fmt.Print].
func (f *DayOfWeek) String() string {
	return string(*f)
}

type DeploymentKind string

const (
	DeploymentKind_Bundle DeploymentKind = "BUNDLE"
)

// String representation for [fmt.Print].
func (f *DeploymentKind) String() string {
	return string(*f)
}

type EventLevel string

const (
	EventLevel_Info    EventLevel = "INFO"
	EventLevel_Warn    EventLevel = "WARN"
	EventLevel_Error   EventLevel = "ERROR"
	EventLevel_Metrics EventLevel = "METRICS"
)

// String representation for [fmt.Print].
func (f *EventLevel) String() string {
	return string(*f)
}

type IngestionSourceType string

const (
	IngestionSourceType_IngestionSourceTypeUnspecified IngestionSourceType = "INGESTION_SOURCE_TYPE_UNSPECIFIED"
	IngestionSourceType_Mysql                          IngestionSourceType = "MYSQL"
	IngestionSourceType_Postgresql                     IngestionSourceType = "POSTGRESQL"
	IngestionSourceType_Redshift                       IngestionSourceType = "REDSHIFT"
	IngestionSourceType_Sqldw                          IngestionSourceType = "SQLDW"
	IngestionSourceType_Sqlserver                      IngestionSourceType = "SQLSERVER"
	IngestionSourceType_Salesforce                     IngestionSourceType = "SALESFORCE"
	IngestionSourceType_Bigquery                       IngestionSourceType = "BIGQUERY"
	IngestionSourceType_Netsuite                       IngestionSourceType = "NETSUITE"
	IngestionSourceType_WorkdayRaas                    IngestionSourceType = "WORKDAY_RAAS"
	IngestionSourceType_Ga4RawData                     IngestionSourceType = "GA4_RAW_DATA"
	IngestionSourceType_Servicenow                     IngestionSourceType = "SERVICENOW"
	IngestionSourceType_ManagedPostgresql              IngestionSourceType = "MANAGED_POSTGRESQL"
	IngestionSourceType_Oracle                         IngestionSourceType = "ORACLE"
	IngestionSourceType_Teradata                       IngestionSourceType = "TERADATA"
	IngestionSourceType_Sharepoint                     IngestionSourceType = "SHAREPOINT"
	IngestionSourceType_Dynamics365                    IngestionSourceType = "DYNAMICS365"
	IngestionSourceType_GoogleDrive                    IngestionSourceType = "GOOGLE_DRIVE"
	IngestionSourceType_Confluence                     IngestionSourceType = "CONFLUENCE"
	IngestionSourceType_MetaMarketing                  IngestionSourceType = "META_MARKETING"
	IngestionSourceType_GoogleAds                      IngestionSourceType = "GOOGLE_ADS"
	IngestionSourceType_TiktokAds                      IngestionSourceType = "TIKTOK_ADS"
	IngestionSourceType_SalesforceMarketingCloud       IngestionSourceType = "SALESFORCE_MARKETING_CLOUD"
	IngestionSourceType_Hubspot                        IngestionSourceType = "HUBSPOT"
	IngestionSourceType_WorkdayHcm                     IngestionSourceType = "WORKDAY_HCM"
	IngestionSourceType_Guidewire                      IngestionSourceType = "GUIDEWIRE"
	IngestionSourceType_Zendesk                        IngestionSourceType = "ZENDESK"
	IngestionSourceType_SlackAuditLogs                 IngestionSourceType = "SLACK_AUDIT_LOGS"
	IngestionSourceType_CrowdstrikeEventStream         IngestionSourceType = "CROWDSTRIKE_EVENT_STREAM"
	IngestionSourceType_WorkdayActivityLogging         IngestionSourceType = "WORKDAY_ACTIVITY_LOGGING"
	IngestionSourceType_AkamaiWaf                      IngestionSourceType = "AKAMAI_WAF"
	IngestionSourceType_Veeva                          IngestionSourceType = "VEEVA"
	IngestionSourceType_VeevaVault                     IngestionSourceType = "VEEVA_VAULT"
	IngestionSourceType_M365AuditLogs                  IngestionSourceType = "M365_AUDIT_LOGS"
	IngestionSourceType_OktaSystemLogs                 IngestionSourceType = "OKTA_SYSTEM_LOGS"
	IngestionSourceType_OnePasswordEventLogs           IngestionSourceType = "ONE_PASSWORD_EVENT_LOGS"
	IngestionSourceType_ProofpointSiem                 IngestionSourceType = "PROOFPOINT_SIEM"
	IngestionSourceType_WizAuditLogs                   IngestionSourceType = "WIZ_AUDIT_LOGS"
	IngestionSourceType_Github                         IngestionSourceType = "GITHUB"
	IngestionSourceType_Outlook                        IngestionSourceType = "OUTLOOK"
	IngestionSourceType_Smartsheet                     IngestionSourceType = "SMARTSHEET"
	IngestionSourceType_MicrosoftTeams                 IngestionSourceType = "MICROSOFT_TEAMS"
	IngestionSourceType_AdobeCampaigns                 IngestionSourceType = "ADOBE_CAMPAIGNS"
	IngestionSourceType_LinkedinAds                    IngestionSourceType = "LINKEDIN_ADS"
	IngestionSourceType_XAds                           IngestionSourceType = "X_ADS"
	IngestionSourceType_BingAds                        IngestionSourceType = "BING_ADS"
	IngestionSourceType_GoogleSearchConsole            IngestionSourceType = "GOOGLE_SEARCH_CONSOLE"
	IngestionSourceType_PinterestAds                   IngestionSourceType = "PINTEREST_ADS"
	IngestionSourceType_RedditAds                      IngestionSourceType = "REDDIT_ADS"
	IngestionSourceType_ForeignCatalog                 IngestionSourceType = "FOREIGN_CATALOG"
)

// String representation for [fmt.Print].
func (f *IngestionSourceType) String() string {
	return string(*f)
}

type MaturityLevel string

const (
	MaturityLevel_Stable     MaturityLevel = "STABLE"
	MaturityLevel_Evolving   MaturityLevel = "EVOLVING"
	MaturityLevel_Deprecated MaturityLevel = "DEPRECATED"
)

// String representation for [fmt.Print].
func (f *MaturityLevel) String() string {
	return string(*f)
}

type PipelineHealthStatus string

const (
	PipelineHealthStatus_Healthy   PipelineHealthStatus = "HEALTHY"
	PipelineHealthStatus_Unhealthy PipelineHealthStatus = "UNHEALTHY"
)

// String representation for [fmt.Print].
func (f *PipelineHealthStatus) String() string {
	return string(*f)
}

type PipelineState string

const (
	PipelineState_Deploying  PipelineState = "DEPLOYING"
	PipelineState_Starting   PipelineState = "STARTING"
	PipelineState_Running    PipelineState = "RUNNING"
	PipelineState_Stopping   PipelineState = "STOPPING"
	PipelineState_Deleted    PipelineState = "DELETED"
	PipelineState_Recovering PipelineState = "RECOVERING"
	PipelineState_Failed     PipelineState = "FAILED"
	PipelineState_Resetting  PipelineState = "RESETTING"
	PipelineState_Idle       PipelineState = "IDLE"
)

// String representation for [fmt.Print].
func (f *PipelineState) String() string {
	return string(*f)
}

type PipelinesAwsAvailability string

const (
	PipelinesAwsAvailability_Spot             PipelinesAwsAvailability = "SPOT"
	PipelinesAwsAvailability_OnDemand         PipelinesAwsAvailability = "ON_DEMAND"
	PipelinesAwsAvailability_SpotWithFallback PipelinesAwsAvailability = "SPOT_WITH_FALLBACK"
)

// String representation for [fmt.Print].
func (f *PipelinesAwsAvailability) String() string {
	return string(*f)
}

type PipelinesAzureAvailability string

const (
	PipelinesAzureAvailability_SpotAzure             PipelinesAzureAvailability = "SPOT_AZURE"
	PipelinesAzureAvailability_OnDemandAzure         PipelinesAzureAvailability = "ON_DEMAND_AZURE"
	PipelinesAzureAvailability_SpotWithFallbackAzure PipelinesAzureAvailability = "SPOT_WITH_FALLBACK_AZURE"
)

// String representation for [fmt.Print].
func (f *PipelinesAzureAvailability) String() string {
	return string(*f)
}

type PipelinesEbsVolumeType string

const (
	PipelinesEbsVolumeType_GeneralPurposeSsd      PipelinesEbsVolumeType = "GENERAL_PURPOSE_SSD"
	PipelinesEbsVolumeType_ThroughputOptimizedHdd PipelinesEbsVolumeType = "THROUGHPUT_OPTIMIZED_HDD"
)

// String representation for [fmt.Print].
func (f *PipelinesEbsVolumeType) String() string {
	return string(*f)
}

type PipelinesGcpAvailability string

const (
	PipelinesGcpAvailability_PreemptibleGcp             PipelinesGcpAvailability = "PREEMPTIBLE_GCP"
	PipelinesGcpAvailability_OnDemandGcp                PipelinesGcpAvailability = "ON_DEMAND_GCP"
	PipelinesGcpAvailability_PreemptibleWithFallbackGcp PipelinesGcpAvailability = "PREEMPTIBLE_WITH_FALLBACK_GCP"
)

// String representation for [fmt.Print].
func (f *PipelinesGcpAvailability) String() string {
	return string(*f)
}

type PublishingMode string

const (
	PublishingMode_PublishingModeUnspecified PublishingMode = "PUBLISHING_MODE_UNSPECIFIED"
	PublishingMode_LegacyPublishingMode      PublishingMode = "LEGACY_PUBLISHING_MODE"
	PublishingMode_DefaultPublishingMode     PublishingMode = "DEFAULT_PUBLISHING_MODE"
)

// String representation for [fmt.Print].
func (f *PublishingMode) String() string {
	return string(*f)
}

type ScdType string

const (
	ScdType_ScdTypeUnspecified ScdType = "SCD_TYPE_UNSPECIFIED"
	ScdType_ScdType1           ScdType = "SCD_TYPE_1"
	ScdType_ScdType2           ScdType = "SCD_TYPE_2"
	ScdType_AppendOnly         ScdType = "APPEND_ONLY"
)

// String representation for [fmt.Print].
func (f *ScdType) String() string {
	return string(*f)
}

type UpdateCause string

const (
	UpdateCause_ApiCall                   UpdateCause = "API_CALL"
	UpdateCause_RetryOnFailure            UpdateCause = "RETRY_ON_FAILURE"
	UpdateCause_ServiceUpgrade            UpdateCause = "SERVICE_UPGRADE"
	UpdateCause_SchemaChange              UpdateCause = "SCHEMA_CHANGE"
	UpdateCause_JobTask                   UpdateCause = "JOB_TASK"
	UpdateCause_UserAction                UpdateCause = "USER_ACTION"
	UpdateCause_InfrastructureMaintenance UpdateCause = "INFRASTRUCTURE_MAINTENANCE"
)

// String representation for [fmt.Print].
func (f *UpdateCause) String() string {
	return string(*f)
}

type UpdateMode string

const (
	UpdateMode_Default    UpdateMode = "DEFAULT"
	UpdateMode_Continuous UpdateMode = "CONTINUOUS"
)

// String representation for [fmt.Print].
func (f *UpdateMode) String() string {
	return string(*f)
}

type FileIngestionOptions_FileFormat string

const (
	FileIngestionOptions_FileFormat_FileFormatUnspecified FileIngestionOptions_FileFormat = "FILE_FORMAT_UNSPECIFIED"
	FileIngestionOptions_FileFormat_Binaryfile            FileIngestionOptions_FileFormat = "BINARYFILE"
	FileIngestionOptions_FileFormat_Json                  FileIngestionOptions_FileFormat = "JSON"
	FileIngestionOptions_FileFormat_Csv                   FileIngestionOptions_FileFormat = "CSV"
	FileIngestionOptions_FileFormat_Xml                   FileIngestionOptions_FileFormat = "XML"
	FileIngestionOptions_FileFormat_Excel                 FileIngestionOptions_FileFormat = "EXCEL"
	FileIngestionOptions_FileFormat_Parquet               FileIngestionOptions_FileFormat = "PARQUET"
	FileIngestionOptions_FileFormat_Avro                  FileIngestionOptions_FileFormat = "AVRO"
	FileIngestionOptions_FileFormat_Orc                   FileIngestionOptions_FileFormat = "ORC"
)

// String representation for [fmt.Print].
func (f *FileIngestionOptions_FileFormat) String() string {
	return string(*f)
}

type FileIngestionOptions_SchemaEvolutionMode string

const (
	FileIngestionOptions_SchemaEvolutionMode_SchemaEvolutionModeUnspecified FileIngestionOptions_SchemaEvolutionMode = "SCHEMA_EVOLUTION_MODE_UNSPECIFIED"
	FileIngestionOptions_SchemaEvolutionMode_AddNewColumnsWithTypeWidening  FileIngestionOptions_SchemaEvolutionMode = "ADD_NEW_COLUMNS_WITH_TYPE_WIDENING"
	FileIngestionOptions_SchemaEvolutionMode_AddNewColumns                  FileIngestionOptions_SchemaEvolutionMode = "ADD_NEW_COLUMNS"
	FileIngestionOptions_SchemaEvolutionMode_Rescue                         FileIngestionOptions_SchemaEvolutionMode = "RESCUE"
	FileIngestionOptions_SchemaEvolutionMode_FailOnNewColumns               FileIngestionOptions_SchemaEvolutionMode = "FAIL_ON_NEW_COLUMNS"
	FileIngestionOptions_SchemaEvolutionMode_None                           FileIngestionOptions_SchemaEvolutionMode = "NONE"
)

// String representation for [fmt.Print].
func (f *FileIngestionOptions_SchemaEvolutionMode) String() string {
	return string(*f)
}

type GoogleDriveOptions_GoogleDriveEntityType string

const (
	GoogleDriveOptions_GoogleDriveEntityType_GoogleDriveEntityTypeUnspecified GoogleDriveOptions_GoogleDriveEntityType = "GOOGLE_DRIVE_ENTITY_TYPE_UNSPECIFIED"
	GoogleDriveOptions_GoogleDriveEntityType_File                             GoogleDriveOptions_GoogleDriveEntityType = "FILE"
	GoogleDriveOptions_GoogleDriveEntityType_FileMetadata                     GoogleDriveOptions_GoogleDriveEntityType = "FILE_METADATA"
	GoogleDriveOptions_GoogleDriveEntityType_Permission                       GoogleDriveOptions_GoogleDriveEntityType = "PERMISSION"
)

// String representation for [fmt.Print].
func (f *GoogleDriveOptions_GoogleDriveEntityType) String() string {
	return string(*f)
}

type SharepointOptions_SharepointEntityType string

const (
	SharepointOptions_SharepointEntityType_SharepointEntityTypeUnspecified SharepointOptions_SharepointEntityType = "SHAREPOINT_ENTITY_TYPE_UNSPECIFIED"
	SharepointOptions_SharepointEntityType_File                            SharepointOptions_SharepointEntityType = "FILE"
	SharepointOptions_SharepointEntityType_FileMetadata                    SharepointOptions_SharepointEntityType = "FILE_METADATA"
	SharepointOptions_SharepointEntityType_Permission                      SharepointOptions_SharepointEntityType = "PERMISSION"
	SharepointOptions_SharepointEntityType_List                            SharepointOptions_SharepointEntityType = "LIST"
)

// String representation for [fmt.Print].
func (f *SharepointOptions_SharepointEntityType) String() string {
	return string(*f)
}

type TikTokAdsOptions_TikTokDataLevel string

const (
	TikTokAdsOptions_TikTokDataLevel_TikTokDataLevelUnspecified TikTokAdsOptions_TikTokDataLevel = "TIK_TOK_DATA_LEVEL_UNSPECIFIED"
	TikTokAdsOptions_TikTokDataLevel_AuctionAdvertiser          TikTokAdsOptions_TikTokDataLevel = "AUCTION_ADVERTISER"
	TikTokAdsOptions_TikTokDataLevel_AuctionCampaign            TikTokAdsOptions_TikTokDataLevel = "AUCTION_CAMPAIGN"
	TikTokAdsOptions_TikTokDataLevel_AuctionAdgroup             TikTokAdsOptions_TikTokDataLevel = "AUCTION_ADGROUP"
	TikTokAdsOptions_TikTokDataLevel_AuctionAd                  TikTokAdsOptions_TikTokDataLevel = "AUCTION_AD"
)

// String representation for [fmt.Print].
func (f *TikTokAdsOptions_TikTokDataLevel) String() string {
	return string(*f)
}

type TikTokAdsOptions_TikTokReportType string

const (
	TikTokAdsOptions_TikTokReportType_TikTokReportTypeUnspecified TikTokAdsOptions_TikTokReportType = "TIK_TOK_REPORT_TYPE_UNSPECIFIED"
	TikTokAdsOptions_TikTokReportType_Basic                       TikTokAdsOptions_TikTokReportType = "BASIC"
	TikTokAdsOptions_TikTokReportType_Audience                    TikTokAdsOptions_TikTokReportType = "AUDIENCE"
	TikTokAdsOptions_TikTokReportType_PlayableAd                  TikTokAdsOptions_TikTokReportType = "PLAYABLE_AD"
	TikTokAdsOptions_TikTokReportType_Dsa                         TikTokAdsOptions_TikTokReportType = "DSA"
	TikTokAdsOptions_TikTokReportType_BusinessCenter              TikTokAdsOptions_TikTokReportType = "BUSINESS_CENTER"
	TikTokAdsOptions_TikTokReportType_GmvMax                      TikTokAdsOptions_TikTokReportType = "GMV_MAX"
)

// String representation for [fmt.Print].
func (f *TikTokAdsOptions_TikTokReportType) String() string {
	return string(*f)
}

type UpdateState_UpdateState string

const (
	UpdateState_UpdateState_Queued              UpdateState_UpdateState = "QUEUED"
	UpdateState_UpdateState_Created             UpdateState_UpdateState = "CREATED"
	UpdateState_UpdateState_WaitingForResources UpdateState_UpdateState = "WAITING_FOR_RESOURCES"
	UpdateState_UpdateState_Initializing        UpdateState_UpdateState = "INITIALIZING"
	UpdateState_UpdateState_Resetting           UpdateState_UpdateState = "RESETTING"
	UpdateState_UpdateState_SettingUpTables     UpdateState_UpdateState = "SETTING_UP_TABLES"
	UpdateState_UpdateState_Running             UpdateState_UpdateState = "RUNNING"
	UpdateState_UpdateState_Stopping            UpdateState_UpdateState = "STOPPING"
	UpdateState_UpdateState_Completed           UpdateState_UpdateState = "COMPLETED"
	UpdateState_UpdateState_Failed              UpdateState_UpdateState = "FAILED"
	UpdateState_UpdateState_Canceled            UpdateState_UpdateState = "CANCELED"
)

// String representation for [fmt.Print].
func (f *UpdateState_UpdateState) String() string {
	return string(*f)
}

// Akamai specific options for ingestion.
type AkamaiOptions struct {
	ConfigIds []int64 `json:"config_ids"`
}

type ApplyEnvironmentRequest struct {
	PipelineId *string `json:"pipeline_id"`
}

type ApplyEnvironmentRequest_Response struct {
}

// Policy for auto full refresh..
type AutoFullRefreshPolicy struct {
	Enabled          *bool `json:"enabled"`
	MinIntervalHours *int  `json:"min_interval_hours"`
}

type ClonePipeline struct {
	PipelineId           *string                             `json:"pipeline_id"`
	ExpectedLastModified *int64                              `json:"expected_last_modified"`
	AllowDuplicateNames  *bool                               `json:"allow_duplicate_names"`
	Id                   *string                             `json:"id"`
	Name                 *string                             `json:"name"`
	Storage              *string                             `json:"storage"`
	Configuration        map[string]string                   `json:"configuration"`
	Clusters             []PipelineCluster                   `json:"clusters"`
	Libraries            []PipelineLibrary                   `json:"libraries"`
	IngestionDefinition  *IngestionPipelineDefinition        `json:"ingestion_definition"`
	GatewayDefinition    *IngestionGatewayPipelineDefinition `json:"gateway_definition"`
	Trigger              *PipelineTrigger                    `json:"trigger"`
	Target               *string                             `json:"target"`
	Schema               *string                             `json:"schema"`
	Filters              *Filters                            `json:"filters"`
	Continuous           *bool                               `json:"continuous"`
	Development          *bool                               `json:"development"`
	Photon               *bool                               `json:"photon"`
	Edition              *string                             `json:"edition"`
	Channel              *string                             `json:"channel"`
	Catalog              *string                             `json:"catalog"`
	Notifications        []Notifications                     `json:"notifications"`
	Serverless           *bool                               `json:"serverless"`
	Deployment           *PipelineDeployment                 `json:"deployment"`
	RestartWindow        *RestartWindow                      `json:"restart_window"`
	BudgetPolicyId       *string                             `json:"budget_policy_id"`
	Tags                 map[string]string                   `json:"tags"`
	EventLog             *EventLogSpec                       `json:"event_log"`
	RootPath             *string                             `json:"root_path"`
	Environment          *PipelinesEnvironment               `json:"environment"`
	UsagePolicyId        *string                             `json:"usage_policy_id"`
	CloneMode            *CloneMode                          `json:"clone_mode"`
}

// Key value pair used to specify configuration parameters to Execution.
type ClonePipeline_ConfigurationEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ClonePipeline_Response struct {
	PipelineId *string `json:"pipeline_id"`
}

// A key-value entry that defines a single pipeline tags..
type ClonePipeline_TagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Confluence specific options for ingestion.
type ConfluenceConnectorOptions struct {
	IncludeConfluenceSpaces []string `json:"include_confluence_spaces"`
}

type ConnectionParameters struct {
	SourceCatalog *string `json:"source_catalog"`
}

// Wrapper message for source-specific options to support multiple connector
// types.
type ConnectorOptions struct {
	GoogleAdsOptions     *GoogleAdsOptions           `json:"google_ads_options"`
	TiktokAdsOptions     *TikTokAdsOptions           `json:"tiktok_ads_options"`
	SharepointOptions    *SharepointOptions          `json:"sharepoint_options"`
	GdriveOptions        *GoogleDriveOptions         `json:"gdrive_options"`
	OutlookOptions       *OutlookOptions             `json:"outlook_options"`
	AkamaiOptions        *AkamaiOptions              `json:"akamai_options"`
	SmartsheetOptions    *SmartsheetOptions          `json:"smartsheet_options"`
	JiraOptions          *JiraConnectorOptions       `json:"jira_options"`
	ConfluenceOptions    *ConfluenceConnectorOptions `json:"confluence_options"`
	MetaMarketingOptions *MetaMarketingOptions       `json:"meta_marketing_options"`
}

type CreatePipeline struct {
	AllowDuplicateNames *bool                               `json:"allow_duplicate_names"`
	DryRun              *bool                               `json:"dry_run"`
	RunAs               *PipelinesJobRunAs                  `json:"run_as"`
	Id                  *string                             `json:"id"`
	Name                *string                             `json:"name"`
	Storage             *string                             `json:"storage"`
	Configuration       map[string]string                   `json:"configuration"`
	Clusters            []PipelineCluster                   `json:"clusters"`
	Libraries           []PipelineLibrary                   `json:"libraries"`
	IngestionDefinition *IngestionPipelineDefinition        `json:"ingestion_definition"`
	GatewayDefinition   *IngestionGatewayPipelineDefinition `json:"gateway_definition"`
	Trigger             *PipelineTrigger                    `json:"trigger"`
	Target              *string                             `json:"target"`
	Schema              *string                             `json:"schema"`
	Filters             *Filters                            `json:"filters"`
	Continuous          *bool                               `json:"continuous"`
	Development         *bool                               `json:"development"`
	Photon              *bool                               `json:"photon"`
	Edition             *string                             `json:"edition"`
	Channel             *string                             `json:"channel"`
	Catalog             *string                             `json:"catalog"`
	Notifications       []Notifications                     `json:"notifications"`
	Serverless          *bool                               `json:"serverless"`
	Deployment          *PipelineDeployment                 `json:"deployment"`
	RestartWindow       *RestartWindow                      `json:"restart_window"`
	BudgetPolicyId      *string                             `json:"budget_policy_id"`
	Tags                map[string]string                   `json:"tags"`
	EventLog            *EventLogSpec                       `json:"event_log"`
	RootPath            *string                             `json:"root_path"`
	Environment         *PipelinesEnvironment               `json:"environment"`
	UsagePolicyId       *string                             `json:"usage_policy_id"`
}

// Key value pair used to specify configuration parameters to Execution.
type CreatePipeline_ConfigurationEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreatePipeline_Response struct {
	PipelineId        *string       `json:"pipeline_id"`
	EffectiveSettings *PipelineSpec `json:"effective_settings"`
}

// A key-value entry that defines a single pipeline tags..
type CreatePipeline_TagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CronTrigger struct {
	QuartzCronSchedule *string `json:"quartz_cron_schedule"`
	TimezoneId         *string `json:"timezone_id"`
}

type DataPlaneId struct {
	Instance *string `json:"instance"`
	SeqNo    *int64  `json:"seq_no"`
}

// Location of staged data storage.
type DataStagingOptions struct {
	CatalogName *string `json:"catalog_name"`
	SchemaName  *string `json:"schema_name"`
	VolumeName  *string `json:"volume_name"`
}

type DeletePipeline struct {
	PipelineId     *string `json:"pipeline_id"`
	Force          *bool   `json:"force"`
	DeleteDatasets *bool   `json:"delete_datasets"`
	Cascade        *bool   `json:"cascade"`
}

type DeletePipeline_Response struct {
}

type EditPipeline struct {
	PipelineId           *string                             `json:"pipeline_id"`
	AllowDuplicateNames  *bool                               `json:"allow_duplicate_names"`
	ExpectedLastModified *int64                              `json:"expected_last_modified"`
	RunAs                *PipelinesJobRunAs                  `json:"run_as"`
	Id                   *string                             `json:"id"`
	Name                 *string                             `json:"name"`
	Storage              *string                             `json:"storage"`
	Configuration        map[string]string                   `json:"configuration"`
	Clusters             []PipelineCluster                   `json:"clusters"`
	Libraries            []PipelineLibrary                   `json:"libraries"`
	IngestionDefinition  *IngestionPipelineDefinition        `json:"ingestion_definition"`
	GatewayDefinition    *IngestionGatewayPipelineDefinition `json:"gateway_definition"`
	Trigger              *PipelineTrigger                    `json:"trigger"`
	Target               *string                             `json:"target"`
	Schema               *string                             `json:"schema"`
	Filters              *Filters                            `json:"filters"`
	Continuous           *bool                               `json:"continuous"`
	Development          *bool                               `json:"development"`
	Photon               *bool                               `json:"photon"`
	Edition              *string                             `json:"edition"`
	Channel              *string                             `json:"channel"`
	Catalog              *string                             `json:"catalog"`
	Notifications        []Notifications                     `json:"notifications"`
	Serverless           *bool                               `json:"serverless"`
	Deployment           *PipelineDeployment                 `json:"deployment"`
	RestartWindow        *RestartWindow                      `json:"restart_window"`
	BudgetPolicyId       *string                             `json:"budget_policy_id"`
	Tags                 map[string]string                   `json:"tags"`
	EventLog             *EventLogSpec                       `json:"event_log"`
	RootPath             *string                             `json:"root_path"`
	Environment          *PipelinesEnvironment               `json:"environment"`
	UsagePolicyId        *string                             `json:"usage_policy_id"`
}

// Key value pair used to specify configuration parameters to Execution.
type EditPipeline_ConfigurationEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type EditPipeline_Response struct {
}

// A key-value entry that defines a single pipeline tags..
type EditPipeline_TagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ErrorDetail struct {
	Exceptions []SerializedException `json:"exceptions"`
	Fatal      *bool                 `json:"fatal"`
}

// Configurable event log parameters..
type EventLogSpec struct {
	Name    *string `json:"name"`
	Schema  *string `json:"schema"`
	Catalog *string `json:"catalog"`
}

type FileFilter struct {
	PathFilter     *string `json:"path_filter"`
	ModifiedBefore *string `json:"modified_before"`
	ModifiedAfter  *string `json:"modified_after"`
}

type FileIngestionOptions struct {
	Format              *FileIngestionOptions_FileFormat          `json:"format"`
	FileFilters         []FileFilter                              `json:"file_filters"`
	InferColumnTypes    *bool                                     `json:"infer_column_types"`
	SchemaEvolutionMode *FileIngestionOptions_SchemaEvolutionMode `json:"schema_evolution_mode"`
	SchemaHints         *string                                   `json:"schema_hints"`
	IgnoreCorruptFiles  *bool                                     `json:"ignore_corrupt_files"`
	CorruptRecordColumn *string                                   `json:"corrupt_record_column"`
	RescuedDataColumn   *string                                   `json:"rescued_data_column"`
	SingleVariantColumn *string                                   `json:"single_variant_column"`
	ReaderCaseSensitive *bool                                     `json:"reader_case_sensitive"`
	FormatOptions       map[string]string                         `json:"format_options"`
}

// Key value pair used to specify configuration parameters Based on public doc
// https://docs.databricks.com/en/ingestion/auto-loader/options.html.
type FileIngestionOptions_FormatOptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type Filters struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type GetPipeline struct {
	PipelineId *string `json:"pipeline_id"`
}

type GetPipeline_Response struct {
	PipelineId              *string               `json:"pipeline_id"`
	Spec                    *PipelineSpec         `json:"spec"`
	State                   *PipelineState        `json:"state"`
	Cause                   *string               `json:"cause"`
	ClusterId               *string               `json:"cluster_id"`
	Name                    *string               `json:"name"`
	Health                  *PipelineHealthStatus `json:"health"`
	CreatorUserName         *string               `json:"creator_user_name"`
	LatestUpdates           []UpdateStateInfo     `json:"latest_updates"`
	LastModified            *int64                `json:"last_modified"`
	RunAsUserName           *string               `json:"run_as_user_name"`
	EffectiveBudgetPolicyId *string               `json:"effective_budget_policy_id"`
	EffectivePublishingMode *PublishingMode       `json:"effective_publishing_mode"`
	RunAs                   *PipelinesJobRunAs    `json:"run_as"`
	EffectiveUsagePolicyId  *string               `json:"effective_usage_policy_id"`
}

type GetUpdate struct {
	PipelineId *string `json:"pipeline_id"`
	UpdateId   *string `json:"update_id"`
}

type GetUpdate_Response struct {
	Update *UpdateInfo `json:"update"`
}

// Google Ads specific options for ingestion.
type GoogleAdsOptions struct {
	ManagerAccountId   *string `json:"manager_account_id"`
	LookbackWindowDays *int    `json:"lookback_window_days"`
	SyncStartDate      *string `json:"sync_start_date"`
}

type GoogleDriveOptions struct {
	Url                  *string                                   `json:"url"`
	EntityType           *GoogleDriveOptions_GoogleDriveEntityType `json:"entity_type"`
	FileIngestionOptions *FileIngestionOptions                     `json:"file_ingestion_options"`
}

type IngestionGatewayPipelineDefinition struct {
	ConnectionName        *string               `json:"connection_name"`
	ConnectionId          *string               `json:"connection_id"`
	GatewayStorageCatalog *string               `json:"gateway_storage_catalog"`
	GatewayStorageSchema  *string               `json:"gateway_storage_schema"`
	GatewayStorageName    *string               `json:"gateway_storage_name"`
	ConnectionParameters  *ConnectionParameters `json:"connection_parameters"`
}

type IngestionPipelineDefinition struct {
	ConnectionName             *string                                          `json:"connection_name"`
	IngestionGatewayId         *string                                          `json:"ingestion_gateway_id"`
	IngestFromUcForeignCatalog *bool                                            `json:"ingest_from_uc_foreign_catalog"`
	Objects                    []IngestionPipelineDefinition_IngestionConfig    `json:"objects"`
	SourceType                 *IngestionSourceType                             `json:"source_type"`
	TableConfiguration         *IngestionPipelineDefinition_TableSpecificConfig `json:"table_configuration"`
	NetsuiteJarPath            *string                                          `json:"netsuite_jar_path"`
	SourceConfigurations       []SourceConfig                                   `json:"source_configurations"`
	FullRefreshWindow          *OperationTimeWindow                             `json:"full_refresh_window"`
	ConnectorType              *ConnectorType                                   `json:"connector_type"`
	DataStagingOptions         *DataStagingOptions                              `json:"data_staging_options"`
}

type IngestionPipelineDefinition_ConfluenceOptions struct {
	IncludeConfluenceSpaces []string `json:"include_confluence_spaces"`
}

type IngestionPipelineDefinition_IngestionConfig struct {
	Schema *IngestionPipelineDefinition_SchemaSpec `json:"schema"`
	Table  *IngestionPipelineDefinition_TableSpec  `json:"table"`
	Report *IngestionPipelineDefinition_ReportSpec `json:"report"`
}

type IngestionPipelineDefinition_JiraOptions struct {
	IncludeJiraSpaces []string `json:"include_jira_spaces"`
}

type IngestionPipelineDefinition_ReportSpec struct {
	SourceUrl          *string                                          `json:"source_url"`
	DestinationCatalog *string                                          `json:"destination_catalog"`
	DestinationSchema  *string                                          `json:"destination_schema"`
	DestinationTable   *string                                          `json:"destination_table"`
	TableConfiguration *IngestionPipelineDefinition_TableSpecificConfig `json:"table_configuration"`
}

type IngestionPipelineDefinition_SchemaSpec struct {
	SourceCatalog      *string                                          `json:"source_catalog"`
	SourceSchema       *string                                          `json:"source_schema"`
	DestinationCatalog *string                                          `json:"destination_catalog"`
	DestinationSchema  *string                                          `json:"destination_schema"`
	TableConfiguration *IngestionPipelineDefinition_TableSpecificConfig `json:"table_configuration"`
	JiraOptions        *IngestionPipelineDefinition_JiraOptions         `json:"jira_options"`
	ConfluenceOptions  *IngestionPipelineDefinition_ConfluenceOptions   `json:"confluence_options"`
	ConnectorOptions   *ConnectorOptions                                `json:"connector_options"`
}

type IngestionPipelineDefinition_TableSpec struct {
	SourceCatalog      *string                                          `json:"source_catalog"`
	SourceSchema       *string                                          `json:"source_schema"`
	SourceTable        *string                                          `json:"source_table"`
	DestinationCatalog *string                                          `json:"destination_catalog"`
	DestinationSchema  *string                                          `json:"destination_schema"`
	DestinationTable   *string                                          `json:"destination_table"`
	TableConfiguration *IngestionPipelineDefinition_TableSpecificConfig `json:"table_configuration"`
	JiraOptions        *IngestionPipelineDefinition_JiraOptions         `json:"jira_options"`
	ConfluenceOptions  *IngestionPipelineDefinition_ConfluenceOptions   `json:"confluence_options"`
	ConnectorOptions   *ConnectorOptions                                `json:"connector_options"`
}

type IngestionPipelineDefinition_TableSpecificConfig struct {
	ScdType                        *ScdType                                                                   `json:"scd_type"`
	PrimaryKeys                    []string                                                                   `json:"primary_keys"`
	SequenceBy                     []string                                                                   `json:"sequence_by"`
	IncludeColumns                 []string                                                                   `json:"include_columns"`
	ExcludeColumns                 []string                                                                   `json:"exclude_columns"`
	SalesforceIncludeFormulaFields *bool                                                                      `json:"salesforce_include_formula_fields"`
	WorkdayReportParameters        *IngestionPipelineDefinition_WorkdayReportParameters                       `json:"workday_report_parameters"`
	RowFilter                      *string                                                                    `json:"row_filter"`
	QueryBasedConnectorConfig      *IngestionPipelineDefinition_TableSpecificConfig_QueryBasedConnectorConfig `json:"query_based_connector_config"`
	AutoFullRefreshPolicy          *AutoFullRefreshPolicy                                                     `json:"auto_full_refresh_policy"`
	TableProperties                map[string]string                                                          `json:"table_properties"`
	EnableAutoClustering           *bool                                                                      `json:"enable_auto_clustering"`
	ClusteringColumns              []string                                                                   `json:"clustering_columns"`
}

// Configurations that are only applicable for query-based ingestion connectors..
type IngestionPipelineDefinition_TableSpecificConfig_QueryBasedConnectorConfig struct {
	CursorColumns                        []string `json:"cursor_columns"`
	DeletionCondition                    *string  `json:"deletion_condition"`
	HardDeletionSyncMinIntervalInSeconds *int64   `json:"hard_deletion_sync_min_interval_in_seconds"`
}

type IngestionPipelineDefinition_TableSpecificConfig_TablePropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type IngestionPipelineDefinition_WorkdayReportParameters struct {
	Incremental      *bool                                                               `json:"incremental"`
	ReportParameters []IngestionPipelineDefinition_WorkdayReportParameters_QueryKeyValue `json:"report_parameters"`
	Parameters       map[string]string                                                   `json:"parameters"`
}

// Key value pair used to specify configuration parameters to Execution.
type IngestionPipelineDefinition_WorkdayReportParameters_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type IngestionPipelineDefinition_WorkdayReportParameters_QueryKeyValue struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Jira specific options for ingestion.
type JiraConnectorOptions struct {
	IncludeJiraSpaces []string `json:"include_jira_spaces"`
}

// The request/response messages for the ListPipelines API. The default behavior
// is to return the 25 newest events in timestamp descending order for the given
// pipeline..
type ListPipelineEvents struct {
	PipelineId *string  `json:"pipeline_id"`
	PageToken  *string  `json:"page_token"`
	MaxResults *int     `json:"max_results"`
	OrderBy    []string `json:"order_by"`
	Filter     *string  `json:"filter"`
}

type ListPipelineEvents_Response struct {
	Events        []PipelineEvent `json:"events"`
	NextPageToken *string         `json:"next_page_token"`
	PrevPageToken *string         `json:"prev_page_token"`
}

// The request/response messages for the ListPipelines API. The default behavior
// is to return the 25 first pipelines in ascending order of pipeline id..
type ListPipelines struct {
	PageToken  *string  `json:"page_token"`
	MaxResults *int     `json:"max_results"`
	OrderBy    []string `json:"order_by"`
	Filter     *string  `json:"filter"`
}

type ListPipelines_Response struct {
	Statuses      []PipelineStateInfo `json:"statuses"`
	NextPageToken *string             `json:"next_page_token"`
}

// The request/response messages for the ListUpdates API. The default behavior
// is to return the 25 most recent updates in timestamp descending order for the
// given pipeline. No custom sorting or filtering is supported..
type ListUpdates struct {
	PipelineId    *string `json:"pipeline_id"`
	PageToken     *string `json:"page_token"`
	MaxResults    *int    `json:"max_results"`
	UntilUpdateId *string `json:"until_update_id"`
}

type ListUpdates_Response struct {
	Updates       []UpdateInfo `json:"updates"`
	NextPageToken *string      `json:"next_page_token"`
	PrevPageToken *string      `json:"prev_page_token"`
}

type ManualTrigger struct {
}

// Meta Marketing (Meta Ads) specific options for ingestion.
type MetaMarketingOptions struct {
	Level                        *string  `json:"level"`
	Breakdowns                   []string `json:"breakdowns"`
	ActionBreakdowns             []string `json:"action_breakdowns"`
	ActionReportTime             *string  `json:"action_report_time"`
	StartDate                    *string  `json:"start_date"`
	CustomInsightsLookbackWindow *int     `json:"custom_insights_lookback_window"`
	TimeIncrement                *string  `json:"time_increment"`
}

type NotebookLibrary struct {
	Path *string `json:"path"`
}

type Notifications struct {
	EmailRecipients []string `json:"email_recipients"`
	Alerts          []string `json:"alerts"`
}

// Proto representing a window.
type OperationTimeWindow struct {
	StartHour  *int        `json:"start_hour"`
	DaysOfWeek []DayOfWeek `json:"days_of_week"`
	TimeZoneId *string     `json:"time_zone_id"`
}

type Origin struct {
	Cloud                         *string `json:"cloud"`
	Region                        *string `json:"region"`
	OrgId                         *int64  `json:"org_id"`
	PipelineId                    *string `json:"pipeline_id"`
	PipelineName                  *string `json:"pipeline_name"`
	ClusterId                     *string `json:"cluster_id"`
	UpdateId                      *string `json:"update_id"`
	MaintenanceId                 *string `json:"maintenance_id"`
	TableId                       *string `json:"table_id"`
	DatasetName                   *string `json:"dataset_name"`
	FlowId                        *string `json:"flow_id"`
	FlowName                      *string `json:"flow_name"`
	BatchId                       *int64  `json:"batch_id"`
	RequestId                     *string `json:"request_id"`
	UcResourceId                  *string `json:"uc_resource_id"`
	Host                          *string `json:"host"`
	MaterializationName           *string `json:"materialization_name"`
	GraphId                       *string `json:"graph_id"`
	IngestionSourceConnectionName *string `json:"ingestion_source_connection_name"`
	IngestionSourceCatalogName    *string `json:"ingestion_source_catalog_name"`
	IngestionSourceSchemaName     *string `json:"ingestion_source_schema_name"`
	IngestionSourceTableName      *string `json:"ingestion_source_table_name"`
	IngestionSourceTableVersion   *string `json:"ingestion_source_table_version"`
}

// Outlook specific options for ingestion.
type OutlookOptions struct {
	FolderFilter  []string `json:"folder_filter"`
	SenderFilter  []string `json:"sender_filter"`
	SubjectFilter []string `json:"subject_filter"`
	StartDate     *string  `json:"start_date"`
}

type PathPattern struct {
	Include *string `json:"include"`
}

type PipelineCluster struct {
	Label                     *string                   `json:"label"`
	ApplyPolicyDefaultValues  *bool                     `json:"apply_policy_default_values"`
	SparkConf                 map[string]string         `json:"spark_conf"`
	AwsAttributes             *PipelinesAwsAttributes   `json:"aws_attributes"`
	AzureAttributes           *PipelinesAzureAttributes `json:"azure_attributes"`
	GcpAttributes             *PipelinesGcpAttributes   `json:"gcp_attributes"`
	NodeTypeId                *string                   `json:"node_type_id"`
	DriverNodeTypeId          *string                   `json:"driver_node_type_id"`
	SshPublicKeys             []string                  `json:"ssh_public_keys"`
	CustomTags                map[string]string         `json:"custom_tags"`
	ClusterLogConf            *PipelinesClusterLogConf  `json:"cluster_log_conf"`
	SparkEnvVars              map[string]string         `json:"spark_env_vars"`
	InitScripts               []PipelinesInitScriptInfo `json:"init_scripts"`
	InstancePoolId            *string                   `json:"instance_pool_id"`
	PolicyId                  *string                   `json:"policy_id"`
	EnableLocalDiskEncryption *bool                     `json:"enable_local_disk_encryption"`
	DriverInstancePoolId      *string                   `json:"driver_instance_pool_id"`
	NumWorkers                *int                      `json:"num_workers"`
	Autoscale                 *PipelinesAutoScale       `json:"autoscale"`
}

// Key value pair used to specify configuration parameters to Execution.
type PipelineCluster_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Key value pair used to specify configuration parameters to Execution.
type PipelineCluster_SparkConfEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Key value pair used to specify configuration parameters to Execution.
type PipelineCluster_SparkEnvVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type PipelineDeployment struct {
	Kind             *DeploymentKind `json:"kind"`
	MetadataFilePath *string         `json:"metadata_file_path"`
}

type PipelineEvent struct {
	Id            *string        `json:"id"`
	Sequence      *Sequencing    `json:"sequence"`
	Origin        *Origin        `json:"origin"`
	Timestamp     *string        `json:"timestamp"`
	Message       *string        `json:"message"`
	Level         *EventLevel    `json:"level"`
	Error         *ErrorDetail   `json:"error"`
	EventType     *string        `json:"event_type"`
	MaturityLevel *MaturityLevel `json:"maturity_level"`
	Truncation    *Truncation    `json:"truncation"`
}

type PipelineLibrary struct {
	Jar      *string                `json:"jar"`
	Maven    *PipelinesMavenLibrary `json:"maven"`
	Whl      *string                `json:"whl"`
	Notebook *NotebookLibrary       `json:"notebook"`
	File     *NotebookLibrary       `json:"file"`
	Glob     *PathPattern           `json:"glob"`
}

type PipelineSpec struct {
	Id                  *string                             `json:"id"`
	Name                *string                             `json:"name"`
	Storage             *string                             `json:"storage"`
	Configuration       map[string]string                   `json:"configuration"`
	Clusters            []PipelineCluster                   `json:"clusters"`
	Libraries           []PipelineLibrary                   `json:"libraries"`
	IngestionDefinition *IngestionPipelineDefinition        `json:"ingestion_definition"`
	GatewayDefinition   *IngestionGatewayPipelineDefinition `json:"gateway_definition"`
	Trigger             *PipelineTrigger                    `json:"trigger"`
	Target              *string                             `json:"target"`
	Schema              *string                             `json:"schema"`
	Filters             *Filters                            `json:"filters"`
	Continuous          *bool                               `json:"continuous"`
	Development         *bool                               `json:"development"`
	Photon              *bool                               `json:"photon"`
	Edition             *string                             `json:"edition"`
	Channel             *string                             `json:"channel"`
	Catalog             *string                             `json:"catalog"`
	Notifications       []Notifications                     `json:"notifications"`
	Serverless          *bool                               `json:"serverless"`
	Deployment          *PipelineDeployment                 `json:"deployment"`
	RestartWindow       *RestartWindow                      `json:"restart_window"`
	BudgetPolicyId      *string                             `json:"budget_policy_id"`
	Tags                map[string]string                   `json:"tags"`
	EventLog            *EventLogSpec                       `json:"event_log"`
	RootPath            *string                             `json:"root_path"`
	Environment         *PipelinesEnvironment               `json:"environment"`
	UsagePolicyId       *string                             `json:"usage_policy_id"`
}

// Key value pair used to specify configuration parameters to Execution.
type PipelineSpec_ConfigurationEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// A key-value entry that defines a single pipeline tags..
type PipelineSpec_TagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type PipelineStateInfo struct {
	PipelineId      *string               `json:"pipeline_id"`
	State           *PipelineState        `json:"state"`
	ClusterId       *string               `json:"cluster_id"`
	Name            *string               `json:"name"`
	LatestUpdates   []UpdateStateInfo     `json:"latest_updates"`
	CreatorUserName *string               `json:"creator_user_name"`
	RunAsUserName   *string               `json:"run_as_user_name"`
	Health          *PipelineHealthStatus `json:"health"`
}

type PipelineTrigger struct {
	Manual *ManualTrigger `json:"manual"`
	Cron   *CronTrigger   `json:"cron"`
}

type PipelinesAutoScale struct {
	MinWorkers *int    `json:"min_workers"`
	MaxWorkers *int    `json:"max_workers"`
	Mode       *string `json:"mode"`
}

// Attributes set during cluster creation which are related to Amazon Web
// Services..
type PipelinesAwsAttributes struct {
	FirstOnDemand       *int                      `json:"first_on_demand"`
	Availability        *PipelinesAwsAvailability `json:"availability"`
	ZoneId              *string                   `json:"zone_id"`
	InstanceProfileArn  *string                   `json:"instance_profile_arn"`
	SpotBidPricePercent *int                      `json:"spot_bid_price_percent"`
	EbsVolumeType       *PipelinesEbsVolumeType   `json:"ebs_volume_type"`
	EbsVolumeCount      *int                      `json:"ebs_volume_count"`
	EbsVolumeSize       *int                      `json:"ebs_volume_size"`
	EbsVolumeIops       *int                      `json:"ebs_volume_iops"`
	EbsVolumeThroughput *int                      `json:"ebs_volume_throughput"`
}

// Attributes set during cluster creation which are related to Azure..
type PipelinesAzureAttributes struct {
	FirstOnDemand   *int                        `json:"first_on_demand"`
	Availability    *PipelinesAzureAvailability `json:"availability"`
	SpotBidMaxPrice *float64                    `json:"spot_bid_max_price"`
}

// Cluster log delivery config.
type PipelinesClusterLogConf struct {
	Dbfs *PipelinesDbfsStorageInfo `json:"dbfs"`
}

// A storage location in DBFS.
type PipelinesDbfsStorageInfo struct {
	Destination *string `json:"destination"`
}

// The environment entity used to preserve serverless environment side panel,
// jobs' environment for non-notebook task, and DLT's environment for classic
// and serverless pipelines. In this minimal environment spec, only pip
// dependencies are supported..
type PipelinesEnvironment struct {
	Dependencies       []string `json:"dependencies"`
	EnvironmentVersion *string  `json:"environment_version"`
}

// Attributes set during cluster creation which are related to Gcp..
type PipelinesGcpAttributes struct {
	GoogleServiceAccount *string                   `json:"google_service_account"`
	BootDiskSize         *int                      `json:"boot_disk_size"`
	Availability         *PipelinesGcpAvailability `json:"availability"`
	ZoneId               *string                   `json:"zone_id"`
	LocalSsdCount        *int                      `json:"local_ssd_count"`
}

// Config for an individual init script.
type PipelinesInitScriptInfo struct {
	Dbfs *PipelinesDbfsStorageInfo `json:"dbfs"`
	S3   *PipelinesS3StorageInfo   `json:"s3"`
}

// Write-only setting, available only in Create/Update calls. Specifies the user
// or service principal that the pipeline runs as. If not specified, the
// pipeline runs as the user who created the pipeline.
//
// Only `user_name` or `service_principal_name` can be specified. If both are
// specified, an error is thrown..
type PipelinesJobRunAs struct {
	UserName             *string `json:"user_name"`
	ServicePrincipalName *string `json:"service_principal_name"`
}

type PipelinesMavenLibrary struct {
	Coordinates *string  `json:"coordinates"`
	Repo        *string  `json:"repo"`
	Exclusions  []string `json:"exclusions"`
}

// A storage location in Amazon S3.
type PipelinesS3StorageInfo struct {
	Destination      *string `json:"destination"`
	Region           *string `json:"region"`
	Endpoint         *string `json:"endpoint"`
	EnableEncryption *bool   `json:"enable_encryption"`
	EncryptionType   *string `json:"encryption_type"`
	KmsKey           *string `json:"kms_key"`
	CannedAcl        *string `json:"canned_acl"`
}

// PG-specific catalog-level configuration parameters.
type PostgresCatalogConfig struct {
	SlotConfig *PostgresSlotConfig `json:"slot_config"`
}

// PostgresSlotConfig contains the configuration for a Postgres logical
// replication slot.
type PostgresSlotConfig struct {
	SlotName        *string `json:"slot_name"`
	PublicationName *string `json:"publication_name"`
}

// Specifies a replace_where predicate override for a replace where flow..
type ReplaceWhereOverride struct {
	FlowName          *string `json:"flow_name"`
	PredicateOverride *string `json:"predicate_override"`
}

type RestartWindow struct {
	StartHour  *int        `json:"start_hour"`
	DaysOfWeek []DayOfWeek `json:"days_of_week"`
	TimeZoneId *string     `json:"time_zone_id"`
}

type RestorePipelineRequest struct {
	PipelineId *string `json:"pipeline_id"`
}

type RestorePipelineRequest_Response struct {
}

// Configuration for rewinding a specific dataset..
type RewindDatasetSpec struct {
	Identifier       *string `json:"identifier"`
	Cascade          *bool   `json:"cascade"`
	ResetCheckpoints *bool   `json:"reset_checkpoints"`
}

// Information about a rewind being requested for this pipeline or some of the
// datasets in it..
type RewindSpec struct {
	RewindTimestamp *string             `json:"rewind_timestamp"`
	DryRun          *bool               `json:"dry_run"`
	Datasets        []RewindDatasetSpec `json:"datasets"`
}

type Sequencing struct {
	DataPlaneId       *DataPlaneId `json:"data_plane_id"`
	ControlPlaneSeqNo *int64       `json:"control_plane_seq_no"`
}

type SerializedException struct {
	ClassName *string      `json:"class_name"`
	Message   *string      `json:"message"`
	Stack     []StackFrame `json:"stack"`
}

type SharepointOptions struct {
	Url                  *string                                 `json:"url"`
	EntityType           *SharepointOptions_SharepointEntityType `json:"entity_type"`
	FileIngestionOptions *FileIngestionOptions                   `json:"file_ingestion_options"`
}

// Smartsheet specific options for ingestion.
type SmartsheetOptions struct {
	EnforceSchema *bool `json:"enforce_schema"`
}

// SourceCatalogConfig contains catalog-level custom configuration parameters
// for each source.
type SourceCatalogConfig struct {
	SourceCatalog *string                `json:"source_catalog"`
	Postgres      *PostgresCatalogConfig `json:"postgres"`
}

type SourceConfig struct {
	Catalog *SourceCatalogConfig `json:"catalog"`
}

type StackFrame struct {
	DeclaringClass *string `json:"declaring_class"`
	MethodName     *string `json:"method_name"`
	FileName       *string `json:"file_name"`
	LineNumber     *int    `json:"line_number"`
}

type StartUpdate struct {
	PipelineId               *string                `json:"pipeline_id"`
	FullRefresh              *bool                  `json:"full_refresh"`
	Cause                    *UpdateCause           `json:"cause"`
	RefreshSelection         []string               `json:"refresh_selection"`
	FullRefreshSelection     []string               `json:"full_refresh_selection"`
	ResetCheckpointSelection []string               `json:"reset_checkpoint_selection"`
	ValidateOnly             *bool                  `json:"validate_only"`
	RewindSpec               *RewindSpec            `json:"rewind_spec"`
	Parameters               map[string]string      `json:"parameters"`
	ReplaceWhereOverrides    []ReplaceWhereOverride `json:"replace_where_overrides"`
}

type StartUpdate_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type StartUpdate_Response struct {
	UpdateId *string `json:"update_id"`
}

type StopPipeline struct {
	PipelineId *string `json:"pipeline_id"`
}

type StopPipeline_Response struct {
}

// TikTok Ads specific options for ingestion.
type TikTokAdsOptions struct {
	LookbackWindowDays *int                               `json:"lookback_window_days"`
	SyncStartDate      *string                            `json:"sync_start_date"`
	Dimensions         []string                           `json:"dimensions"`
	Metrics            []string                           `json:"metrics"`
	ReportType         *TikTokAdsOptions_TikTokReportType `json:"report_type"`
	DataLevel          *TikTokAdsOptions_TikTokDataLevel  `json:"data_level"`
	QueryLifetime      *bool                              `json:"query_lifetime"`
}

// Information about truncations applied to this event..
type Truncation struct {
	TruncatedFields []Truncation_TruncationDetail `json:"truncated_fields"`
}

// Details about a specific field that was truncated..
type Truncation_TruncationDetail struct {
	FieldName *string `json:"field_name"`
}

type UpdateInfo struct {
	PipelineId           *string                  `json:"pipeline_id"`
	UpdateId             *string                  `json:"update_id"`
	Config               *PipelineSpec            `json:"config"`
	Cause                *UpdateCause             `json:"cause"`
	State                *UpdateState_UpdateState `json:"state"`
	ClusterId            *string                  `json:"cluster_id"`
	CreationTime         *int64                   `json:"creation_time"`
	FullRefresh          *bool                    `json:"full_refresh"`
	RefreshSelection     []string                 `json:"refresh_selection"`
	FullRefreshSelection []string                 `json:"full_refresh_selection"`
	ValidateOnly         *bool                    `json:"validate_only"`
	Mode                 *UpdateMode              `json:"mode"`
	Parameters           map[string]string        `json:"parameters"`
}

type UpdateInfo_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UpdateState struct {
}

type UpdateStateInfo struct {
	UpdateId     *string                  `json:"update_id"`
	State        *UpdateState_UpdateState `json:"state"`
	CreationTime *string                  `json:"creation_time"`
}
