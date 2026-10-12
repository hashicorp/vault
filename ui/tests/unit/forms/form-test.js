/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { module, test } from 'qunit';

module('Unit | forms | form', function (hooks) {
  hooks.beforeEach(function () {
    this.data = {
      foo: 'bar',
    };
    this.options = { isNew: true };
    this.validations = {
      foo: [{ type: 'presence', message: 'Foo is required' }],
    };
  });

  test('it should set data, options and validations on class', async function (assert) {
    const form = new Form(this.data, this.options, this.validations);
    assert.deepEqual(form.data, this.data, 'Data is set correctly');
    assert.true(form.isNew, 'isNew is set correctly');
    assert.deepEqual(form.validations, this.validations, 'Validations are set correctly');
  });

  test('it should shim Ember object set method', async function (assert) {
    const form = new Form(this.data);
    form.set('data.foo', 'baz');
    assert.strictEqual(form.data.foo, 'baz', 'Property is updated using set method');
  });

  test('it should return data as valid without validations from toJSON method', async function (assert) {
    const expected = { isValid: true, state: {}, invalidFormMessage: '', data: this.data };
    const form = new Form(this.data);
    const json = form.toJSON();
    assert.deepEqual(json, expected, 'toJSON returns correct data and state');
  });

  test('it should validate and return data from toJSON method', async function (assert) {
    const state = { foo: { isValid: true, errors: [], warnings: [] } };
    const expected = { isValid: true, state, invalidFormMessage: '', data: this.data };

    const form = new Form(this.data, this.options, this.validations);
    const json = form.toJSON();

    assert.deepEqual(json, expected, 'toJSON returns correct data and validation state');
  });

  test('it should allow for custom data to be passed in toJSON method', async function (assert) {
    this.data.foo = 'string with whitespace';
    this.validations.foo = [{ type: 'containsWhiteSpace', message: 'Whitespace is not allowed' }];

    const serializedData = { foo: this.data.foo.replace(/\s/g, '_') };
    const state = { foo: { isValid: true, errors: [], warnings: [] } };
    const expected = { isValid: true, state, invalidFormMessage: '', data: serializedData };

    const form = new Form(this.data, this.options, this.validations);
    const json = form.toJSON(serializedData);

    assert.deepEqual(json, expected, 'toJSON allows data to be overridden');
  });

  test('it should proxy get/set to data object for form fields', async function (assert) {
    class TestForm extends Form {
      formFields = [new FormField('foo', 'string')];
    }
    const form = new TestForm(this.data);

    assert.strictEqual(
      form.foo,
      'bar',
      'data object value returned when accessing form field property on class'
    );
    form.foo = 'baz';
    assert.strictEqual(
      form.data.foo,
      'baz',
      'data object property is set when setting form field property on class'
    );
    assert.strictEqual(form.foo, 'baz', 'updated data object value is returned');
  });

  module('max duration', function (hooks) {
    hooks.beforeEach(function () {
      this.tooLarge = 'Rotation period is too large. Maximum allowed value is 106,751 days.';
      this.TtlForm = class extends Form {
        formFields = [
          new FormField('name', 'string'),
          new FormField('rotation_period', 'string', { editType: 'ttl', label: 'Rotation period' }),
        ];
      };
    });

    test('it should invalidate a TTL field longer than Vault supports', async function (assert) {
      const form = new this.TtlForm({ name: 'role', rotation_period: '863999999913600s' });
      const { isValid, state, invalidFormMessage } = form.toJSON();

      assert.false(isValid, 'form is invalid');
      assert.deepEqual(state.rotation_period.errors, [this.tooLarge], 'error is reported on the TTL field');
      assert.strictEqual(invalidFormMessage, 'There is an error with this form.', 'error count is reported');
    });

    test('it should accept empty, zero and maximum TTL values', async function (assert) {
      ['', '0s', '9223372036s', '106751d'].forEach((rotation_period) => {
        const { isValid } = new this.TtlForm({ rotation_period }).toJSON();
        assert.true(isValid, `"${rotation_period}" is valid`);
      });
    });

    test('it should merge with validations defined on the form', async function (assert) {
      const form = new this.TtlForm(
        { rotation_period: '' },
        {},
        {
          rotation_period: [{ type: 'presence', message: 'Rotation period is required' }],
        }
      );
      assert.deepEqual(
        form.toJSON().state.rotation_period.errors,
        ['Rotation period is required'],
        'existing rule still runs'
      );

      form.data.rotation_period = '106752d';
      assert.deepEqual(
        form.toJSON().state.rotation_period.errors,
        [this.tooLarge],
        'max rule runs alongside it'
      );
    });

    test('it should not change the validations defined on the form', async function (assert) {
      const validations = { name: [{ type: 'presence', message: 'Name is required' }] };
      const form = new this.TtlForm({ name: 'role', rotation_period: '106752d' }, {}, validations);

      form.toJSON();
      const { state } = form.toJSON();

      assert.deepEqual(
        form.validations,
        { name: [{ type: 'presence', message: 'Name is required' }] },
        'form validations are unchanged'
      );
      assert.deepEqual(
        state.rotation_period.errors,
        [this.tooLarge],
        'error is not duplicated on repeat calls'
      );
    });

    test('it should ignore fields that are not TTL fields', async function (assert) {
      const { isValid, state } = new this.TtlForm({ name: '106752d' }).toJSON();

      assert.true(isValid, 'form is valid');
      assert.strictEqual(state.name, undefined, 'no rule is added to non-TTL fields');
    });

    test('it should validate nested and grouped TTL fields', async function (assert) {
      class GroupedForm extends Form {
        formFieldGroups = [
          new FormFieldGroup('default', [
            new FormField('config.max_lease_ttl', 'string', { editType: 'ttl', label: 'Max Lease TTL' }),
          ]),
        ];
      }
      const { isValid, state } = new GroupedForm({ config: { max_lease_ttl: '9223372037s' } }).toJSON();

      assert.false(isValid, 'form is invalid');
      assert.deepEqual(
        state['config.max_lease_ttl'].errors,
        ['Max Lease TTL is too large. Maximum allowed value is 106,751 days.'],
        'error is reported on the nested field path'
      );
    });

    test('it should add the rule once when a TTL field is in both formFields and formFieldGroups', async function (assert) {
      // PkiConfigCrlForm and OpenApiForm list the same fields in both collections.
      const ttl = new FormField('expiry', 'string', { editType: 'ttl', label: 'Expiry' });
      class DuplicatedForm extends Form {
        formFields = [ttl];
        formFieldGroups = [new FormFieldGroup('default', [ttl])];
      }
      const { state, invalidFormMessage } = new DuplicatedForm({ expiry: '106752d' }).toJSON();

      assert.deepEqual(
        state.expiry.errors,
        ['Expiry is too large. Maximum allowed value is 106,751 days.'],
        'error is reported once'
      );
      assert.strictEqual(invalidFormMessage, 'There is an error with this form.', 'error is counted once');
    });

    test('it should name unlabeled TTL fields the way the form field label does', async function (assert) {
      class UnlabeledForm extends Form {
        formFields = [new FormField('verification_ttl', 'string', { editType: 'ttl' })];
      }
      const { state } = new UnlabeledForm({ verification_ttl: '106752d' }).toJSON();

      assert.deepEqual(
        state.verification_ttl.errors,
        ['Verification ttl is too large. Maximum allowed value is 106,751 days.'],
        'message uses the humanized field name'
      );
    });
  });
});
