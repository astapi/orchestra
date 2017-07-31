package ecs

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"strings"
)

func (client *ECS) ListClusters() ([]string, error) {
	input := &ecs.ListClustersInput{}
	result, err := client.Client.ListClusters(input)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(result.ClusterArns))
	for i, v := range result.ClusterArns {
		ret[i] = strings.Split(*v, "/")[1]
	}
	return ret, nil
}
