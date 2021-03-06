package iam

import (
	"fmt"
	"strings"

	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/genevieve/leftovers/aws/common"
)

type instanceProfilesClient interface {
	ListInstanceProfiles(*awsiam.ListInstanceProfilesInput) (*awsiam.ListInstanceProfilesOutput, error)
	RemoveRoleFromInstanceProfile(*awsiam.RemoveRoleFromInstanceProfileInput) (*awsiam.RemoveRoleFromInstanceProfileOutput, error)
	DeleteInstanceProfile(*awsiam.DeleteInstanceProfileInput) (*awsiam.DeleteInstanceProfileOutput, error)
}

type InstanceProfiles struct {
	client instanceProfilesClient
	logger logger
}

func NewInstanceProfiles(client instanceProfilesClient, logger logger) InstanceProfiles {
	return InstanceProfiles{
		client: client,
		logger: logger,
	}
}

func (i InstanceProfiles) ListAll(filter string) ([]common.Deletable, error) {
	return i.get(filter)
}

func (i InstanceProfiles) List(filter string) ([]common.Deletable, error) {
	resources, err := i.get(filter)
	if err != nil {
		return nil, err
	}

	var delete []common.Deletable
	for _, r := range resources {
		proceed := i.logger.PromptWithDetails(r.Type(), r.Name())
		if !proceed {
			continue
		}

		delete = append(delete, r)
	}

	return delete, nil
}

func (i InstanceProfiles) get(filter string) ([]common.Deletable, error) {
	profiles, err := i.client.ListInstanceProfiles(&awsiam.ListInstanceProfilesInput{})
	if err != nil {
		return nil, fmt.Errorf("Listing instance profiles: %s", err)
	}

	var resources []common.Deletable
	for _, p := range profiles.InstanceProfiles {
		resource := NewInstanceProfile(i.client, p.InstanceProfileName, p.Roles, i.logger)

		if !strings.Contains(resource.Name(), filter) {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
