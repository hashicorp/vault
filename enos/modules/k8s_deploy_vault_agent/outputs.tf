output "transport" {
  description = "The transport configuration to use when connecting to the deployed vault agent pod."
  value = {
    kubernetes = data.enos_kubernetes_pods.vault_agent_pods.transports[0]
  }
}
