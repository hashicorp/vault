output "address" {
    value = "${aws_elb.vault.dns_name}"
}

// Can be used to add additional SG rules to Vault instances.
output "vault_security_group" {
    value = "${aws_security_group.vault.id}"
}

// Can be used to add additional SG rules to the Vault ELB.
output "elb_security_group" {
    value = "${aws_security_group.elb.id}"
}
