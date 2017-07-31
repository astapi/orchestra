package ecs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func New(sess *session.Session) *ECS {
	return &ECS{Client: ecs.New(sess)}
}

func (client *ECS) ListTaskDefinitionFamilies() ([]string, error) {
	input := &ecs.ListTaskDefinitionFamiliesInput{}
	result, err := client.Client.ListTaskDefinitionFamilies(input)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(result.Families))
	for i, v := range result.Families {
		ret[i] = *v
	}
	return ret, nil
}

func (client *ECS) RunTask(clusterName string, taskName string) error {
	input := &ecs.RunTaskInput{
		Cluster:        aws.String(clusterName),
		TaskDefinition: aws.String(taskName),
	}
	result, err := client.Client.RunTask(input)
	fmt.Println(result)
	return err
}

func (client *ECS) UpdateTask(taskName string) (*ecs.TaskDefinition, error) {
	taskDefinition, err := client.DescriveTaskDefinition(taskName)
	if err != nil {
		return nil, err
	}

	containerDefinitions := taskDefinition.TaskDefinition.ContainerDefinitions
	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: containerDefinitions,
		Family:               aws.String(taskName),
	}
	result, err := client.Client.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return result.TaskDefinition, nil
}

func (client *ECS) DescriveTaskDefinition(taskName string) (*ecs.DescribeTaskDefinitionOutput, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskName),
	}
	return client.Client.DescribeTaskDefinition(input)
}
