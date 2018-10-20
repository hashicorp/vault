//-------------------------------------------------------------------
// Vault settings
//-------------------------------------------------------------------

variable "download-url" {
    default = "https://releases.hashicorp.com/vault/0.11.3/vault_0.11.3_linux_amd64.zip"
    description = "URL to download Vault"
}

variable "config" {
    description = "Configuration (text) for Vault"
}

variable "extra-install" {
    default = ""
    description = "Extra commands to run in the install script"
}

//-------------------------------------------------------------------
// AWS settings
//-------------------------------------------------------------------

variable "ami" {
    default = "ami-7eb2a716"
    description = "AMI for Vault instances"
}

variable "availability-zones" {
    default = "us-east-1a,us-east-1b"
    description = "Availability zones for launching the Vault instances"
}

variable "elb-health-check" {
    default = "HTTP:8200/v1/sys/health"
    description = "Health check for Vault servers"
}

variable "instance_type" {
    default = "m3.medium"
    description = "Instance type for Vault instances"
}

variable "key-name" {
    default = "default"
    description = "SSH key name for Vault instances"
}

variable "nodes" {
    default = "2"
    description = "number of Vault instances"
}

variable "subnets" {
    description = "list of subnets to launch Vault within"
}

variable "vpc-id" {
    description = "VPC ID"
}

variable "source-networks" {
    type = "list"
    #not secure by default, get your IP from checkip.amazonaws.com and set to something like:
    #default = ["108.92.149.99/32","34.222.225.77/32","18.237.144.183/32"]
    default = ["0.0.0.0/0"]
    description = "List of source networks from which to allow access"
}
