package v1

import "time"

type CreateCustomOAuthAppIntegration struct {
	AccountId            *string            `json:"account_id"`
	RedirectUrls         []string           `json:"redirect_urls"`
	Name                 *string            `json:"name"`
	Confidential         *bool              `json:"confidential"`
	TokenAccessPolicy    *TokenAccessPolicy `json:"token_access_policy"`
	Scopes               []string           `json:"scopes"`
	UserAuthorizedScopes []string           `json:"user_authorized_scopes"`
}

type CreatePublishedOAuthAppIntegration struct {
	AccountId         *string            `json:"account_id"`
	AppId             *string            `json:"app_id"`
	TokenAccessPolicy *TokenAccessPolicy `json:"token_access_policy"`
}

type CreatePublishedOAuthAppIntegration_Response struct {
	IntegrationId *string `json:"integration_id"`
}

type CustomOAuthAppIntegration struct {
	IntegrationId        *string            `json:"integration_id"`
	ClientId             *string            `json:"client_id"`
	RedirectUrls         []string           `json:"redirect_urls"`
	Name                 *string            `json:"name"`
	Confidential         *bool              `json:"confidential"`
	TokenAccessPolicy    *TokenAccessPolicy `json:"token_access_policy"`
	Scopes               []string           `json:"scopes"`
	CreatedBy            *int64             `json:"created_by"`
	CreateTime           *string            `json:"create_time"`
	CreatorUsername      *string            `json:"creator_username"`
	UserAuthorizedScopes []string           `json:"user_authorized_scopes"`
	PrincipalId          *int64             `json:"principal_id"`
}

type CustomOAuthAppIntegrationSecret struct {
	IntegrationId          *string    `json:"integration_id"`
	ClientId               *string    `json:"client_id"`
	ClientSecret           *string    `json:"client_secret"`
	PrincipalId            *int64     `json:"principal_id"`
	ClientSecretExpireTime *time.Time `json:"client_secret_expire_time"`
}

type DeleteCustomOAuthAppIntegration struct {
	AccountId     *string `json:"account_id"`
	IntegrationId *string `json:"integration_id"`
}

type DeleteCustomOAuthAppIntegration_Response struct {
}

type DeletePublishedOAuthAppIntegration struct {
	AccountId     *string `json:"account_id"`
	IntegrationId *string `json:"integration_id"`
}

type DeletePublishedOAuthAppIntegration_Response struct {
}

type GetCustomOAuthAppIntegration struct {
	AccountId     *string `json:"account_id"`
	IntegrationId *string `json:"integration_id"`
}

type GetPublishedOAuthAppIntegration struct {
	AccountId     *string `json:"account_id"`
	IntegrationId *string `json:"integration_id"`
}

type ListCustomOAuthAppIntegrations struct {
	AccountId              *string `json:"account_id"`
	PageToken              *string `json:"page_token"`
	PageSize               *int    `json:"page_size"`
	IncludeCreatorUsername *bool   `json:"include_creator_username"`
}

type ListCustomOAuthAppIntegrations_Response struct {
	Apps          []CustomOAuthAppIntegration `json:"apps"`
	NextPageToken *string                     `json:"next_page_token"`
}

type ListPublishedOAuthAppIntegrations struct {
	AccountId *string `json:"account_id"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListPublishedOAuthAppIntegrations_Response struct {
	Apps          []PublishedOAuthAppIntegration `json:"apps"`
	NextPageToken *string                        `json:"next_page_token"`
}

type PublishedOAuthAppIntegration struct {
	AppId             *string            `json:"app_id"`
	IntegrationId     *string            `json:"integration_id"`
	Name              *string            `json:"name"`
	TokenAccessPolicy *TokenAccessPolicy `json:"token_access_policy"`
	CreatedBy         *int64             `json:"created_by"`
	CreateTime        *string            `json:"create_time"`
}

type TokenAccessPolicy struct {
	AccessTokenTtlInMinutes          *int  `json:"access_token_ttl_in_minutes"`
	RefreshTokenTtlInMinutes         *int  `json:"refresh_token_ttl_in_minutes"`
	EnableSingleUseRefreshTokens     *bool `json:"enable_single_use_refresh_tokens"`
	AbsoluteSessionLifetimeInMinutes *int  `json:"absolute_session_lifetime_in_minutes"`
}

type UpdateCustomOAuthAppIntegration struct {
	AccountId            *string            `json:"account_id"`
	IntegrationId        *string            `json:"integration_id"`
	RedirectUrls         []string           `json:"redirect_urls"`
	TokenAccessPolicy    *TokenAccessPolicy `json:"token_access_policy"`
	Scopes               []string           `json:"scopes"`
	UserAuthorizedScopes []string           `json:"user_authorized_scopes"`
}

type UpdateCustomOAuthAppIntegration_Response struct {
}

type UpdatePublishedOAuthAppIntegration struct {
	AccountId         *string            `json:"account_id"`
	IntegrationId     *string            `json:"integration_id"`
	TokenAccessPolicy *TokenAccessPolicy `json:"token_access_policy"`
}

type UpdatePublishedOAuthAppIntegration_Response struct {
}
