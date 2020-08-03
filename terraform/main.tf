// iam stuff
// github action secret

resource "aws_iam_user" "publisher" {
  name = "truss-cli-github-actions"
  path = "/github/"
}

resource "aws_iam_user_policy" "publisher" {
  name = "truss-cli-github-actions-policy"
  user = aws_iam_user.publisher.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_access_key" "publisher" {
  user = aws_iam_user.publisher.name
}

resource "github_actions_secret" "aws_access_key" {
  repository      = "truss-cli"
  secret_name     = "AWS_ACCESS_KEY_ID"
  plaintext_value = aws_iam_access_key.publisher.id
}

resource "github_actions_secret" "aws_secret_key" {
  repository      = "truss-cli"
  secret_name     = "AWS_SECRET_ACCESS_KEY"
  plaintext_value = aws_iam_access_key.publisher.secret
}
