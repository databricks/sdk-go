package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/databricks/options"
	dataqualityv1 "github.com/databricks/sdk-go/dataquality/v1"
	qualitymonitorv2 "github.com/databricks/sdk-go/qualitymonitor/v2"
)

func createQualityMonitor() {
	tokenProvider := auth.TokenProviderFn(func(ctx context.Context) (*auth.Token, error) {
		return &auth.Token{
			Value: "dapi1234567890",
		}, nil
	})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client, err := qualitymonitorv2.NewClient(context.Background(),
		options.WithHost("https://dbc-<your-workspace-id>.cloud.databricks.com/"),
		options.WithCredentials(auth.NewTokenCredentials(tokenProvider)),
		options.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	resp, err := client.CreateQualityMonitor(context.Background(), &qualitymonitorv2.CreateQualityMonitorRequest{
		QualityMonitor: &qualitymonitorv2.QualityMonitor{
			ObjectType: toPtr("schema"),
			ObjectId:   toPtr("<your-schema-id>"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create quality monitor: %v", err)
	}
	fmt.Printf("Quality monitor created: %+v", resp)
}

func deleteQualityMonitor() {
	tokenProvider := auth.TokenProviderFn(func(ctx context.Context) (*auth.Token, error) {
		return &auth.Token{
			Value: "dapi1234567890",
		}, nil
	})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client, err := qualitymonitorv2.NewClient(context.Background(),
		options.WithHost("https://dbc-<your-workspace-id>.cloud.databricks.com/"),
		options.WithCredentials(auth.NewTokenCredentials(tokenProvider)),
		options.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	err = client.DeleteQualityMonitor(context.Background(), &qualitymonitorv2.DeleteQualityMonitorRequest{
		ObjectType: toPtr("schema"),
		ObjectId:   toPtr("<your-schema-id>"),
	})
	if err != nil {
		log.Fatalf("Failed to delete quality monitor: %v", err)
	}
	fmt.Printf("Quality monitor deleted")
}

func listQualityMonitors() {
	tokenProvider := auth.TokenProviderFn(func(ctx context.Context) (*auth.Token, error) {
		return &auth.Token{
			Value: "dapi1234567890",
		}, nil
	})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client, err := qualitymonitorv2.NewClient(context.Background(),
		options.WithHost("https://dbc-<your-workspace-id>.cloud.databricks.com/"),
		options.WithCredentials(auth.NewTokenCredentials(tokenProvider)),
		options.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	iter := client.ListQualityMonitorIter(context.Background(), &qualitymonitorv2.ListQualityMonitorRequest{})
	for item, err := range iter {
		if err != nil {
			log.Fatalf("Failed to get quality monitor: %v", err)
		}
		fmt.Printf("Quality monitor: %+v", item)
	}
}

func createDataQualityMonitor() {
	tokenProvider := auth.TokenProviderFn(func(ctx context.Context) (*auth.Token, error) {
		return &auth.Token{
			Value: "dapi1234567890",
		}, nil
	})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client, err := dataqualityv1.NewClient(context.Background(),
		options.WithHost("https://dbc-<your-workspace-id>.cloud.databricks.com/"),
		options.WithCredentials(auth.NewTokenCredentials(tokenProvider)),
		options.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	resp, err := client.CreateMonitor(context.Background(), &dataqualityv1.CreateMonitorRequest{
		Monitor: &dataqualityv1.Monitor{
			ObjectType: toPtr("schema"),
			ObjectId:   toPtr("<your-schema-id>"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create quality monitor: %v", err)
	}
	fmt.Printf("Quality monitor created: %+v", resp)
}

func listDataQualityMonitors() {
	tokenProvider := auth.TokenProviderFn(func(ctx context.Context) (*auth.Token, error) {
		return &auth.Token{
			Value: "dapi1234567890",
		}, nil
	})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	client, err := dataqualityv1.NewClient(context.Background(),
		options.WithHost("https://dbc-<your-workspace-id>.cloud.databricks.com/"),
		options.WithCredentials(auth.NewTokenCredentials(tokenProvider)),
		options.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	iter := client.ListRefreshIter(context.Background(), &dataqualityv1.ListRefreshRequest{
		ObjectType: toPtr("schema"),
		ObjectId:   toPtr("<your-schema-id>"),
	})
	for item, err := range iter {
		if err != nil {
			log.Fatalf("Failed to get data quality monitor: %v", err)
		}
		fmt.Printf("Data quality monitor: %s, %s\n", *item.ObjectType, *item.ObjectId)
	}
}

func main() {
	createQualityMonitor()
	deleteQualityMonitor()
	// listQualityMonitors()
	// createDataQualityMonitor()
	// listDataQualityMonitors()
}

func toPtr[T any](v T) *T {
	return &v
}
