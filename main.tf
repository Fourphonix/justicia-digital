terraform {
  required_version = ">= 1.8.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region                      = "us-east-1"
  skip_credentials_validation = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
}

# Bucket de almacenamiento WORM inalterable
resource "aws_s3_bucket" "expedientes_worm" {
  bucket        = "justicia-digital-expedientes-worm-store"
  force_destroy = false

  object_lock_enabled = true
}

resource "aws_s3_bucket_object_lock_configuration" "worm_lock" {
  bucket = aws_s3_bucket.expedientes_worm.id

  rule {
    default_retention {
      mode = "COMPLIANCE"
      days = 3650
    }
  }
}
