resource "template_file" "install" {
    filename = "${path.module}/scripts/install.sh.tpl"

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

    tag {
        key = "Name"
        value = "vault"
        propagate_at_launch = true
    }
}

resource "aws_launch_configuration" "vault" {
    image_id = "${var.ami}"
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

    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        from_port = 8200
        to_port = 8200
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}
