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
import { CLIENT_LIST_RESPONSE, SELECTORS } from 'vault/tests/helpers/oidc-config';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import OidcKeyForm from 'vault/forms/oidc/key';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import { clickTrigger } from 'ember-power-select/test-support/helpers';

module('Integration | Component | oidc/key-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    oidcConfigHandlers(this.server);

    const api = this.owner.lookup('service:api');
    sinon.stub(api.identity, 'oidcReadClient').resolves({ data: CLIENT_LIST_RESPONSE });
    this.writeStub = sinon.stub(api.identity, 'oidcWriteKey').resolves();

    this.onSave = sinon.spy();
    this.onCancel = sinon.spy();
    this.clients = [{ name: 'app-1' }];

    this.renderComponent = (client = {}) => {
      const defaultValues = {
        algorithm: 'RS256',
        rotation_period: '24h',
        verification_ttl: '24h',
      };
      const data = { ...defaultValues, ...client };
      this.form = new OidcKeyForm(data, { isNew: !Object.keys(client).length });
      return render(hbs`
        <Oidc::KeyForm
          @form={{this.form}}
          @clients={{this.clients}}
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
      },
    });
  });

  test('it should save new key', async function (assert) {
    assert.expect(8);

    await this.renderComponent();
    assert.dom(SELECTORS.keySaveButton).hasText('Create', 'Save button has correct text');
    assert.dom('[data-test-input="algorithm"]').hasValue('RS256', 'default algorithm is correct');
    assert.strictEqual(findAll('[data-test-field]').length, 4, 'renders all input fields');

    // check validation errors
    await fillIn('[data-test-input="name"]', ' ');
    await click(SELECTORS.keySaveButton);

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Name is required. Name cannot contain whitespace.', 'Validation messages are shown for name');
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There are 2 errors with this form.', 'Renders form error count');

    assert.dom('[data-test-oidc-radio="limited"] input').isDisabled('limit radio button disabled on create');
    await fillIn('[data-test-input="name"]', 'test-key');

    await click(SELECTORS.keySaveButton);
    assert.true(this.onSave.calledOnce, 'onSave callback fires on save success');
    assert.true(this.writeStub.calledWith('test-key'), 'API called to save key with correct parameters');
  });

  test('it should update key and limit access to selected applications', async function (assert) {
    assert.expect(11);

    await this.renderComponent({ name: 'test-key', allowed_client_ids: ['*'] });
    assert.dom(SELECTORS.keySaveButton).hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test-key', 'Name input is populated with model value');
    assert.dom('[data-test-oidc-radio="allow-all"] input').isChecked('Allow all radio button is selected');

    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom('[data-test-component="search-select"]#oidc-key-form-client-select')
      .exists('Limited radio button shows clients search select');
    await clickTrigger();
    assert.strictEqual(findAll('li.ember-power-select-option').length, 1, 'dropdown only renders one option');
    assert
      .dom('li.ember-power-select-option')
      .hasTextContaining('app-1', 'dropdown contains client that references key');
    assert.dom('[data-test-smaller-id]').exists('renders smaller client id in dropdown');

    await click('[data-test-oidc-radio="allow-all"]');
    assert
      .dom('[data-test-component="search-select"]#oidc-key-form-client-select')
      .doesNotExist('Allow all radio button hides search select');

    await click(SELECTORS.keySaveButton);
    assert.true(this.onSave.calledOnce, 'onSave callback fires on save success');
    assert.true(this.writeStub.calledWith('test-key'), 'API called to save key with correct parameters');
  });

  test('it should fire callback on cancel', async function (assert) {
    assert.expect(1);

    await this.renderComponent({ name: 'test-key' });
    await click(SELECTORS.keyCancelButton);
    assert.true(this.onCancel.calledOnce, 'onCancel callback fires on cancel');
  });

  test('it should render fallback for search select', async function (assert) {
    assert.expect(1);

    this.clients = [];
    await this.renderComponent({ name: 'test-key', allowed_client_ids: ['*'] });

    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom(
        '[data-test-component="search-select"]#oidc-key-form-client-select [data-test-component="string-list"]'
      )
      .exists('Radio toggle shows client string-list input');
  });

  test('it should render error alerts when API returns an error', async function (assert) {
    assert.expect(2);

    this.writeStub.rejects(getErrorResponse());

    await this.renderComponent();
    await fillIn('[data-test-input="name"]', 'test-app');
    await click(SELECTORS.keySaveButton);
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom('[data-test-message-error]').exists('alert banner renders');
  });
});
