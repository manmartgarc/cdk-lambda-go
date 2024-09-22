import { GoFunction } from '@aws-cdk/aws-lambda-go-alpha';
import * as cdk from 'aws-cdk-lib';
import { AttributeType, Table } from 'aws-cdk-lib/aws-dynamodb';
import { Construct } from 'constructs';

export class LambdaGoStack extends cdk.Stack {
  private readonly goFunc: GoFunction;
  private readonly ddbTable: Table;
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.ddbTable = new Table(this, 'HelloGoTable', {
      partitionKey: {name: 'Name', type: AttributeType.STRING}
    });

    this.goFunc = new GoFunction(this, 'HelloGo', {
      entry: 'go-lambda',
      environment: {
        'TABLE_NAME': this.ddbTable.tableName
      }
    });

    this.ddbTable.grantReadWriteData(this.goFunc);
  }
}
