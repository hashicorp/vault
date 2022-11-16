terraform "bootstrap" {
  required_providers {
    tfe = {
      source = "hashicorp/tfe"
    }
  }

  cloud {
    organization = var.organization
    hostname     = "app.terraform.io"
    token        = var.tfc_api_token

    workspaces {
      name = "${var.product_line}-ci-enos-boostrap"
    }
  }
}

terraform "us_east_1" {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  cloud {
    organization = var.organization
    hostname     = "app.terraform.io"
    token        = var.tfc_api_token

    workspaces {
      name = "${var.product_line}-ci-enos-boostrap-us-east-1"
    }
  }
}

terraform "us_east_2" {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  cloud {
    organization = var.organization
    hostname     = "app.terraform.io"
    token        = var.tfc_api_token

    workspaces {
      name = "${var.product_line}-ci-enos-boostrap-us-east-2"
    }
  }
}

terraform "us_west_1" {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  cloud {
    organization = var.organization
    hostname     = "app.terraform.io"
    token        = var.tfc_api_token

    workspaces {
      name = "${var.product_line}-ci-enos-boostrap-us-west-1"
    }
  }
}

terraform "us_west_2" {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  cloud {
    organization = var.organization
    hostname     = "app.terraform.io"
    token        = var.tfc_api_token

    workspaces {
      name = "${var.product_line}-ci-enos-boostrap-us-west-2"
    }
  }
}
