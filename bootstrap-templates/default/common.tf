module "truss-tenant" {
  source   = "git::ssh://git@github.com/instructure/truss//modules/truss-tenant"
  name     = "{{ .Params.name }}"
  istio    = true
  starlord = {{ .Params.starlord }}
  apps = [{
    name = "{{ .Params.role }}"
    vault = [{
      path         = "secret/data/bridge/{env}/{region}/shared/*"
      capabilities = ["read", "list"]
      }, {
      path         = "secret/data/bridge/{env}/{region}/{{ .Params.name }}/*"
      capabilities = ["read", "list"]
    }]
  }]
}
