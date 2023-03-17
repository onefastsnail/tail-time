terraform {
  required_version = "~> 1.3"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">=4.56.0"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      creator = "terraform"
    }
  }
}
