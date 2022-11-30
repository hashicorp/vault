output "ci_role" {
  value = {
    name = aws_iam_role.role.name
    arn  = aws_iam_role.role.arn
  }
}

output "ci_role_policy" {
  value = {
    name   = aws_iam_role_policy.role_policy.name
    policy = aws_iam_role_policy.role_policy.policy
  }
}
