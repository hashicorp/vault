terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.3.2"
    }
  }
}

data "aws_vpc" "vpc" {
  id = var.vpc_id
}

data "aws_subnets" "vpc" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}

data "aws_kms_key" "kms_key" {
  key_id = var.awskms_unseal_key_arn
}

data "aws_iam_policy_document" "target" {
  statement {
    resources = ["*"]

    actions = [
      "ec2:DescribeInstances",
      "secretsmanager:*"
    ]
  }

  statement {
    resources = [var.awskms_unseal_key_arn]

    actions = [
      "kms:DescribeKey",
      "kms:ListKeys",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:GenerateDataKey"
    ]
  }
}

data "aws_iam_policy_document" "target_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "fleet" {
  statement {
    resources = ["*"]

    actions = [
      "ec2:DescribeImages",
      "ec2:DescribeSubnets",
      "ec2:RequestSpotInstances",
      "ec2:TerminateInstances",
      "ec2:DescribeInstanceStatus",
      "ec2:CancelSpotFleetRequests",
      "ec2:CreateTags",
      "ec2:RunInstances",
      "ec2:StartInstances",
      "ec2:StopInstances",
    ]
  }

  statement {
    effect = "Deny"

    resources = [
      "arn:aws:ec2:*:*:instance/*",
    ]

    actions = [
      "ec2:RunInstances",
    ]

    condition {
      test     = "StringNotEquals"
      variable = "ec2:InstanceMarketType"
      values   = ["spot"]
    }
  }

  statement {
    resources = ["*"]

    actions = [
      "iam:PassRole",
    ]

    condition {
      test     = "StringEquals"
      variable = "iam:PassedToService"
      values = [
        "ec2.amazonaws.com",
      ]
    }
  }

  statement {
    resources = [
      "arn:aws:elasticloadbalancing:*:*:loadbalancer/*",
    ]

    actions = [
      "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
    ]
  }

  statement {
    resources = [
      "arn:aws:elasticloadbalancing:*:*:*/*"
    ]

    actions = [
      "elasticloadbalancing:RegisterTargets"
    ]
  }
}

data "aws_iam_policy_document" "fleet_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["spotfleet.amazonaws.com"]
    }
  }
}

data "enos_environment" "localhost" {}

resource "random_string" "cluster_name" {
  length  = 8
  lower   = true
  upper   = false
  numeric = false
  special = false
}

resource "random_string" "unique_id" {
  length  = 4
  lower   = true
  upper   = false
  numeric = false
  special = false
}

locals {
  instances    = toset([for idx in range(var.instance_count) : tostring(idx)])
  cluster_name = coalesce(var.cluster_name, random_string.cluster_name.result)
  name_prefix  = "${var.project_name}-${local.cluster_name}-${random_string.unique_id.result}"
  fleet_tag    = "${local.name_prefix}-spot-fleet-target"
  fleet_tags = {
    Name      = "${local.name_prefix}-target"
    Type      = local.cluster_name
    SpotFleet = local.fleet_tag
  }
}

resource "aws_iam_role" "target" {
  name               = "${local.name_prefix}-target-role"
  assume_role_policy = data.aws_iam_policy_document.target_role.json
}

resource "aws_iam_instance_profile" "target" {
  name = "${local.name_prefix}-target-profile"
  role = aws_iam_role.target.name
}

resource "aws_iam_role_policy" "target" {
  name   = "${local.name_prefix}-target-policy"
  role   = aws_iam_role.target.id
  policy = data.aws_iam_policy_document.target.json
}

resource "aws_iam_role" "fleet" {
  name               = "${local.name_prefix}-fleet-role"
  assume_role_policy = data.aws_iam_policy_document.fleet_role.json
}

resource "aws_iam_role_policy" "fleet" {
  name   = "${local.name_prefix}-fleet-policy"
  role   = aws_iam_role.fleet.id
  policy = data.aws_iam_policy_document.fleet.json
}

