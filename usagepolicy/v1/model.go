package v1

type SortSpec_Field string

const (
	SortSpec_Field_FieldUnspecified SortSpec_Field = "FIELD_UNSPECIFIED"
	SortSpec_Field_PolicyName       SortSpec_Field = "POLICY_NAME"
)

// String representation for [fmt.Print].
func (f *SortSpec_Field) String() string {
	return string(*f)
}

// A request to create a UsagePolicy..
type CreateUsagePolicyRequest struct {
	RequestId *string      `json:"request_id"`
	AccountId *string      `json:"account_id"`
	Policy    *UsagePolicy `json:"policy"`
}

type CustomPolicyTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Deletes a usage policy.
type DeleteUsagePolicyRequest struct {
	PolicyId  *string `json:"policy_id"`
	AccountId *string `json:"account_id"`
}

// Structured representation of a filter to be applied to a list of policies.
// All specified filters will be applied in conjunction..
type Filter struct {
	PolicyName      *string `json:"policy_name"`
	CreatorUserId   *int64  `json:"creator_user_id"`
	CreatorUserName *string `json:"creator_user_name"`
}

// Request to get a usage policy.
type GetUsagePolicyRequest struct {
	PolicyId  *string `json:"policy_id"`
	AccountId *string `json:"account_id"`
}

// The limit configuration of the policy. Limit configuration provide a budget
// policy level cost control by enforcing the limit..
type LimitConfig struct {
}

// Request to list usage policies.
type ListUsagePoliciesRequest struct {
	PageSize  *int      `json:"page_size"`
	PageToken *string   `json:"page_token"`
	FilterBy  *Filter   `json:"filter_by"`
	SortSpec  *SortSpec `json:"sort_spec"`
	AccountId *string   `json:"account_id"`
}

// A list of usage policies..
type ListUsagePoliciesResponse struct {
	Policies          []UsagePolicy `json:"policies"`
	NextPageToken     *string       `json:"next_page_token"`
	PreviousPageToken *string       `json:"previous_page_token"`
}

type SortSpec struct {
	Field      *SortSpec_Field `json:"field"`
	Descending *bool           `json:"descending"`
}

// Updates a UsagePolicy..
type UpdateUsagePolicyRequest struct {
	Policy      *UsagePolicy `json:"policy"`
	AccountId   *string      `json:"account_id"`
	LimitConfig *LimitConfig `json:"limit_config"`
}

// Contains the UsagePolicy details (same structure as BudgetPolicy).
type UsagePolicy struct {
	PolicyId            *string           `json:"policy_id"`
	PolicyName          *string           `json:"policy_name"`
	CustomTags          []CustomPolicyTag `json:"custom_tags"`
	BindingWorkspaceIds []int64           `json:"binding_workspace_ids"`
}
