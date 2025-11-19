#!/bin/sh
set -e

# Load S3 credentials if available
if [ -f /shared/s3-credentials.env ]; then
    echo "Loading S3 credentials..."
    export $(cat /shared/s3-credentials.env | xargs)
    echo "✓ S3 credentials loaded"
else
    echo "Warning: S3 credentials not found"
fi

# Install dependencies
echo "Installing dependencies..."
go mod download
bun install

# Start air for hot-reloading
exec air

