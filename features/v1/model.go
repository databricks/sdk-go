package v1

import "time"

type Function_FunctionType string

const (
	Function_FunctionType_FunctionTypeUnspecified Function_FunctionType = "FUNCTION_TYPE_UNSPECIFIED"
	Function_FunctionType_Avg                     Function_FunctionType = "AVG"
	Function_FunctionType_Count                   Function_FunctionType = "COUNT"
	Function_FunctionType_Sum                     Function_FunctionType = "SUM"
	Function_FunctionType_Min                     Function_FunctionType = "MIN"
	Function_FunctionType_Max                     Function_FunctionType = "MAX"
	Function_FunctionType_First                   Function_FunctionType = "FIRST"
	Function_FunctionType_Last                    Function_FunctionType = "LAST"
	Function_FunctionType_ApproxCountDistinct     Function_FunctionType = "APPROX_COUNT_DISTINCT"
	Function_FunctionType_ApproxPercentile        Function_FunctionType = "APPROX_PERCENTILE"
	Function_FunctionType_StddevPop               Function_FunctionType = "STDDEV_POP"
	Function_FunctionType_StddevSamp              Function_FunctionType = "STDDEV_SAMP"
	Function_FunctionType_VarPop                  Function_FunctionType = "VAR_POP"
	Function_FunctionType_VarSamp                 Function_FunctionType = "VAR_SAMP"
)

// String representation for [fmt.Print].
func (f *Function_FunctionType) String() string {
	return string(*f)
}

type MaterializedFeature_PipelineScheduleState string

const (
	MaterializedFeature_PipelineScheduleState_PipelineScheduleStateUnspecified MaterializedFeature_PipelineScheduleState = "PIPELINE_SCHEDULE_STATE_UNSPECIFIED"
	MaterializedFeature_PipelineScheduleState_Snapshot                         MaterializedFeature_PipelineScheduleState = "SNAPSHOT"
	MaterializedFeature_PipelineScheduleState_Active                           MaterializedFeature_PipelineScheduleState = "ACTIVE"
	MaterializedFeature_PipelineScheduleState_Paused                           MaterializedFeature_PipelineScheduleState = "PAUSED"
)

// String representation for [fmt.Print].
func (f *MaterializedFeature_PipelineScheduleState) String() string {
	return string(*f)
}

// An aggregation function applied over a time window..
type AggregationFunction struct {
	Avg                 *AvgFunction                 `json:"avg"`
	CountFunction       *CountFunction               `json:"count_function"`
	Sum                 *SumFunction                 `json:"sum"`
	Min                 *MinFunction                 `json:"min"`
	Max                 *MaxFunction                 `json:"max"`
	First               *FirstFunction               `json:"first"`
	Last                *LastFunction                `json:"last"`
	ApproxCountDistinct *ApproxCountDistinctFunction `json:"approx_count_distinct"`
	ApproxPercentile    *ApproxPercentileFunction    `json:"approx_percentile"`
	StddevPop           *StddevPopFunction           `json:"stddev_pop"`
	StddevSamp          *StddevSampFunction          `json:"stddev_samp"`
	VarPop              *VarPopFunction              `json:"var_pop"`
	VarSamp             *VarSampFunction             `json:"var_samp"`
	TimeWindow          *TimeWindow                  `json:"time_window"`
}

// Computes the approximate count of distinct values..
type ApproxCountDistinctFunction struct {
	Input      *string  `json:"input"`
	RelativeSd *float64 `json:"relative_sd"`
}

// Computes the approximate percentile of values..
type ApproxPercentileFunction struct {
	Input      *string  `json:"input"`
	Percentile *float64 `json:"percentile"`
	Accuracy   *int64   `json:"accuracy"`
}

type AuthConfig struct {
	UcServiceCredentialName *string `json:"uc_service_credential_name"`
}

// Computes the average of values..
type AvgFunction struct {
	Input *string `json:"input"`
}

type BackfillSource struct {
	DeltaTableSource *DeltaTableSource `json:"delta_table_source"`
}

type BatchCreateMaterializedFeaturesRequest struct {
	Requests []CreateMaterializedFeatureRequest `json:"requests"`
}

type BatchCreateMaterializedFeaturesResponse struct {
	MaterializedFeatures []MaterializedFeature `json:"materialized_features"`
}

type ColumnIdentifier struct {
	VariantExprPath *string `json:"variant_expr_path"`
}

