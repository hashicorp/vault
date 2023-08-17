locals {
  // This is the code of the service quota to request a change for. Each adjustable limit has a
  // unique code. See, https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/servicequotas_service_quota#quota_code
  subnets_per_vps_quota = "L-F678F1CE"
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_east_1" {
  provider     = aws.us_east_2
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_east_2" {
  provider     = aws.us_east_2
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_west_1" {
  provider     = aws.us_west_1
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}

resource "aws_servicequotas_service_quota" "vpcs_per_region_us_west_2" {
  provider     = aws.us_west_2
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}
