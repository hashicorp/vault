/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import oidcConfigHandlers from 'vault/mirage/handlers/oidc-config';
import { SELECTORS, CLIENT_LIST_RESPONSE } from 'vault/tests/helpers/oidc-config';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import OidcProviderForm from 'vault/forms/oidc/provider';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

const ISSUER_URL = 'http://127.0.0.1:8200/v1/identity/oidc/provider/test-provider';

module('Integration | Component | oidc/provider-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    oidcConfigHandlers(this.server);
    this.api = this.owner.lookup('service:api');
    this.apiStub = sinon.stub(this.api.identity, 'oidcWriteProvider').resolves();

    this.scopes = [{ id: 'test-scope' }];
    this.clients = this.api.keyInfoToArray(CLIENT_LIST_RESPONSE, 'name');
    this.onCancel = sinon.spy();
    this.onSave = sinon.spy();

    this.renderComponent = (data) => {
      this.form = new OidcProviderForm(data || {}, { isNew: !data });
      return render(hbs`
        <Oidc::ProviderForm
          @form={{this.form}}
          @scopes={{this.scopes}}
          @clients={{this.clients}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
        />
      `);
    };

    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should save new provider', async function (assert) {
    assert.expect(13);

    await this.renderComponent();

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create Provider', 'Form title renders correct text');
    assert.dom(SELECTORS.providerSaveButton).hasText('Create', 'Save button has correct text');
    assert
      .dom('[data-test-input="issuer"]')
      .hasAttribute('placeholder', 'e.g. https://example.com:8200', 'issuer placeholder text is correct');
    assert.strictEqual(findAll('[data-test-field]').length, 3, 'renders all input fields');
    await click(
      '[data-test-component="search-select"]#oidc-provider-form-scope-select .ember-basic-dropdown-trigger'
    );
    assert.dom('li.ember-power-select-option').hasText('test-scope', 'dropdown renders scopes');

    // check validation errors
    await fillIn('[data-test-input="name"]', ' ');
    await click(SELECTORS.providerSaveButton);

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Name is required. Name cannot contain whitespace.', 'Validation messages are shown for name');
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There are 2 errors with this form.', 'Renders form error count');

    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom('[data-test-component="search-select"]#oidc-provider-form-client-select')
      .exists('Limited radio button shows clients search select');
    await click(
      '[data-test-component="search-select"]#oidc-provider-form-client-select .ember-basic-dropdown-trigger'
    );
    assert.dom('li.ember-power-select-option').hasTextContaining('test-app', 'dropdown renders client name');
    assert.dom('[data-test-smaller-id]').exists('renders smaller client id in dropdown');

    await click('[data-test-oidc-radio="allow-all"]');
    assert
      .dom('[data-test-component="search-select"]#oidc-provider-form-client-select')
      .doesNotExist('Allow all radio button hides search select');

    await fillIn('[data-test-input="name"]', 'test-provider');
    await click(SELECTORS.providerSaveButton);

    assert.true(this.onSave.called, 'onSave callback fires on save success');
    assert.true(this.apiStub.calledWith('test-provider'), 'Request made to save provider');
  });

  test('it should update provider', async function (assert) {
    assert.expect(8);

    const provider = {
      name: 'test-provider',
      allowed_client_ids: ['*'],
      issuer: ISSUER_URL,
      scopes_supported: ['test-scope'],
    };

    await this.renderComponent(provider);

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Edit Provider', 'Title renders correct text');
    assert.dom(SELECTORS.providerSaveButton).hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert
      .dom('[data-test-input="name"]')
      .hasValue('test-provider', 'Name input is populated with model value');
    assert.dom('[data-test-selected-option]').hasText('test-scope', 'model scope is selected');
    assert.dom('[data-test-oidc-radio="allow-all"] input').isChecked('Allow all radio button is selected');
    await click(SELECTORS.providerSaveButton);

    assert.true(this.onSave.called, 'onSave callback fires on save success');
    const { name, ...payload } = provider;
    assert.true(this.apiStub.calledWith(name, payload), 'Request made to save provider');
  });

  test('it should fire callback on cancel', async function (assert) {
    assert.expect(1);

    await this.renderComponent();
    await click(SELECTORS.providerCancelButton);
    assert.true(this.onCancel.called, 'onCancel callback fires on cancel');
  });

  test('it should render fallback for search select', async function (assert) {
    assert.expect(2);

    this.scopes = [];
    this.clients = [];
    await this.renderComponent();

    assert
      .dom(
        '[data-test-component="search-select"]#oidc-provider-form-scope-select [data-test-component="string-list"]'
      )
      .exists('renders fall back for scopes search select');
    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom(
        '[data-test-component="search-select"]#oidc-provider-form-client-select [data-test-component="string-list"]'
      )
      .exists('Radio toggle shows assignments string-list input');
  });

  test('it should render error alerts when API returns an error', async function (assert) {
    assert.expect(2);

    this.apiStub.rejects(getErrorResponse());
    await this.renderComponent();

    await fillIn('[data-test-input="name"]', 'some-provider');
    await click(SELECTORS.providerSaveButton);
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom('[data-test-message-error]').exists('alert banner renders');
  });
});
