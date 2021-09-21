import * as cdk from '@aws-cdk/core';
import { OrangeStackOpenApiServices } from '@orangestack-cdk/openapi-lambda';

export class {{NAME_PROJECT|to_camel}}Stack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    // Create OrangeStackOpenApiServices pointing to spec file
    // @ts-ignore
    new OrangeStackOpenApiServices(this, '{{NAME_PROJECT|to_camel}}Api', {
      specPath: 'spec/{{NAME_PROJECT|to_kebab}}.yaml'
    });
  }
}
