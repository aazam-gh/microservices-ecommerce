data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

resource "aws_cloudwatch_log_group" "app" {
  count             = var.create_example_resources ? 1 : 0
  name              = "/aws/${local.name_prefix}/application"
  retention_in_days = 14
}
