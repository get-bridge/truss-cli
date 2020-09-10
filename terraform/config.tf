terraform {
  backend "s3" {
    bucket   = "bridge-shared-terraform-us-east-2"
    key      = "truss-cli/terraform/terraform.tfstate"
    region   = "us-east-2"
    acl      = "bucket-owner-full-control"
    encrypt  = true
    role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
  }
}

provider "aws" {
  region              = "us-east-2"
  allowed_account_ids = ["127178877223"]

  assume_role {
    role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
  }
}

provider "github" {
  organization = "instructure-bridge"
}
