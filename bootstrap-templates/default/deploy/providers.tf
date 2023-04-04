provider "aws" {
  region              = module.truss-metadata.aws_region[local.region]
  allowed_account_ids = ["127178877223"]

  assume_role {
    role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
  }
}

module "kubeconfig" {
  source = "git@github.com:get-bridge/terraform-truss-kubeconfig"

  environment = local.truss_env
  region      = local.region
}

provider "kustomization" {
  kubeconfig_raw = module.kubeconfig.body
}

provider "spinnaker" {
  address   = "https://api-prod.spinnaker.bridgeops.sh"
  cert_path = "~/.spin/shared-prod/spinnaker-client.crt"
  key_path  = "~/.spin/shared-prod/spinnaker-client.key"
}
