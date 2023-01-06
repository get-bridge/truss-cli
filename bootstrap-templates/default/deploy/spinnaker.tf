# data "terraform_remote_state" "common" {
#   backend = "s3"
#   config = {
#     bucket   = "bridge-tfstate"
#     key      = "{{ .Params.name }}/truss/common/terraform.tfstate"
#     region   = "us-west-2"
#     role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
#   }
# }

module "deploy-pipeline" {
  source = "git@github.com:get-bridge/bridge-terraform-modules//bridge-spinnaker-eks-pipeline"

  service     = "{{ .Params.name }}"
  role        = "{{ .Params.role }}"
  environment = local.app_env
  region      = local.region

  app_resources = {
    requests = {
      cpu    = "200m"
      memory = "512Mi"
    }
    limits = {
      cpu    = null
      memory = "512Mi"
    }
  }

  smoketest_job_id = "$${execution['id']}" # TODO: default in the module so I don't have to set it here
  smoketest_image  = "hello-world"         # TODO: replace with real smoketests
  deploy_image     = "{{ .Params.image }}:$${trigger['parameters']['sha']}"
  version_id       = "$${trigger['parameters']['sha']}"

  # We set replicas to null here so that HPA can manage
  deploy_replicas_count = null

  restrict_execution_during_time_window = false

  additional_labels = {
    service-name = "{{ .Params.name }}-{{ .Params.role }}"
  }

  service_port      = {{ .Params.httpPort }}
  health_check_path = "{{ .Params.healthCheckPath }}"

  # If migrations/boot begin taking too long, consider increasing this until
  # migrations are broken out into a separate step.
  health_check_initial_delay = 0

  # trigger_webhook_source = "{{ .Params.name }}-github"
  # trigger_webhook_payload_constraints = {
  #   token = data.terraform_remote_state.common.outputs.spinnaker_webhook_token
  # }

  slack_channel      = "{{ .Params.slackChannel }}"
  enable_manual_gate = false

  failed_message   = "`${local.region}-${local.app_env}` pipeline has failed. $${trigger['parameters']['committer_name']} ($${trigger['parameters']['committer_email']}) `$${trigger['parameters']['message']}` ($${trigger['parameters']['sha']})"
  complete_message = "Finished `$${trigger['parameters']['message']}` ($${trigger['parameters']['sha']})."
  starting_message = "Deploying `$${trigger['parameters']['message']}` ($${trigger['parameters']['sha']})."


  parameters = {
    sha = {
      default     = null
      description = "Git commit SHA"
      label       = null
      required    = true
    }
    message = {
      default     = null
      description = "Git commit message"
      label       = null
      required    = true
    }
    committer_name = {
      default     = null
      description = "Git committer name"
      label       = null
      required    = true
    }
    committer_email = {
      default     = null
      description = "Git committer email"
      label       = null
      required    = true
    }
  }
}
