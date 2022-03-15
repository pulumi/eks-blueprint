package provider

import (
	"fmt"
	"time"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-eks/sdk/go/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (s *Ssp) getManagedNodeGroups() eks.ManagedNodeGroupArray {
	array := eks.ManagedNodeGroupArray{}

	for _, nodeGroup := range s.ManagedNodeGroups {
		array = append(array, nodeGroup)
	}

	return array
}

type AddManagedNodeGroupArgs struct {
	Name              string   `pulumi:"name"`
	InstanceTypes     []string `pulumi:"instanceTypes"`
	MinSize           int      `pulumi:"minSize"`
	MaxSize           int      `pulumi:"maxSize"`
	DesiredSize       int      `pulumi:"desiredSize"`
	KubernetesVersion string   `pulumi:"kubernetesVersion"`
}

type AddManagedNodeGroupResult struct {
	Result *eks.ManagedNodeGroup `pulumi:"result"`
}

var managedPolicyArns = []string{
	"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
	"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
	"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
}

func (s *Ssp) createRole(ctx *pulumi.Context, name string) (*iam.Role, error) {
	version := "2012-10-17"
	statementId := "AllowAssumeRole"
	effect := "Allow"
	instanceAssumeRolePolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Version: &version,
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Sid:    &statementId,
				Effect: &effect,
				Actions: []string{
					"sts:AssumeRole",
				},
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type: "Service",
						Identifiers: []string{
							"ec2.amazonaws.com",
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	role, err := iam.NewRole(ctx, name, &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(instanceAssumeRolePolicy.Json),
	}, pulumi.Parent(s))
	if err != nil {
		return nil, err
	}

	for i, policy := range managedPolicyArns {
		_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-policy-%d", name, i), &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String(policy),
			Role:      role.ID(),
		}, pulumi.Parent(s))
		if err != nil {
			return nil, err
		}
	}
	return role, nil
}

func (s *Ssp) addManagedNodeGroup(ctx *pulumi.Context, args *AddManagedNodeGroupArgs) (*AddManagedNodeGroupResult, error) {
	var opts []pulumi.ResourceOption

	opts = append(opts, pulumi.DependsOn([]pulumi.Resource{s.Cluster}))
	opts = append(opts, pulumi.Parent(s))

	managedNodeGroup, err := eks.NewManagedNodeGroup(ctx,
		args.Name,
		&eks.ManagedNodeGroupArgs{
			Cluster:       s.Cluster.Core,
			NodeGroupName: pulumi.String(args.Name),
			InstanceTypes: pulumi.ToStringArray(args.InstanceTypes),
			NodeRole:      s.InstanceRole,
			// ScalingConfig: &awsEks.NodeGroupScalingConfigArgs{
			// 	MinSize:     pulumi.Int(args.MinSize),
			// 	MaxSize:     pulumi.Int(args.MaxSize),
			// 	DesiredSize: pulumi.Int(args.DesiredSize),
			// },
			// ScalingConfig: &awsEks.NodeGroupScalingConfigArgs{
			// 	MinSize: pulumi.Int(args.MinSize),
			// },
		}, opts...)

	if err != nil {
		return nil, err
	}

	s.ManagedNodeGroups = append(s.ManagedNodeGroups, managedNodeGroup)

	memoCache := getMemoCache()
	err = memoCache.Storage.Replace("ssp", s, time.Hour*1)
	if err != nil {
		return nil, err
	}

	return &AddManagedNodeGroupResult{
		Result: managedNodeGroup,
	}, nil
}
