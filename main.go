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
	}

	return session.NewSession(&cfg)
}

func main() {
	green := color.New(color.FgGreen).SprintFunc()

	app := &cli.App{
		Usage:   "Helper functions for common aws actions",
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:        "invalidate",
				Usage:       "invalidate <name>",
				Description: "invalidate a CloudFront distribution's cache",
				Action: func(c *cli.Context) (err error) {
					awsSession, err := createAWSSession(c)
					if err != nil {
						return
					}

					name := c.Args().First()
					if name == "" {
						return errors.New("missing distribution name")
					}

					distributionId := os.Getenv("cloudfront_" + name)
					if distributionId == "" {
						return errors.New("distribution alias not found in configuration")
					}

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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: ".env",
				Usage: "env configuration file",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
