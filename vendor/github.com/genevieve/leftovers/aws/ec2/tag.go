package ec2

import (
	"fmt"

	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
)

type Tag struct {
	client     tagsClient
	key        *string
	value      *string
	resourceId *string
	identifier string
	rtype      string
}

func NewTag(client tagsClient, key, value, resourceId *string) Tag {
	return Tag{
		client:     client,
		key:        key,
		value:      value,
		resourceId: resourceId,
		identifier: *value,
		rtype:      "EC2 Tag",
	}
}

func (t Tag) Delete() error {
	_, err := t.client.DeleteTags(&awsec2.DeleteTagsInput{
		Tags:      []*awsec2.Tag{{Key: t.key, Value: t.value}},
		Resources: []*string{t.resourceId},
	})

	if err != nil {
		return fmt.Errorf("FAILED deleting %s %s: %s", t.rtype, t.identifier, err)
	}

	return nil
}

func (t Tag) Name() string {
	return t.identifier
}

func (t Tag) Type() string {
	return t.rtype
}
