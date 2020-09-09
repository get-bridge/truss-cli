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

resource "spinnaker_application" "application" {
  name          = "{{ .Params.name }}"
  email         = "bridge-eng@instructure.com"
  instance_port = {{ .Params.httpPort }}

  permissions {
    read    = ["bridge-engineering-all"]
    write   = ["bridge-engineering-all"]
    execute = ["bridge-engineering-all"]
  }
}