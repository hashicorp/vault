provider "aws" "us_east_1" {
  region = "us-east-1"
}

provider "aws" "us_east_2" {
  region = "us-east-2"
}

provider "aws" "us_west_1" {
  region = "us-west-1"
}

provider "aws" "us_west_2" {
  region = "us-west-2"
}

provider "tfe" "bootstrap" {
  hostname = "app.terraform.io"
  token    = var.tfc_api_token
}
