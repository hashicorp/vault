/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  _getPathParam,
  combineOpenApiAttrs,
  expandOpenApiProps,
  getHelpUrlForModel,
  pathToHelpUrlSegment,
} from 'vault/utils/openapi-helpers';
import Model, { attr } from '@ember-data/model';
import { setupTest } from 'ember-qunit';
import { camelize } from '@ember/string';

module('Unit | Utility | OpenAPI helper utils', function (hooks) {
  setupTest(hooks);

  test(`pathToHelpUrlSegment`, function (assert) {
    [
      { path: '/auth/{username}', result: '/auth/example' },
      { path: '{username}/foo', result: 'example/foo' },
      { path: 'foo/{username}/bar', result: 'foo/example/bar' },
      { path: '', result: '' },
      { path: undefined, result: '' },
    ].forEach((test) => {
      assert.strictEqual(pathToHelpUrlSegment(test.path), test.result, `translates ${test.path}`);
    });
  });

  test(`_getPathParam`, function (assert) {
    [
      { path: '/auth/{username}', result: 'username' },
      { path: '{unicorn}/foo', result: 'unicorn' },
      { path: 'foo/{bigfoot}/bar', result: 'bigfoot' },
      { path: '{alphabet}/bowl/{soup}', result: 'alphabet' },
      { path: 'no/params', result: false },
      { path: '', result: false },
      { path: undefined, result: false },
    ].forEach((test) => {
      assert.strictEqual(_getPathParam(test.path), test.result, `returns first param for ${test.path}`);
    });
  });

  test(`getHelpUrlForModel`, function (assert) {
    [
      { modelType: 'kmip/config', result: '/v1/foobar/config?help=1' },
      { modelType: 'does-not-exist', result: null },
      { modelType: 4, result: null },
      { modelType: '', result: null },
      { modelType: undefined, result: null },
    ].forEach((test) => {
      assert.strictEqual(
        getHelpUrlForModel(test.modelType, 'foobar'),
        test.result,
        `returns first param for ${test.path}`
      );
    });
  });

  test('combineOpenApiAttrs should combine attributes correctly', async function (assert) {
    class FooModel extends Model {
      @attr('string', {
        label: 'Foo',
        subText: 'A form field',
      })
      foo;
      @attr('boolean', {
        label: 'Bar',
        subText: 'Maybe a checkbox',
      })
      bar;
      @attr('number', {
        label: 'Baz',
        subText: 'A number field',
      })
      baz;
    }
    this.owner.register('model:foo', FooModel);
    const myModel = this.owner.lookup('service:store').modelFor('foo');
    const newProps = {
      foo: {
        editType: 'ttl',
      },
      baz: {
        type: 'number',
        editType: 'slider',
        label: 'Old label',
      },
      foobar: {
        type: 'string',
        label: 'Foo-bar',
      },
    };
    const expected = [
      {
        name: 'foo',
        type: 'string',
        options: {
          label: 'Foo',
          subText: 'A form field',
          editType: 'ttl',
        },
      },
      {
        name: 'bar',
        type: 'boolean',
        options: {
          label: 'Bar',
          subText: 'Maybe a checkbox',
        },
      },
      {
        name: 'baz',
        type: 'number',
        options: {
          label: 'Baz', // uses the value we set on the model
          editType: 'slider',
          subText: 'A number field',
        },
      },
      {
        name: 'foobar',
        type: 'string',
        options: {
          label: 'Foo-bar',
        },
      },
    ];
    const { attrs, newFields } = combineOpenApiAttrs(myModel.attributes, newProps);
    assert.deepEqual(newFields, ['foobar'], 'correct newFields added');

    // When combineOpenApiAttrs
    assert.strictEqual(attrs.length, 4, 'correct number of attributes returned');
    expected.forEach((exp) => {
      const name = exp.name;
      const attr = attrs.find((a) => a.name === name);
      assert.deepEqual(attr, exp, `${name} combined properly`);
    });
  });

  module('expandopenApiProps', function () {
    const OPENAPI_RESPONSE_PROPS = {
      ttl: {
        type: 'string',
        format: 'seconds',
        description: 'this is a TTL!',
        'x-vault-displayAttrs': {
          name: 'TTL',
        },
      },
      'awesome-people': {
        type: 'array',
        items: {
          type: 'string',
        },
        'x-vault-displayAttrs': {
          value: 'Grace Hopper,Lady Ada',
        },
      },
      'favorite-ice-cream': {
        type: 'string',
        enum: ['vanilla', 'chocolate', 'strawberry'],
      },
      'default-value': {
        default: 30,
        'x-vault-displayAttrs': {
          value: 300,
        },
        type: 'integer',
      },
      default: {
        'x-vault-displayAttrs': {
          value: 30,
        },
        type: 'integer',
      },
      'super-secret': {
        type: 'string',
        'x-vault-displayAttrs': {
          sensitive: true,
        },
        description: 'A really secret thing',
      },
    };
    const EXPANDED_PROPS = {
      ttl: {
        helpText: 'this is a TTL!',
        editType: 'ttl',
        label: 'TTL',
        fieldGroup: 'default',
      },
      awesomePeople: {
        editType: 'stringArray',
        defaultValue: 'Grace Hopper,Lady Ada',
        fieldGroup: 'default',
      },
      favoriteIceCream: {
        editType: 'string',
        type: 'string',
        possibleValues: ['vanilla', 'chocolate', 'strawberry'],
        fieldGroup: 'default',
      },
      defaultValue: {
        editType: 'number',
        type: 'number',
        defaultValue: 300,
        fieldGroup: 'default',
      },
      default: {
        editType: 'number',
        type: 'number',
        defaultValue: 30,
        fieldGroup: 'default',
      },
      superSecret: {
        type: 'string',
        editType: 'string',
        sensitive: true,
        helpText: 'A really secret thing',
        fieldGroup: 'default',
      },
    };
    const OPENAPI_DESCRIPTIONS = {
      token_bound_cidrs: {
        type: 'array',
        description:
          'Comma separated string or JSON list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
        items: {
          type: 'string',
        },
        'x-vault-displayAttrs': {
          description:
            'List of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
          name: "Generated Token's Bound CIDRs",
          group: 'Tokens',
        },
      },
      blah_blah: {
        type: 'array',
        description: 'Comma-separated list of policies',
        items: {
          type: 'string',
        },
        'x-vault-displayAttrs': {
          name: "Generated Token's Policies",
          group: 'Tokens',
        },
      },
      only_display_description: {
        type: 'array',
        items: {
          type: 'string',
        },
        'x-vault-displayAttrs': {
          description: 'Hello there, you look nice today',
        },
      },
    };

    const STRING_ARRAY_DESCRIPTIONS = {
      token_bound_cidrs: {
        helpText:
          'List of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
      },
      blah_blah: {
        helpText: 'Comma-separated list of policies',
      },
      only_display_description: {
        helpText: 'Hello there, you look nice today',
      },
    };
    test('it creates objects from OpenAPI schema props', function (assert) {
      assert.expect(6);
      const generatedProps = expandOpenApiProps(OPENAPI_RESPONSE_PROPS);
      for (const propName in EXPANDED_PROPS) {
        assert.deepEqual(EXPANDED_PROPS[propName], generatedProps[propName], `correctly expands ${propName}`);
      }
    });
    test('it uses the description from the display attrs block if it exists', async function (assert) {
      assert.expect(3);
      const generatedProps = expandOpenApiProps(OPENAPI_DESCRIPTIONS);
      for (const propName in STRING_ARRAY_DESCRIPTIONS) {
        assert.strictEqual(
          generatedProps[camelize(propName)].helpText,
          STRING_ARRAY_DESCRIPTIONS[propName].helpText,
          `correctly updates helpText for ${propName}`
        );
      }
    });
  });
});
