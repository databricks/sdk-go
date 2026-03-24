package v1

// Request message for CreateFeatureTag..
type CreateFeatureTagRequest struct {
	TableName   *string     `json:"table_name"`
	FeatureName *string     `json:"feature_name"`
	FeatureTag  *FeatureTag `json:"feature_tag"`
}

// Request message for DeleteFeatureTag..
type DeleteFeatureTagRequest struct {
	TableName   *string `json:"table_name"`
	FeatureName *string `json:"feature_name"`
	Key         *string `json:"key"`
}

type FeatureLineage struct {
	Models         []FeatureLineage_Model         `json:"models"`
	FeatureSpecs   []FeatureLineage_FeatureSpec   `json:"feature_specs"`
	OnlineFeatures []FeatureLineage_OnlineFeature `json:"online_features"`
}

type FeatureLineage_FeatureSpec struct {
	Name *string `json:"name"`
}

type FeatureLineage_Model struct {
	Name    *string `json:"name"`
	Version *int64  `json:"version"`
}

type FeatureLineage_OnlineFeature struct {
	FeatureName *string `json:"feature_name"`
	TableName   *string `json:"table_name"`
}

// Represents a tag on a feature in a feature table..
type FeatureTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type GetFeatureLineageRequest struct {
	FeatureName *string `json:"feature_name"`
	TableName   *string `json:"table_name"`
}

// Request message for GetFeatureTag..
type GetFeatureTagRequest struct {
	TableName   *string `json:"table_name"`
	FeatureName *string `json:"feature_name"`
	Key         *string `json:"key"`
}

// Request message for ListFeatureTags..
type ListFeatureTagsRequest struct {
	TableName   *string `json:"table_name"`
	FeatureName *string `json:"feature_name"`
	PageToken   *string `json:"page_token"`
	PageSize    *int    `json:"page_size"`
}

// Response message for ListFeatureTag..
type ListFeatureTagsResponse struct {
	FeatureTags   []FeatureTag `json:"feature_tags"`
	NextPageToken *string      `json:"next_page_token"`
}

// Request message for UpdateFeatureTag..
type UpdateFeatureTagRequest struct {
	TableName   *string     `json:"table_name"`
	FeatureName *string     `json:"feature_name"`
	FeatureTag  *FeatureTag `json:"feature_tag"`
	UpdateMask  *string     `json:"update_mask"`
}