// A ColumnSelection function, equivalent to the LAST() record of an entity over
// a lifetime ContinuousWindow.
type ColumnSelection struct {
	Column *string `json:"column"`
}

type ContinuousWindow struct {
	WindowDuration *time.Duration `json:"window_duration"`
	Offset         *time.Duration `json:"offset"`
}

// Computes the count of values..
type CountFunction struct {
	Input *string `json:"input"`
}

type CreateFeatureRequest struct {
	Feature *Feature `json:"feature"`
}

type CreateKafkaConfigRequest struct {
	KafkaConfig *KafkaConfig `json:"kafka_config"`
}

type CreateMaterializedFeatureRequest struct {
	MaterializedFeature *MaterializedFeature `json:"materialized_feature"`
}

type DataSource struct {
	DeltaTableSource *DeltaTableSource `json:"delta_table_source"`
	KafkaSource      *KafkaSource      `json:"kafka_source"`
}

type DeleteFeatureRequest struct {
	FullName *string `json:"full_name"`
}

type DeleteKafkaConfigRequest struct {
	Name *string `json:"name"`
}

type DeleteMaterializedFeatureRequest struct {
	MaterializedFeatureId *string `json:"materialized_feature_id"`
}

type DeltaTableSource struct {
	FullName          *string  `json:"full_name"`
	EntityColumns     []string `json:"entity_columns"`
	TimeseriesColumn  *string  `json:"timeseries_column"`
	FilterCondition   *string  `json:"filter_condition"`
	TransformationSql *string  `json:"transformation_sql"`
	DataframeSchema   *string  `json:"dataframe_schema"`
}

type EntityColumn struct {
	Name *string `json:"name"`
}

type Feature struct {
	FullName         *string           `json:"full_name"`
	Source           *DataSource       `json:"source"`
	Inputs           []string          `json:"inputs"`
	Function         *Function         `json:"function"`
	TimeWindow       *TimeWindow       `json:"time_window"`
	Description      *string           `json:"description"`
	FilterCondition  *string           `json:"filter_condition"`
	LineageContext   *LineageContext   `json:"lineage_context"`
	Entities         []EntityColumn    `json:"entities"`
	TimeseriesColumn *TimeseriesColumn `json:"timeseries_column"`
}

// Returns the first value..
type FirstFunction struct {
	Input *string `json:"input"`
}

type Function struct {
	FunctionType        *Function_FunctionType    `json:"function_type"`
	ExtraParameters     []Function_ExtraParameter `json:"extra_parameters"`
	AggregationFunction *AggregationFunction      `json:"aggregation_function"`
	ColumnSelection     *ColumnSelection          `json:"column_selection"`
}

// Deprecated: Use typed fields on function-specific messages (e.g.
// ApproxPercentileFunction.percentile) or AggregationFunction.ExtraParameter
// instead. Kept for backwards compatibility..
type Function_ExtraParameter struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type GetFeatureRequest struct {
	FullName *string `json:"full_name"`
}

type GetKafkaConfigRequest struct {
	Name *string `json:"name"`
}

type GetMaterializedFeatureRequest struct {
	MaterializedFeatureId *string `json:"materialized_feature_id"`
}

type JobContext struct {
	JobId    *int64 `json:"job_id"`
	JobRunId *int64 `json:"job_run_id"`
}

type KafkaConfig struct {
	Name             *string           `json:"name"`
	BootstrapServers *string           `json:"bootstrap_servers"`
	SubscriptionMode *SubscriptionMode `json:"subscription_mode"`
	AuthConfig       *AuthConfig       `json:"auth_config"`
	KeySchema        *SchemaConfig     `json:"key_schema"`
	ValueSchema      *SchemaConfig     `json:"value_schema"`
	ExtraOptions     map[string]string `json:"extra_options"`
	BackfillSource   *BackfillSource   `json:"backfill_source"`
}

type KafkaConfig_ExtraOptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type KafkaSource struct {
	Name                       *string            `json:"name"`
	EntityColumnIdentifiers    []ColumnIdentifier `json:"entity_column_identifiers"`
	TimeseriesColumnIdentifier *ColumnIdentifier  `json:"timeseries_column_identifier"`
	FilterCondition            *string            `json:"filter_condition"`
}

// Returns the last value..
type LastFunction struct {
	Input *string `json:"input"`
}

