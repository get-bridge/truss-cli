provider "spinnaker" {
  address   = "https://api-prod.spinnaker.bridgeops.sh"
  cert_path = "~/.spin/shared-prod/spinnaker-client.crt"
  key_path  = "~/.spin/shared-prod/spinnaker-client.key"
}

module "deploy-pipeline" {
  source = "git::ssh://gerrit.instructure.com:29418/bridge-terraform-modules//bridge-spinnaker-eks-pipeline"

  service     = "{{ .Params.name }}"
  role        = "{{ .Params.role }}"
  environment = var.app_env
  region      = var.region_code

  smoketest_image       = "hello-world"
  deploy_image          = "{{ .Params.image }}:$${trigger['properties']['sha']}"
  version_id            = "$${trigger['properties']['sha']}"
  deploy_replicas_count = 3

  additional_labels = {
    service-name = "{{ .Params.name }}-{{ .Params.role }}"
  }
  service_port      = {{ .Params.httpPort }}
  health_check_path = "/health-check"

  trigger_jenkins_job = "{{ .Params.trigger_jenkins_job }}"

  # slack_channel      = "#bridge-noisy"
  enable_manual_gate = var.app_env == "prod"
}
