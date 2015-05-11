output "address" {
    value = "${aws_elb.vault.dns_name}"
}
