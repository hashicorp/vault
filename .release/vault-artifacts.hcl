schema = 1
artifacts {
  zip = [
    "vault_${version}_linux_amd64.zip",
    "vault_${version}_linux_arm64.zip",
  ]
  rpm = [
    "vault-${version_linux}-1.aarch64.rpm",
    "vault-${version_linux}-1.x86_64.rpm",
  ]
  deb = [
    "vault_${version_linux}-1_amd64.deb",
    "vault_${version_linux}-1_arm64.deb",
  ]
  container = [
    "vault_default_linux_amd64_${version}_${commit_sha}.docker.tar",
    "vault_default_linux_arm64_${version}_${commit_sha}.docker.tar",
    "vault_ubi_linux_amd64_${version}_${commit_sha}.docker.redhat.tar",
  ]
}
