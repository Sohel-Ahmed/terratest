package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-s3-example using Terratest.
func TestTerraformAwsS3Example(t *testing.T) {
	t.Parallel()

	var LocalEndpoints = map[string]string{
		"apigateway":     "http://localhost:4566",
		"cloudformation": "http://localhost:4566",
		"cloudwatch":     "http://localhost:4566",
		"dynamodb":       "http://localhost:4566",
		"es":             "http://localhost:4566",
		"firehose":       "http://localhost:4566",
		"iam":            "http://localhost:4566",
		"kinesis":        "http://localhost:4566",
		"lambda":         "http://localhost:4566",
		"route53":        "http://localhost:4566",
		"redshift":       "http://localhost:4566",
		"s3":             "http://localhost:4566",
		"secretsmanager": "http://localhost:4566",
		"ses":            "http://localhost:4566",
		"sns":            "http://localhost:4566",
		"sqs":            "http://localhost:4566",
		"ssm":            "http://localhost:4566",
		"stepfunctions":  "http://localhost:4566",
		"sts":            "http://localhost:4566",
	}
	aws.SetAwsEndpointsOverrides(LocalEndpoints)

	// Give this S3 Bucket a unique ID for a name tag so we can distinguish it from any other Buckets provisioned
	// in your AWS account
	expectedName := fmt.Sprintf("terratest-aws-s3-example-%s", strings.ToLower(random.UniqueId()))

	// Give this S3 Bucket an environment to operate as a part of for the purposes of resource tagging
	expectedEnvironment := "Automated Testing"

	// This usually picks a random AWS region to test in. This helps ensure your code works in all regions.
	// I used us-east-1 only since the default region for LocalStack is us-east-1
	awsRegion := aws.GetRandomStableRegion(t, []string{"us-east-1"}, nil)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/s3_bucket",
				// Variables to pass to our Terraform code using -var options
				Vars: map[string]interface{}{
					"tag_bucket_name":        expectedName,
					"tag_bucket_environment": expectedEnvironment,
					"with_policy":            "true",
				},
		
				// Environment variables to set when running Terraform
				EnvVars: map[string]string{
					"AWS_DEFAULT_REGION": awsRegion,
				},
			}
		
			// At the end of the test, run `terraform destroy` to clean up any resources that were created
			defer terraform.Destroy(t, terraformOptions)
		
			// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
			terraform.InitAndApply(t, terraformOptions)
		
			// Run `terraform output` to get the value of an output variable
			bucketID := terraform.Output(t, terraformOptions, "bucket_id")
		
			// Verify that our Bucket has versioning enabled
			actualStatus := aws.GetS3BucketVersioning(t, awsRegion, bucketID)
			expectedStatus := "Enabled"
			assert.Equal(t, expectedStatus, actualStatus)
		
			// Verify that our Bucket has a policy attached
			aws.AssertS3BucketPolicyExists(t, awsRegion, bucketID)
		}