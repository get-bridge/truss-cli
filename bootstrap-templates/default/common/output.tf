output "spinnaker_webhook_token" {
  value     = random_password.webhook_token.result
  sensitive = true
}
