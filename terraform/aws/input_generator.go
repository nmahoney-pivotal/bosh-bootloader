package aws

import (
	"crypto/sha1"
	"fmt"

	"github.com/cloudfoundry/bosh-bootloader/aws"
	"github.com/cloudfoundry/bosh-bootloader/storage"
)

type InputGenerator struct {
	availabilityZoneRetriever aws.AvailabilityZoneRetriever
}

const terraformNameCharLimit = 18

func NewInputGenerator(availabilityZoneRetriever aws.AvailabilityZoneRetriever) InputGenerator {
	return InputGenerator{
		availabilityZoneRetriever: availabilityZoneRetriever,
	}
}

func (i InputGenerator) Generate(state storage.State) (map[string]interface{}, error) {
	azs, err := i.availabilityZoneRetriever.RetrieveAvailabilityZones(state.AWS.Region)
	if err != nil {
		return map[string]interface{}{}, err
	}

	shortEnvID := state.EnvID
	if len(shortEnvID) > terraformNameCharLimit {
		sha1 := fmt.Sprintf("%x", sha1.Sum([]byte(state.EnvID)))
		shortEnvID = fmt.Sprintf("%s-%s", shortEnvID[:terraformNameCharLimit-8], sha1[:terraformNameCharLimit-11])
	}

	inputs := map[string]interface{}{
		"env_id":             state.EnvID,
		"short_env_id":       shortEnvID,
		"region":             state.AWS.Region,
		"availability_zones": azs,
	}

	if state.LB.Type == "cf" {
		inputs["ssl_certificate"] = state.LB.Cert
		inputs["ssl_certificate_private_key"] = state.LB.Key
		inputs["ssl_certificate_chain"] = state.LB.Chain

		if state.LB.Domain != "" {
			inputs["system_domain"] = state.LB.Domain
		}
	}

	return inputs, nil
}

func (i InputGenerator) Credentials(state storage.State) map[string]string {
	return map[string]string{
		"access_key": state.AWS.AccessKeyID,
		"secret_key": state.AWS.SecretAccessKey,
	}
}
