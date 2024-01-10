/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | mfa-login-enforcement-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('mfa-login-enforcement');
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
      },
    });
  });

  test('it should render correct fields', async function (assert) {
    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @model={{this.model}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave)}}
      />
    `);

    const fields = {
      name: {
        label: 'Name',
        subText:
          'The name for this enforcement. Giving it a name means that you can refer to it again later. This name will not be editable later.',
      },
      methods: {
        label: 'MFA methods',
        subText: 'The MFA method(s) that this enforcement will apply to.',
      },
      targets: {
        label: 'Targets',
        subText:
          'The list of authentication types, authentication mounts, groups, and/or entities that will require this MFA configuration.',
      },
    };

    const subTexts = this.element.querySelectorAll('[data-test-label-subtext]');
    Object.keys(fields).forEach((field, index) => {
      const { label, subText } = fields[field];
      assert.dom(`[data-test-mlef-label="${field}"]`).hasText(label, `${field} field label renders`);
      assert.dom(subTexts[index]).hasText(subText, `${subText} field label sub text renders`);
    });
    assert.dom('[data-test-mlef-input="name"]').exists(`Name field input renders`);
    assert.dom('[data-test-mlef-search="methods"]').exists('MFA method search select renders');
    assert.dom('[data-test-mlef-select="target-type"]').exists('Target type selector renders');
    assert.dom('[data-test-mlef-select="accessor"]').exists('Auth mount target selector renders by default');
  });

  test('it should render inline', async function (assert) {
    this.errors = this.model.validate().state;
    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @model={{this.model}}
        @isInline={{true}}
        @modelErrors={{this.errors}}
      />
    `);

    assert.dom('[data-test-mlef-input="name"]').exists(`Name field input renders`);
    assert.dom('[data-test-mlef-search="methods"]').doesNotExist('MFA method search select does not render');
    assert.dom('[data-test-mlef-select="target-type"]').exists('Target type selector renders');
    assert
      .dom('[data-test-inline-error-message]')
      .exists({ count: 2 }, 'External validation errors are displayed');
  });

  test('it should display field validation errors on save', async function (assert) {
    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @model={{this.model}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave)}}
      />
    `);

    await click('[data-test-mlef-save]');
    const errors = this.element.querySelectorAll('[data-test-inline-error-message]');
    assert.dom(errors[0]).hasText('Name is required', 'Name error message renders');
    assert.dom(errors[1]).hasText('At least one MFA method is required', 'Methods error message renders');
    assert
      .dom(errors[2])
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
        @model={{this.model}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave) true}}
      />
    `);

    await fillIn('[data-test-mlef-input="name"]', 'bar');
    await click('.ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await fillIn('[data-test-mlef-select="accessor"] select', 'auth_userpass_1234');
    await click('[data-test-mlef-add-target]');
    await click('[data-test-mlef-save]');
    assert.true(this.didSave, 'onSave callback triggered');
    assert.strictEqual(this.model.name, 'bar', 'Name property set on model');
    assert.strictEqual(this.model.mfa_methods.firstObject.id, '123456', 'Mfa method added to model');
    assert.deepEqual(
      this.model.auth_method_accessors,
      ['auth_userpass_1234'],
      'Target saved to correct model property'
    );
  });

  test('it should populate fields with model data', async function (assert) {
    this.model.name = 'foo';
    const [method] = (await this.store.query('mfa-method', {})).toArray();
    this.model.mfa_methods.addObject(method);
    this.model.auth_method_accessors.addObject('auth_userpass_1234');

    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @model={{this.model}}
        @onClose={{fn (mut this.didClose)}}
        @onSave={{fn (mut this.didSave) true}}
      />
    `);

    assert.dom('[data-test-mlef-input="name"]').hasValue('foo', 'Name input is populated');
    assert.dom('.search-select-list-item').includesText('TOTP', 'MFA method type renders in selected option');
    assert
      .dom('.search-select-list-item small')
      .hasText('123456', 'MFA method id renders in selected option');
    assert
      .dom('[data-test-row-label="Authentication mount"]')
      .hasText('Authentication mount', 'Selected target type renders');
    assert
      .dom('[data-test-value-div="Authentication mount"]')
      .hasText('auth_userpass_1234', 'Selected target value renders');

    await click('[data-test-mlef-remove-target]');
    await click('[data-test-mlef-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .includesText('At least one target is required', 'Target is removed');
    assert.notOk(this.model.auth_method_accessors.length, 'Target is removed from appropriate model prop');

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
    this.model.auth_method_accessors.addObject('auth_userpass_1234');
    this.model.auth_method_types.addObject('userpass');
    const [entity] = (await this.store.query('identity/entity', {})).toArray();
    this.model.identity_entities.addObject(entity);
    const [group] = (await this.store.query('identity/group', {})).toArray();
    this.model.identity_groups.addObject(group);

    await render(hbs`
      <Mfa::MfaLoginEnforcementForm
        @model={{this.model}}
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
      { label: 'Group', value: 'bar group 1234', key: 'identity_groups', type: 'identity/group' },
      { label: 'Entity', value: 'foo entity 1234', key: 'identity_entities', type: 'identity/entity' },
    ];

    for (const [index, target] of targets.entries()) {
      // target populated from model
      assert
        .dom(`[data-test-row-label="${target.label}"]`)
        .hasText(target.label, `${target.label} target populated with correct type label`);
      assert
        .dom(`[data-test-value-div="${target.label}"]`)
        .hasText(target.value, `${target.label} target populated with correct value`);
      // remove target
      await click(`[data-test-mlef-remove-target="${target.label}"]`);
      assert
        .dom('[data-test-mlef-target]')
        .exists({ count: targets.length - (index + 1) }, `${target.label} target removed`);
      assert.notOk(this.model[target.key].length, `${target.label} removed from correct model prop`);
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
      assert.ok(this.model[target.key].length, `${target.label} added to correct model prop`);
    }
    assert.dom('[data-test-mlef-target]').exists({ count: 4 }, 'All targets were added back');
  });
});
