provider "aws" {
  region              = module.bridge_lookups.region_lookup[var.region_code]
  allowed_account_ids = ["127178877223"]

  assume_role {
    role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
  }
}

provider "spinnaker" {
  address   = "https://api-prod.spinnaker.bridgeops.sh"
  cert_path = "~/.spin/shared-prod/spinnaker-client.crt"
  key_path  = "~/.spin/shared-prod/spinnaker-client.key"
}

module "bridge_lookups" {
  source = "git@github.com:instructure/truss.git//modules/lookups"
}

# provider "aws" {
#   alias               = "cmh"
#   region              = "us-east-2"
#   allowed_account_ids = ["127178877223"]

#   assume_role {
#     role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
#   }
# }

# data "aws_s3_bucket_object" "kubeconfig" {
#   bucket   = "truss-kubeconfig-us-east-2"
#   key      = "kubeconfig-truss-${var.truss_env}-${var.region_code}"
#   provider = aws.cmh
# }
