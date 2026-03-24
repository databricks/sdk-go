package v1

import "time"

// Policy that determines how to resolve conflicts when multiple upstream
// sources have different tag values..
type ConflictResolutionPolicy struct {
	DefaultValueOverride *DefaultValueOverridePolicy `json:"default_value_override"`
}

type CreateTagPolicyRequest struct {
	TagPolicy *TagPolicy `json:"tag_policy"`
}

// Policy that specifies a default value to use when resolving tag conflicts
// during propagation..
type DefaultValueOverridePolicy struct {
	DefaultValue *string `json:"default_value"`
}

type DeleteTagPolicyRequest struct {
	TagKey *string `json:"tag_key"`
}

type GetTagPolicyRequest struct {
	TagKey *string `json:"tag_key"`
}

type ListTagPoliciesRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type ListTagPoliciesResponse struct {
	TagPolicies   []TagPolicy `json:"tag_policies"`
	NextPageToken *string     `json:"next_page_token"`
}

// Configuration that controls how tags are automatically propagated through
// data lineage..
type PropagationConfig struct {
	Enabled            *bool                     `json:"enabled"`
	ConflictResolution *ConflictResolutionPolicy `json:"conflict_resolution"`
}

type TagPolicy struct {
	TagKey            *string            `json:"tag_key"`
	Id                *string            `json:"id"`
	Description       *string            `json:"description"`
	Values            []Value            `json:"values"`
	CreateTime        *time.Time         `json:"create_time"`
	UpdateTime        *time.Time         `json:"update_time"`
	PropagationConfig *PropagationConfig `json:"propagation_config"`
	AccountId         *string            `json:"account_id"`
}

type UpdateTagPolicyRequest struct {
	TagPolicy  *TagPolicy `json:"tag_policy"`
	UpdateMask *string    `json:"update_mask"`
}

type Value struct {
	Name *string `json:"name"`
}
