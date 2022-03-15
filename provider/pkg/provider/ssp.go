// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-eks/sdk/go/eks"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Ssp struct {
	pulumi.ResourceState

	Cluster           *eks.Cluster            `pulumi:"cluster"`
	ClusterName       string                  `json:"clusterName"`
	PlatformNamespace corev1.NamespaceOutput  `pulumi:"platformNamespace"`
	InstanceRole      *iam.Role               `pulumi:"instanceRole"`
	ManagedNodeGroups []*eks.ManagedNodeGroup `pulumi:"managedNodeGroups"`
	NodeGroups        []*eks.NodeGroup        `pulumi:"nodeGroups"`
	FargateProfiles   []*eks.FargateProfile   `pulumi:"fargateProfiles"`
}

type ClusterArgs struct {
	Region            pulumi.StringInput    `pulumi:"region"`
	KubernetesVersion pulumi.StringInput    `pulumi:"kubernetesVersion"`
	EncryptionConfig  pulumi.StringInput    `pulumi:"encryptionConfig"`
	Tags              pulumi.StringMapInput `pulumi:"tags"`
}

func NewSharedServicesPlatform(ctx *pulumi.Context,
	name string, args *ClusterArgs, opts ...pulumi.ResourceOption) (*Ssp, error) {
	if args == nil {
		args = &ClusterArgs{}
	}

	component := &Ssp{
		InstanceRole:      &iam.Role{},
		ManagedNodeGroups: []*eks.ManagedNodeGroup{},
		NodeGroups:        []*eks.NodeGroup{},
		FargateProfiles:   []*eks.FargateProfile{},
	}

	err := ctx.RegisterComponentResource("ssp:index:SharedServicesPlatform", name, component, opts...)
	if err != nil {
		return nil, err
	}

	instanceRole, err := component.createRole(ctx, name)
	if err != nil {
		return nil, err
	}

	component.ClusterName = name

	cluster, err := eks.NewCluster(ctx, name, &eks.ClusterArgs{
		Name:                 pulumi.String(name),
		SkipDefaultNodeGroup: pulumi.Bool(true),
		Tags:                 args.Tags,
		Version:              args.KubernetesVersion,
		UseDefaultVpcCni:     pulumi.Bool(true),
		InstanceRole:         instanceRole,
	}, append(opts, pulumi.Parent(component))...)
	if err != nil {
		fmt.Printf("Error during NewCluster: %w\n", err)
		return nil, err
	}

	platformNamespace := cluster.Provider.ApplyT(func(kubernetesProvider *kubernetes.Provider) (*corev1.Namespace, error) {
		opts = append(opts, pulumi.DependsOn([]pulumi.Resource{cluster}))
		opts = append(opts, pulumi.Parent(component))
		opts = append(opts, pulumi.Provider(kubernetesProvider))

		platformNamespace, err := corev1.NewNamespace(ctx, "platform", &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("platform"),
			},
		}, opts...)

		if err != nil {
			return nil, err
		}

		return platformNamespace, nil
	}).(corev1.NamespaceOutput)

	component.Cluster = cluster
	component.InstanceRole = instanceRole
	component.PlatformNamespace = platformNamespace

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"kubeconfig": cluster.Kubeconfig,
	}); err != nil {
		return nil, err
	}

	return component, nil
}
