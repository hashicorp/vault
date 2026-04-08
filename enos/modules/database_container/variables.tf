# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

variable "database_type" {
  description = "Type of database to create (postgres, mongodb, mysql)"
  type        = string
  validation {
    condition     = contains(["postgres", "mongodb", "mysql"], var.database_type)
    error_message = "database_type must be one of: postgres, mongodb, mysql"
  }
}

variable "db_version" {
  description = "Database version to use"
  type        = string
}

variable "username" {
  description = "Database username"
  type        = string
}

variable "password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "database" {
  description = "Database name"
  type        = string
}

variable "port" {
  description = "Database port"
  type        = number
}

variable "host" {
  description = "Host configuration with public_ip"
  type = object({
    public_ip  = string
    private_ip = string
    ipv6       = optional(string)
  })
}

variable "instance_name" {
  description = "Unique instance name for the container (defaults to 'default')"
  type        = string
  default     = "default"
}

variable "depends_on_modules" {
  description = "List of modules this depends on"
  type        = list(any)
  default     = []
}
