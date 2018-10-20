
output "address" {
    value = "${aws_lb.vault.dns_name}"
}

output "vault_profile" {
    value = <<EOF
    
    export VAULT_ADDR=\"https://${aws_lb.vault.dns_name}\""
    vault status -tls-skip-verify
EOF
}

// Can be used to add additional SG rules to Vault instances.
output "vault_security_group" {
    value = "${aws_security_group.vault.id}"
}
