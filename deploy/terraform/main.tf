terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Latest Ubuntu 24.04 AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# Security group: allow SSH (portfolio) + admin SSH
resource "aws_security_group" "portfolio" {
  name_prefix = "${var.project_name}-"
  description = "SSH portfolio server"

  # Portfolio SSH (port 22)
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH portfolio access"
  }

  # Admin SSH (port 2222) — restrict to your IP in production
  ingress {
    from_port   = 2222
    to_port     = 2222
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Admin SSH access"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = var.project_name
  }
}

resource "aws_instance" "portfolio" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [aws_security_group.portfolio.id]

  root_block_device {
    volume_size = 8
    volume_type = "gp3"
  }

  user_data = <<-EOF
    #!/bin/bash
    set -e

    # Move default SSH to port 2222 so portfolio gets port 22
    sed -i 's/^#Port 22/Port 2222/' /etc/ssh/sshd_config
    sed -i 's/^Port 22/Port 2222/' /etc/ssh/sshd_config
    systemctl restart sshd

    # Install Docker
    apt-get update -y
    apt-get install -y docker.io
    systemctl enable docker
    systemctl start docker

    # Create app directory
    mkdir -p /opt/ssh-portfolio/data
    mkdir -p /opt/ssh-portfolio/.ssh

    # Pull and run (update with your Docker Hub/GHCR image)
    # docker pull ghcr.io/mayurathavale18/ssh-portfolio:latest
    # docker run -d \
    #   --name ssh-portfolio \
    #   --restart always \
    #   -p 22:22 \
    #   -v /opt/ssh-portfolio/.ssh:/app/.ssh \
    #   -v /opt/ssh-portfolio/data:/app/data \
    #   ghcr.io/mayurathavale18/ssh-portfolio:latest
  EOF

  tags = {
    Name = var.project_name
  }
}

# Elastic IP for stable address
resource "aws_eip" "portfolio" {
  instance = aws_instance.portfolio.id
  domain   = "vpc"

  tags = {
    Name = var.project_name
  }
}
