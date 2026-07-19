#!/bin/bash
# Build the Go backend as a Lambda-compatible ZIP package
set -e

echo "Building Lambda binary (linux/arm64)..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go

echo "Creating deployment ZIP..."
zip -j function.zip bootstrap

echo "Done! Upload function.zip to AWS Lambda."
echo "  Runtime:  provided.al2023"
echo "  Handler:  bootstrap"
echo "  Arch:     arm64"

# Cleanup intermediate binary
rm -f bootstrap
