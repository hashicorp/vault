provider "aws" "default" {
  region = var.aws_region
}

provider "aws" "west1" {
  region = "us-west-1"
}

provider "aws" "west2" {
  region = "us-west-2"
}

provider "enos" "rhel" {
  transport = {
    ssh = {
      user             = "ec2-user"
      private_key_path = abspath(var.aws_ssh_private_key_path)
    }
  }
}

provider "enos" "ubuntu" {
  transport = {
    ssh = {
      user             = "ubuntu"
      private_key_path = abspath(var.aws_ssh_private_key_path)
    }
  }
}
