/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { EXTERNAL_TABS } from 'vault/tests/helpers/pki/assertion-helpers';
import sinon from 'sinon';

const SELECTORS = {
  pathByContainer: (idx) => `${GENERAL.cardContainer(idx)} ${GENERAL.inputByAttr('path')}`,
};

module('Integration | Component | pki | external-pki | ExternalPki::HeaderTabs', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.backend = new SecretsEngineResource({
      accessor: 'pki-external-ca_e158c567',
      type: 'pki-external-ca',
      path: 'my-pki-external-ca/',
    });

    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::HeaderTabs
        @backend={{this.backend}}
        @roleName={{this.roleName}}
        @showConfigSnippets={{this.showConfigSnippets}} 
        @hideTabs={{this.hideTabs}}
        @breadcrumbs={{array (hash label="Crumb!")}} 
        />`,
        { owner: this.engine }
      );
  });

  test('it renders ManageDropdown when @backend is provided', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.dropdownToggle('Manage')).exists('ManageDropdown renders when @backend is present');
  });

  test('it does not render ManageDropdown without @backend', async function (assert) {
    this.backend = undefined;
    await this.renderComponent();

    assert
      .dom(GENERAL.dropdownToggle('Manage'))
      .doesNotExist('ManageDropdown is not rendered when @backend is absent');
  });

  test('it renders all default tabs when @backend is set', async function (assert) {
    await this.renderComponent();

    EXTERNAL_TABS.forEach((tab) => {
      assert.dom(GENERAL.linkTo(tab)).exists(`"${tab}" tab renders`);
    });
  });

  test('it renders only Overview tab when showConfigSnippets is true', async function (assert) {
    this.showConfigSnippets = true;
    await this.renderComponent();

    assert.dom(GENERAL.linkTo('Overview')).exists('"Overview" tab renders');

    EXTERNAL_TABS.filter((t) => t !== 'Overview').forEach((tab) => {
      assert.dom(GENERAL.linkTo(tab)).doesNotExist(`"${tab}" tab is hidden when showConfigSnippets is true`);
    });
  });

  test('it hides all tabs when @hideTabs is true', async function (assert) {
    this.hideTabs = true;
    await this.renderComponent();

    EXTERNAL_TABS.forEach((tab) => {
      assert.dom(GENERAL.linkTo(tab)).doesNotExist(`"${tab}" tab is hidden when @hideTabs is true`);
    });
  });

  test('it renders custom tabs block when provided', async function (assert) {
    await render(
      hbs`
          <ExternalPki::HeaderTabs @backend={{this.backend}}  @breadcrumbs={{array (hash label="Crumb!")}} >
            <:tabs>
              <li><a data-test-link-to="Custom tab">Custom tab</a></li>
            </:tabs>
          </ExternalPki::HeaderTabs>
        `,
      { owner: this.engine }
    );

    assert.dom(GENERAL.linkTo('Custom tab')).exists('custom tabs block renders');
    EXTERNAL_TABS.forEach((tab) => {
      assert
        .dom(GENERAL.linkTo(tab))
        .doesNotExist(`default "${tab}" tab is replaced by the custom tabs block`);
    });
  });

  module('policy flyout pre-population', function (hooks) {
    hooks.beforeEach(function () {
      this.roleName = 'my-role';
      // The Generate policy button only renders for enterprise
      this.owner.lookup('service:version').type = 'enterprise';

      this.router = this.owner.lookup('service:router');
      this.currentRouteNameStub = sinon.stub(this.router, 'currentRouteName');
    });

    test('it pre-populates flyout stanzas for the overview route', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.overview');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/config/acme-account');
      assert.dom(SELECTORS.pathByContainer(1)).hasValue('my-pki-external-ca/config/dns');
      assert.dom(SELECTORS.pathByContainer(2)).hasValue('my-pki-external-ca/role');
      assert.dom(SELECTORS.pathByContainer(3)).hasValue('my-pki-external-ca/lookup/orders');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 4 }, 'four stanzas pre-populated for overview route');
    });

    test('it pre-populates flyout stanzas for the recent orders route', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.orders.index');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/lookup/orders/recent');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 1 }, 'one stanza pre-populated for recent orders route');
    });

    test('it pre-populates flyout stanzas for the orders order sub-route with roleName', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.orders.order');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/lookup/order/*');
      assert.dom(SELECTORS.pathByContainer(1)).hasValue('my-pki-external-ca/role/my-role/order/+/fetch-cert');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 2 }, 'two stanzas pre-populated for orders order route');
    });

    test('it pre-populates flyout stanzas for the certificate sub-route with roleName', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.certificates.certificate');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/lookup/cert/*');
      assert.dom(SELECTORS.pathByContainer(1)).hasValue('my-pki-external-ca/role/my-role/order/+/fetch-cert');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 2 }, 'two stanzas pre-populated for certificate route');
    });

    test('it pre-populates flyout stanzas for the role overview sub-route with roleName', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.roles.role.overview');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/role/my-role');
      assert.dom(SELECTORS.pathByContainer(1)).hasValue('my-pki-external-ca/role/my-role/active-orders');
      assert.dom(SELECTORS.pathByContainer(2)).hasValue('my-pki-external-ca/role/my-role/cached');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 3 }, 'three stanzas pre-populated for role overview route');
    });

    test('it pre-populates flyout stanzas for the role order sub-route with roleName', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.roles.role.order');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/role/my-role/order/+/status');
      assert.dom(SELECTORS.pathByContainer(1)).hasValue('my-pki-external-ca/role/my-role/order/+/fetch-cert');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 2 }, 'two stanzas pre-populated for role order route');
    });

    test('it pre-populates flyout stanzas for the role active-orders sub-route with roleName', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.roles.role.active-orders');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/role/my-role/active-orders');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 1 }, 'one stanza pre-populated for role active-orders route');
    });

    test('it uses fallback when role name does not exist', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.roles.role.active-orders');
      this.roleName = undefined;
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/role/:role_name/active-orders');
    });

    test('flyout opens with blank stanza when current route is not in the map', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.unknown');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('', 'no path pre-populated for unmapped route');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 1 }, 'single blank stanza shown for unmapped route');
    });

    test('flyout opens with blank stanza when current route does not have matching prefix', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.auth.backend.pki.external.overview');
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('', 'no path pre-populated for unmapped route');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 1 }, 'single blank stanza shown for unmapped route');
    });

    test('it pre-populates flyout stanzas for the error route via URL fallback', async function (assert) {
      this.currentRouteNameStub.value('vault.cluster.secrets.backend.pki.external.error');
      // matchesCurrentUrl uses router.recognize() + router.currentURL to find the active route
      sinon
        .stub(this.router, 'currentURL')
        .value('/ui/vault/secrets/my-pki-external-ca/pki/external/dns-providers');
      sinon.stub(this.router, 'rootURL').value('/ui/');
      sinon.stub(this.router, 'recognize').callsFake(() => ({
        name: 'vault.cluster.secrets.backend.pki.external.dns-providers',
      }));
      await this.renderComponent();
      await click(GENERAL.button('Generate policy'));

      assert.dom(SELECTORS.pathByContainer(0)).hasValue('my-pki-external-ca/config/dns');
      assert
        .dom(GENERAL.cardContainer())
        .exists({ count: 1 }, 'one stanza pre-populated for error route via URL fallback');
    });
  });
});
