package addons

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AddMetricsServerArgs struct {
	KubernetesProvider *kubernetes.ProviderOutput `pulumi:"kubernetesProvider"`
	Namespace          *corev1.NamespaceOutput    `pulumi:"namespace"`
	Version            string                     `pulumi:"version"`
}

type AddMetricsServerResult struct {
	Result helm.ReleaseOutput `pulumi:"result"`
}

func AddMetricsServer(ctx *pulumi.Context, args *AddMetricsServerArgs) (*AddMetricsServerResult, error) {
	var opts []pulumi.ResourceOption

	helmRelease := pulumi.All(args.KubernetesProvider, args.Namespace).ApplyT(func(iArgs []interface{}) (*helm.Release, error) {
		kubernetesProvider := iArgs[0].(*kubernetes.Provider)
		namespace := iArgs[1].(*corev1.Namespace)

		opts = append(opts, pulumi.DependsOn([]pulumi.Resource{kubernetesProvider}))
		opts = append(opts, pulumi.Parent(kubernetesProvider))
		opts = append(opts, pulumi.Provider(kubernetesProvider))

		if args.Version == "" {
			args.Version = "3.7.0"
		}

		helmRelease, err := helm.NewRelease(ctx, "ssp-addon-metrics-server", &helm.ReleaseArgs{
			Chart:     pulumi.String("metrics-server"),
			Version:   pulumi.String(args.Version),
			Namespace: namespace.Metadata.Name(),
			Name:      pulumi.String("ssp-addon-metrics-servers"),
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://kubernetes-sigs.github.io/metrics-server"),
			},
		}, opts...)
		if err != nil {
			return nil, err
		}

		return helmRelease, nil
	})

	return &AddMetricsServerResult{
		Result: helmRelease.(helm.ReleaseOutput),
	}, nil
}
