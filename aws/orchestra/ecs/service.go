package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
)

func (client *ECS) ListServices(clusterName string) ([]string, error) {
	input := &ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
	}
	result, err := client.Client.ListServices(input)
	if err != nil {
		return nil, err
	}
	arns := result.ServiceArns
	ret := make([]string, len(arns))
	for i, v := range arns {
		ret[i] = strings.Split(*v, "/")[1]
	}
	return ret, nil
}

func (client *ECS) UpdateService(clusterName string, serviceName string) error {
	taskName, err := client.CurrentTaskDefinition(clusterName, serviceName)
	if err != nil {
		return err
	}

	taskDefinition, err := client.UpdateTask(taskName)
	if err != nil {
		return err
	}

	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(clusterName),
		Service:        aws.String(serviceName),
		TaskDefinition: aws.String(*taskDefinition.TaskDefinitionArn),
	}
	_, err = client.Client.UpdateService(input)
	if err != nil {
		return err
	}

	//fmt.Println(result)
	return nil
}

func (client *ECS) CurrentTaskDefinition(clusterName string, serviceName string) (string, error) {
	input := &ecs.DescribeServicesInput{
		Cluster: aws.String(clusterName),
		Services: []*string{
			aws.String(serviceName),
		},
	}
	result, err := client.Client.DescribeServices(input)
	if err != nil {
		return "", err
	}
	taskRevision := strings.Split(*result.Services[0].TaskDefinition, "/")[1]
	return strings.Split(taskRevision, ":")[0], nil
}

func (client *ECS) DescribeService(serviceName string) {
}
