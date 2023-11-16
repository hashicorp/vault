/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import Model, { attr } from '@ember-data/model';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

// create class using decorator
const createClass = () => {
  @withExpandedAttributes()
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

    get fieldGroups() {
      return [{ default: ['baz'] }, { 'Other options': ['foo', 'bar'] }];
    }
  }
  return new Foo();
};

module('Unit | Decorators | model-expanded-attributes', function (hooks) {
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
    @withExpandedAttributes()
    class Foo {} // eslint-disable-line
    const message =
      'withExpandedAttributes decorator must be used on instance of ember-data Model class. Decorator not applied to returned class';
    assert.ok(this.spy.calledWith(message), 'Error is printed to console');
  });

  test('it adds allByKey value to model', function (assert) {
    assert.expect(1);
    const model = createClass();
    assert.deepEqual(
      { foo: this.fooField, bar: this.barField, baz: this.bazField },
      model.allByKey,
      'allByKey set on Model class'
    );
  });

  test('_expandGroups helper works correctly', function (assert) {
    const model = createClass();
    const result = model._expandGroups(model.fieldGroups);
    assert.deepEqual(result, [
      { default: [this.bazField] },
      { 'Other options': [this.fooField, this.barField] },
    ]);
  });

  test('_expandGroups throws assertion when incorrect inputs', function (assert) {
    assert.expect(1);
    const model = createClass();
    try {
      model._expandGroups({ foo: ['bar'] });
    } catch (e) {
      assert.strictEqual(e.message, '_expandGroups expects an array of objects');
    }
  });
});
