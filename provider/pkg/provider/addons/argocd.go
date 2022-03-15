package addons

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AddArgoCDArgs struct {
	KubernetesProvider *kubernetes.Provider `pulumi:"kubernetesProvider"`
	Namespace          *corev1.Namespace    `pulumi:"namespace"`
	Version            string               `pulumi:"version"`
}

type AddArgoCDResult struct {
	Result helm.ReleaseOutput `pulumi:"result"`
}

func AddArgoCD(ctx *pulumi.Context, args *AddArgoCDArgs) (*AddArgoCDResult, error) {
	var opts []pulumi.ResourceOption

	helmRelease := pulumi.All(args.KubernetesProvider, args.Namespace).ApplyT(func(iArgs []interface{}) (*helm.Release, error) {
		kubernetesProvider := iArgs[0].(*kubernetes.Provider)
		namespace := iArgs[1].(*corev1.Namespace)

		opts = append(opts, pulumi.DependsOn([]pulumi.Resource{kubernetesProvider, namespace}))
		opts = append(opts, pulumi.Parent(kubernetesProvider))
		opts = append(opts, pulumi.Provider(kubernetesProvider))

		if args.Version == "" {
			args.Version = "3.33.5"
		}

		helmRelease, err := helm.NewRelease(ctx, "ssp-addon-argocd", &helm.ReleaseArgs{
			Chart:     pulumi.String("argo-cd"),
			Version:   pulumi.String(args.Version),
			Namespace: namespace.Metadata.Name(),
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://argoproj.github.io/argo-helm"),
			},
		}, opts...)
		if err != nil {
			return &helm.Release{}, err
		}

		return helmRelease, nil
	})

	return &AddArgoCDResult{
		Result: helmRelease.(helm.ReleaseOutput),
	}, nil
}
