module "{{ .Params.name }}" {
  source = "../../modules/{{ .Params.name }}-{{ .Params.role }}"

  truss_env   = "nonprod"
  app_env     = "edge"
  region_code = "cmh"
}
