/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  _getPathParam,
  combineOpenApiAttrs,
  getHelpUrlForModel,
  pathToHelpUrlSegment,
} from 'vault/utils/openapi-helpers';
import Model, { attr } from '@ember-data/model';
import { setupTest } from 'ember-qunit';

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
});
