/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { withFormFields } from 'vault/decorators/model-form-fields';
import sinon from 'sinon';

module('Unit | Decorators | ModelFormFields', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy(console, 'error');
    this.store = this.owner.lookup('service:store');
  });
  hooks.afterEach(function () {
    this.spy.restore();
  });

  test('it should warn when applying decorator to class that does not extend Model', function (assert) {
    @withFormFields()
    class Foo {} // eslint-disable-line
    const message =
      'withFormFields decorator must be used on instance of ember-data Model class. Decorator not applied to returned class';
    assert.ok(this.spy.calledWith(message), 'Error is printed to console');
  });

  test('it return allFields when arguments not provided', function (assert) {
    assert.expect(1);
    // test by instantiating a record that uses this decorator
    const record = this.store.createRecord('kv/data');
    assert.deepEqual(
      record.allFields,
      [
        {
          name: 'backend',
          options: {},
          type: 'string',
        },
        {
          name: 'path',
          options: {
            label: 'Path for this secret',
            subText: 'Names with forward slashes define hierarchical path structures.',
          },
          type: 'string',
        },
        {
          name: 'secretData',
          options: {},
          type: 'object',
        },
        {
          name: 'createdTime',
          options: {},
          type: 'string',
        },
        {
          name: 'customMetadata',
          options: {},
          type: 'object',
        },
        {
          name: 'deletionTime',
          options: {},
          type: 'string',
        },
        {
          name: 'destroyed',
          options: {},
          type: 'boolean',
        },
        {
          name: 'version',
          options: {},
          type: 'number',
        },
        {
          name: 'failReadErrorCode',
          options: {},
          type: 'number',
        },
        {
          name: 'casVersion',
          options: {},
          type: 'number',
        },
      ],
      'allFields set on Model class'
    );
  });

  test('it should set formFields prop on Model class', function (assert) {
    // this model uses withFormFields
    const record = this.store.createRecord('clients/config');
    assert.deepEqual(
      record.formFields,
      [
        {
          name: 'enabled',
          options: {},
          type: 'string',
        },
        {
          name: 'retentionMonths',
          options: {
            label: 'Retention period',
            subText: 'The number of months of activity logs to maintain for client tracking.',
          },
          type: 'number',
        },
      ],
      'formFields set on Model class'
    );
  });

  test('it should set formFieldGroups on Model class', function (assert) {
    // this model uses withFormFields with groups
    const record = this.store.createRecord('ldap/config');
    const groups = record.formFieldGroups.map((group) => Object.keys(group)[0]);
    assert.deepEqual(
      groups,
      ['default', 'TLS options', 'More options'],
      'formFieldGroups set on Model class with correct group labels'
    );
  });
});
