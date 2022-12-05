terraform {
  backend "s3" {
    bucket   = "bridge-tfstate"
    key      = "{{ .Params.name }}/truss/deploy/terraform.tfstate"
    region   = "us-west-2"
    acl      = "bucket-owner-full-control"
    encrypt  = true
    role_arn = "arn:aws:iam::127178877223:role/xacct/ops-admin"
  }
}
