# Uncomment to generate a random token to use when triggering Spinnaker
# pipelines from Github Actions. Make sure `repository` is set correctly
# in the `github_actions_secret`. If not, you may overwrite another repo's
# secret by accident.
#
# resource "random_password" "webhook_token" {
#   length  = 32
#   special = false
# }

# resource "github_actions_secret" "spinnaker_token" {
#   #
#   repository      = "{{ .Params.name }}"
#   secret_name     = "SPINNAKER_TRIGGER_TOKEN"
#   plaintext_value = random_password.webhook_token.result
# }

resource "spinnaker_application" "application" {
  name  = "{{ .Params.name }}"
  email = "bridge-engineering-all@getbridge.com"

  permissions {
    read    = ["bridge-engineering-all"]
    write   = ["bridge-engineering-all"]
    execute = ["bridge-engineering-all"]
  }
}
