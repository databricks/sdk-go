package v2

// Returns the details of a policy family at a specific version.
type GetPolicyFamily struct {
	PolicyFamilyId *string `json:"policy_family_id"`
	Version        *int64  `json:"version"`
}

// Returns the list of policy families available to use at their latest version.
type ListPolicyFamilies struct {
	MaxResults *int64  `json:"max_results"`
	PageToken  *string `json:"page_token"`
}

type ListPolicyFamilies_Response struct {
	PolicyFamilies []PolicyFamily `json:"policy_families"`
	NextPageToken  *string        `json:"next_page_token"`
}

type PolicyFamily struct {
	PolicyFamilyId *string `json:"policy_family_id"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Definition     *string `json:"definition"`
}
