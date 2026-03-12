package v1

type ListPublishedOAuthApps struct {
	AccountId *string `json:"account_id"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListPublishedOAuthApps_Response struct {
	Apps          []PublishedOAuthApp `json:"apps"`
	NextPageToken *string             `json:"next_page_token"`
}

type PublishedOAuthApp struct {
	AppId                *string  `json:"app_id"`
	ClientId             *string  `json:"client_id"`
	Name                 *string  `json:"name"`
	Description          *string  `json:"description"`
	IsConfidentialClient *bool    `json:"is_confidential_client"`
	RedirectUrls         []string `json:"redirect_urls"`
	Scopes               []string `json:"scopes"`
}
