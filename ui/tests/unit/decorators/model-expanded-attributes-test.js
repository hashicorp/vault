/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import sinon from 'sinon';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

module('Unit | Decorators | model-expanded-attributes', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.spy = sinon.spy(console, 'error');
  });
  hooks.afterEach(function () {
    this.spy.restore();
  });

  test('it should warn when applying decorator to class that does not extend Model', function (assert) {
    @withExpandedAttributes()
    class Foo {} // eslint-disable-line
    const message =
      'withExpandedAttributes decorator must be used on instance of ember-data Model class. Decorator not applied to returned class';
    assert.ok(this.spy.calledWith(message), 'Error is printed to console');
  });

  test('it adds allByKey value to model', function (assert) {
    assert.expect(1);
    const model = this.store.modelFor('namespace');
    assert.deepEqual(
      model.prototype.allByKey,
      {
        path: {
          name: 'path',
          options: {},
          type: 'string',
        },
      },
      'allByKey set on Model class'
    );
  });

  test('_expandGroups helper works correctly', function (assert) {
    const model = this.store.modelFor('aws-credential');
    const result = model.prototype._expandGroups([
      { default: ['roleArn'] },
      { 'Other options': ['ttl', 'leaseId'] },
    ]);
    assert.deepEqual(result, [
      {
        default: [
          {
            name: 'roleArn',
            options: {
              helpText:
                'The ARN of the role to assume if credential_type on the Vault role is assumed_role. Optional if the role has a single role ARN; required otherwise.',
              label: 'Role ARN',
            },
            type: 'string',
          },
        ],
      },
      {
        'Other options': [
          {
            name: 'ttl',
            options: {
              defaultValue: '3600s',
              editType: 'ttl',
              helpText:
                'Specifies the TTL for the use of the STS token. Valid only when credential_type is assumed_role, federation_token, or session_token.',
              label: 'TTL',
              setDefault: true,
              ttlOffValue: '',
            },
            type: undefined,
          },
          {
            name: 'leaseId',
            options: {},
            type: 'string',
          },
        ],
      },
    ]);
  });

  test('_expandGroups throws assertion when incorrect inputs', function (assert) {
    assert.expect(1);
    const model = this.store.modelFor('aws-credential');
    try {
      model.prototype._expandGroups({ foo: ['bar'] });
    } catch (e) {
      assert.strictEqual(e.message, '_expandGroups expects an array of objects');
    }
  });
});
