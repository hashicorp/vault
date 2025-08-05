/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { fillIn, render, typeIn } from '@ember/test-helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import config from 'vault/config/environment';
import { camelize } from '@ember/string';

const SELECTORS = {
  container: '[data-test-swagger-ui]',
  searchInput: 'input.operation-filter-input',
  apiPathBlock: '.opblock-post',
  operationId: '.opblock-summary-operation-id',
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
    // in testing only the input is not filling correctly except after the second time
    await fillIn(SELECTORS.searchInput, 'moot');
    await typeIn(SELECTORS.searchInput, 'token');
    // for some reason search results are not rendered immediately in tests,
    // so asserting that the search input has the value we expect is the best we can do here
    // if the search fn breaks, this test will fail
    assert.dom(SELECTORS.searchInput).hasValue('token', 'search input has value');
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
});
