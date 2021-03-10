package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func getInvalidationInput(distributionId string) (input *cloudfront.CreateInvalidationInput) {
	invalidationPathItem := "/*"

	var invalidationPathItems []*string
	invalidationPathItems = append(invalidationPathItems, &invalidationPathItem)

	invalidationPaths := &cloudfront.Paths{}
	invalidationPaths.SetItems(invalidationPathItems)
	invalidationPaths.SetQuantity(1)

	invalidationBatch := &cloudfront.InvalidationBatch{}
	invalidationBatch.SetPaths(invalidationPaths)
	invalidationBatch.SetCallerReference(fmt.Sprintf("%s-%d", distributionId, time.Now().Unix()))

	input = &cloudfront.CreateInvalidationInput{}
	input.SetDistributionId(distributionId)
	input.SetInvalidationBatch(invalidationBatch)

	return
}

func getRegisterTaskDefinitionInput(taskDefinitionStr []byte) (input *ecs.RegisterTaskDefinitionInput, err error) {
	err = json.Unmarshal(taskDefinitionStr, &input)
	if err != nil {
		return
	}

	return
}

func getUpdateServiceInput(cluster string, service string, taskDefinition string) (input *ecs.UpdateServiceInput) {
	input = &ecs.UpdateServiceInput{
		Cluster:            aws.String(cluster),
		Service:            aws.String(service),
		TaskDefinition:     aws.String(taskDefinition),
		ForceNewDeployment: aws.Bool(true),
	}
	return
}

func invalidate(awsSession *session.Session, distributionId string) (err error) {
	cloudfrontInstance := cloudfront.New(awsSession)
	result, err := cloudfrontInstance.CreateInvalidation(getInvalidationInput(distributionId))
	if err != nil {
		return
	}

	err = cloudfrontInstance.WaitUntilInvalidationCompleted(&cloudfront.GetInvalidationInput{
		DistributionId: &distributionId,
		Id:             result.Invalidation.Id,
	})
	if err != nil {
		return
	}

	return
}

func createAWSSession(c *cli.Context) (s *session.Session, err error) {
	envPath := c.String("env")

	err = godotenv.Load(envPath)
	if err != nil {
		return
	}

	cfg := aws.Config{
		Credentials: credentials.NewSharedCredentials(
			os.Getenv("AWS_CREDENTIALS_FILE"),
			os.Getenv("AWS_PROFILE"),
		),
		Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	}

	return session.NewSession(&cfg)
}
