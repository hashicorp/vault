locals {
  subnets_per_vps_quota = "L-F678F1CE"
}

resource "aws_servicequotas_service_quota" "vpcs_pre_region_us_east_1" {
  provider     = aws.us_east_2
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}

resource "aws_servicequotas_service_quota" "vpcs_pre_region_us_east_2" {
  provider     = aws.us_east_2
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}

resource "aws_servicequotas_service_quota" "vpcs_pre_region_us_west_1" {
  provider     = aws.us_west_1
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}

resource "aws_servicequotas_service_quota" "vpcs_pre_region_us_west_2" {
  provider     = aws.us_west_2
  quota_code   = local.subnets_per_vps_quota
  service_code = "vpc"
  value        = 50
}
