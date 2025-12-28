package jobs

import (
	"context"
	"log/slog"
	"testing"

	"github.com/databricks/sdk-go/databricks/options"
	"github.com/databricks/sdk-go/internal/fixture"
	httpfixture "github.com/databricks/sdk-go/internal/fixture/http"
	"github.com/google/go-cmp/cmp"
)

func TestCreateJob(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{
			Method:       "POST",
			Path:         "/api/2.1/jobs/create",
			ExpectedBody: []byte(`{"job":{"id":"","name":"test-job","description":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`),
			Response:     []byte(`{"id":"12345","name":"test-job","description":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`),
		},
	)

	transport := httpfixture.NewTransport(seq)
	ctx := context.Background()
	client, err := NewClient(ctx,
		options.WithHost("https://example.com/api/2.1/jobs/create"),
		options.WithHTTPClient(transport.Client()),
		options.WithLogger(slog.Default()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	job, err := client.CreateJob(ctx, &CreateJobRequest{Job: Job{Name: "test-job"}})
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}

	want := &Job{ID: "12345", Name: "test-job"}
	if diff := cmp.Diff(want, job); diff != "" {
		t.Errorf("CreateJob mismatch (-want +got):\n%s", diff)
	}

	if err := transport.Verify(); err != nil {
		t.Errorf("fixture verification failed: %v", err)
	}
}

func TestListJobsIter(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{
			Method:   "GET",
			Path:     "/api/2.1/jobs/list?page_size=2",
			Response: []byte(`{"jobs":[{"id":"1","name":"job-1"},{"id":"2","name":"job-2"}],"next_page_token":"page2"}`),
		},
		fixture.Fixture{
			Method:   "GET",
			Path:     "/api/2.1/jobs/list?page_size=2&page_token=page2",
			Response: []byte(`{"jobs":[{"id":"3","name":"job-3"}],"next_page_token":""}`),
		},
	)

	transport := httpfixture.NewTransport(seq)
	ctx := context.Background()
	client, err := NewClient(ctx,
		options.WithHost("https://example.com/api/2.1/jobs/list"),
		options.WithHTTPClient(transport.Client()),
		options.WithLogger(slog.Default()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var jobs []*Job
	for job, err := range client.ListJobsIter(ctx, &ListJobsRequest{PageSize: 2}) {
		if err != nil {
			t.Fatalf("ListJobsIter failed: %v", err)
		}
		jobs = append(jobs, job)
	}

	want := []*Job{
		{ID: "1", Name: "job-1"},
		{ID: "2", Name: "job-2"},
		{ID: "3", Name: "job-3"},
	}
	if diff := cmp.Diff(want, jobs); diff != "" {
		t.Errorf("ListJobsIter mismatch (-want +got):\n%s", diff)
	}

	if err := transport.Verify(); err != nil {
		t.Errorf("fixture verification failed: %v", err)
	}
}
