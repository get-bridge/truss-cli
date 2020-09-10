module "{{ .Params.name }}" {
  source = "../../modules/{{ .Params.name }}-{{ .Params.role }}"

  truss_env   = "prod"
  app_env     = "prod"
  region_code = "cmh"
}
