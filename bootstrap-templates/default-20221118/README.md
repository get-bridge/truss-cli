# Deploying {{ .Params.name }}

This is how you deploy {{ .Params.name }} on Truss! Spinnaker pipelines are available at https://prod.spinnaker.bridgeops.sh/#/applications/{{ .Params.name }}/executions

## Directory Structure

```
{{ .TrussDir }}/
  # The common/ directory contains Terraform to provision resources that are
  # not environment-specific. For example, this directory contains your Truss
  # Tenant, ECR image repository, Spinnaker application, and Spinnaker webhook
  # token, if you are using one.
  common/

  # The deploy/ directory contains environment-specific Terraform and uses Terraform
  # workspaces to provision each environment (e.g. edge-cmh, staging-dub, prod-syd).
  # This includes your application's Kubernetes manifests (e.g. ConfigMap, Service,
  # Ingress) and Spinnaker pipelines for the individual environment, which are responsible
  # for creating and updating your application's Deployment.
  #
  # If your application requires per-environment AWS or other resources, like an RDS
  # database or an S3 bucket, those should also be specified here.
  deploy/

    # The deploy/kustomize/ directory contains a Kustomize base, or set of Kubernetes
    # manifests, that can be customized and applied per environment via the Terraform
    # located in kustomize.tf. This pattern provides a convenient method for using data
    # from Terraform resources in your application's Kubernetes manifests.
    kustomize/
    kustomize.tf

  # secrets/ and secrets.yaml hold configuration and encrypted secret data to be managed
  # via truss-cli's `truss secrets` command. Once pushed to Vault in each cluster, secrets
  # can be easily accessed for use as environment variables in applications running in Truss.
  secrets/
  secrets.yaml
```

## Runbook

### Getting Started

If you are using (tfenv)[https://github.com/tfutils/tfenv], it will automagically use
the correct version of Terraform as specified in `.terraform-version`. If you are not,
please make sure you are using the correct version of Terraform.

First, provision your shared resources:

```shell
cd {{ .TrussDir }}/common

aws-vault exec bridge -- terraform init
aws-vault exec bridge -- terraform plan

# Review the plan, and if everything looks good...
aws-vault exec bridge -- terraform apply
```

Then, for each environment you want to deploy to, provision environment-specific resources:

```shell
cd {{ .TrussDir }}/deploy

aws-vault exec bridge -- terraform init
aws-vault exec bridge -- terraform workspace new <env>-<region> # e.g. edge-cmh
aws-vault exec bridge -- terraform plan

# Review the plan, and if everything looks good...
aws-vault exec bridge -- terraform apply
```

Finally, before deploying, you'll also need to push secrets.

To push secrets to a single environment:

```shell
cd {{ .TrussDir }}
aws-vault exec bridge -- truss secrets push <env>-<region> # e.g. edge-cmh
```

To push secrets to all environments:
```shell
cd {{ .TrussDir }}
aws-vault exec bridge -- truss secrets push --all
```

### Adding an RDS database

TODO

### Adding a Promotion Pipeline

TODO

### Triggering Spinnaker Pipelines from Github Actions

TODO