resource "aws_security_group" "target" {
  name        = "${local.name_prefix}-target"
  description = "Target instance security group"
  vpc_id      = var.vpc_id

  # SSH traffic
  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"
    cidr_blocks = [
      "${data.enos_environment.localhost.public_ip_address}/32",
      join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block),
    ]
  }

  # Vault traffic
  ingress {
    from_port = 8200
    to_port   = 8201
    protocol  = "tcp"
    cidr_blocks = flatten([
      "${data.enos_environment.localhost.public_ip_address}/32",
      join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block),
      formatlist("%s/32", var.ssh_allow_ips)
    ])
  }

  # Consul traffic
  ingress {
    from_port = 8301
    to_port   = 8301
    protocol  = "tcp"
    cidr_blocks = [
      "${data.enos_environment.localhost.public_ip_address}/32",
      join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block),
    ]
  }

  ingress {
    from_port = 8301
    to_port   = 8301
    protocol  = "udp"
    cidr_blocks = [
      "${data.enos_environment.localhost.public_ip_address}/32",
      join(",", data.aws_vpc.vpc.cidr_block_associations.*.cidr_block),
    ]
  }

  # Internal traffic
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  # External traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${local.name_prefix}-sg"
    },
  )
}

resource "aws_launch_template" "target" {
  name     = "${local.name_prefix}-target"
  image_id = var.ami_id
  key_name = var.ssh_keypair

  iam_instance_profile {
    name = aws_iam_instance_profile.target.name
  }

  network_interfaces {
    associate_public_ip_address = true
    delete_on_termination       = true
    security_groups             = [aws_security_group.target.id]
  }

  tag_specifications {
    resource_type = "instance"

    tags = merge(
      var.common_tags,
      local.fleet_tags,
    )
  }
}

# There are three primary knobs we can turn to try and optimize our costs by
# using a spot fleet: our min and max instance requirements, our max bid
# price, and the allocation strategy to use when fulfilling the spot request.
# We've currently configured our instance requirements to allow for anywhere
# from 2-4 vCPUs and 4-16GB of RAM. We intentionally have a wide range
# to allow for a large instance size pool to be considered. Our next knob is our
# max bid price. As we're using spot fleets to save on instance cost, we never
# want to pay more for an instance than we were on-demand. We've set the max price
# to equal what we pay for t3.medium instances on-demand, which are the smallest
# reliable size for Vault scenarios. The final knob is the allocation strategy
# that AWS will use when looking for instances that meet our resource and cost
# requirements. We're using the "lowestPrice" strategy to get the absolute
# cheapest machines that will fit the requirements, but it comes with a slightly
# higher capacity risk than say, "capacityOptimized" or "priceCapacityOptimized".
# Unless we see capacity issues or instances being shut down then we ought to
# stick with that strategy.
resource "aws_spot_fleet_request" "targets" {
  allocation_strategy = "lowestPrice"
  fleet_type          = "request"
  iam_fleet_role      = aws_iam_role.fleet.arn
  // Set this to zero so re-runs don't plan for replacement
  instance_pools_to_use_count   = 0
  target_capacity               = var.instance_count
  terminate_instances_on_delete = true
  wait_for_fulfillment          = true

  launch_template_config {
    launch_template_specification {
      id      = aws_launch_template.target.id
      version = aws_launch_template.target.latest_version
    }

    overrides {
      spot_price = var.spot_price_max
      subnet_id  = data.aws_subnets.vpc.ids[0]

      instance_requirements {
        burstable_performance = "included"

        memory_mib {
          min = var.instance_mem_min
          max = var.instance_mem_max
        }

        vcpu_count {
          min = var.instance_cpu_min
          max = var.instance_cpu_max
        }
      }
    }
  }

  tags = merge(
    var.common_tags,
    local.fleet_tags,
  )
}

data "aws_instances" "targets" {
  depends_on = [
    aws_spot_fleet_request.targets,
  ]

  instance_tags = local.fleet_tags
  instance_state_names = [
    "pending",
    "running",
  ]

  filter {
    name   = "image-id"
    values = [var.ami_id]
  }

  filter {
    name   = "iam-instance-profile.arn"
    values = [aws_iam_instance_profile.target.arn]
  }
}

data "aws_instance" "targets" {
  depends_on = [
    aws_spot_fleet_request.targets,
    data.aws_instances.targets
  ]
  for_each = local.instances

  instance_id = data.aws_instances.targets.ids[each.key]
}
