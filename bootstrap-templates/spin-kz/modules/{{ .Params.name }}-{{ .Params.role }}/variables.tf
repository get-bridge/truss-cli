locals {
  app = "{{ .Params.name }}"
}

variable "truss_env" {
  type        = string
  description = "Truss environment, i.e. nonprod, prod, dev"
}

variable "app_env" {
  type        = string
  description = "App environment, i.e. edge, staging, perf, prod"
}

variable "region_code" {
  type        = string
  description = "Short region code, i.e. cmh, iad, syd, dub"
}

variable "account" {
  type    = string
  default = "bridge-shared"
}
