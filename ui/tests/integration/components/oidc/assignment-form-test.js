/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { overrideMirageResponse } from 'vault/tests/helpers/oidc-config';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | oidc/assignment-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'oidcConfig';
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', () => {});
    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
      },
    });
  });

  test('it should save new assignment', async function (assert) {
    assert.expect(5);
    this.model = this.store.createRecord('oidc/assignment');
    this.server.post('/identity/oidc/assignment/test', (schema, req) => {
      assert.ok(true, 'Request made to save assignment');
      return JSON.parse(req.requestBody);
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-oidc-assignment-save]').hasText('Create', 'Save button has correct label');
    await click('[data-test-oidc-assignment-save]');
    assert
      .dom('[data-test-inline-alert]')
      .hasText('Name is required.', 'Validation message is shown for name');
    assert.strictEqual(
      findAll('[data-test-inline-error-message]').length,
      2,
      `there are two validations errors.`
    );
    await fillIn('[data-test-input="name"]', 'test');
    await click('[data-test-component="search-select"]#entities .ember-basic-dropdown-trigger');
    await click('.ember-power-select-option');
    await click('[data-test-oidc-assignment-save]');
  });

  test('it should populate fields with model data on edit view and update an assignment', async function (assert) {
    assert.expect(5);

    this.store.pushPayload('oidc/assignment', {
      modelName: 'oidc/assignment',
      name: 'test',
      entity_ids: ['1234-12345'],
      group_ids: ['abcdef-123'],
    });
    this.model = this.store.peekRecord('oidc/assignment', 'test');

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-oidc-assignment-save]').hasText('Update', 'Save button has correct label');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test', 'Name input is populated with model value');
    assert
      .dom('[data-test-search-select="entities"] [data-test-smaller-id]')
      .hasText('1234-12345', 'entity id renders in selected option');
    assert
      .dom('[data-test-search-select="groups"] [data-test-smaller-id]')
      .hasText('abcdef-123', 'group id renders in selected option');
  });

  test('it should use fallback component on create if no permissions for entities or groups', async function (assert) {
    assert.expect(2);
    this.model = this.store.createRecord('oidc/assignment');
    this.server.get('/identity/entity/id', () => overrideMirageResponse(403));
    this.server.get('/identity/group/id', () => overrideMirageResponse(403));

    await render(hbs`
      <Oidc::AssignmentForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert
      .dom('[data-test-component="search-select"]#entities [data-test-component="string-list"]')
      .exists('entities string list fallback component exists');
    assert
      .dom('[data-test-component="search-select"]#groups [data-test-component="string-list"]')
      .exists('groups string list fallback component exists');
  });

  test('it should use fallback component on edit if no permissions for entities or groups', async function (assert) {
    assert.expect(8);
    this.store.pushPayload('oidc/assignment', {
      modelName: 'oidc/assignment',
      name: 'test',
      entity_ids: ['1234-12345'],
      group_ids: ['abcdef-123'],
    });
    this.model = this.store.peekRecord('oidc/assignment', 'test');
    this.server.get('/identity/entity/id', () => overrideMirageResponse(403));
    this.server.get('/identity/group/id', () => overrideMirageResponse(403));

    await render(hbs`
    <Oidc::AssignmentForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
  `);

    assert
      .dom('[data-test-component="search-select"]#entities [data-test-component="string-list"]')
      .exists('entities string list fallback component exists');
    assert
      .dom('[data-test-component="search-select"]#entities [data-test-string-list-input="0"]')
      .hasValue('1234-12345', 'first row pre-populated with model entity');
    assert
      .dom(
        '[data-test-component="search-select"]#entities [data-test-string-list-row="0"] [data-test-string-list-button="delete"]'
      )
      .exists('first row renders delete icon');
    assert
      .dom(
        '[data-test-component="search-select"]#entities [data-test-string-list-row="1"] [data-test-string-list-button="add"]'
      )
      .exists('second row renders add icon');

    assert
      .dom('[data-test-component="search-select"]#groups [data-test-component="string-list"]')
      .exists('groups string list fallback component exists');
    assert
      .dom('[data-test-component="search-select"]#groups [data-test-string-list-input="0"]')
      .hasValue('abcdef-123', 'first row pre-populated with model group');
    assert
      .dom(
        '[data-test-component="search-select"]#groups [data-test-string-list-row="0"] [data-test-string-list-button="delete"]'
      )
      .exists('first row renders delete icon');
    assert
      .dom(
        '[data-test-component="search-select"]#groups [data-test-string-list-row="1"] [data-test-string-list-button="add"]'
      )
      .exists('second row renders add icon');
  });
});
