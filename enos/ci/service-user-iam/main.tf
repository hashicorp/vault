# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

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
  policy   = data.aws_iam_policy_document.role_policy.json
}

data "aws_iam_policy_document" "role_policy" {
  source_policy_documents = [
    data.aws_iam_policy_document.enos_scenario.json,
    data.aws_iam_policy_document.aws_nuke.json,
  ]
}

data "aws_iam_policy_document" "aws_nuke" {
  provider = aws.us_east_1

  statement {
    effect = "Allow"
    actions = [
      "ec2:DescribeInternetGateways",
      "ec2:DescribeNatGateways",
      "ec2:DescribeRegions",
      "ec2:DescribeVpnGateways",
      "iam:DeleteAccessKey",
      "iam:DeleteUser",
      "iam:DeleteUserPolicy",
      "iam:GetUser",
      "iam:ListAccessKeys",
      "iam:ListAccountAliases",
      "iam:ListGroupsForUser",
      "iam:ListUserPolicies",
      "iam:ListUserTags",
      "iam:ListUsers",
      "iam:UntagUser",
      "servicequotas:ListServiceQuotas"
    ]

    resources = ["*"]
  }
}

data "aws_iam_policy_document" "enos_scenario" {
  provider = aws.us_east_1

  statement {
    effect = "Allow"
    actions = [
      "ec2:AssociateRouteTable",
      "ec2:AttachInternetGateway",
      "ec2:AuthorizeSecurityGroupEgress",
      "ec2:AuthorizeSecurityGroupIngress",
      "ec2:CancelSpotFleetRequests",
      "ec2:CancelSpotInstanceRequests",
      "ec2:CreateEgressOnlyInternetGateway",
      "ec2:CreateInternetGateway",
      "ec2:CreateKeyPair",
      "ec2:CreateFleet",
      "ec2:CreateLaunchTemplate",
      "ec2:CreateLaunchTemplateVersion",
      "ec2:CreateRoute",
      "ec2:CreateRouteTable",
      "ec2:CreateSecurityGroup",
      "ec2:CreateSpotDatafeedSubscription",
      "ec2:CreateSubnet",
      "ec2:CreateTags",
      "ec2:CreateVolume",
      "ec2:CreateVPC",
      "ec2:DeleteEgressOnlyInternetGateway",
      "ec2:DeleteFleets",
      "ec2:DeleteInternetGateway",
      "ec2:DeleteLaunchTemplate",
      "ec2:DeleteLaunchTemplateVersions",
      "ec2:DeleteKeyPair",
      "ec2:DeleteRoute",
      "ec2:DeleteRouteTable",
      "ec2:DeleteSecurityGroup",
      "ec2:DeleteSpotDatafeedSubscription",
      "ec2:DeleteSubnet",
      "ec2:DeleteTags",
      "ec2:DeleteVolume",
      "ec2:DeleteVPC",
      "ec2:DescribeAccountAttributes",
      "ec2:DescribeAvailabilityZones",
      "ec2:DescribeEgressOnlyInternetGateways",
      "ec2:DescribeFleets",
      "ec2:DescribeFleetHistory",
      "ec2:DescribeFleetInstances",
      "ec2:DescribeImages",
      "ec2:DescribeInstanceAttribute",
      "ec2:DescribeInstanceCreditSpecifications",
      "ec2:DescribeInstances",
      "ec2:DescribeInstanceTypeOfferings",
      "ec2:DescribeInstanceTypes",
      "ec2:DescribeInternetGateways",
      "ec2:DescribeKeyPairs",
      "ec2:DescribeLaunchTemplates",
      "ec2:DescribeLaunchTemplateVersions",
      "ec2:DescribeNatGateways",
      "ec2:DescribeNetworkAcls",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DescribeRegions",
      "ec2:DescribeRouteTables",
      "ec2:DescribeSecurityGroups",
      "ec2:DescribeSpotDatafeedSubscription",
      "ec2:DescribeSpotFleetInstances",
      "ec2:DescribeSpotFleetInstanceRequests",
      "ec2:DescribeSpotFleetRequests",
      "ec2:DescribeSpotFleetRequestHistory",
      "ec2:DescribeSpotInstanceRequests",
      "ec2:DescribeSpotPriceHistory",
      "ec2:DescribeSubnets",
      "ec2:DescribeTags",
      "ec2:DescribeVolumes",
      "ec2:DescribeVpcAttribute",
      "ec2:DescribeVpcClassicLink",
      "ec2:DescribeVpcClassicLinkDnsSupport",
      "ec2:DescribeVpcs",
      "ec2:DescribeVpnGateways",
      "ec2:DetachInternetGateway",
      "ec2:DisassociateRouteTable",
      "ec2:GetLaunchTemplateData",
      "ec2:GetSpotPlacementScores",
      "ec2:ImportKeyPair",
      "ec2:ModifyFleet",
      "ec2:ModifyInstanceAttribute",
      "ec2:ModifyLaunchTemplate",
      "ec2:ModifySpotFleetRequest",
      "ec2:ModifySubnetAttribute",
      "ec2:ModifyVPCAttribute",
      "ec2:RequestSpotInstances",
      "ec2:RequestSpotFleet",
      "ec2:ResetInstanceAttribute",
      "ec2:RevokeSecurityGroupEgress",
      "ec2:RevokeSecurityGroupIngress",
      "ec2:RunInstances",
      "ec2:SendSpotInstanceInterruptions",
      "ec2:TerminateInstances",
      "elasticloadbalancing:DescribeLoadBalancers",
      "elasticloadbalancing:DescribeTargetGroups",
      "iam:AddRoleToInstanceProfile",
      "iam:AttachRolePolicy",
      "iam:CreateInstanceProfile",
      "iam:CreatePolicy",
      "iam:CreateRole",
      "iam:CreateServiceLinkedRole",
      "iam:DeleteInstanceProfile",
      "iam:DeletePolicy",
      "iam:DeleteRole",
      "iam:DeleteRolePolicy",
      "iam:DetachRolePolicy",
      "iam:GetInstanceProfile",
      "iam:GetRole",
      "iam:GetRolePolicy",
      "iam:ListAccountAliases",
      "iam:ListAttachedRolePolicies",
      "iam:ListInstanceProfiles",
      "iam:ListInstanceProfilesForRole",
      "iam:ListPolicies",
      "iam:ListRolePolicies",
      "iam:ListRoles",
      "iam:PassRole",
      "iam:PutRolePolicy",
      "iam:RemoveRoleFromInstanceProfile",
      "kms:CreateAlias",
      "kms:CreateKey",
      "kms:Decrypt",
      "kms:DeleteAlias",
      "kms:DescribeKey",
      "kms:Encrypt",
      "kms:GetKeyPolicy",
      "kms:GetKeyRotationStatus",
      "kms:ListAliases",
      "kms:ListKeys",
      "kms:ListResourceTags",
      "kms:ScheduleKeyDeletion",
      "kms:TagResource",
      "servicequotas:ListServiceQuotas"
    ]

    resources = ["*"]
  }
}
