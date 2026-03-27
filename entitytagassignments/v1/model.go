package v1

import "time"

type TagAssignmentSourceType string

const (
	TagAssignmentSourceType_TagAssignmentSourceTypeUnspecified              TagAssignmentSourceType = "TAG_ASSIGNMENT_SOURCE_TYPE_UNSPECIFIED"
	TagAssignmentSourceType_TagAssignmentSourceTypeSystemDataClassification TagAssignmentSourceType = "TAG_ASSIGNMENT_SOURCE_TYPE_SYSTEM_DATA_CLASSIFICATION"
)

// String representation for [fmt.Print].
func (f *TagAssignmentSourceType) String() string {
	return string(*f)
}

// Request to create a new entity tag assignment.
type CreateEntityTagAssignmentRequest struct {
	TagAssignment *EntityTagAssignment `json:"tag_assignment"`
}

// Request to delete an entity tag assignment.
type DeleteEntityTagAssignmentRequest struct {
	EntityName *string `json:"entity_name"`
	TagKey     *string `json:"tag_key"`
	EntityType *string `json:"entity_type"`
}

// Represents a tag assignment to an entity.
type EntityTagAssignment struct {
	EntityName *string                  `json:"entity_name"`
	TagKey     *string                  `json:"tag_key"`
	TagValue   *string                  `json:"tag_value"`
	EntityType *string                  `json:"entity_type"`
	UpdateTime *time.Time               `json:"update_time"`
	UpdatedBy  *string                  `json:"updated_by"`
	SourceType *TagAssignmentSourceType `json:"source_type"`
	Inherited  *bool                    `json:"inherited"`
}

// Request to get an entity tag assignment.
type GetEntityTagAssignmentRequest struct {
	EntityName       *string `json:"entity_name"`
	TagKey           *string `json:"tag_key"`
	EntityType       *string `json:"entity_type"`
	IncludeInherited *bool   `json:"include_inherited"`
}

// Request to list entity tag assignments.
type ListEntityTagAssignmentsRequest struct {
	EntityName       *string `json:"entity_name"`
	MaxResults       *int    `json:"max_results"`
	PageToken        *string `json:"page_token"`
	EntityType       *string `json:"entity_type"`
	IncludeInherited *bool   `json:"include_inherited"`
}

type ListEntityTagAssignmentsResponse struct {
	TagAssignments []EntityTagAssignment `json:"tag_assignments"`
	NextPageToken  *string               `json:"next_page_token"`
}

// Request to update an entity tag assignment.
type UpdateEntityTagAssignmentRequest struct {
	TagAssignment *EntityTagAssignment `json:"tag_assignment"`
	UpdateMask    *string              `json:"update_mask"`
}
