package v1

type GetAssignableRolesForResourceRequest struct {
	AccountId *string `json:"account_id"`
	Resource  *string `json:"resource"`
}

type GetAssignableRolesForResourceResponse struct {
	Roles []Role `json:"roles"`
}

type GetRuleSetRequest struct {
	AccountId *string `json:"account_id"`
	Name      *string `json:"name"`
	Etag      *string `json:"etag"`
}

type GrantRule struct {
	Principals []string `json:"principals"`
	Role       *string  `json:"role"`
}

type Role struct {
	Name *string `json:"name"`
}

type RuleSet struct {
	Name       *string     `json:"name"`
	Etag       *string     `json:"etag"`
	GrantRules []GrantRule `json:"grant_rules"`
}

type RuleSetUpdateRequest struct {
	Name       *string     `json:"name"`
	Etag       *string     `json:"etag"`
	GrantRules []GrantRule `json:"grant_rules"`
}

type UpdateRuleSetRequest struct {
	AccountId *string               `json:"account_id"`
	Name      *string               `json:"name"`
	RuleSet   *RuleSetUpdateRequest `json:"rule_set"`
}
