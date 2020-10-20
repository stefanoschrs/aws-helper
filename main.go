package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/go-sharp/color"
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

func main() {
	green := color.New(color.FgGreen).SprintFunc()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg := aws.Config{
		Credentials: credentials.NewSharedCredentials(
			os.Getenv("AWS_CREDENTIALS_FILE"),
			os.Getenv("AWS_PROFILE"),
		),
	}
	awsSession, err := session.NewSession(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "invalidate",
				Usage:       "invalidate <name>",
				Description: "invalidate a CloudFront distribution's cache",
				Action: func(c *cli.Context) (err error) {
					name := c.Args().First()
					if name == "" {
						return errors.New("missing distribution name")
					}

					distributionId := os.Getenv("cloudfront_" + name)
					fmt.Printf("Invalidating %s (%s).. ", green(name), distributionId)

					err = invalidate(awsSession, distributionId)
					if err != nil {
						return
					}

					fmt.Println(green("COMPLETE"))
					return
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
