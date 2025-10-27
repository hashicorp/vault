# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

provider "aws" {
  region = "us-east-1"
  alias  = "us_east_1"
}

provider "aws" {
  region = "us-east-2"
  alias  = "us_east_2"
}

provider "aws" {
  region = "us-west-1"
  alias  = "us_west_1"
}

provider "aws" {
  region = "us-west-2"
  alias  = "us_west_2"
}
