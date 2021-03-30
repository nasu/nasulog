import * as cdk from '@aws-cdk/core'
import * as ecs from '@aws-cdk/aws-ecs'
import * as ec2 from '@aws-cdk/aws-ec2'
import * as ecr from '@aws-cdk/aws-ecr'
import { FargateService } from '@aws-cdk/aws-ecs'
import { ServicePrincipal } from '@aws-cdk/aws-iam'

interface ECSServiceStackProps extends cdk.StackProps {
  vpc: ec2.IVpc
  ecsCluster: ecs.ICluster
  repository: ecr.IRepository
}

export class NasulogECSServiceStack extends cdk.Stack {
  public readonly ecsLBTarget: ecs.IEcsLoadBalancerTarget

  constructor(scope: cdk.Construct, id: string, props: ECSServiceStackProps) {
    super(scope, id, props)

    const image = ecs.ContainerImage.fromEcrRepository(props.repository)
    const taskDefinition = new ecs.FargateTaskDefinition(this, id + 'TaskDefinition', {})
    taskDefinition.addContainer(id + 'TaskDefintionContainer', {
      image,
      environment: {
      }
    }).addPortMappings({
      containerPort: 8080,
      hostPort: 8080,
      protocol: ecs.Protocol.TCP,
    })
    const service = new ecs.FargateService(this, id + 'FargateService', {
      serviceName: id + 'FargateService',
      cluster: props.ecsCluster,
      taskDefinition: taskDefinition,
      vpcSubnets: props.vpc.selectSubnets({subnetType: ec2.SubnetType.PRIVATE}),
    })
    this.ecsLBTarget = service.loadBalancerTarget({
      containerName: taskDefinition.defaultContainer!.containerName,
      containerPort: taskDefinition.defaultContainer!.containerPort,
    })
  }
}
