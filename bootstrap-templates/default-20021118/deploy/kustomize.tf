locals {
  vault_path = "vault:secret/data/bridge/${local.app_env}/${local.region}/{{ .Params.name }}/{{ .Params.role }}/default"
  hostname   = "hello-jdharrington-${local.app_env}.${local.truss_env}-${local.region}.truss.bridgeops.sh"
}

data "kustomization_overlay" "{{ .Params.name }}" {
  common_labels = {
    app    = "{{ .Params.name }}"
    env    = local.app_env
    region = local.region
    role   = "{{ .Params.role }}"
  }

  name_prefix = "{{ .Params.name }}-"

  namespace = "{{ .Params.name }}-${local.app_env}"

  resources = ["${path.root}/kustomize"]

  patches {
    target {
      kind = "Ingress"
      name = "api"
    }

    patch = <<-EOF
    - op: replace
      path: /spec/rules/0/host
      value: ${local.hostname}
    - op: replace
      path: /spec/tls/0/hosts
      value: [${local.hostname}]
    EOF
  }

  patches {
    target {
      kind = "ConfigMap"
      name = "{{ .Params.role }}"
    }

    patch = <<-EOF
      kind: ConfigMap
      metadata:
        name: {{ .Params.role }}
      data:
        EXAMPLE_ENV_SPECIFIC_VAR: ${local.app_config[terraform.workspace]["EXAMPLE_ENV_SPECIFIC_VAR"]}
        EXAMPLE_SECRET: ${local.vault_path}#EXAMPLE_SECRET
    EOF
  }
}

resource "kustomization_resource" "{{ .Params.name }}" {
  for_each = data.kustomization_overlay.{{ .Params.name }}.ids
  manifest = data.kustomization_overlay.{{ .Params.name }}.manifests[each.value]
}
