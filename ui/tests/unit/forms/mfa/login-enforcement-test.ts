/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import MfaLoginEnforcementForm from 'vault/forms/mfa/login-enforcement';

module('Unit | Form | mfa/login-enforcement', function (hooks) {
  setupTest(hooks);

  test('it rejects empty name, mfa_methods, and targets', function (assert) {
    const form = new MfaLoginEnforcementForm({}, { isNew: true });
    const { isValid, state } = form.toJSON();

    assert.false(isValid, 'empty form is invalid');
    assert.deepEqual(state['name']?.errors, ['Name is required'], 'name presence error');
    assert.deepEqual(
      state['mfa_methods']?.errors,
      ['At least one MFA method is required'],
      'mfa_methods presence error'
    );
    assert.deepEqual(
      state['targets']?.errors,
      [
        "At least one target is required. If you've selected one, click 'Add' to make sure it's added to this enforcement.",
      ],
      'targets validation error'
    );
  });

  test('it accepts valid form with auth_method_accessors', function (assert) {
    const form = new MfaLoginEnforcementForm(
      {
        name: 'test-enforcement',
        mfa_methods: ['method-id-1'],
        auth_method_accessors: ['auth_userpass_123'],
      },
      { isNew: true }
    );
    const { isValid, data } = form.toJSON();

    assert.true(isValid, 'form with accessor target is valid');
    assert.strictEqual(data.name, 'test-enforcement');
    assert.deepEqual(data.mfa_methods, ['method-id-1']);
    assert.deepEqual(data.auth_method_accessors, ['auth_userpass_123']);
  });

  test('it accepts valid form with auth_method_types', function (assert) {
    const form = new MfaLoginEnforcementForm(
      {
        name: 'test-enforcement',
        mfa_methods: ['method-id-1'],
        auth_method_types: ['userpass'],
      },
      { isNew: true }
    );
    const { isValid } = form.toJSON();

    assert.true(isValid, 'form with auth method type target is valid');
  });

  test('it accepts valid form with identity_entity_ids', function (assert) {
    const form = new MfaLoginEnforcementForm(
      {
        name: 'test-enforcement',
        mfa_methods: ['method-id-1'],
        identity_entity_ids: ['entity-id-1'],
      },
      { isNew: true }
    );
    const { isValid } = form.toJSON();

    assert.true(isValid, 'form with entity target is valid');
  });

  test('it accepts valid form with identity_group_ids', function (assert) {
    const form = new MfaLoginEnforcementForm(
      {
        name: 'test-enforcement',
        mfa_methods: ['method-id-1'],
        identity_group_ids: ['group-id-1'],
      },
      { isNew: true }
    );
    const { isValid } = form.toJSON();

    assert.true(isValid, 'form with group target is valid');
  });

  test('it disables name edit when isNew is false', function (assert) {
    const newForm = new MfaLoginEnforcementForm({}, { isNew: true });
    const nameFieldNew = newForm.formFieldGroups[0]?.['name']?.find((f) => f.name === 'name');
    assert.false(nameFieldNew?.options?.editDisabled, 'name is editable on new form');

    const editForm = new MfaLoginEnforcementForm({ name: 'existing-enforcement' }, { isNew: false });
    const nameFieldEdit = editForm.formFieldGroups[0]?.['name']?.find((f) => f.name === 'name');
    assert.true(nameFieldEdit?.options?.editDisabled, 'name is not editable on edit form');
  });
});
