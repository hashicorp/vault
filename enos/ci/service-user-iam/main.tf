terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  cloud {
    hostname     = "app.terraform.io"
    organization = "hashicorp-qti"
    // workspace must be exported in the environment as: TF_WORKSPACE=<vault|vault-enterprise>-ci-enos-service-user-iam
  }
}

locals {
  enterprise_repositories = ["vault-enterprise"]
  is_ent                  = contains(local.enterprise_repositories, var.repository)
  ci_account_prefix       = local.is_ent ? "vault_enterprise" : "vault"
  service_user            = "github_actions-${local.ci_account_prefix}_ci"
  aws_account_id          = local.is_ent ? "505811019928" : "040730498200"
}

resource "aws_iam_role" "role" {
  provider           = aws.us_east_1
  name               = local.service_user
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy_document.json
}

data "aws_iam_policy_document" "assume_role_policy_document" {
  provider = aws.us_east_1
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${local.aws_account_id}:user/${local.service_user}"]
    }
  }
}

resource "aws_iam_role_policy" "role_policy" {
  provider = aws.us_east_1
  role     = aws_iam_role.role.name
  name     = "${local.service_user}_policy"
  policy   = data.aws_iam_policy_document.iam_policy_document.json
}

data "aws_iam_policy_document" "iam_policy_document" {
  provider = aws.us_east_1
  statement {
    effect = "Allow"
    actions = [
      "iam:ListRoles",
      "iam:CreateRole",
      "iam:GetRole",
      "iam:DeleteRole",
      "iam:ListInstanceProfiles",
      "iam:ListInstanceProfilesForRole",
      "iam:CreateInstanceProfile",
      "iam:GetInstanceProfile",
      "iam:DeleteInstanceProfile",
      "iam:ListPolicies",
      "iam:CreatePolicy",
      "iam:DeletePolicy",
      "iam:ListRoles",
      "iam:CreateRole",
      "iam:AddRoleToInstanceProfile",
      "iam:PassRole",
      "iam:RemoveRoleFromInstanceProfile",
      "iam:DeleteRole",
      "iam:ListRolePolicies",
      "iam:ListAttachedRolePolicies",
      "iam:AttachRolePolicy",
      "iam:GetRolePolicy",
      "iam:PutRolePolicy",
      "iam:DetachRolePolicy",
      "iam:DeleteRolePolicy",
      "ec2:DescribeAccountAttributes",
      "ec2:DescribeInstanceTypes",
      "ec2:DescribeInstanceCreditSpecifications",
      "ec2:DescribeImages",
      "ec2:DescribeRegions",
      "ec2:DescribeTags",
      "ec2:DescribeVpcClassicLink",
      "ec2:DescribeVpcClassicLinkDnsSupport",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DescribeAvailabilityZones",
      "ec2:DescribeSecurityGroups",
      "ec2:CreateSecurityGroup",
      "ec2:AuthorizeSecurityGroupIngress",
      "ec2:AuthorizeSecurityGroupEgress",
      "ec2:DeleteSecurityGroup",
      "ec2:RevokeSecurityGroupIngress",
      "ec2:RevokeSecurityGroupEgress",
      "ec2:DescribeInstances",
      "ec2:DescribeInstanceAttribute",
      "ec2:CreateTags",
      "ec2:RunInstances",
      "ec2:ModifyInstanceAttribute",
      "ec2:TerminateInstances",
      "ec2:ResetInstanceAttribute",
      "ec2:DeleteTags",
      "ec2:DescribeVolumes",
      "ec2:CreateVolume",
      "ec2:DeleteVolume",
      "ec2:DescribeVpcs",
      "ec2:DescribeVpcAttribute",
      "ec2:CreateVPC",
      "ec2:ModifyVPCAttribute",
      "ec2:DeleteVPC",
      "ec2:DescribeSubnets",
      "ec2:CreateSubnet",
      "ec2:ModifySubnetAttribute",
      "ec2:DeleteSubnet",
      "ec2:DescribeInternetGateways",
      "ec2:CreateInternetGateway",
      "ec2:AttachInternetGateway",
      "ec2:DetachInternetGateway",
      "ec2:DeleteInternetGateway",
      "ec2:DescribeRouteTables",
      "ec2:CreateRoute",
      "ec2:CreateRouteTable",
      "ec2:AssociateRouteTable",
      "ec2:DisassociateRouteTable",
      "ec2:DeleteRouteTable",
      "ec2:CreateKeyPair",
      "ec2:ImportKeyPair",
      "ec2:DeleteKeyPair",
      "ec2:DescribeKeyPairs",
      "kms:ListKeys",
      "kms:ListResourceTags",
      "kms:GetKeyPolicy",
      "kms:GetKeyRotationStatus",
      "kms:DescribeKey",
      "kms:CreateKey",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ScheduleKeyDeletion",
      "kms:ListAliases",
      "kms:CreateAlias",
      "kms:DeleteAlias",
      "servicequotas:ListServiceQuotas",
      "ec2:DescribeInternetGateways",
      "ec2:DescribeNatGateways",
      "ec2:DescribeVpnGateways",
      "iam:ListAccountAliases",
      "elasticloadbalancing:DescribeLoadBalancers",
      "elasticloadbalancing:DescribeTargetGroups"
    ]
    resources = ["*"]
  }
}
