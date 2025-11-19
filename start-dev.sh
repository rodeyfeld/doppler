#!/bin/sh
set -e

# Load S3 credentials if available
if [ -f /shared/s3-credentials.env ]; then
    echo "Loading S3 credentials from /shared/s3-credentials.env..."
    # Read and export credentials, handling special characters
    export S3_ACCESS_KEY_ID="$(grep '^S3_ACCESS_KEY_ID=' /shared/s3-credentials.env | cut -d'=' -f2-)"
    export S3_SECRET_ACCESS_KEY="$(grep '^S3_SECRET_ACCESS_KEY=' /shared/s3-credentials.env | cut -d'=' -f2-)"
    echo "  Loaded: S3_ACCESS_KEY_ID=${S3_ACCESS_KEY_ID}"
    echo "  Loaded: S3_SECRET_ACCESS_KEY=***[hidden]"
    echo "✓ S3 credentials loaded"
else
    echo "⚠ Warning: S3 credentials not found at /shared/s3-credentials.env"
fi

# Install dependencies
echo "Installing dependencies..."
go mod download
bun install

# Verify credentials are set
if [ -z "$S3_ACCESS_KEY_ID" ] || [ -z "$S3_SECRET_ACCESS_KEY" ]; then
    echo "ERROR: S3 credentials not loaded properly!"
    exit 1
fi

# Start air for hot-reloading
exec air

