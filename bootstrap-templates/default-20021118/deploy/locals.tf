locals {
  workspace_args = split("-", terraform.workspace)

  app_env   = local.workspace_args[0]
  region    = local.workspace_args[1]
  truss_env = local.app_env == "prod" ? "prod" : "nonprod"

  # This can be organized however you like, just provided as an example.
  app_config = {
    "edge-cmh" = {
      EXAMPLE_ENV_SPECIFIC_VAR = "value-for-edge-cmh"
    }
  }
}
