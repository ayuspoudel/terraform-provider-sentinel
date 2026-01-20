#!/usr/bin/env bash

# Author: @ayuspoudel


set -euo pipefail
build_version="0.3.8"
echo "=== Sentinel Terraform Provider: Local Build Script $build_version==="

echo "1) Checking for package name consistency..."

failed=0

while IFS= read -r dir; do
  gofiles=$(find "$dir" -maxdepth 1 -type f -name "*.go")
  [ -z "$gofiles" ] && continue

  pkgs=$(echo "$gofiles" | xargs grep -h "^package " | awk '{print $2}' | sort -u)
  count=$(echo "$pkgs" | wc -l | tr -d ' ')

  if [ "$count" -gt 1 ]; then
    echo
    echo "ERROR: multiple packages detected in $dir"
    echo "$pkgs"
    failed=1
  fi
done < <(find internal -type d)

if [ "$failed" -eq 1 ]; then
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
PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/ayuspoudel/sentinel/0.3.8/$(go env GOOS)_$(go env GOARCH)"
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
      version = "0.3.8"
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
  terraform plan

)

echo
echo "SUCCESS: Provider builds and Terraform initializes correctly with $build_version."
