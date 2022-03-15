package provider

import (
	"fmt"
	"time"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-ssp/provider/pkg/provider/addons"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type OnboardTeamArgs struct {
	Name       string `pulumi:"name"`
	Repository string `pulumi:"repository"`
	Controller string `pulumi:"controller"`
}

type OnboardTeamResult struct {
	Result pulumi.Output `pulumi:"result"`
}

func (s *Ssp) onboardTeam(ctx *pulumi.Context, args *OnboardTeamArgs) (*OnboardTeamResult, error) {
	var opts []pulumi.ResourceOption

	opts = append(opts, pulumi.DependsOn([]pulumi.Resource{s.Cluster}))
	opts = append(opts, pulumi.Parent(s))

	team := pulumi.All(s.Cluster.Provider).ApplyT(func(iArgs []interface{}) (string, error) {
		kubernetesProvider := iArgs[0].(*kubernetes.Provider)

		opts = append(opts, pulumi.Provider(kubernetesProvider))

		namespace, err := corev1.NewNamespace(ctx, args.Name, &corev1.NamespaceArgs{}, opts...)
		if err != nil {
			return "", err
		}

		switch args.Controller {
		case "pulumi":
			addons.AddPulumiOperator(ctx, &addons.AddPulumiOperatorArgs{
				KubernetesProvider: kubernetesProvider,
				Namespace:          namespace,
				Version:            "1.5.0",
			})

		case "argocd":
			addons.AddArgoCD(ctx, &addons.AddArgoCDArgs{
				KubernetesProvider: kubernetesProvider,
				Namespace:          namespace,
				Version:            "1.5.0",
			})

		default:
			return "", fmt.Errorf("Couldn't find a team controller %s", args.Controller)
		}

		return "done", nil
	})

	memoCache := getMemoCache()
	err := memoCache.Storage.Replace("ssp", s, time.Hour*1)
	if err != nil {
		return nil, err
	}

	return &OnboardTeamResult{
		Result: team,
	}, nil
}
