import * as cdk from '@aws-cdk/core'
import * as ec2 from '@aws-cdk/aws-ec2'

export class NasulogNetworkStack extends cdk.Stack {
  public readonly vpc: ec2.IVpc

  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)

    this.vpc = new ec2.Vpc(this, id, {
      cidr: '10.0.0.0/16',
      maxAzs: 3,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      subnetConfiguration: [
        { cidrMask: 24, name: 'Public',   subnetType: ec2.SubnetType.PUBLIC, },
        { cidrMask: 24, name: 'Private',  subnetType: ec2.SubnetType.PRIVATE, },
        //{ cidrMask: 24, name: 'Isolated', subnetType: ec2.SubnetType.ISOLATED, },
      ]
    })
    this.vpc.addGatewayEndpoint('DynamoDBVPCEndpoint', {
      service: ec2.GatewayVpcEndpointAwsService.DYNAMODB,
      subnets: [{subnetType:ec2.SubnetType.PRIVATE}],
    })
    this.vpc.addInterfaceEndpoint('ECRDockerEndpoint', {
      service: ec2.InterfaceVpcEndpointAwsService.ECR_DOCKER,
      subnets: {subnetType:ec2.SubnetType.PRIVATE},
    })
  }
}