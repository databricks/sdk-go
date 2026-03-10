package v1

type CreateTagAssignmentRequest struct {
	TagAssignment *TagAssignment `json:"tag_assignment"`
}

type DeleteTagAssignmentRequest struct {
	EntityType *string `json:"entity_type"`
	EntityId   *string `json:"entity_id"`
	TagKey     *string `json:"tag_key"`
}

type GetTagAssignmentRequest struct {
	EntityType *string `json:"entity_type"`
	EntityId   *string `json:"entity_id"`
	TagKey     *string `json:"tag_key"`
}

type ListTagAssignmentsRequest struct {
	EntityType *string `json:"entity_type"`
	EntityId   *string `json:"entity_id"`
	PageSize   *int    `json:"page_size"`
	PageToken  *string `json:"page_token"`
}

type ListTagAssignmentsResponse struct {
	TagAssignments []TagAssignment `json:"tag_assignments"`
	NextPageToken  *string         `json:"next_page_token"`
}

type TagAssignment struct {
	EntityType *string `json:"entity_type"`
	EntityId   *string `json:"entity_id"`
	TagKey     *string `json:"tag_key"`
	TagValue   *string `json:"tag_value"`
}

type UpdateTagAssignmentRequest struct {
	TagAssignment *TagAssignment `json:"tag_assignment"`
	UpdateMask    *string        `json:"update_mask"`
}
