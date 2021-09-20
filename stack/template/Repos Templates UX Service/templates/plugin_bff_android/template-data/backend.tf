terraform {
  backend "s3" {
    bucket = "{{BUCKET_TFSTATE_NAME}}"
    key    = "sqs/sqs-{{NAME|unify}}.tfstate"
    region = "{{BUCKET_TFSTATE_REGION}}"
  }
  required_providers {
    aws = {
      version = "~> 3.2"
      source = "hashicorp/aws"
    }
  }
}