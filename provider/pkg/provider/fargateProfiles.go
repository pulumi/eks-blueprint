package provider

import (
	"time"

	awseks "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AddFargateProfileArgs struct {
	Name      string                             `pulumi:"name"`
	Selectors awseks.FargateProfileSelectorInput `pulumi:"selectors"`
}

type AddFargateProfileResult struct {
	Result *awseks.FargateProfile `pulumi:"result"`
}

func (s *Ssp) addFargateProfile(ctx *pulumi.Context, args *AddFargateProfileArgs) (*AddFargateProfileResult, error) {
	var opts []pulumi.ResourceOption

	opts = append(opts, pulumi.DependsOn([]pulumi.Resource{s.Cluster}))
	opts = append(opts, pulumi.Parent(s))

	fargateProfile, err := awseks.NewFargateProfile(ctx, args.Name, &awseks.FargateProfileArgs{
		ClusterName: pulumi.String(s.ClusterName),
		Selectors: &awseks.FargateProfileSelectorArray{
			args.Selectors,
		},
	})

	// eks.FargateProfile(ctx, args.Name, eks.FargateProfileArgs{})
	if err != nil {
		return nil, err
	}

	// s.FargateProfiles = append(s.FargateProfiles, fargateProfile)

	memoCache := getMemoCache()
	err = memoCache.Storage.Replace("ssp", s, time.Hour*1)
	if err != nil {
		return nil, err
	}

	return &AddFargateProfileResult{
		Result: fargateProfile,
	}, nil
}
