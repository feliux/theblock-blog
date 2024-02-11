terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.36.0"
    }
  }
}

provider "aws" {
  region                   = var.region
  shared_credentials_files = [var.credentials]
  profile                  = var.profile // develop
}
