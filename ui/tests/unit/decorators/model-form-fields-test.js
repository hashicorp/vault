/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { withFormFields } from 'vault/decorators/model-form-fields';
import sinon from 'sinon';
import Model, { attr } from '@ember-data/model';

// create class using decorator
const createClass = (propertyNames, groups) => {
  @withFormFields(propertyNames, groups)
  class Foo extends Model {
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
  return new Foo();
};

module('Unit | Decorators | ModelFormFields', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy(console, 'error');
    this.fooField = {
      name: 'foo',
      options: { label: 'Foo', subText: 'A form field' },
      type: 'string',
    };
    this.barField = {
      name: 'bar',
      options: { label: 'Bar', subText: 'Maybe a checkbox' },
      type: 'boolean',
    };
    this.bazField = {
      name: 'baz',
      options: { label: 'Baz', subText: 'A number field' },
      type: 'number',
    };
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
    const model = createClass();
    assert.deepEqual(
      [this.fooField, this.barField, this.bazField],
      model.allFields,
      'allFields set on Model class'
    );
  });

  test('it should set formFields prop on Model class', function (assert) {
    const model = createClass(['foo']);
    assert.deepEqual([this.fooField], model.formFields, 'formFields set on Model class');
  });

  test('it should set formFieldGroups on Model class', function (assert) {
    const groups = [{ default: ['foo'] }, { subgroup: ['bar'] }];
    const model = createClass(null, groups);
    const fieldGroups = [{ default: [this.fooField] }, { subgroup: [this.barField] }];
    assert.deepEqual(fieldGroups, model.formFieldGroups, 'formFieldGroups set on Model class');
  });
});
