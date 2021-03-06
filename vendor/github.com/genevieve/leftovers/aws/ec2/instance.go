package ec2

import (
	"fmt"
	"strings"

	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
)

type Instance struct {
	client     instancesClient
	id         *string
	identifier string
	rtype      string
}

func NewInstance(client instancesClient, id, keyName *string, tags []*awsec2.Tag) Instance {
	identifier := *id

	extra := []string{}
	for _, t := range tags {
		extra = append(extra, fmt.Sprintf("%s:%s", *t.Key, *t.Value))
	}

	if keyName != nil && *keyName != "" {
		extra = append(extra, fmt.Sprintf("KeyPairName:%s", *keyName))
	}

	if len(extra) > 0 {
		identifier = fmt.Sprintf("%s (%s)", *id, strings.Join(extra, ", "))
	}

	return Instance{
		client:     client,
		id:         id,
		identifier: identifier,
		rtype:      "EC2 Instance",
	}
}

func (i Instance) Delete() error {
	_, err := i.client.TerminateInstances(&awsec2.TerminateInstancesInput{
		InstanceIds: []*string{i.id},
	})

	if err != nil {
		return fmt.Errorf("FAILED deleting %s %s: %s", i.rtype, i.identifier, err)
	}

	return nil
}

func (i Instance) Name() string {
	return i.identifier
}

func (i Instance) Type() string {
	return i.rtype
}
