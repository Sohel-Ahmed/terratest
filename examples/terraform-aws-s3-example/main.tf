# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------
# provider "aws" {
#   region = var.region
# }


provider "aws" {
  region                      = "us-east-1"
  s3_force_path_style         = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway     = "http://localhost:4566"
    cloudformation = "http://localhost:4566"
    cloudwatch     = "http://localhost:4566"
    dynamodb       = "http://localhost:4566"
    es             = "http://localhost:4566"
    firehose       = "http://localhost:4566"
    iam            = "http://localhost:4566"
    kinesis        = "http://localhost:4566"
    lambda         = "http://localhost:4566"
    route53        = "http://localhost:4566"
    redshift       = "http://localhost:4566"
    s3             = "http://localhost:4566"
    secretsmanager = "http://localhost:4566"
    ses            = "http://localhost:4566"
    sns            = "http://localhost:4566"
    sqs            = "http://localhost:4566"
    ssm            = "http://localhost:4566"
    stepfunctions  = "http://localhost:4566"
    sts            = "http://localhost:4566"
    ec2            = "http://localhost:4566"
  }
}

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A S3 BUCKET WITH VERSIONING ENABLED INCLUDING TAGS
# See test/terraform_aws_s3_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

# data "aws_iam_policy_document" "s3_bucket_policy" {
#   statement {
#     effect = "Allow"
#     principals {
#       # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
#       # force an interpolation expression to be interpreted as a list by wrapping it
#       # in an extra set of list brackets. That form was supported for compatibility in
#       # v0.11, but is no longer supported in Terraform v0.12.
#       #
#       # If the expression in the following list itself returns a list, remove the
#       # brackets to avoid interpretation as a list of lists. If the expression
#       # returns a single list item then leave it as-is and remove this TODO comment.
#       identifiers = [local.aws_account_id]
#       type        = "AWS"
#     }
#     actions   = ["*"]
#     resources = ["${aws_s3_bucket.test_bucket.arn}/*"]
#   }

#   statement {
#     effect = "Deny"
#     principals {
#       identifiers = ["*"]
#       type        = "AWS"
#     }
#     actions   = ["*"]
#     resources = ["${aws_s3_bucket.test_bucket.arn}/*"]

#     condition {
#       test     = "Bool"
#       variable = "aws:SecureTransport"
#       values = [
#         "false",
#       ]
#     }
#   }
# }

resource "aws_s3_bucket" "test_bucket_logs" {
  bucket = "${local.aws_account_id}-${var.tag_bucket_name}-logs"
  acl    = "log-delivery-write"

  tags = {
    Name        = "${local.aws_account_id}-${var.tag_bucket_name}-logs"
    Environment = var.tag_bucket_environment
  }

  force_destroy = true
}

# resource "aws_s3_bucket" "test_bucket" {
#   bucket = "${local.aws_account_id}-${var.tag_bucket_name}"
#   acl    = "private"

#   versioning {
#     enabled = true
#   }

#   logging {
#     target_bucket = aws_s3_bucket.test_bucket_logs.id
#     target_prefix = "TFStateLogs/"
#   }

#   tags = {
#     Name        = var.tag_bucket_name
#     Environment = var.tag_bucket_environment
#   }
# }

# resource "aws_s3_bucket_policy" "bucket_access_policy" {
#   count  = var.with_policy ? 1 : 0
#   bucket = aws_s3_bucket.test_bucket.id
#   policy = data.aws_iam_policy_document.s3_bucket_policy.json
# }

# ---------------------------------------------------------------------------------------------------------------------
# LOCALS
# Used to represent any data that requires complex expressions/interpolations
# ---------------------------------------------------------------------------------------------------------------------

data "aws_caller_identity" "current" {
}

locals {
  aws_account_id = data.aws_caller_identity.current.account_id
}

