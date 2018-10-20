resource "template_file" "install" {
    template = "${file("${path.module}/scripts/install.sh.tpl")}"

    vars {
        download_url  = "${var.download-url}"
        config        = "${var.config}"
        extra-install = "${var.extra-install}"
    }
}

// We launch Vault into an ASG so that it can properly bring them up for us.
resource "aws_autoscaling_group" "vault" {
    name = "vault - ${aws_launch_configuration.vault.name}"
    launch_configuration = "${aws_launch_configuration.vault.name}"
    availability_zones = ["${split(",", var.availability-zones)}"]
    min_size = "${var.nodes}"
    max_size = "${var.nodes}"
    desired_capacity = "${var.nodes}"
    health_check_grace_period = 15
    health_check_type = "EC2"
    vpc_zone_identifier = ["${split(",", var.subnets)}"]
    /*
    load_balancers = ["${aws_lb.vault.id}"]
    */
    
    tag {
        key = "Name"
        value = "vault"
        propagate_at_launch = true
    }
}

resource "aws_launch_configuration" "vault" {
    #image_id = "${var.ami}"
    image_id = "${data.aws_ami.latest_ubuntu_image.id}"
    instance_type = "${var.instance_type}"
    key_name = "${var.key-name}"
    security_groups = ["${aws_security_group.vault.id}"]
    user_data = "${template_file.install.rendered}"

}

// Security group for Vault allows SSH and HTTP access (via "tcp" in
// case TLS is used)
resource "aws_security_group" "vault" {
    name = "vault"
    description = "Vault servers"
    vpc_id = "${var.vpc-id}"
    
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "vault-ssh" {
    security_group_id = "${aws_security_group.vault.id}"
    type = "ingress"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["${var.source-networks}"]

  lifecycle {
    create_before_destroy = true
  }
}

// This rule allows Vault HTTP API access to individual nodes, since each will
// need to be addressed individually for unsealing.
resource "aws_security_group_rule" "vault-http-api" {
    security_group_id = "${aws_security_group.vault.id}"
    type = "ingress"
    from_port = 8200
    to_port = 8200
    protocol = "tcp"
    source_security_group_id = "${aws_security_group.elb.id}"
}

resource "aws_security_group_rule" "vault-egress" {
    security_group_id = "${aws_security_group.vault.id}"
    type = "egress"
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
}


resource "aws_lb" "vault" {
  name               = "vault-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = ["${aws_security_group.elb.id}"]
  subnets            = ["${split(",", var.subnets)}"]

  enable_deletion_protection = true
/*
  access_logs {
    bucket  = "${aws_s3_bucket.lb_logs.bucket}"
    prefix  = "test-lb"
    enabled = true
  }

  tags {
    Environment = "production"
  }
*/
}

resource "aws_lb_target_group" "vault-lb-tgt-group" {
  name     = "vault-lb-tgt"
  port     = 8200
  protocol = "HTTP"
  vpc_id = "${var.vpc-id}"
}

resource "aws_lb_listener" "vault-http" {
  load_balancer_arn = "${aws_lb.vault.arn}"
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.vault-lb-tgt-group.arn}"
  }
}

resource "aws_lb_listener" "vault-https" {
  load_balancer_arn = "${aws_lb.vault.arn}"
  port              = "443"
  protocol          = "HTTPS"
  
  ssl_policy        = "ELBSecurityPolicy-2015-05"
  certificate_arn   = "${aws_iam_server_certificate.vault_cert.arn}"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.vault-lb-tgt-group.arn}"
  }
}

# Create a new ALB Target Group attachment
resource "aws_autoscaling_attachment" "asg_attachment_bar" {
  autoscaling_group_name = "${aws_autoscaling_group.vault.id}"
  alb_target_group_arn   = "${aws_lb_target_group.vault-lb-tgt-group.arn}"
}

#https://www.terraform.io/docs/providers/tls/r/self_signed_cert.html
resource "tls_private_key" "example" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

# run vault cmd with -tls-skip-verify flag for now
resource "tls_self_signed_cert" "example" {
  key_algorithm   = "${tls_private_key.example.algorithm}"
  private_key_pem = "${tls_private_key.example.private_key_pem}"

  subject {
    common_name  = "testvault"
    organization = "BAD Examples, Inc"
  }

  validity_period_hours = 6

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth"
  ]
  
  
  lifecycle {
    create_before_destroy = true
  }
}
resource "aws_iam_server_certificate" "vault_cert" {
  name      = "vault-cert"
  certificate_body = "${tls_self_signed_cert.example.cert_pem}"
  private_key      = "${tls_private_key.example.private_key_pem}"
  
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "elb" {
    name = "vault-elb"
    description = "Vault ELB"
    vpc_id = "${var.vpc-id}"
}

resource "aws_security_group_rule" "vault-elb-http" {
    security_group_id = "${aws_security_group.elb.id}"
    type = "ingress"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["${var.source-networks}"]
}

resource "aws_security_group_rule" "vault-elb-https" {
    security_group_id = "${aws_security_group.elb.id}"
    type = "ingress"
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["${var.source-networks}"]
}

resource "aws_security_group_rule" "vault-elb-egress" {
    security_group_id = "${aws_security_group.elb.id}"
    type = "egress"
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
}


# automatic lookup based on   #https://aws.amazon.com/amazon-linux-ami/
# aws ec2 describe-images --owners 099720109477 --filters 'Name=name,Values=ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-????????' 'Name=state,Values=available' --output json | jq -r '.Images | sort_by(.CreationDate) | last(.[]).ImageId'

data "aws_ami" "latest_ubuntu_image" {
  most_recent = true

  owners = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-*"]
  }

}
