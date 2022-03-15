"""A Python Pulumi program"""
import pulumi
import pulumi_ssp
from pulumi_ssp import SharedServicesPlatform, ClusterArgsArgs
from pulumi_aws.eks import FargateProfileSelectorArgs

ssp = SharedServicesPlatform("python-platform", cluster_args=ClusterArgsArgs(kubernetes_version="1.21.0", region="us-east-1"))
ssp.add_managed_node_group(name="my-python-ng", desired_size=1, min_size=1, max_size=1, instance_types=["t3.medium"])
# ssp.add_fargate_profile(name="fargate", selectors=[FargateProfileSelectorArgs(namespace="fargate", labels={"namespace": "fargate"})])
ssp.onboard_team(name="team-python", repository="github.com/pulumi/pulumi", controller="pulumi")
