variable "azure_ad_application_id" {
  description = "The application id for the azure active directory app."
  type        = string
}

output "chart_values" {
  description = "The additional chart values that are required to enable and configure Azure WIF for the Vault servers."
  value = {
    "server" = {
      "extraLabels" = {
        "azure.workload.identity/use" = "true"
      }
      "serviceAccount" = {
        "annotations" = {
          "azure.workload.identity/client-id" = var.azure_ad_application_id
        }
      }
    }
  }
}
