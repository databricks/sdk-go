package v1

type DestinationType string

const (
	DestinationType_Slack          DestinationType = "SLACK"
	DestinationType_Email          DestinationType = "EMAIL"
	DestinationType_Webhook        DestinationType = "WEBHOOK"
	DestinationType_Pagerduty      DestinationType = "PAGERDUTY"
	DestinationType_MicrosoftTeams DestinationType = "MICROSOFT_TEAMS"
)

// String representation for [fmt.Print].
func (f *DestinationType) String() string {
	return string(*f)
}

type Config struct {
	Slack          *SlackConfig          `json:"slack"`
	Email          *EmailConfig          `json:"email"`
	GenericWebhook *GenericWebhookConfig `json:"generic_webhook"`
	Pagerduty      *PagerdutyConfig      `json:"pagerduty"`
	MicrosoftTeams *MicrosoftTeamsConfig `json:"microsoft_teams"`
}

type CreateNotificationDestinationRequest struct {
	DisplayName *string `json:"display_name"`
	Config      *Config `json:"config"`
}

type DeleteNotificationDestinationRequest struct {
	Id *string `json:"id"`
}

type EmailConfig struct {
	Addresses []string `json:"addresses"`
}

type Empty struct {
}

type GenericWebhookConfig struct {
	Url         *string `json:"url"`
	UrlSet      *bool   `json:"url_set"`
	Username    *string `json:"username"`
	UsernameSet *bool   `json:"username_set"`
	Password    *string `json:"password"`
	PasswordSet *bool   `json:"password_set"`
}

type GetNotificationDestinationRequest struct {
	Id *string `json:"id"`
}

type ListNotificationDestinationsRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int64  `json:"page_size"`
}

type ListNotificationDestinationsResponse struct {
	Results       []ListNotificationDestinationsResult `json:"results"`
	NextPageToken *string                              `json:"next_page_token"`
}

type ListNotificationDestinationsResult struct {
	Id              *string          `json:"id"`
	DisplayName     *string          `json:"display_name"`
	DestinationType *DestinationType `json:"destination_type"`
	Config          *Config          `json:"config"`
}

type MicrosoftTeamsConfig struct {
	Url           *string `json:"url"`
	UrlSet        *bool   `json:"url_set"`
	AppId         *string `json:"app_id"`
	AppIdSet      *bool   `json:"app_id_set"`
	AuthSecret    *string `json:"auth_secret"`
	AuthSecretSet *bool   `json:"auth_secret_set"`
	ChannelUrl    *string `json:"channel_url"`
	ChannelUrlSet *bool   `json:"channel_url_set"`
	TenantId      *string `json:"tenant_id"`
	TenantIdSet   *bool   `json:"tenant_id_set"`
}

type NotificationDestination struct {
	Id              *string          `json:"id"`
	DisplayName     *string          `json:"display_name"`
	DestinationType *DestinationType `json:"destination_type"`
	Config          *Config          `json:"config"`
}

type PagerdutyConfig struct {
	IntegrationKey    *string `json:"integration_key"`
	IntegrationKeySet *bool   `json:"integration_key_set"`
}

type SlackConfig struct {
	Url           *string `json:"url"`
	UrlSet        *bool   `json:"url_set"`
	OauthToken    *string `json:"oauth_token"`
	OauthTokenSet *bool   `json:"oauth_token_set"`
	ChannelId     *string `json:"channel_id"`
	ChannelIdSet  *bool   `json:"channel_id_set"`
}

type UpdateNotificationDestinationRequest struct {
	Id          *string `json:"id"`
	DisplayName *string `json:"display_name"`
	Config      *Config `json:"config"`
}
