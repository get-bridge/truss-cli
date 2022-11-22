resource "aws_ecr_repository" "{{ .Params.name }}" {
  name = "{{ .Params.name }}"

  image_scanning_configuration {
    scan_on_push = true
  }
}
