variable "project_name" {
  description = "Project name used for naming and tagging."
  type        = string
  default     = "ecommerce-platform"
}

variable "environment" {
  description = "Deployment environment name."
  type        = string
  default     = "dev"
}

variable "aws_region" {
  description = "AWS region where resources are managed."
  type        = string
  default     = "us-east-1"
}

variable "tags" {
  description = "Additional tags applied to all managed resources."
  type        = map(string)
  default     = {}
}

variable "create_example_resources" {
  description = "Create optional starter resources when true."
  type        = bool
  default     = false
}
