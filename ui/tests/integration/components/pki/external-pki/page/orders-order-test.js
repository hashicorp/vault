/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const SELECTORS = {
  pathByContainer: (idx) => `${GENERAL.cardContainer(idx)} ${GENERAL.inputByAttr('path')}`,
};

module('Integration | Component | pki | external-pki | ExternalPki::Page::OrdersOrder', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');

    this.model = {
      engine: { id: 'my-pki-external-ca' },
      order_id: '019faf78-0d52-71e8-abfa-71bc5de7cc9a',
      order: { details: { order_status: 'completed', role_name: 'my-role' } },
      certificate: undefined,
      responseTimestamp: new Date('2026-07-14T21:00:00Z'),
    };

    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::Page::OrdersOrder @model={{this.model}} @breadcrumbs={{array (hash label="View order")}} />`,
        { owner: this.engine }
      );
  });

  test('it renders the last refreshed timestamp', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.textBody('Last refreshed')).hasTextContaining('Last refreshed July 14, 2026');
  });

  test('it calls router.refresh with the order route when Refresh button is clicked', async function (assert) {
    const refreshStub = sinon.stub(this.router, 'refresh');
    await this.renderComponent();
    await click(GENERAL.button('Refresh'));
    assert.true(refreshStub.calledOnce, 'refresh was called once');
    const [route] = refreshStub.lastCall.args;
    assert.strictEqual(
      route,
      'vault.cluster.secrets.backend.pki.external.orders.order',
      'refresh was called with the order route'
    );
  });

  module('policy flyout pre-population', function (hooks) {
    hooks.beforeEach(function () {
      // The Generate policy button only renders for enterprise
      this.owner.lookup('service:version').type = 'enterprise';
      this.currentRouteNameStub = sinon.stub(this.router, 'currentRouteName');
    });

    test('it pre-populates flyout stanzas when the role name is unavailable', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.orders.order');
      this.model.order = undefined;
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/lookup/order/*');
      assert
        .dom(SELECTORS.pathByContainer(1))
        .hasValue('my-pki-external-ca/role/:role_name/order/+/fetch-cert');
      assert.dom(GENERAL.cardContainer()).exists({ count: 2 });
    });

    test('it pre-populates flyout stanzas when the role name is provided', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.orders.order');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/lookup/order/*');
      assert.dom(SELECTORS.pathByContainer(1)).hasValue('my-pki-external-ca/role/my-role/order/+/fetch-cert');
      assert.dom(GENERAL.cardContainer()).exists({ count: 2 });
    });
  });
});
