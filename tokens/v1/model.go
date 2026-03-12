package v1

type CreateToken struct {
	LifetimeSeconds *int64   `json:"lifetime_seconds"`
	Comment         *string  `json:"comment"`
	Scopes          []string `json:"scopes"`
}

type CreateToken_Response struct {
	TokenValue *string          `json:"token_value"`
	TokenInfo  *PublicTokenInfo `json:"token_info"`
}

type ListTokens struct {
}

type ListTokens_Response struct {
	TokenInfos []PublicTokenInfo `json:"token_infos"`
}

type PublicTokenInfo struct {
	TokenId          *string  `json:"token_id"`
	CreationTime     *int64   `json:"creation_time"`
	ExpiryTime       *int64   `json:"expiry_time"`
	Comment          *string  `json:"comment"`
	Scopes           []string `json:"scopes"`
	LastAccessedTime *int64   `json:"last_accessed_time"`
}

type RevokeToken struct {
	TokenId *string `json:"token_id"`
}

type RevokeToken_Response struct {
}

type UpdateToken struct {
	TokenId    *string          `json:"token_id"`
	Token      *PublicTokenInfo `json:"token"`
	UpdateMask *string          `json:"update_mask"`
}

type UpdateTokenResponse struct {
}
