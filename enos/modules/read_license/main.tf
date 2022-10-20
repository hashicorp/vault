variable "file_name" {}

output "license" {
  value = file(var.file_name)
}
