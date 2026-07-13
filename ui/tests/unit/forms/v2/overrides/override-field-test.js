/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { configBuilder, overrideFieldsInSection } from 'vault/forms/v2/overrides/override-field';

module('Unit | forms/v2/overrides | override-field', function (hooks) {
  hooks.beforeEach(function () {
    this.generatedFormConfig = {
      name: 'test-form',
      path: '/test/path',
      title: 'Test Form',
      payload: {
        username: '',
        email: '',
        enabled: false,
      },
      submit: async () => ({ success: true }),
      sections: [
        {
          name: 'section1',
          title: 'User Information',
          description: 'Enter user details',
          fields: [
            {
              name: 'username',
              type: 'TextInput',
              label: 'Username',
              helperText: 'Enter your username',
            },
            {
              name: 'email',
              type: 'TextInput',
              label: 'Email',
            },
          ],
        },
        {
          name: 'section2',
          title: 'Settings',
          fields: [
            {
              name: 'enabled',
              type: 'Toggle',
              label: 'Enable feature',
            },
          ],
        },
      ],
    };
  });

  test('it adds a field to the config: addField', function (assert) {
    const originalSectionCount = this.generatedFormConfig.sections.length;
    const originalFieldCount = this.generatedFormConfig.sections[0].fields.length;

    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .addField('section1', {
        name: 'password',
        type: 'MaskedInput',
        label: 'Password',
      })
      .build();

    assert.strictEqual(
      this.generatedFormConfig.sections.length,
      originalSectionCount,
      'original config section count unchanged'
    );
    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields.length,
      originalFieldCount,
      'original config field count unchanged after adding field'
    );
    assert.notOk(
      this.generatedFormConfig.sections[0].fields.find((f) => f.name === 'password'),
      'password field does not exist in original config'
    );

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    assert.strictEqual(section1.fields.length, 3, 'section1 has 3 fields after adding');
    assert.ok(
      section1.fields.find((f) => f.name === 'password'),
      'password field exists in section1'
    );
  });

  test('it removes a field from the config: removeField', function (assert) {
    const originalSectionCount = this.generatedFormConfig.sections.length;
    const originalFieldCount = this.generatedFormConfig.sections[0].fields.length;

    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .removeField('section1', 'username')
      .build();

    assert.strictEqual(
      this.generatedFormConfig.sections.length,
      originalSectionCount,
      'original config section count unchanged'
    );
    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields.length,
      originalFieldCount,
      'original config field count unchanged after removing field'
    );
    assert.ok(
      this.generatedFormConfig.sections[0].fields.find((f) => f.name === 'username'),
      'username field still exists in original config'
    );

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    assert.strictEqual(section1.fields.length, 1, 'section1 has 1 field after removing');
    assert.notOk(
      section1.fields.find((f) => f.name === 'username'),
      'username field removed from section1'
    );
  });

  test('it updates a field in the config: updateField', function (assert) {
    const originalLabel = this.generatedFormConfig.sections[0].fields[0].label;
    const originalType = this.generatedFormConfig.sections[0].fields[0].type;

    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .updateField('section1', 'username', {
        label: 'Updated Username',
        type: 'TextArea',
      })
      .build();

    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields[0].label,
      originalLabel,
      'original field label unchanged'
    );
    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields[0].type,
      originalType,
      'original field type unchanged'
    );

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    const usernameField = section1.fields.find((f) => f.name === 'username');
    assert.strictEqual(usernameField.label, 'Updated Username', 'field label updated');
    assert.strictEqual(usernameField.type, 'TextArea', 'field type updated');
  });

  test('it moves a field between sections in the config: moveField', function (assert) {
    const originalSection1FieldCount = this.generatedFormConfig.sections[0].fields.length;
    const originalSection2FieldCount = this.generatedFormConfig.sections[1].fields.length;

    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .moveField('username', 'section1', 'section2')
      .build();

    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields.length,
      originalSection1FieldCount,
      'original section1 field count unchanged'
    );
    assert.strictEqual(
      this.generatedFormConfig.sections[1].fields.length,
      originalSection2FieldCount,
      'original section2 field count unchanged'
    );
    assert.ok(
      this.generatedFormConfig.sections[0].fields.find((f) => f.name === 'username'),
      'username field still in section1 in original config'
    );
    assert.notOk(
      this.generatedFormConfig.sections[1].fields.find((f) => f.name === 'username'),
      'username field not in section2 in original config'
    );

    const section1 = overriddenConfig.sections[0];
    const section2 = overriddenConfig.sections[1];

    assert.strictEqual(section1.fields.length, 1, 'section1 has 1 field after moving');
    assert.notOk(
      section1.fields.find((f) => f.name === 'username'),
      'username removed from section1'
    );

    assert.strictEqual(section2.fields.length, 2, 'section2 has 2 fields after moving');
    assert.ok(
      section2.fields.find((f) => f.name === 'username'),
      'username added to section2'
    );
  });

  test('it reorders fields within a section in the config: reorderFields', function (assert) {
    const originalFirstFieldName = this.generatedFormConfig.sections[0].fields[0].name;
    const originalSecondFieldName = this.generatedFormConfig.sections[0].fields[1].name;

    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .reorderFields('section1', ['email', 'username'])
      .build();

    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields[0].name,
      originalFirstFieldName,
      'original first field unchanged'
    );
    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields[1].name,
      originalSecondFieldName,
      'original second field unchanged'
    );

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    assert.strictEqual(section1.fields[0].name, 'email', 'email is first field');
    assert.strictEqual(section1.fields[1].name, 'username', 'username is second field');
  });

  test('it supports method chaining', function (assert) {
    const overriddenConfig = configBuilder(this.generatedFormConfig)
      .addField('section1', { name: 'password', type: 'MaskedInput', label: 'Password' })
      .updateField('section1', 'username', { label: 'Updated Username' })
      .removeField('section2', 'enabled')
      .build();

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    const section2 = overriddenConfig.sections.find((s) => s.name === 'section2');

    assert.strictEqual(section1.fields.length, 3, 'section1 has 3 fields');
    assert.ok(
      section1.fields.find((f) => f.name === 'password'),
      'password field added'
    );
    assert.strictEqual(
      section1.fields.find((f) => f.name === 'username').label,
      'Updated Username',
      'username label updated'
    );
    assert.strictEqual(section2.fields.length, 0, 'section2 has 0 fields');
  });

  test('it throws error when section not found', function (assert) {
    assert.throws(
      () => {
        configBuilder(this.generatedFormConfig).addField('nonexistent', { name: 'test', type: 'TextInput' });
      },
      /Section "nonexistent" not found/,
      'throws error for nonexistent section'
    );
  });

  test('it throws error when field not found', function (assert) {
    assert.throws(
      () => {
        configBuilder(this.generatedFormConfig).updateField('section1', 'nonexistent', { label: 'Test' });
      },
      /Field "nonexistent" not found/,
      'throws error for nonexistent field'
    );
  });

  test('overrideFieldsInSection: it overrides multiple fields in a section', function (assert) {
    const overriddenConfig = overrideFieldsInSection(this.generatedFormConfig, 'section1', {
      username: {
        label: 'User Name',
        helperText: 'Your unique username',
      },
      email: {
        type: 'TextArea',
      },
    });

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    const usernameField = section1.fields.find((f) => f.name === 'username');
    const emailField = section1.fields.find((f) => f.name === 'email');

    assert.strictEqual(usernameField.label, 'User Name', 'username label updated');
    assert.strictEqual(usernameField.helperText, 'Your unique username', 'username helperText updated');
    assert.strictEqual(usernameField.type, 'TextInput', 'username type preserved');

    assert.strictEqual(emailField.type, 'TextArea', 'email type updated');
    assert.strictEqual(emailField.label, 'Email', 'email label preserved');
  });

  test('overrideFieldsInSection: it preserves unmodified fields', function (assert) {
    const overriddenConfig = overrideFieldsInSection(this.generatedFormConfig, 'section1', {
      username: {
        label: 'Updated Username',
      },
    });

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    const emailField = section1.fields.find((f) => f.name === 'email');

    assert.strictEqual(emailField.type, 'TextInput', 'email type unchanged');
    assert.strictEqual(emailField.label, 'Email', 'email label unchanged');
  });

  test('overrideFieldsInSection: it preserves other sections', function (assert) {
    const overriddenConfig = overrideFieldsInSection(this.generatedFormConfig, 'section1', {
      username: { label: 'New Label' },
    });

    const section2 = overriddenConfig.sections.find((s) => s.name === 'section2');
    assert.ok(section2, 'section2 exists');
    assert.strictEqual(section2.fields.length, 1, 'section2 has 1 field');
    assert.strictEqual(section2.fields[0].name, 'enabled', 'section2 field unchanged');
  });

  test('overrideFieldsInSection: it does not mutate the original config', function (assert) {
    const originalUsername = this.generatedFormConfig.sections[0].fields[0].label;

    overrideFieldsInSection(this.generatedFormConfig, 'section1', {
      username: { label: 'Changed Label' },
    });

    assert.strictEqual(
      this.generatedFormConfig.sections[0].fields[0].label,
      originalUsername,
      'original config unchanged'
    );
  });

  test('overrideFieldsInSection: it throws error for nonexistent section', function (assert) {
    assert.throws(
      () => {
        overrideFieldsInSection(this.generatedFormConfig, 'nonexistent', { username: { label: 'Test' } });
      },
      /Section "nonexistent" not found/,
      'throws error for nonexistent section'
    );
  });

  test('overrideFieldsInSection: it throws error for nonexistent field', function (assert) {
    assert.throws(
      () => {
        overrideFieldsInSection(this.generatedFormConfig, 'section1', { nonexistent: { label: 'Test' } });
      },
      /Field "nonexistent" not found/,
      'throws error for nonexistent field'
    );
  });

  test('overrideFieldsInSection: it handles empty overrides object', function (assert) {
    const overriddenConfig = overrideFieldsInSection(this.generatedFormConfig, 'section1', {});

    const section1 = overriddenConfig.sections.find((s) => s.name === 'section1');
    assert.strictEqual(section1.fields.length, 2, 'section1 still has 2 fields');
    assert.strictEqual(section1.fields[0].label, 'Username', 'username label unchanged');
  });
});
