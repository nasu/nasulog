import * as cdk from '@aws-cdk/core'
import * as ecs from '@aws-cdk/aws-ecs'
import * as ec2 from '@aws-cdk/aws-ec2'

interface ECSClusterStackProps extends cdk.StackProps {
  vpc: ec2.IVpc
}

export class NasulogECSCluseterStack extends cdk.Stack {
  public readonly ecsCluster: ecs.ICluster

  constructor(scope: cdk.Construct, id: string, props: ECSClusterStackProps) {
    super(scope, id, props)

    this.ecsCluster = new ecs.Cluster(this, id, {
      clusterName: id,
      vpc: props.vpc,
      capacity: {
        instanceType: new ec2.InstanceType("t3.nano"),
        minCapacity: 1,
        maxCapacity: 1,
      },
    })
  }
}