module "truss-tenant" {
  source = "git@github.com:get-bridge/truss.git//modules/truss-tenant"

  name = "{{ .Params.name }}"

  # TODO: default istio = false in the tenant module
  istio = false

  apps = [{
    name = "{{ .Params.role }}"
    vault = [{
      path         = "secret/data/bridge/{env}/{region}/{{ .Params.name }}/{{ .Params.role }}/*"
      capabilities = ["read", "list"]
    }]
  }]
}
