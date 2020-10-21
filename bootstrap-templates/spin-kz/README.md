# Deploying {{ .Params.name }}

This is how you deploy {{ .Params.name }} on Truss! Spinnaker pipelins are available at https://prod.spinnaker.bridgeops.sh/#/applications/{{ .Params.name }}/executions

## Directory Structure

- `{{ .TrussDir }}/` - Holds your tenant configuration and Spinnaker Application
  - `/modules/{{ .Params.name }}-{{ .Params.role }}` - Holds the deployment configuration for a single instance of {{ .Params.name }}
    - `/kustomize` - Kubernetes manifests for {{ .Params.name }}'s common infrastructure components
      - `/deployment` -  Base config for deployment
      - `/postdeploy` - Base config for post-deploy job
      - `/predeploy` - Base config for pre-deploy job
  - `/{edge|staging|prod}/{cmh|dub|syd}/` - Holds deployment configuration for a given environment/region of {{ .Params.name }}
    - `/kustomize` - Environment-specific infra overrides
      - `/deployment` - Environment-speicifc deployment overrides
      - `/postdeploy` - Environment-specific post-deploy overrides
      - `/predeploy` - Environment-specific pre-deploy overrides

## Runbook

- Provision your tenant: `cd {{ .TrussDir }} && terraform init && terraform apply`
- Provision a given environment (i.e. edge-cmh): `cd {{ .TrussDir }}/edge/cmh && terraform init && terraform apply`
- Retrieve realtime logs (i.e. edge-cmh): `truss wrap -e cmh-edge -- kubectl -n {{ .Params.name }}-edge logs -c {{ .Params.name }}-{{ .Params.role }} deployment/{{ .Params.name }}-{{ .Params.role }}`