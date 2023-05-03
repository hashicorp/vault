output "application_id" {
  description = "The application id of the vault app"
  value       = azuread_application.vault_app.application_id
}

output "client_secret" {
  description = "The password for the the application service principal, to be used as the client_secret when configuring the vault azure auth method."
  value       = azuread_service_principal_password.password.value
  sensitive   = true
}
