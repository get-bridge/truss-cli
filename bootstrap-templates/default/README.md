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

If you are using [tfenv](https://github.com/tfutils/tfenv), it will automagically use
the correct version of Terraform as specified in `.terraform-version`. If you are not,
please make sure you are using the correct version of Terraform.

You'll also need to make sure you have Spinnaker credentials generated and configured
properly to create pipelines. If you're not sure that you have credentials, you can
sign into Google and generate them [here](https://spinnaker-x509.nonprod-cmh.truss.bridgeops.sh/x509/prod).

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

TODO, but will use [Pier's Aurora TF module](https://github.com/get-bridge/bridge-pier-aurora-module).

### Adding a Promotion Pipeline

After creating each workspace's pipeline in the `{{ .TrussDir }}/deploy/` directory,
you will likely want to add a promotion pipeline, which will allow you to trigger a
single pipeline that will in turn deploy to Edge, then Staging, then Production while
validating that each environment is behaving correctly by running your smoketests.

To add the promotion pipeline, in `{{ .TrussDir }}/common/spinnaker.tf`, uncomment the
`promotion-pipeline` module at the bottom of the file.

Then, re-apply the Terraform as follows:

```shell
cd {{ .TrussDir }}/common

aws-vault exec bridge -- terraform init
aws-vault exec bridge -- terraform plan

# Review the plan, and if everything looks good...
aws-vault exec bridge -- terraform apply
```

You will now find your promotion pipeline located at:
https://prod.spinnaker.bridgeops.sh/#/applications/{{ .Params.name }}/executions?pipeline=Promote%20{{ .Params.name }}

See the [bridge-spinnaker-promote-pipeline module docs](https://github.com/get-bridge/bridge-terraform-modules/tree/master/bridge-spinnaker-promote-pipeline) for more information about available options including manual judgements, notifications, etc.

### Triggering Spinnaker Pipelines from Github Actions

Whether you want to add a promotion pipeline or simply trigger a single pipeline,
such as edge-cmh for a new service not yet deployed to Staging or Production, the
process is similar.

If you are NOT using a promotion pipeline, you'll need to enable a webhook trigger
on the single pipeline(s) that you'll be calling. In `{{ .TrussDir }}/deploy/spinnaker.tf`,
you'll need to uncomment the `terraform_remote_state` data resource at the top of the
file, which will be used to obtain the webhook token that was generated in the set of
`common` Terraform. You'll also need to uncomment the `trigger_webhook_source` and
`trigger_webhook_payload_constraints` arguments to the `deploy-pipeline` module. Once
you've made these changes, you'll need to apply them in any workspace whose pipeline you'd
like to trigger:

```shell
cd {{ .TrussDir }}/deploy

aws-vault exec bridge -- terraform init
aws-vault exec bridge -- terraform workspace select <env>-<region> # e.g. edge-cmh
aws-vault exec bridge -- terraform plan

# Review the plan, and if everything looks good...
aws-vault exec bridge -- terraform apply
```

No matter which pipeline you'd like to trigger, you'll need to add the webhook
token as a secret to your Github repository and update your Github workflow to
send a webhook to Spinnaker.

First, make sure that you have a Github Personal Access Token with read-write permission
to manage repository secrets set in a `GITHUB_TOKEN` environment variable.

To add the token, in `{{ .TrussDir }}/common/spinnaker.tf`, uncomment the `github_actions_secret`
resource above the `promotion-pipeline` module. **NOTE**: Please be certain that the `repository`
indicated in this resource is the one you want to set the secret on. If the repo name is wrong,
you will overwrite the wrong secret and break CD for someone else! Once uncommented, re-apply
the `common` Terraform:

```shell
# Ensure that GITHUB_TOKEN is set...
echo $GITHUB_TOKEN

cd {{ .TrussDir }}/common

aws-vault exec bridge -- terraform init
aws-vault exec bridge -- terraform plan

# Review the plan, and if everything looks good...
aws-vault exec bridge -- terraform apply
```

The details of building a Docker image for your service are up to you, but you will want
a Github Actions Workflow that builds your image, pushes it to EC2 Container Registry, and
then triggers a deployment in Spinnaker via webhook. An example of such a workflow is:

```yaml
# .github/workflows/build.yaml
name: Build & Deploy

on:
  push:
    branches:
      - main

jobs:
  build-image:
    name: Build & Deploy
    runs-on: ubuntu-latest
    env:
      ECR_AWS_ACCESS_KEY_ID: {{ "${{ secrets.TRUSS_AWS_ACCESS_KEY_ID }}" }}
      ECR_AWS_SECRET_ACCESS_KEY: {{ "${{ secrets.TRUSS_AWS_SECRET_ACCESS_KEY }}" }}
      ECR_AWS_DEFAULT_REGION: us-east-2
    steps:
      - name: Checkout repo
        uses: actions/checkout@v3

      - name: Checkout actions repo
        uses: actions/checkout@v3
        with:
          repository: get-bridge/actions
          token: {{ "${{ secrets.GIT_HUB_TOKEN }}" }}
          path: .github/actions

      - name: Login to ECR
        uses: ./.github/actions/ecr-auth

      - uses: docker/build-push-action@v3
        env:
          IMAGE_REPO: {{ "${{ env.ECR_REGISTRY }}" }}/{{ .Params.name }}
        with:
          tags:  |
            {{ "${{ env.IMAGE_REPO }}:${{ github.sha }}" }}
            {{ "${{ env.IMAGE_REPO }}:latest" }}
          push: true

      - name: Trigger Spinnaker Deploy
        uses: get-bridge/spinnaker-webhook@v2
        env:
          SPINNAKER_WEBHOOK_HOST: {{ "${{ secrets.SPINNAKER_HOST }}" }}
          SPINNAKER_WEBHOOK_TOKEN: {{ "${{ secrets.SPINNAKER_TRIGGER_TOKEN }}" }}
          SPINNAKER_WEBHOOK_NAME: {{ .Params.name }}-service-github
```

With such a workflow in place, merging to your repo's `main` branch should now build your
application's container image as specified by your `Dockerfile`, push it to EC2 Container
Registry, and send a webhook to Spinnaker, handing it off to be deployed by your pipeline
of choice.

See the [spinnaker-webhook Action's docs](https://github.com/get-bridge/spinnaker-webhook) for more options.
