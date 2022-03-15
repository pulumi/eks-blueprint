package provider

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi-ssp/provider/pkg/provider/addons"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AddClusterAddOnArgs struct {
	Name    string `pulumi:"name"`
	Version string `pulumi:"version"`
}

type AddClusterAddOnResult struct {
	Result *helm.ReleaseOutput `pulumi:"result"`
}

func (s *Ssp) clusterAddon(ctx *pulumi.Context, args *AddClusterAddOnArgs) (*AddClusterAddOnResult, error) {
	switch args.Name {
	case "argocd":
		h, err := addons.AddArgoCD(ctx, &addons.AddArgoCDArgs{
			KubernetesProvider: &s.Cluster.Provider,
			Namespace:          &s.PlatformNamespace,
			Version:            args.Version,
		})
		if err != nil {
			return nil, err
		}

		return &AddClusterAddOnResult{
			Result: &h.Result,
		}, nil

	case "metrics-server":
		h, err := addons.AddMetricsServer(ctx, &addons.AddMetricsServerArgs{
			KubernetesProvider: &s.Cluster.Provider,
			Namespace:          &s.PlatformNamespace,
			Version:            args.Version,
		})
		if err != nil {
			return nil, err
		}

		return &AddClusterAddOnResult{
			Result: &h.Result,
		}, nil

	default:
		return nil, errors.Errorf("unknown addon %s", args.Name)
	}
}
