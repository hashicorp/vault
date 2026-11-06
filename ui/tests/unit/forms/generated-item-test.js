/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import GeneratedItemForm from 'vault/forms/generated-item';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

module('Unit | Form | generated item', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    const username = new FormField('username', 'string', { readOnly: true });
    const password = new FormField('password', 'string', { sensitive: true });
    const numUses = new FormField('token_num_uses', 'number');
    const noDefault = new FormField('token_no_default_policy', 'boolean');
    const policies = new FormField('token_policies', undefined, { editType: 'stringArray' });
    const cidrs = new FormField('token_bound_cidrs', undefined, { editType: 'stringArray' });
    const metadata = new FormField('alias_metadata', 'object');
    const tokenType = new FormField('token_type', 'string', { defaultValue: 'default-service' });
    const formFields = [username, password, numUses, noDefault, policies, cidrs, metadata, tokenType];
    this.metadata = {
      formFields,
      formFieldGroups: [
        new FormFieldGroup('default', [username, password]),
        new FormFieldGroup('Tokens', [numUses, noDefault, policies, cidrs, metadata, tokenType]),
      ],
      identifier: 'username',
    };
  });

  test('create payload omits the identifier and untouched blanks, and coerces number fields', function (assert) {
    const form = new GeneratedItemForm(this.metadata, {}, { isNew: true });
    form.set('username', 'payload-user');
    form.set('password', 'disposable-password');
    form.set('token_num_uses', '5');

    assert.strictEqual(form.id, 'payload-user', 'id is the identifier value');
    assert.deepEqual(
      form.toJSON().data,
      { password: 'disposable-password', token_num_uses: 5, token_type: 'default-service' },
      'schema defaults are sent; the identifier and untouched fields are not'
    );
  });

  test('edit payload sends every non-blank value, including zero, false and empty objects', function (assert) {
    const form = new GeneratedItemForm(
      this.metadata,
      {
        username: 'existing',
        token_num_uses: 0,
        token_no_default_policy: false,
        token_policies: ['default'],
        token_bound_cidrs: [],
        alias_metadata: {},
        token_type: 'default',
      },
      { isNew: false }
    );
    form.set('token_num_uses', '7');

    assert.deepEqual(
      form.toJSON().data,
      {
        token_num_uses: 7,
        token_no_default_policy: false,
        token_policies: ['default'],
        alias_metadata: {},
        token_type: 'default',
      },
      'empty arrays and unread sensitive values are omitted'
    );
  });

  test('a blank value the user cleared is still sent', function (assert) {
    const form = new GeneratedItemForm(this.metadata, { username: 'existing', token_type: 'batch' });
    form.set('token_type', '');

    assert.strictEqual(
      form.toJSON().data.token_type,
      '',
      'cleared value is sent so the backend can reset it'
    );
  });

  test('field descriptions render as tooltips, as these pages did before', function (assert) {
    const password = new FormField('password', 'string', { subText: 'Password for this user.' });
    const form = new GeneratedItemForm({
      formFields: [password],
      formFieldGroups: [new FormFieldGroup('default', [password])],
    });
    const [field] = form.formFields;

    assert.strictEqual(field.options.helpText, 'Password for this user.', 'description is a tooltip');
    assert.strictEqual(field.options.subText, undefined, 'description is not repeated inline');
    assert.strictEqual(form.formFieldGroups[0].default[0], field, 'groups render the same field instance');
  });

  test('read-only field groups leave out the identifier', function (assert) {
    const form = new GeneratedItemForm(this.metadata, { username: 'existing' });
    const [defaultGroup] = form.fieldGroups;

    assert.deepEqual(
      defaultGroup.default.map(({ name }) => name),
      ['password'],
      'identifier is shown as the title rather than a row'
    );
  });
});
