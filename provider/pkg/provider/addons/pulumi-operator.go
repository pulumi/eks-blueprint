package addons

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AddPulumiOperatorArgs struct {
	KubernetesProvider *kubernetes.Provider `pulumi:"kubernetesProvider"`
	Namespace          *corev1.Namespace    `pulumi:"namespace"`
	Version            string               `pulumi:"version"`
}

type AddPulumiOperatorResult struct {
	Result appsv1.Deployment `pulumi:"result"`
}

func AddPulumiOperator(ctx *pulumi.Context, args *AddPulumiOperatorArgs) (*AddPulumiOperatorResult, error) {
	var resourceOptions []pulumi.ResourceOption

	resourceOptions = append(resourceOptions, pulumi.Provider(args.KubernetesProvider))

	crds, err := yaml.NewConfigFile(ctx, "crds", &yaml.ConfigFileArgs{
		File: fmt.Sprintf("https://raw.githubusercontent.com/pulumi/pulumi-kubernetes-operator/v%s/deploy/crds/pulumi.com_stacks.yaml", args.Version),
	}, resourceOptions...)

	if err != nil {
		return nil, err
	}

	operatorServiceAccount, err := corev1.NewServiceAccount(ctx, "operator-service-account", &corev1.ServiceAccountArgs{}, resourceOptions...)
	if err != nil {
		return nil, err
	}

	operatorRole, err := rbacv1.NewRole(ctx, "operator-role", &rbacv1.RoleArgs{
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("pods"),
					pulumi.String("services"),
					pulumi.String("services/finalizers"),
					pulumi.String("endpoints"),
					pulumi.String("persistentvolumeclaims"),
					pulumi.String("events"),
					pulumi.String("configmaps"),
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apps"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("deployments"),
					pulumi.String("daemonsets"),
					pulumi.String("replicasets"),
					pulumi.String("statefulsets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("monitoring.coreos.com"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("servicemonitors"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("get"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apps"),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("pulumi-kubernetes-operator"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("deployments/finalizers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("pods"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apps"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("replicasets"),
					pulumi.String("deployments"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("pulumi.com"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("*"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("coordination.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("leases"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("update"),
				},
			},
		},
	}, resourceOptions...)
	if err != nil {
		return nil, err
	}

	_, err = rbacv1.NewRoleBinding(ctx, "operator-role-binding", &rbacv1.RoleBindingArgs{
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind: pulumi.String("ServiceAccount"),
				Name: operatorServiceAccount.Metadata.Name().Elem(),
			},
		},
		RoleRef: &rbacv1.RoleRefArgs{
			Kind:     pulumi.String("Role"),
			Name:     operatorRole.Metadata.Name().Elem(),
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
		},
	}, resourceOptions...)
	if err != nil {
		return nil, err
	}

	resourceOptions = append(resourceOptions, pulumi.DependsOn([]pulumi.Resource{crds}))

	deployment, err := appsv1.NewDeployment(ctx, "pulumi-kubernetes-operator", &appsv1.DeploymentArgs{
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"name": pulumi.String("pulumi-kubernetes-operator"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"name": pulumi.String("pulumi-kubernetes-operator"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: operatorServiceAccount.Metadata.Name(),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("pulumi-kubernetes-operator"),
							Image: pulumi.String("pulumi/pulumi-kubernetes-operator:v1.5.0"),
							Command: pulumi.StringArray{
								pulumi.String("pulumi-kubernetes-operator"),
							},
							Args: pulumi.StringArray{
								pulumi.String("--zap-level=debug"),
								pulumi.String("--zap-time-encoding=iso8601"),
							},
							ImagePullPolicy: pulumi.String("Always"),
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("WATCH_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAME"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.name"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("OPERATOR_NAME"),
									Value: pulumi.String("pulumi-kubernetes-operator"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("GRACEFUL_SHUTDOWN_TIMEOUT_DURATION"),
									Value: pulumi.String("5m"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("MAX_CONCURRENT_RECONCILES"),
									Value: pulumi.String("10"),
								},
							},
						},
					},
					// Should be same or larger than GRACEFUL_SHUTDOWN_TIMEOUT_DURATION
					TerminationGracePeriodSeconds: pulumi.Int(300),
				},
			},
		},
	}, resourceOptions...)
	if err != nil {
		return nil, err
	}

	return &AddPulumiOperatorResult{
		Result: *deployment,
	}, nil
}
