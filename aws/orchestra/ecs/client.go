package ecs

import (
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ECS struct {
	Client *ecs.ECS
}
