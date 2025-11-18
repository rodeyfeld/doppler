package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	bucketName = "doppler-photos"
	keyLabel   = "doppler-dev"
)

func main() {
	log.SetFlags(0)
	ctx := context.Background()

	// Wait for Garage to be ready
	waitForGarage(ctx)

	// Get node ID and setup cluster
	nodeID := parseField(garage(ctx, "node", "id"), "", 0)[:16]
	log.Printf("Node ID: %s", nodeID)

	log.Printf("Assigning layout to node %s...", nodeID)
	garageOk(ctx, "layout", "assign", nodeID, "-z", "local", "-c", "1G")
	garageOk(ctx, "layout", "apply", "--version", "1")

	log.Printf("Ensuring %s bucket exists...", bucketName)
	garageOk(ctx, "bucket", "create", bucketName)

	// Create or find key
	keyID := findOrCreateKey(ctx)
	log.Printf("Granting key %s access to %s...", keyID, bucketName)
	garageOk(ctx, "bucket", "allow", "--key", keyID, "--read", "--write", bucketName)

	// Export credentials for the web service
	exportCredentials(ctx, keyID)

	log.Printf("Garage bootstrap finished.")
}

func exportCredentials(ctx context.Context, keyID string) {
	log.Printf("Exporting S3 credentials...")

	// Get full key details including secret
	out := garage(ctx, "key", "info", keyID)

	// Parse the secret key
	secretKey := parseField(out, "Secret key:", 2)

	// Write to shared volume that web service can read
	envContent := fmt.Sprintf("S3_ACCESS_KEY_ID=%s\nS3_SECRET_ACCESS_KEY=%s\n", keyID, secretKey)

	if err := os.WriteFile("/shared/s3-credentials.env", []byte(envContent), 0644); err != nil {
		log.Printf("Warning: Could not write credentials file: %v", err)
		return
	}

	log.Printf("✓ Exported S3 credentials to /shared/s3-credentials.env")
}

func waitForGarage(ctx context.Context) {
	log.Printf("Waiting for Garage...")
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Fatal("Timeout waiting for Garage")
		case <-ticker.C:
			if _, err := garageRun(ctx, "status"); err == nil {
				log.Printf("Garage ready!")
				return
			}
		}
	}
}

func findOrCreateKey(ctx context.Context) string {
	out, _ := garageRun(ctx, "key", "list")
	if strings.Contains(out, keyLabel) {
		log.Printf("Key '%s' exists", keyLabel)
		return parseField(out, keyLabel, 0)
	}

	log.Printf("Creating key '%s'...", keyLabel)
	out = garage(ctx, "key", "create", keyLabel)
	return parseField(out, "Key ID:", 2)
}

// garage runs a command and exits on error
func garage(ctx context.Context, args ...string) string {
	out, err := garageRun(ctx, args...)
	if err != nil {
		log.Fatalf("garage %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	if out != "" {
		log.Print(out)
	}
	return out
}

// garageOk runs a command and logs warnings on error (idempotent ops)
func garageOk(ctx context.Context, args ...string) {
	out, err := garageRun(ctx, args...)
	if out != "" {
		log.Print(out)
	}
	if err != nil {
		log.Printf("warning: garage %s: %v", strings.Join(args, " "), err)
	}
}

// garageRun executes a garage command
func garageRun(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "/garage", append([]string{"-c", "/etc/garage.toml"}, args...)...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// parseField extracts a field from output
func parseField(output, lineMatch string, fieldIdx int) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if lineMatch == "" || strings.Contains(line, lineMatch) {
			fields := strings.Fields(line)
			if fieldIdx < len(fields) {
				return fields[fieldIdx]
			}
		}
	}
	log.Fatalf("Could not parse field %d from output matching '%s'", fieldIdx, lineMatch)
	return ""
}
