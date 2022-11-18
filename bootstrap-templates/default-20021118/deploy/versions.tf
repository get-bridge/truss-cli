terraform {
  required_version = "~> 1.3"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4"
    }

    kustomization = {
      source  = "kbst/kustomization"
      version = "~> 0.9"
    }

    spinnaker = {
      source  = "get-bridge/spinnaker"
      version = "~> 4"
    }
  }
}
