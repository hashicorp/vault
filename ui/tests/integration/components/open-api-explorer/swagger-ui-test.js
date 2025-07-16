/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render, waitFor, waitUntil } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import config from 'vault/config/environment';
import { camelize } from '@ember/string';

const SELECTORS = {
  container: '[data-test-swagger-ui]',
  searchInput: 'input.operation-filter-input',
  apiPathBlock: '.opblock',
  operationId: '.opblock-summary-operation-id',
  controlArrowButton: '.opblock-control-arrow',
  copyButton: '.copy-to-clipboard',
  tryItOutButton: '.try-out button',
};

// for some reason search filtering does not update with ember test helpers
// possibly due to swagger-ui event implementation
// using native window fn to workaround
const setNativeInputValue = (value) => {
  const input = document.querySelector(SELECTORS.searchInput);
  const nativeInputValueSetter = Object.getOwnPropertyDescriptor(
    window.HTMLInputElement.prototype,
    'value'
  ).set;
  nativeInputValueSetter.call(input, value);
  input.dispatchEvent(new Event('input', { bubbles: true }));
};

module('Integration | Component | open-api-explorer | swagger-ui', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'open-api-explorer');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');

    this.openApiResponse = this.server.create('open-api-explorer');
    this.server.get('sys/internal/specs/openapi', () => {
      return this.openApiResponse;
    });

    this.totalApiPaths = Object.keys(this.openApiResponse.paths).length;

    this.renderComponent = () => render(hbs`<SwaggerUi/>`, { owner: this.engine });
  });

  test('it renders', async function (assert) {
    await this.renderComponent();

    assert.dom(SELECTORS.container).exists('renders component');
    assert.dom(SELECTORS.apiPathBlock).exists({ count: this.totalApiPaths }, 'renders all api paths');
  });

  test('it can search', async function (assert) {
    await this.renderComponent();

    setNativeInputValue('token');
    assert.dom(SELECTORS.searchInput).hasValue('token', 'search input has value');
    assert.dom(SELECTORS.apiPathBlock).exists({ count: 2 }, 'renders filtered api paths');
  });

  test('it should render camelized operation ids', async function (assert) {
    const envStub = sinon.stub(config, 'environment').value('development');

    await this.renderComponent();

    const id = this.openApiResponse.paths['/auth/token/create'].post.operationId;
    assert.dom(SELECTORS.operationId).hasText(camelize(id), 'renders camelized operation id');

    envStub.restore();
  });

  test('it should not render operation ids in production', async function (assert) {
    const envStub = sinon.stub(config, 'environment').value('production');

    await this.renderComponent();

    assert.dom(SELECTORS.operationId).doesNotExist('operation ids are hidden in production environment');

    envStub.restore();
  });

  test('it contains a11y fixes', async function (assert) {
    const envStub = sinon.stub(config, 'environment').value('development');

    await this.renderComponent();

    await waitUntil(() => {
      return document.querySelector(SELECTORS.controlArrowButton).getAttribute('tabindex') === '0';
    });
    assert.dom(SELECTORS.controlArrowButton).hasAttribute('tabindex', '0');

    await waitUntil(() => {
      return document.querySelector(SELECTORS.copyButton).getAttribute('tabindex') === '0';
    });
    assert.dom(SELECTORS.copyButton).hasAttribute('tabindex', '0');

    const controlArrowButton = document.querySelectorAll(SELECTORS.controlArrowButton)[1];
    await click(controlArrowButton);
    await waitFor(SELECTORS.tryItOutButton);

    const input = document.querySelector('.parameters input:read-only');
    assert.dom(input).exists('parameter input is readonly');

    assert
      .dom(SELECTORS.tryItOutButton)
      .hasAttribute(
        'aria-description',
        'Caution: This will make requests to the Vault server on your behalf which may create or delete items.'
      );

    envStub.restore();
  });

  test('it retains a11y fixes after filtering', async function (assert) {
    const envStub = sinon.stub(config, 'environment').value('development');

    await this.renderComponent();

    setNativeInputValue('token');

    await waitUntil(() => {
      return document.querySelector(SELECTORS.controlArrowButton).getAttribute('tabindex') === '0';
    });
    assert.dom(SELECTORS.controlArrowButton).hasAttribute('tabindex', '0');

    await waitUntil(() => {
      return document.querySelector(SELECTORS.copyButton).getAttribute('tabindex') === '0';
    });
    assert.dom(SELECTORS.copyButton).hasAttribute('tabindex', '0');

    const controlArrowButton = document.querySelectorAll(SELECTORS.controlArrowButton)[1];
    await click(controlArrowButton);
    await waitFor(SELECTORS.tryItOutButton);

    const input = document.querySelector('.parameters input:read-only');
    assert.dom(input).exists('parameter input is readonly');
    assert
      .dom(SELECTORS.tryItOutButton)
      .hasAttribute(
        'aria-description',
        'Caution: This will make requests to the Vault server on your behalf which may create or delete items.'
      );

    envStub.restore();
  });
});
