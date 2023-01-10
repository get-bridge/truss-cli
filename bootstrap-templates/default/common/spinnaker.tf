resource "spinnaker_application" "application" {
  name  = "{{ .Params.name }}"
  email = "bridge-engineering-all@getbridge.com"

  permissions {
    read    = ["bridge-engineering-all"]
    write   = ["bridge-engineering-all"]
    execute = ["bridge-engineering-all"]
  }
}

resource "random_password" "webhook_token" {
  length  = 32
  special = false
}

# Uncomment this resource if you will be triggering Spinnaker pipelines
# from Github Actions
#
# resource "github_actions_secret" "spinnaker_token" {
#   repository      = "{{ .Params.name }}"
#   secret_name     = "SPINNAKER_TRIGGER_TOKEN"
#   plaintext_value = random_password.webhook_token.result
# }

# Uncomment the module below to add a promotion pipeline
#
# module "promotion-pipeline" {
#   source = "git@github.com:get-bridge/bridge-terraform-modules.git//bridge-spinnaker-promote-pipeline"

#   service = "{{ .Params.name }}"

#   parameters = {
#     sha = {
#       default     = null
#       description = "Git commit SHA"
#       label       = null
#       required    = true
#     }
#     message = {
#       default     = null
#       description = "Git commit message"
#       label       = null
#       required    = true
#     }
#     committer_name = {
#       default     = null
#       description = "Git committer name"
#       label       = null
#       required    = true
#     }
#     committer_email = {
#       default     = null
#       description = "Git committer email"
#       label       = null
#       required    = true
#     }
#   }

#   # Configure Trigger
#   trigger_webhook_source              = "{{ .Params.name }}-github"
#   trigger_webhook_payload_constraints = { token = random_password.webhook_token.result }

#   # Configure Edge
#   edge_pipelines        = [for r in ["cmh"] : "Deploy api ${r}-edge"]
#   edge_manual_judgement = false

#   # Configure Staging
#   staging_pipelines        = [for r in ["cmh", "dub", "syd"] : "Deploy api ${r}-staging"]
#   staging_manual_judgement = false

#   # Configure Prod
#   prod_pipelines        = [for r in ["cmh", "dub", "syd"] : "Deploy api ${r}-prod"]
#   prod_manual_judgement = false

#   # Configure Slack Notifications
#   slack_channel         = "{{ .Params.slackChannel }}"
#   slack_pipeline_alerts = true
#   complete_message      = "Promotion Completed"
#   failed_message        = "Promotion Failed"
#   starting_message      = "Starting Promotion"

#   # Configure Manual Judgement Notifications
#   slack_judgement_prompts = false
# }
