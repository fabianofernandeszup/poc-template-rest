terraform {
  backend "s3" {
    bucket = "notification65465498746546"
    key    = "sqs/sqs-notification.tfstate"
    region = "us-east-1"
  }
  required_providers {
    aws = {
      version = "~> 3.2"
      source = "hashicorp/aws"
    }
  }
}