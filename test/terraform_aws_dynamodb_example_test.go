package test

import (
	"fmt"
	"testing"
	"time"

	awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-dynamodb-example using Terratest.
func TestTerraformAwsDynamoDBExample(t *testing.T) {
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
		"kms":            "http://localhost:4566",
	}
	aws.SetAwsEndpointsOverrides(LocalEndpoints)

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	// awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	// awsRegion := "us-east-1"
	// // Set up expected values to be checked later
	expectedTableName := fmt.Sprintf("terratest-aws-dynamodb-example-table-test") // random.UniqueId())
	// expectedKmsKeyArn := aws.GetCmkArn(t, awsRegion, "alias/aws/dynamodb")
	expectedKeySchema := []*dynamodb.KeySchemaElement{
		{AttributeName: awsSDK.String("userId"), KeyType: awsSDK.String("HASH")},
		{AttributeName: awsSDK.String("department"), KeyType: awsSDK.String("RANGE")},
	}
	expectedTags := []*dynamodb.Tag{
		{Key: awsSDK.String("Environment"), Value: awsSDK.String("production")},
	}

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-dynamodb-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"table_name": expectedTableName,
			// "region":     awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	time.Sleep(5 * time.Second)
	// Look up the DynamoDB table by name
	// table := aws.GetDynamoDBTable(t, awsRegion, expectedTableName)

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sessionOptions := session.Options{
		Config: awsSDK.Config{
			Region:           awsSDK.String("us-east-1"),
			Endpoint:         awsSDK.String("http://localhost:4566"),
			S3ForcePathStyle: awsSDK.Bool(true),
		}}
	sess := session.Must(session.NewSessionWithOptions(sessionOptions))
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n, "- this is the table name")
			// assert.Equal(t, "terratest-aws-dynamodb-example-table-test", *n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	// create the input configuration instance
	input_new := &dynamodb.DescribeTableInput{
		TableName: awsSDK.String("terratest-aws-dynamodb-example-table-test"),
	}

	// fmt.Printf("Describe Table:\n")

	// Get the list of tables
	table_desc, err := svc.DescribeTable(input_new)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	// fmt.Println(table_desc, "This is the table description")

	// for _, n := range table_desc.Table.AttributeDefinitions {
	// 	//fmt.Println(*n, "- this is the table Attributes")
	// 	// assert.Equal(t, "userId", *n)
	// }
	assert.Equal(t, "ACTIVE", awsSDK.StringValue(table_desc.Table.TableStatus))
	assert.ElementsMatch(t, expectedKeySchema, table_desc.Table.KeySchema)

	// // Verify server-side encryption configuration
	// assert.Equal(t, expectedKmsKeyArn, awsSDK.StringValue(table_desc.Table.SSEDescription.KMSMasterKeyArn))
	assert.Equal(t, "ENABLED", awsSDK.StringValue(table_desc.Table.SSEDescription.Status))
	assert.Equal(t, "KMS", awsSDK.StringValue(table_desc.Table.SSEDescription.SSEType))

	// Verify TTL configuration
	ttl_input := &dynamodb.DescribeTimeToLiveInput{
		TableName: awsSDK.String("terratest-aws-dynamodb-example-table-test"),
	}

	table_ttl, err := svc.DescribeTimeToLive(ttl_input)

	fmt.Println("TTL of table is ", table_ttl)

	// ttl := aws.GetDynamoDBTableTimeToLive(t, awsRegion, expectedTableName)
	assert.Equal(t, "expires", awsSDK.StringValue(table_ttl.TimeToLiveDescription.AttributeName))
	assert.Equal(t, "ENABLED", awsSDK.StringValue(table_ttl.TimeToLiveDescription.TimeToLiveStatus))

	// // Verify resource tags
	tags_input := &dynamodb.ListTagsOfResourceInput{
		ResourceArn: table_desc.Table.TableArn,
	}
	table_tags, err := svc.ListTagsOfResource(tags_input)
	fmt.Println("Table tags are", table_tags)

	// tags := aws.GetDynamoDbTableTags(t, awsRegion, expectedTableName)
	assert.ElementsMatch(t, expectedTags, table_tags.Tags)
}
