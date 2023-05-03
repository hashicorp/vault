variable "node_count" {
  description = "The number of nodes in the AKS cluster"
  type        = number
  default     = 3
}

variable "kubeconfig_path" {
  description = "The path to the kubeconfig file that should be written"
  type        = string
}
