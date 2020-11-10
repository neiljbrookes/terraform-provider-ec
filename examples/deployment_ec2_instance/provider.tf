terraform {
  required_version = ">= 0.12"

  required_providers {
    ec = {
      source = "elastic/ec"
    }

    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "ec" {
  apikey = var.ec_api_key
}

provider "aws" {
  region  = var.region
  profile = var.aws_profile
}