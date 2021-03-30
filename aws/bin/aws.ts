#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { NasulogECSCluseterStack } from '../lib/ecscluster-stack';
import { NasulogECRStack } from '../lib/ecr-stack';
import { NasulogNetworkStack } from '../lib/network-stack';
import { NasulogDynamoDBStack } from '../lib/dynamodb-stack';
import { NasulogECSServiceStack } from '../lib/ecsservice-stack';
import { NasulogELBStack } from '../lib/elb-stack';

const app = new cdk.App();

new NasulogDynamoDBStack(app, NasulogDynamoDBStack.name)
const ecr = new NasulogECRStack(app, NasulogECRStack.name)
const network = new NasulogNetworkStack(app, NasulogNetworkStack.name)
const ecscluster = new NasulogECSCluseterStack(app, NasulogECSCluseterStack.name, { vpc: network.vpc })
const ecsservice = new NasulogECSServiceStack(app, NasulogECSServiceStack.name, {
  vpc: network.vpc,
  repository: ecr.repository,
  ecsCluster: ecscluster.ecsCluster,
})
const elb = new NasulogELBStack(app, NasulogELBStack.name, {
  vpc: network.vpc,
  lbTarget: ecsservice.ecsLBTarget,
})