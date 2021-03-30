import * as cdk from '@aws-cdk/core'
import * as dynamodb from '@aws-cdk/aws-dynamodb'
import * as iam from '@aws-cdk/aws-iam'

export class NasulogDynamoDBStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)

    const table = new dynamodb.Table(this, id, {
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
  }
}