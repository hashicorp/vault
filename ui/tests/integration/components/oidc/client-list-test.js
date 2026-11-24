/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | oidc/client-list', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = (isLimited) => {
      const clients = [{ name: 'first-client' }, { name: 'second-client' }];

      const capabilities = clients.reduce((capabilities, client) => {
        const path = this.owner.lookup('service:capabilities').pathFor('oidcClient', client);
        const isFirstClient = client.name === 'first-client';
        const canRead = isLimited ? isFirstClient : true;
        const canUpdate = !isLimited;
        capabilities[path] = { canRead, canUpdate };
        return capabilities;
      }, {});

      this.model = {
        clients,
        capabilities,
      };
      return render(hbs`<Oidc::ClientList @model={{this.model}} />`);
    };
  });

  test('it renders list of clients', async function (assert) {
    await this.renderComponent(false);

    assert.dom('[data-test-oidc-client-linked-block]').exists({ count: 2 }, 'Two clients are rendered');
    assert.dom('[data-test-oidc-client-linked-block="first-client"]').exists('First client is rendered');
    assert.dom('[data-test-oidc-client-linked-block="second-client"]').exists('Second client is rendered');

    await click('[data-test-oidc-client-linked-block="first-client"] [data-test-popup-menu-trigger]');
    assert.dom('[data-test-oidc-client-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-client-menu-link="edit"]').exists('Edit link is rendered');
  });

  test('it renders popup menu based on permissions', async function (assert) {
    await this.renderComponent(true);

    assert.dom('[data-test-popup-menu-trigger]').exists({ count: 1 }, 'Only one popup menu is rendered');
    await click(GENERAL.menuTrigger);
    assert.dom('[data-test-oidc-client-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-client-menu-link="edit"]').doesNotExist('Edit link is not rendered');
  });
});
