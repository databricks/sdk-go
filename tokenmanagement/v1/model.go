package v1

type AdminTokenInfo struct {
	TokenId           *string  `json:"token_id"`
	CreationTime      *int64   `json:"creation_time"`
	ExpiryTime        *int64   `json:"expiry_time"`
	Comment           *string  `json:"comment"`
	CreatedById       *int64   `json:"created_by_id"`
	CreatedByUsername *string  `json:"created_by_username"`
	OwnerId           *int64   `json:"owner_id"`
	WorkspaceId       *int64   `json:"workspace_id"`
	LastUsedDay       *int64   `json:"last_used_day"`
	Scopes            []string `json:"scopes"`
}

// Configuration details for creating on-behalf tokens..
type CreateOnBehalfOfToken struct {
	ApplicationId   *string  `json:"application_id"`
	LifetimeSeconds *int64   `json:"lifetime_seconds"`
	Comment         *string  `json:"comment"`
	Scopes          []string `json:"scopes"`
}

// An on-behalf token was successfully created for the service principal..
type CreateOnBehalfOfToken_Response struct {
	TokenValue *string         `json:"token_value"`
	TokenInfo  *AdminTokenInfo `json:"token_info"`
}

// !! KEEP THIS IN-SYNC WITH THE WORKSPACE PROTO DEFINITIONS IN SERVICE.PROTO !!
//
// The only differences should be: 1. The OpenAPI labels. 2. The account_id
// request parameter..
type GetToken struct {
	TokenId *string `json:"token_id"`
}

// Token with specified Token ID was successfully returned..
type GetToken_Response struct {
	TokenInfo *AdminTokenInfo `json:"token_info"`
}

// !! KEEP THIS IN-SYNC WITH THE ACCOUNT PROTO DEFINITIONS IN
// ACCOUNT_SERVICE.PROTO !!
//
// The only differences should be: 1. The OpenAPI labels. 2. The account_id
// request parameter. 3. The string filter parameter instead of hard-coded
// filters..
type ListTokens struct {
	CreatedById       *int64  `json:"created_by_id"`
	CreatedByUsername *string `json:"created_by_username"`
}

// Tokens were successfully returned..
type ListTokens_Response struct {
	TokenInfos []AdminTokenInfo `json:"token_infos"`
}

type RevokeToken struct {
	TokenId *string `json:"token_id"`
}

// The token was successfully deleted..
type RevokeToken_Response struct {
}
