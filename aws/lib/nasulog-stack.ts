import * as cdk from '@aws-cdk/core'
import * as dynamodb from "@aws-cdk/aws-dynamodb"
import * as ecr from "@aws-cdk/aws-ecr"

export class NasulogInitStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)

    new dynamodb.Table(this, 'nasulogDynamoDB', {
      tableName: 'nasulog',
      partitionKey: {
        name: 'parition_key',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'sort_key',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    })

    new ecr.Repository(this, 'nasulogRepository', {
      repositoryName: 'nasulog',
    })
  }
}