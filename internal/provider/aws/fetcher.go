package aws

import (
	"context"
	"fmt"

	"github.com/user/driftwatch/internal/provider"
)

// defaultFetcher is the real AWS fetcher using the SDK (stubbed for now).
type defaultFetcher struct {
	region  string
	profile string
}

// FetchEC2Instances retrieves EC2 instance resources from AWS.
// In production this would use the AWS SDK; here it returns a descriptive error
// until the SDK dependency is wired in.
func (f *defaultFetcher) FetchEC2Instances(ctx context.Context, region string) ([]provider.Resource, error) {
	if region == "" {
		return nil, fmt.Errorf("aws region must not be empty")
	}
	// TODO: replace with real AWS SDK call:
	//   cfg, _ := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	//   client := ec2.NewFromConfig(cfg)
	//   ...
	return nil, fmt.Errorf("aws SDK not yet wired: set up EC2 client for region %s", region)
}
