#!/usr/bin/env bash

# Author: @ayuspoudel


set -euo pipefail

echo "=== Sentinel Terraform Provider: Local Build Script ==="

echo
echo "1) Checking for package name consistency..."
pkg_issues=$(grep -R "^package " internal | awk '{print $2}' | sort | uniq -c | awk '$1 > 1')
if [ -n "$pkg_issues" ]; then
  echo "ERROR: Multiple package names detected in same directories:"
  echo "$pkg_issues"
  exit 1
fi
echo "OK: Package names are consistent."

echo
echo "2) Cleaning Go build and module cache..."
go clean -cache -testcache -modcache
echo "OK: Cache cleaned."

echo
echo "3) Building provider binary..."
go build -o terraform-provider-sentinel
echo "OK: Binary built."

echo
echo "4) Verifying binary..."
ls -lh terraform-provider-sentinel
file terraform-provider-sentinel

echo
echo "5) Installing binary into local Terraform plugin directory..."
PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/ayuspoudel/sentinel/0.0.0/$(go env GOOS)_$(go env GOARCH)"
mkdir -p "$PLUGIN_DIR"
cp terraform-provider-sentinel "$PLUGIN_DIR/"
chmod +x "$PLUGIN_DIR/terraform-provider-sentinel"
echo "Installed to $PLUGIN_DIR"

echo
echo "6) Terraform smoke test..."
TMP_DIR=$(mktemp -d)
cat <<EOF > "$TMP_DIR/main.tf"
terraform {
  required_providers {
    sentinel = {
      source  = "ayuspoudel/sentinel"
      version = "0.0.0"
    }
  }
}

provider "sentinel" {
  endpoint = "http://localhost:9000"
}
EOF

(
  cd "$TMP_DIR"
  terraform init -input=false
)

echo
echo "SUCCESS: Provider builds and Terraform initializes correctly."
