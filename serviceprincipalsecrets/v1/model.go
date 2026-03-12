package v1

import "time"

type CreateServicePrincipalSecret struct {
	AccountId        *string        `json:"account_id"`
	ServicePrincipal *string        `json:"service_principal"`
	Lifetime         *time.Duration `json:"lifetime"`
}

type CreateServicePrincipalSecretResponse struct {
	Id         *string    `json:"id"`
	Secret     *string    `json:"secret"`
	SecretHash *string    `json:"secret_hash"`
	CreateTime *string    `json:"create_time"`
	UpdateTime *string    `json:"update_time"`
	Status     *string    `json:"status"`
	ExpireTime *time.Time `json:"expire_time"`
}

type DeleteServicePrincipalSecret struct {
	AccountId        *string `json:"account_id"`
	ServicePrincipal *string `json:"service_principal"`
	SecretId         *string `json:"secret_id"`
}

type DeleteServicePrincipalSecret_Response struct {
}

type ListServicePrincipalSecrets struct {
	AccountId        *string `json:"account_id"`
	ServicePrincipal *string `json:"service_principal"`
	PageToken        *string `json:"page_token"`
	PageSize         *int    `json:"page_size"`
}

type ListServicePrincipalSecrets_Response struct {
	Secrets       []ServicePrincipalSecret `json:"secrets"`
	NextPageToken *string                  `json:"next_page_token"`
}

type ServicePrincipalSecret struct {
	Id         *string    `json:"id"`
	Secret     *string    `json:"secret"`
	SecretHash *string    `json:"secret_hash"`
	CreateTime *string    `json:"create_time"`
	UpdateTime *string    `json:"update_time"`
	Status     *string    `json:"status"`
	ExpireTime *time.Time `json:"expire_time"`
}
