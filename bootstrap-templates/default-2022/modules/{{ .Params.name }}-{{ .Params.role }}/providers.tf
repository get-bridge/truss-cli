provider "aws" {
  region              = module.lookups.region_lookup[var.region_code]
  allowed_account_ids = ["127178877223"]

  assume_role {
    role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
  }
}

provider "aws" {
  alias               = "cmh"
  region              = "us-east-2"
  allowed_account_ids = ["127178877223"]

  assume_role {
    role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
  }
}

data "aws_s3_bucket_object" "kubeconfig" {
  bucket   = "truss-kubeconfig-us-east-2"
  key      = "kubeconfig-truss-${var.truss_env}-${var.region_code}"
  provider = aws.cmh
}

provider "kustomization" {
  kubeconfig_raw = data.aws_s3_bucket_object.kubeconfig.body
}
