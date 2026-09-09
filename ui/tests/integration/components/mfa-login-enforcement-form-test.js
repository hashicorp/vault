/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import MfaLoginEnforcementForm from 'vault/forms/mfa/login-enforcement';

module('Integration | Component | mfa-login-enforcement-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.methods = [{ id: '123456', type: 'totp' }];
    this.form = new MfaLoginEnforcementForm({}, { isNew: true });
    this.server.get('/sys/auth', () => ({
      data: { 'userpass/': { type: 'userpass', accessor: 'auth_userpass_1234' } },
    }));
    this.server.get('/identity/mfa/method', () => ({
      data: {
        key_info: {
          123456: { type: 'totp' },
        },
        keys: ['123456'],
      },
    }));
    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
        // TODO: add labels to enforcement targets key/value style inputs
        'select-name': { enabled: false },
        'aria-prohibited-attr': { enabled: false },
      },
    });
  });

  test('it should render correct fields', async function (assert) {
    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @form={{this.form}}
        @methods={{this.methods}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave)}}
      />
    `);

    const subTexts = this.element.querySelectorAll('[data-test-label-subtext]');

    assert.dom(GENERAL.inputByAttr('name')).exists('Name field input renders');
    assert.dom(GENERAL.labelFor('name')).hasText('Name', 'name field label renders');
    assert
      .dom('#helper-text-name')
      .hasText(
        'The name for this enforcement. Giving it a name means that you can refer to it again later. This name will not be editable later.',
        'name field label sub text renders'
      );
    assert.dom(GENERAL.labelFor('mfa_methods')).hasText('MFA methods', 'methods field label renders');
    assert
      .dom(subTexts[0])
      .hasText(
        'The MFA method(s) that this enforcement will apply to.',
        'methods field label sub text renders'
      );
    assert.dom('label[for="targets"]').hasText('Targets', 'targets field label renders');
    assert
      .dom(subTexts[1])
      .hasText(
        'The list of authentication types, authentication mounts, groups, and/or entities that will require this MFA configuration.',
        'targets field label sub text renders'
      );
    assert.dom('[data-test-mlef-search="methods"]').exists('MFA method search select renders');
    assert.dom('[data-test-mlef-select="target-type"]').exists('Target type selector renders');
    assert.dom('[data-test-mlef-select="accessor"]').exists('Auth mount target selector renders by default');
  });

  test('it should render inline', async function (assert) {
    const { state } = this.form.toJSON();
    this.errors = state;
    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @form={{this.form}}
        @methods={{this.methods}}
        @isInline={{true}}
        @modelErrors={{this.errors}}
      />
    `);

    assert.dom(GENERAL.inputByAttr('name')).exists(`Name field input renders`);
    assert.dom('[data-test-mlef-search="methods"]').doesNotExist('MFA method search select does not render');
    assert.dom('[data-test-mlef-select="target-type"]').exists('Target type selector renders');

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .exists({ count: 1 }, 'Name field external validation errors are displayed');
    assert
      .dom('[data-test-inline-error-message]')
      .exists({ count: 1 }, 'Targets field external validation errors are displayed');
  });

  test('it should display field validation errors on save', async function (assert) {
    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @form={{this.form}}
        @methods={{this.methods}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave)}}
      />
    `);

    await click('[data-test-mlef-save]');
    const errors = this.element.querySelectorAll('[data-test-inline-error-message]');
    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Name is required', 'Name error message renders');
    assert.dom(errors[0]).hasText('At least one MFA method is required', 'Methods error message renders');
    assert
      .dom(errors[1])
      .hasText(
        "At least one target is required. If you've selected one, click 'Add' to make sure it's added to this enforcement.",
        'Targets error message renders'
      );
  });

  test('it should save new enforcement', async function (assert) {
    assert.expect(5);

    this.server.post('/identity/mfa/login-enforcement/bar', () => {
      assert.ok(true, 'save request sent to server');
      return {};
    });

    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @form={{this.form}}
        @methods={{this.methods}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave) true}}
      />
    `);

    await fillIn(GENERAL.inputByAttr('name'), 'bar');
    await click('.ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await fillIn('[data-test-mlef-select="accessor"] select', 'auth_userpass_1234');
    await click('[data-test-mlef-add-target]');
    await click('[data-test-mlef-save]');
    assert.true(this.didSave, 'onSave callback triggered');
    assert.strictEqual(this.form.name, 'bar', 'Name property set on form');
    assert.deepEqual(this.form.mfa_methods, ['123456'], 'Mfa method added to form');
    assert.deepEqual(
      this.form.auth_method_accessors,
      ['auth_userpass_1234'],
      'Target saved to correct model property'
    );
  });

  test('it should populate fields with model data', async function (assert) {
    this.form = new MfaLoginEnforcementForm(
      {
        name: 'foo',
        mfa_methods: [{ id: '123456', type: 'totp', displayName: 'TOTP' }],
        auth_method_accessors: ['auth_userpass_1234'],
      },
      { isNew: true }
    );

    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @form={{this.form}}
        @methods={{this.methods}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave) true}}
      />
    `);

    assert.dom(GENERAL.inputByAttr('name')).hasValue('foo', 'Name input is populated');

    assert
      .dom('.search-select-list-item')
      .includesText('TOTP 123456', 'MFA method type renders in selected option');
    assert
      .dom('.search-select-list-item small')
      .hasText('123456', 'MFA method id renders in selected option');
    assert
      .dom(GENERAL.infoRowLabel('Authentication mount'))
      .hasText('Authentication mount', 'Selected target type renders');
    assert
      .dom(GENERAL.infoRowValue('Authentication mount'))
      .hasText('auth_userpass_1234', 'Selected target value renders');

    await click('[data-test-mlef-remove-target]');
    await click('[data-test-mlef-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('At least one target is required', 'Target is removed');
    assert.notOk(this.form.auth_method_accessors.length, 'Target is removed from appropriate form prop');

    await fillIn('[data-test-mlef-select="accessor"] select', 'auth_userpass_1234');
    await click('[data-test-mlef-add-target]');
    await click('[data-test-selected-list-button="delete"]');
    await click('[data-test-mlef-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('At least one MFA method is required', 'Target is removed');
  });

  test('it should add and remove targets', async function (assert) {
    assert.expect();

    this.server.get('/identity/entity/id', () => ({
      data: {
        key_info: { 1234: { name: 'foo entity' } },
        keys: ['1234'],
      },
    }));
    this.server.get('/identity/group/id', () => ({
      data: {
        key_info: { 1234: { name: 'bar group' } },
        keys: ['1234'],
      },
    }));
    this.form = new MfaLoginEnforcementForm(
      {
        auth_method_accessors: ['auth_userpass_1234'],
        auth_method_types: ['userpass'],
        identity_entity_ids: ['1234'],
        identity_group_ids: ['1234'],
      },
      { isNew: true }
    );

    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @form={{this.form}}
        @methods={{this.methods}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave) true}}
      />
    `);

    const targets = [
      {
        label: 'Authentication mount',
        value: 'auth_userpass_1234',
        key: 'auth_method_accessors',
        type: 'accessor',
      },
      { label: 'Authentication method', value: 'userpass', key: 'auth_method_types', type: 'method' },
      { label: 'Group', value: 'bar group 1234', key: 'identity_group_ids', type: 'identity/group' },
      { label: 'Entity', value: 'foo entity 1234', key: 'identity_entity_ids', type: 'identity/entity' },
    ];

    for (const [index, target] of targets.entries()) {
      // target populated from model
      assert
        .dom(GENERAL.infoRowLabel(target.label))
        .hasText(target.label, `${target.label} target populated with correct type label`);
      assert
        .dom(GENERAL.infoRowValue(target.label))
        .hasText(target.value, `${target.label} target populated with correct value`);
      // remove target
      await click(`[data-test-mlef-remove-target="${target.label}"]`);
      assert
        .dom('[data-test-mlef-target]')
        .exists({ count: targets.length - (index + 1) }, `${target.label} target removed`);
      assert.notOk(this.form[target.key].length, `${target.label} removed from correct form prop`);
    }
    // add targets
    for (const target of targets) {
      await fillIn('[data-test-mlef-select="target-type"] select', target.type);
      if (['Group', 'Entity'].includes(target.label)) {
        await click(`[data-test-mlef-search="${target.type}"] .ember-basic-dropdown-trigger`);
        await click('.ember-power-select-option');
      } else {
        const key = target.label === 'Authentication method' ? 'auth-method' : 'accessor';
        const value = target.label === 'Authentication method' ? 'userpass' : 'auth_userpass_1234';
        await fillIn(`[data-test-mlef-select="${key}"] select`, value);
      }
      await click('[data-test-mlef-add-target]');
      assert.ok(this.form[target.key].length, `${target.label} added to correct form prop`);
    }
    assert.dom('[data-test-mlef-target]').exists({ count: 4 }, 'All targets were added back');
  });
});
