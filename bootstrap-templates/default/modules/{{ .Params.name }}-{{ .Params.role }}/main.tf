data "kustomization" "app" {
  path = "${path.root}/kustomize"
}

resource "kustomization_resource" "app" {
  for_each = data.kustomization.app.ids
  manifest = data.kustomization.app.manifests[each.value]
}

module "lookups" {
  source = "git@github.com:instructure/truss.git//modules/lookups"
}
