variable "type" {
  description = "The type of transport configuration block to create, must be one of: [ssh|kubernetes|nomad]"
  type        = string
  validation {
    condition     = contains(["ssh", "kubernetes", "nomad"], var.type)
    error_message = "The transport type must be one of [ssh|kubernetes|nomad]"
  }
}

variable "configuration" {
  description = "The transport configuration, excluding the type"
  type        = any # unfortunately attempting to use an object type (even with optional properties) does not work
}

output "transport_config" {
  description = "The transport specific configuration block for a transport, i.e.  { host = ??, user = ?? } for ssh."
  value = {
    (var.type) = var.configuration
  }
}
