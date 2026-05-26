/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { withFormFields } from 'vault/decorators/model-form-fields';
import Model, { attr } from '@ember-data/model';
import { run } from '@ember/runloop';
import sinon from 'sinon';

module('Unit | Decorators | ModelFormFields', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy(console, 'error');
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

  test('it sets formFields and allFields when applied to a Model subclass', function (assert) {
    @withFormFields(['name', 'role'])
    class TestModel extends Model {
      @attr('string', { label: 'Name' }) name;
      @attr('string', { label: 'Role' }) role;
    }
    this.owner.register('model:test-with-form-fields', TestModel);

    const model = run(() => this.owner.lookup('service:store').createRecord('test-with-form-fields'));

    assert.ok(Array.isArray(model.formFields), 'formFields is set as an array');
    assert.strictEqual(model.formFields.length, 2, 'formFields contains the specified fields');
    assert.strictEqual(model.formFields[0].name, 'name', 'first formField is name');
    assert.strictEqual(model.formFields[1].name, 'role', 'second formField is role');

    assert.ok(Array.isArray(model.allFields), 'allFields is set as an array');
    assert.strictEqual(model.allFields.length, 2, 'allFields contains all model attributes');
  });

  test('it sets formFieldGroups when groupPropertyNames are provided', function (assert) {
    @withFormFields(['name'], [{ default: ['name'] }, { Options: ['role'] }])
    class GroupedModel extends Model {
      @attr('string') name;
      @attr('string') role;
    }
    this.owner.register('model:test-with-form-field-groups', GroupedModel);

    const model = run(() => this.owner.lookup('service:store').createRecord('test-with-form-field-groups'));

    assert.ok(Array.isArray(model.formFieldGroups), 'formFieldGroups is set as an array');
    assert.ok(
      model.formFieldGroups.some((g) => Object.prototype.hasOwnProperty.call(g, 'default')),
      'formFieldGroups contains a default group'
    );
    assert.ok(
      model.formFieldGroups.some((g) => Object.prototype.hasOwnProperty.call(g, 'Options')),
      'formFieldGroups contains the Options group'
    );
  });
});
