import * as cdk from '@aws-cdk/core'
import * as elb from '@aws-cdk/aws-elasticloadbalancingv2'
import * as ec2 from '@aws-cdk/aws-ec2'
import * as ecs from '@aws-cdk/aws-ecs'

interface ELBStackProps extends cdk.StackProps {
  vpc: ec2.IVpc
  lbTarget: ecs.IEcsLoadBalancerTarget
}

export class NasulogELBStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: ELBStackProps) {
    super(scope, id, props)

    const alb = new elb.ApplicationLoadBalancer(this, id, {
      loadBalancerName: id,
      internetFacing: true,
      vpc: props.vpc,
      vpcSubnets: {subnets: props.vpc.publicSubnets},
    })
    const albListner = alb.addListener(id + 'Listner', {
      port: 80,
    })
    albListner.addTargets(id + 'ListnerTarget', {
      protocol: elb.ApplicationProtocol.HTTP,
      port: 80,
      targets: [ props.lbTarget ],
    })
  }
}
