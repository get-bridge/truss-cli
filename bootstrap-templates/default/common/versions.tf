terraform {
  required_version = "~> 1.3"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4"
    }

    github = {
      source  = "hashicorp/github"
      version = "~> 5"
    }

    spinnaker = {
      source  = "get-bridge/spinnaker"
      version = "~> 4"
    }
  }
}
