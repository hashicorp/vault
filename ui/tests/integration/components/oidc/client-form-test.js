/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import {
  OIDC_BASE_URL,
  SELECTORS,
  overrideMirageResponse,
  overrideCapabilities,
} from 'vault/tests/helpers/oidc-config';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const searchSelect = create(ss);

module('Integration | Component | oidc/client-form', function (hooks) {
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
    this.server.get('/identity/oidc/key', () => {
      return {
        request_id: 'key-list-id',
        data: {
          keys: ['default'],
        },
      };
    });
    this.server.get('/identity/oidc/assignment', () => {
      return {
        request_id: 'assignment-list-id',
        data: {
          keys: ['allow_all', 'assignment-1'],
        },
      };
    });
    this.server.get('/identity/oidc/assignment/assignment-1', () => {
      return {
        request_id: 'assignment-1-id',
        data: {
          entity_ids: ['1234-12345'],
          group_ids: ['abcdef-123'],
        },
      };
    });
    setRunOptions({
      rules: {
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
        'color-contrast': { enabled: false },
      },
    });
  });

  test('it should save new client', async function (assert) {
    assert.expect(14);

    this.server.post('/identity/oidc/client/test-app', (schema, req) => {
      assert.ok(true, 'Request made to save client');
      return JSON.parse(req.requestBody);
    });
    this.model = this.store.createRecord('oidc/client');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);
    await click('[data-test-toggle-group="More options"]');
    assert
      .dom('[data-test-oidc-client-title]')
      .hasText('Create Application', 'Form title renders correct text');
    assert.dom(SELECTORS.clientSaveButton).hasText('Create', 'Save button has correct text');
    assert.strictEqual(findAll('[data-test-field]').length, 6, 'renders all attribute fields');
    assert
      .dom('[data-test-oidc-radio="allow-all"] input')
      .isChecked('Allow all radio button selected by default');
    assert.dom('[data-test-ttl-value="ID Token TTL"]').hasValue('1', 'ttl defaults to 24h');
    assert.dom('[data-test-ttl-value="Access Token TTL"]').hasValue('1', 'ttl defaults to 24h');
    assert.dom('[data-test-selected-option]').hasText('default', 'Search select has default key selected');

    // check validation errors
    await fillIn('[data-test-input="name"]', ' ');
    await click('[data-test-selected-list-button="delete"]');
    await click(SELECTORS.clientSaveButton);

    const validationErrors = findAll(SELECTORS.inlineAlert);
    assert
      .dom(validationErrors[0])
      .hasText('Name is required. Name cannot contain whitespace.', 'Validation messages are shown for name');
    assert.dom(validationErrors[1]).hasText('Key is required.', 'Validation message is shown for key');
    assert.dom(validationErrors[2]).hasText('There are 3 errors with this form.', 'Renders form error count');

    // fill out form with valid inputs
    await clickTrigger();
    await fillIn('.ember-power-select-search input', 'default');
    await searchSelect.options.objectAt(0).click();

    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom('[data-test-search-select-with-modal]')
      .exists('Limited radio button shows assignments search select');

    await clickTrigger();
    assert.dom('li.ember-power-select-option').hasText('assignment-1', 'dropdown renders assignments');
    await fillIn('[data-test-input="name"]', 'test-app');
    await click(SELECTORS.clientSaveButton);
  });

  test('it should update client', async function (assert) {
    assert.expect(11);

    this.server.post('/identity/oidc/client/test-app', (schema, req) => {
      assert.ok(true, 'Request made to save client');
      return JSON.parse(req.requestBody);
    });

    this.store.pushPayload('oidc/client', {
      modelName: 'oidc/client',
      name: 'test-app',
      clientType: 'public',
    });

    this.model = this.store.peekRecord('oidc/client', 'test-app');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);
    await click('[data-test-toggle-group="More options"]');
    assert.dom('[data-test-oidc-client-title]').hasText('Edit Application', 'Title renders correct text');
    assert.dom(SELECTORS.clientSaveButton).hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test-app', 'Name input is populated with model value');
    assert.dom('[data-test-input="key"]').isDisabled('Signing key input is disabled');
    assert.dom('[data-test-input="key"]').hasValue('default', 'Key input populated with default');
    assert.dom('[data-test-input="clientType"] input').isDisabled('client type input is disabled on edit');
    assert
      .dom('[data-test-input="clientType"] input#confidential')
      .isChecked('Correct radio button is selected');
    assert.dom('[data-test-oidc-radio="allow-all"] input').isChecked('Allow all radio button is selected');
    await click(SELECTORS.clientSaveButton);
  });

  test('it should rollback attributes or unload record on cancel', async function (assert) {
    assert.expect(4);
    this.model = this.store.createRecord('oidc/client');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click(SELECTORS.clientCancelButton);
    assert.true(this.model.isDestroyed, 'New model is unloaded on cancel');

    this.store.pushPayload('oidc/client', {
      modelName: 'oidc/client',
      name: 'test-app',
      assignments: ['allow_all'],
      redirectUris: [],
    });
    this.model = this.store.peekRecord('oidc/client', 'test-app');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await fillIn('[data-test-input="redirectUris"] [data-test-string-list-input="0"]', 'some-url.com');
    await click('[data-test-string-list-button="add"]');
    await click(SELECTORS.clientCancelButton);
    assert.strictEqual(this.model.redirectUris, undefined, 'Model attributes rolled back on cancel');
  });

  test('it should show create assignment modal', async function (assert) {
    assert.expect(3);
    this.model = this.store.createRecord('oidc/client');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
          `);
    await click('[data-test-oidc-radio="limited"]');
    await clickTrigger();
    await fillIn('.ember-power-select-search input', 'test-new');
    await searchSelect.options.objectAt(0).click();
    assert.dom('#search-select-modal').exists('modal with form opens');
    assert.dom('[data-test-modal-title]').hasText('Create new assignment', 'Create assignment modal renders');
    await click(SELECTORS.assignmentCancelButton);
    assert.dom('#search-select-modal').doesNotExist('modal disappears onCancel');
  });

  test('it should render fallback for search select', async function (assert) {
    assert.expect(1);
    this.model = this.store.createRecord('oidc/client');
    this.server.get('/identity/oidc/assignment', () => overrideMirageResponse(403));
    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom('[data-test-component="string-list"]')
      .exists('Radio toggle shows assignments string-list input');
  });

  test('it should render error alerts when API returns an error', async function (assert) {
    assert.expect(2);
    this.model = this.store.createRecord('oidc/client');
    this.server.post('/sys/capabilities-self', () => overrideCapabilities(OIDC_BASE_URL + '/clients'));
    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);
    await fillIn('[data-test-input="name"]', 'test-app');
    await click(SELECTORS.clientSaveButton);
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom('[data-test-message-error]').exists('alert banner renders');
  });
});
