/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import { setupMirage } from 'ember-cli-mirage/test-support';
import oidcConfigHandlers from 'vault/mirage/handlers/oidc-config';
import { SELECTORS } from 'vault/tests/helpers/oidc-config';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import OidcClientForm from 'vault/forms/oidc/client';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

const searchSelect = create(ss);

module('Integration | Component | oidc/client-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    oidcConfigHandlers(this.server);
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

    const api = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(api.identity, 'oidcWriteClient').resolves();

    this.renderComponent = (client) => {
      const data = {
        key: 'default',
        id_token_ttl: '24h',
        access_token_ttl: '24h',
        client_type: 'confidential',
        ...(client || {}),
      };
      this.form = new OidcClientForm(data, { isNew: !client });
      this.keys = [{ id: 'default' }];
      this.onCancel = sinon.spy();
      this.onSave = sinon.spy();

      return render(hbs`
        <Oidc::ClientForm
          @form={{this.form}}
          @keys={{this.keys}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
        />
      `);
    };

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

    await this.renderComponent();
    await click(GENERAL.button('More options'));
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create Application', 'Form title renders correct text');
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
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Name is required. Name cannot contain whitespace.', 'Validation messages are shown for name');
    assert
      .dom(GENERAL.validationErrorByAttr('key'))
      .hasText('Key is required.', 'Validation message is shown for key');
    assert.dom(validationErrors[1]).hasText('There are 3 errors with this form.', 'Renders form error count');

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
    assert.true(this.apiStub.calledWith('test-app'), 'API called with correct parameters');
    assert.true(this.onSave.called, 'onSave callback is called on successful save');
  });

  test('it should update client', async function (assert) {
    assert.expect(11);

    await this.renderComponent({ name: 'test-app', client_type: 'public' });
    await click(GENERAL.button('More options'));
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Edit Application', 'Title renders correct text');
    assert.dom(SELECTORS.clientSaveButton).hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test-app', 'Name input is populated with model value');
    assert.dom('[data-test-input="key"]').isDisabled('Signing key input is disabled');
    assert.dom('[data-test-input="key"]').hasValue('default', 'Key input populated with default');
    assert
      .dom('[data-test-input-group="client_type"] input')
      .isDisabled('client type input is disabled on edit');
    assert
      .dom('[data-test-input-group="client_type"] input#public')
      .isChecked('Correct radio button is selected');
    assert.dom('[data-test-oidc-radio="allow-all"] input').isChecked('Allow all radio button is selected');
    await click(SELECTORS.clientSaveButton);
    assert.true(this.apiStub.calledWith('test-app'), 'API called with correct parameters');
    assert.true(this.onSave.called, 'onSave callback is called on successful save');
  });

  test('it should show create assignment modal', async function (assert) {
    assert.expect(3);

    await this.renderComponent();
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

    this.server.get('/identity/oidc/assignment', () => overrideResponse(403));

    await this.renderComponent();

    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom('[data-test-component="string-list"]')
      .exists('Radio toggle shows assignments string-list input');
  });

  test('it should render error alerts when API returns an error', async function (assert) {
    assert.expect(2);

    this.apiStub.rejects(getErrorResponse());

    await this.renderComponent();
    await fillIn('[data-test-input="name"]', 'test-app');
    await click(SELECTORS.clientSaveButton);
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom('[data-test-message-error]').exists('alert banner renders');
  });
});
