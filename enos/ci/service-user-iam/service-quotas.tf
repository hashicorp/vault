# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

locals {
  // This is the code of the service quota to request a change for. Each adjustable limit has a
  // unique code. See, https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/servicequotas_service_quota#quota_code
  subnets_per_vpcs_quota                = "L-F678F1CE"
  standard_spot_instance_requests_quota = "L-34B43A08"
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_east_1" {
  provider     = aws.us_east_1
  quota_code   = local.subnets_per_vpcs_quota
  service_code = "vpc"
  value        = 100
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_east_2" {
  provider     = aws.us_east_2
  quota_code   = local.subnets_per_vpcs_quota
  service_code = "vpc"
  value        = 100
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_west_1" {
  provider     = aws.us_west_1
  quota_code   = local.subnets_per_vpcs_quota
  service_code = "vpc"
  value        = 100
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_west_2" {
  provider     = aws.us_west_2
  quota_code   = local.subnets_per_vpcs_quota
  service_code = "vpc"
  value        = 100
}

resource "aws_servicequotas_service_quota" "spot_requests_per_region_us_east_1" {
  provider     = aws.us_east_1
  quota_code   = local.standard_spot_instance_requests_quota
  service_code = "ec2"
  value        = 640
}

resource "aws_servicequotas_service_quota" "spot_requests_per_region_us_east_2" {
  provider     = aws.us_east_2
  quota_code   = local.standard_spot_instance_requests_quota
  service_code = "ec2"
  value        = 640
}

resource "aws_servicequotas_service_quota" "spot_requests_per_region_us_west_1" {
  provider     = aws.us_west_1
  quota_code   = local.standard_spot_instance_requests_quota
  service_code = "ec2"
  value        = 640
}

resource "aws_servicequotas_service_quota" "spot_requests_per_region_us_west_2" {
  provider     = aws.us_west_2
  quota_code   = local.standard_spot_instance_requests_quota
  service_code = "ec2"
  value        = 640
}
