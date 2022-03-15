using Pulumi;
using Pulumi.Ssp;

class MyStack : Stack
{
  public MyStack()
  {
    var ssp = new SharedServicesPlatform("csharp", new SharedServicesPlatformArgs
    {
      ClusterArgs = new Pulumi.Ssp.Inputs.ClusterArgsArgs
      {
        KubernetesVersion = "1.21.0",
        Region = "us-east-1",
      }
    });

    ssp.OnboardTeam(new SharedServicesPlatformOnboardTeamArgs
    {
      Name = "payments",
      Controller = "pulumi",
      Repository = "github.com/awesome-org/payments-team/infra"
    });
  }
}