// Lineage context information for tracking where an API was invoked. This will
// allow us to track lineage, which currently uses caller entity information for
// use across the Lineage Client and Observability in Lumberjack..
type LineageContext struct {
	NotebookId *int64      `json:"notebook_id"`
	JobContext *JobContext `json:"job_context"`
}

type ListFeaturesRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListFeaturesResponse struct {
	Features      []Feature `json:"features"`
	NextPageToken *string   `json:"next_page_token"`
}

type ListKafkaConfigsRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListKafkaConfigsResponse struct {
	KafkaConfigs  []KafkaConfig `json:"kafka_configs"`
	NextPageToken *string       `json:"next_page_token"`
}

type ListMaterializedFeaturesRequest struct {
	FeatureName *string `json:"feature_name"`
	PageToken   *string `json:"page_token"`
	PageSize    *int    `json:"page_size"`
}

type ListMaterializedFeaturesResponse struct {
	MaterializedFeatures []MaterializedFeature `json:"materialized_features"`
	NextPageToken        *string               `json:"next_page_token"`
}

// A materialized feature represents a feature that is continuously computed and
// stored..
type MaterializedFeature struct {
	MaterializedFeatureId   *string                                    `json:"materialized_feature_id"`
	FeatureName             *string                                    `json:"feature_name"`
	OfflineStoreConfig      *OfflineStoreConfig                        `json:"offline_store_config"`
	OnlineStoreConfig       *OnlineStoreConfig                         `json:"online_store_config"`
	TableName               *string                                    `json:"table_name"`
	PipelineScheduleState   *MaterializedFeature_PipelineScheduleState `json:"pipeline_schedule_state"`
	LastMaterializationTime *time.Time                                 `json:"last_materialization_time"`
	CronSchedule            *string                                    `json:"cron_schedule"`
}

// Computes the maximum value..
type MaxFunction struct {
	Input *string `json:"input"`
}

// Computes the minimum value..
type MinFunction struct {
	Input *string `json:"input"`
}

// Configuration for offline store destination..
type OfflineStoreConfig struct {
	CatalogName     *string `json:"catalog_name"`
	SchemaName      *string `json:"schema_name"`
	TableNamePrefix *string `json:"table_name_prefix"`
}

// Configuration for online store destination..
type OnlineStoreConfig struct {
	CatalogName     *string `json:"catalog_name"`
	SchemaName      *string `json:"schema_name"`
	TableNamePrefix *string `json:"table_name_prefix"`
	OnlineStoreName *string `json:"online_store_name"`
}

type SchemaConfig struct {
	JsonSchema *string `json:"json_schema"`
}

type SlidingWindow struct {
	WindowDuration *time.Duration `json:"window_duration"`
	SlideDuration  *time.Duration `json:"slide_duration"`
}

// Computes the population standard deviation..
type StddevPopFunction struct {
	Input *string `json:"input"`
}

// Computes the sample standard deviation..
type StddevSampFunction struct {
	Input *string `json:"input"`
}

type SubscriptionMode struct {
	Assign           *string `json:"assign"`
	Subscribe        *string `json:"subscribe"`
	SubscribePattern *string `json:"subscribe_pattern"`
}

// Computes the sum of values..
type SumFunction struct {
	Input *string `json:"input"`
}

type TimeWindow struct {
	Continuous *ContinuousWindow `json:"continuous"`
	Tumbling   *TumblingWindow   `json:"tumbling"`
	Sliding    *SlidingWindow    `json:"sliding"`
}

type TimeseriesColumn struct {
	Name *string `json:"name"`
}

type TumblingWindow struct {
	WindowDuration *time.Duration `json:"window_duration"`
}

type UpdateFeatureRequest struct {
	Feature    *Feature `json:"feature"`
	UpdateMask *string  `json:"update_mask"`
}

type UpdateKafkaConfigRequest struct {
	KafkaConfig *KafkaConfig `json:"kafka_config"`
	UpdateMask  *string      `json:"update_mask"`
}

type UpdateMaterializedFeatureRequest struct {
	MaterializedFeature *MaterializedFeature `json:"materialized_feature"`
	UpdateMask          *string              `json:"update_mask"`
}

// Computes the population variance..
type VarPopFunction struct {
	Input *string `json:"input"`
}

// Computes the sample variance..
type VarSampFunction struct {
	Input *string `json:"input"`
}
