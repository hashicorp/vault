#!/usr/bin/env bash
# Copyright IBM Corp. 2025, 2026
# SPDX-License-Identifier: BUSL-1.1

# Script to capture AWS resource counts for debugging and leak detection
# Usage: capture-aws-resources.sh <region> <scenario> <output_file> <stage>

set -euo pipefail

REGION="${1:-}"
SCENARIO="${2:-}"
OUTPUT_FILE="${3:-}"
STAGE="${4:-Unknown}"

if [[ -z "$REGION" || -z "$SCENARIO" || -z "$OUTPUT_FILE" ]]; then
  echo "ERROR: Missing required arguments"
  echo "Usage: $0 <region> <scenario> <output_file> <stage>"
  exit 1
fi

# Disable command echo for security (prevent credential exposure in logs)
set +x

# Function to safely query AWS and handle errors
query_aws() {
  local service="$1"
  local resource="$2"
  local query="$3"
  
  local result
  result=$(aws "$service" "$resource" \
    --region "$REGION" \
    --query "$query" \
    --output text \
    --no-cli-pager 2>&1) || {
    echo "ERROR: Failed to query $resource" >&2
    echo "0"
    return
  }
  
  echo "${result:-0}"
}

# Capture resource counts
{
  echo "=== AWS Resource Counts $STAGE ==="
  echo "Timestamp: $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
  echo "Region: $REGION"
  echo "Scenario: $SCENARIO"
  echo ""
  echo "## Resource Counts"
  echo "VPCs: $(query_aws ec2 describe-vpcs 'Vpcs | length(@)')"
  echo "Subnets: $(query_aws ec2 describe-subnets 'Subnets | length(@)')"
  echo "Internet Gateways: $(query_aws ec2 describe-internet-gateways 'InternetGateways | length(@)')"
  echo "NAT Gateways: $(query_aws ec2 describe-nat-gateways 'NatGateways | length(@)')"
  echo "Route Tables: $(query_aws ec2 describe-route-tables 'RouteTables | length(@)')"
  echo "Security Groups: $(query_aws ec2 describe-security-groups 'SecurityGroups | length(@)')"
  echo "Network ACLs: $(query_aws ec2 describe-network-acls 'NetworkAcls | length(@)')"
  echo "Elastic IPs: $(query_aws ec2 describe-addresses 'Addresses | length(@)')"
  echo "EC2 Instances (All): $(query_aws ec2 describe-instances 'Reservations[].Instances | length(@)')"
  echo "EC2 Instances (Running): $(aws ec2 describe-instances --region "$REGION" --filters "Name=instance-state-name,Values=running" --query 'Reservations[].Instances | length(@)' --output text --no-cli-pager 2>&1 || echo "0")"
  echo "Network Interfaces: $(query_aws ec2 describe-network-interfaces 'NetworkInterfaces | length(@)')"
  echo "Volumes: $(query_aws ec2 describe-volumes 'Volumes | length(@)')"
  echo "Key Pairs: $(query_aws ec2 describe-key-pairs 'KeyPairs | length(@)')"
  echo "Load Balancers (ELB): $(query_aws elb describe-load-balancers 'LoadBalancerDescriptions | length(@)')"
  echo "Load Balancers (ALB/NLB): $(query_aws elbv2 describe-load-balancers 'LoadBalancers | length(@)')"
  echo "Target Groups: $(query_aws elbv2 describe-target-groups 'TargetGroups | length(@)')"
  echo "Launch Templates: $(query_aws ec2 describe-launch-templates 'LaunchTemplates | length(@)')"
  echo "Auto Scaling Groups: $(query_aws autoscaling describe-auto-scaling-groups 'AutoScalingGroups | length(@)')"
} > "$OUTPUT_FILE"

echo "Resource counts captured to $OUTPUT_FILE"
cat "$OUTPUT_FILE"
