import * as cdk from '@aws-cdk/core'
import * as dynamodb from "@aws-cdk/aws-dynamodb"
import * as ecr from "@aws-cdk/aws-ecr"
import * as ecs from "@aws-cdk/aws-ecs"
import * as ec2 from "@aws-cdk/aws-ec2"
import * as elb from "@aws-cdk/aws-elasticloadbalancingv2"
import { ApplicationProtocol } from '@aws-cdk/aws-elasticloadbalancingv2'
import { Protocol, RESERVED_TUNNEL_INSIDE_CIDR } from '@aws-cdk/aws-ec2'

export class NasulogInitStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)

    new dynamodb.Table(this, 'nasulogDynamoDB', {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      tableName: 'nasulog',
      partitionKey: {
        name: 'partition_key',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'sort_key',
        type: dynamodb.AttributeType.STRING,
      },
    })

    const repository = new ecr.Repository(this, 'nasulogRepository', {
      repositoryName: 'nasulog',
    })

    const vpc = new ec2.Vpc(this, 'nasulogVpc', {
      cidr: '10.0.0.0/16',
      maxAzs: 3,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      subnetConfiguration: [
        { cidrMask: 24, name: 'Public',   subnetType: ec2.SubnetType.PUBLIC, },
        { cidrMask: 24, name: 'Private',  subnetType: ec2.SubnetType.PRIVATE, },
        { cidrMask: 24, name: 'Isolated', subnetType: ec2.SubnetType.ISOLATED, },
      ]
    })
    vpc.addGatewayEndpoint("nasulogDynamoDBVpcEndpoint", {
      service: ec2.GatewayVpcEndpointAwsService.DYNAMODB,
      subnets: [{subnetType:ec2.SubnetType.PRIVATE}],
    })

    const cluster = new ecs.Cluster(this, 'nasulogCluster', {
      clusterName: 'nasulogCluster',
      vpc: vpc,
      capacity: {
        instanceType: new ec2.InstanceType("t3.nano"),
        minCapacity: 1,
        maxCapacity: 1,
      },
    })
    const image = ecs.ContainerImage.fromEcrRepository(repository)
    const taskDefinition = new ecs.FargateTaskDefinition(this, 'nasulogTaskDefinition', {})
    taskDefinition.addContainer('taskDefintionContainerDefinition', {
      image,
      environment: {
        "DYNAMODB_URL": "http://dynamodb:9000",
        "SQS_URL": "http://elasticmq:9324",
        "AWS_ACCESS_KEY_ID": "AKIA0000000000000000",
        "AWS_SECRET_ACCESS_KEY": "s7C0000000000000000000000000000000000000",
      }
    }).addPortMappings({
      containerPort: 8080,
      hostPort: 8080,
      protocol: ecs.Protocol.TCP,
    })
    const fargate = new ecs.FargateService(this, 'nasulogFargateService2', {
      serviceName: 'nasulogFargateService2',
      cluster,
      taskDefinition: taskDefinition,
      vpcSubnets: vpc.selectSubnets({subnetType: ec2.SubnetType.PRIVATE}),
    })

    const alb = new elb.ApplicationLoadBalancer(this, 'nasulogALB', {
      loadBalancerName: 'nasulogALB',
      internetFacing: true,
      vpc: vpc,
      vpcSubnets: {subnets: vpc.publicSubnets},
    })
    const albListner = alb.addListener('nasulogALBListner', {
      port: 80,
    })
    albListner.addTargets('nasulogALBListnerTarget', {
      protocol: ApplicationProtocol.HTTP,
      port: 80,
      targets: [
        fargate.loadBalancerTarget({
          containerName: taskDefinition.defaultContainer!.containerName,
          containerPort: taskDefinition.defaultContainer!.containerPort,
        }),
      ],
    })
  }
}