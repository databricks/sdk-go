package v1

type DownloadRequest struct {
	AccountId    *string `json:"account_id"`
	StartMonth   *string `json:"start_month"`
	EndMonth     *string `json:"end_month"`
	PersonalData *bool   `json:"personal_data"`
}

// Billable usage data was returned successfully..
type DownloadResponse struct {
	Content *string `json:"content"`
}
