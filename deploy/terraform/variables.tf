variable "aws_region" {
  description = "AWS region to deploy in"
  type        = string
  default     = "ap-south-1"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "key_name" {
  description = "Name of the SSH key pair for EC2 access"
  type        = string
}

variable "domain_name" {
  description = "Domain name for the portfolio (optional)"
  type        = string
  default     = ""
}

variable "project_name" {
  description = "Name tag for resources"
  type        = string
  default     = "ssh-portfolio"
}
