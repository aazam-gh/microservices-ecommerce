output "aws_account_id" {
  description = "AWS account ID used by the current credentials."
  value       = data.aws_caller_identity.current.account_id
}

output "aws_region" {
  description = "AWS region in use by the provider."
  value       = data.aws_region.current.name
}

output "name_prefix" {
  description = "Common name prefix for resources."
  value       = local.name_prefix
}

output "cloudwatch_log_group_name" {
  description = "Name of the optional starter CloudWatch log group."
  value       = var.create_example_resources ? aws_cloudwatch_log_group.app[0].name : null
}
