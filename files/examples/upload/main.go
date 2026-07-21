// Command upload demonstrates uploading a file to a Unity Catalog volume with
// the files/v2 client, then downloading it back to verify the round trip.
//
// The client credentials and host are read from the default Databricks profile.
// For more information on setting up a Databricks profile, refer to the
// [Databricks Auth] documentation.
//
// [Databricks Auth]: https://pkg.go.dev/github.com/databricks/sdk-go/auth
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"

	files "github.com/databricks/sdk-go/files/v2"
	"github.com/databricks/sdk-go/options/client"
)

// remotePath is the Unity Catalog volume path the example writes to. Override it
// when running against a real workspace with a volume you own.
const remotePath = "/Volumes/main/default/my_volume/example-upload.bin"

func main() {
	ctx := context.Background()
	logger := slog.Default()

	// Create a new client from the default Databricks profile.
	c, err := files.NewClient(ctx, client.WithLogger(logger))
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}

	// Generate a dummy payload of 40 KiB.
	payload := bytes.Repeat([]byte("databricks"), 4096)

	// Upload the payload to the remote path.
	_, err = c.Upload(ctx, remotePath, bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("Upload: %v", err)
	}

	fmt.Printf("Uploaded %d bytes to %s\n", len(payload), remotePath)

	// Download the payload from the remote path.
	resp, err := c.DownloadFile(ctx, &files.DownloadFileRequest{FilePath: new(remotePath)})
	if err != nil {
		log.Fatalf("DownloadFile: %v", err)
	}
	defer resp.Contents.Close()

	// Verify the downloaded payload matches the original payload.
	got, err := io.ReadAll(resp.Contents)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}
	if !bytes.Equal(got, payload) {
		log.Fatalf("round trip mismatch: got %d bytes, want %d", len(got), len(payload))
	}

	fmt.Printf("Downloaded %d bytes; contents match\n", len(got))
}
