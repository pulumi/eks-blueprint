import * as ssp from "@pulumi/ssp";

// import { FargateProfile } from "@pulumi/aws/eks";

// const fp = new FargateProfile("my-fargate-profile", {
//   clusterName: "",
//   podExecutionRoleArn: "",
//   selectors: [{ namespace: "serverless" }],
// });

const sharedServicesPlatform = new ssp.SharedServicesPlatform("my-platform", {
  clusterArgs: {
    kubernetesVersion: "1.21.0",
    region: "us-east-1",
  },
});

sharedServicesPlatform.addManagedNodeGroup({
  name: "my-managed-group",
  desiredSize: 1,
  minSize: 1,
  kubernetesVersion: "1.21.0",
  maxSize: 2,
  instanceTypes: ["t3.medium"],
});

sharedServicesPlatform.clusterAddon({
  name: "argocd",
});

sharedServicesPlatform.onboardTeam({
  name: "team-a",
  repository: "github.com/pulumi/pulumi",
  controller: "argocd",
});

// sharedServicesPlatform.clusterAddon({
//   name: "metrics-server",
// });

// sharedServicesPlatform.addFargateNodeGroup({
//   profiles: [
//     {
//       selectors: [{ namespace: "dynatrace" }],
//     },
//   ],
// });
