module "deploy-pipeline" {
  source = "git::ssh://gerrit.instructure.com:29418/bridge-terraform-modules//bridge-spinnaker-eks-pipeline"

  service     = "{{ .Params.name }}"
  role        = "{{ .Params.role }}"
  environment = var.app_env
  region      = var.region_code

  version_id      = var.app_env == "prod" ? local.prod_version_id : local.nonprod_version_id
  commit_message  = var.app_env == "prod" ? local.prod_commit_message : local.nonprod_commit_message
  committer_name  = var.app_env == "prod" ? local.prod_committer_name : local.nonprod_committer_name
  committer_email = var.app_env == "prod" ? local.prod_committer_email : local.nonprod_committer_email

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

  infra_kustomize = {
    artifact_account    = "inst-bridge-github"
    github_repo         = "{{ .Params.githubRepo }}"
    github_branch       = "master"
    checkout_subpath    = "{{ .TrussDir }}"
    kustomize_file_path = "{{ .TrussDir }}/${var.app_env}/${var.region_code}/kustomization.yaml"
  }

  # predeploy_kustomize = {
  #   artifact_account    = "inst-bridge-github"
  #   github_repo         = "{{ .Params.githubRepo }}"
  #   github_branch       = "master"
  #   checkout_subpath    = "{{ .TrussDir }}"
  #   kustomize_file_path = "{{ .TrussDir }}/${var.app_env}/${var.region_code}/predeploy/kustomization.yaml"
  # }

  deploy_kustomize = {
    artifact_account    = "inst-bridge-github"
    github_repo         = "{{ .Params.githubRepo }}"
    github_branch       = "master"
    checkout_subpath    = "{{ .TrussDir }}"
    kustomize_file_path = "{{ .TrussDir }}/${var.app_env}/${var.region_code}/deployment/kustomization.yaml"
  }

  # postdeploy_kustomize = {
  #   artifact_account    = "inst-bridge-github"
  #   github_repo         = "{{ .Params.githubRepo }}"
  #   github_branch       = "master"
  #   checkout_subpath    = "{{ .TrussDir }}"
  #   kustomize_file_path = "{{ .TrussDir }}/${var.app_env}/${var.region_code}/postdeploy/kustomization.yaml"
  # }

  smoketest_image = "{{ .Params.smoketestImage }}"

  # slack_channel      = "#bridge_noisy"
  enable_manual_gate = var.app_env == "prod"
}


locals {
  nonprod_version_id      = "$${trigger['parameters']['sha']}"
  nonprod_commit_message  = "$${trigger['parameters']['message']}"
  nonprod_committer_name  = "$${trigger['parameters']['committer_name']}"
  nonprod_committer_email = "$${trigger['parameters']['committer_email']}"
  prod_version_id         = "$${trigger['parentExecution']['trigger']['parameters']['sha']}"
  prod_commit_message     = "$${trigger['parentExecution']['trigger']['parameters']['message']}"
  prod_committer_name     = "$${trigger['parentExecution']['trigger']['parameters']['committer_name']}"
  prod_committer_email    = "$${trigger['parentExecution']['trigger']['parameters']['committer_email']}"
}
