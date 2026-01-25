package jobs

import "time"

type Job struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateJobRequest struct {
	Job Job `json:"job"`
}

type ListJobsRequest struct {
	PageSize  int    `json:"page_size"`
	PageToken string `json:"page_token"`
}

type ListJobsResponse struct {
	Jobs          []Job  `json:"jobs"`
	NextPageToken string `json:"next_page_token"`
}
