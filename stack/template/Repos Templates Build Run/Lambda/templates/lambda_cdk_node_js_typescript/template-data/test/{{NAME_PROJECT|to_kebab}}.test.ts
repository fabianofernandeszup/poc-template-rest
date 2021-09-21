import { expect as expectCDK, matchTemplate, MatchStyle } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import * as {{NAME_PROJECT|to_camel}} from '../lib/{{NAME_PROJECT|to_kebab}}-stack';

test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new {{NAME_PROJECT|to_camel}}.{{NAME_PROJECT|to_camel}}Stack(app, '{{NAME_PROJECT}}');
    // THEN
    expectCDK(stack).to(matchTemplate({
      "Resources": {}
    }, MatchStyle.EXACT))
});
