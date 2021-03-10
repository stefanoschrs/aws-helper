package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/go-sharp/color"
	"github.com/urfave/cli/v2"
)

func ActionInvalidate(c *cli.Context) (err error) {
	green := color.New(color.FgGreen).SprintFunc()

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

	// Invalidate
	fmt.Printf("Invalidating %s (%s).. ", green(name), distributionId)

	err = invalidate(awsSession, distributionId)
	if err != nil {
		return
	}

	fmt.Println(green("COMPLETE"))
	return
}

func ActionECSDeploy(c *cli.Context) (err error) {
	green := color.New(color.FgGreen).SprintFunc()

	awsSession, err := createAWSSession(c)
	if err != nil {
		return
	}

	name := c.Args().First()
	if name == "" {
		return errors.New("missing service name")
	}

	cluster := os.Getenv("ecs_" + name + "_cluster")
	if cluster == "" {
		return errors.New("cluster alias not found in configuration")
	}
	service := os.Getenv("ecs_" + name + "_service")
	if service == "" {
		return errors.New("service alias not found in configuration")
	}
	taskDefinitionPath := os.Getenv("ecs_" + name + "_task-definition")
	if cluster == "" {
		return errors.New("taskDefinitionPath alias not found in configuration")
	}

	// Load task definition json
	taskDefinition, err := ioutil.ReadFile(taskDefinitionPath)
	if err != nil {
		err = fmt.Errorf("load task definition: %w", err)
		return
	}

	ecsInstance := ecs.New(awsSession)

	// RegisterTaskDefinition
	rtdInput, err := getRegisterTaskDefinitionInput(taskDefinition)
	if err != nil {
		return
	}

	fmt.Printf("Registering new task definition for %s.. ", green(*rtdInput.Family))

	rtdOutput, err := ecsInstance.RegisterTaskDefinition(rtdInput)
	if err != nil {
		return
	}

	fmt.Println(green("COMPLETE"))

	// UpdateService
	taskRevision := fmt.Sprintf("%s:%d", *rtdInput.Family, *rtdOutput.TaskDefinition.Revision)
	fmt.Printf("Updating service %s.. (%s) ", green(service), taskRevision)

	usOutput, err := ecsInstance.UpdateService(getUpdateServiceInput(cluster, service, taskRevision))
	if err != nil {
		return
	}

	fmt.Println(green("COMPLETE"))

	for _, set := range usOutput.Service.TaskSets {
		fmt.Printf("Task: %s\tDefinition: %s\tStatus: %s\n", *set.Id, *set.TaskDefinition, *set.Status)
	}
	return
}
