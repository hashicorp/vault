/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { allowAllCapabilitiesStub, capabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | oidc/client-list', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.createRecord('oidc/client', { name: 'first-client' });
    this.store.createRecord('oidc/client', { name: 'second-client' });
    this.model = this.store.peekAll('oidc/client');
  });

  test('it renders list of clients', async function (assert) {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub(['read', 'update']));
    await render(hbs`<Oidc::ClientList @model={{this.model}} />`);

    assert.dom('[data-test-oidc-client-linked-block]').exists({ count: 2 }, 'Two clients are rendered');
    assert.dom('[data-test-oidc-client-linked-block="first-client"]').exists('First client is rendered');
    assert.dom('[data-test-oidc-client-linked-block="second-client"]').exists('Second client is rendered');

    await click('[data-test-oidc-client-linked-block="first-client"] [data-test-popup-menu-trigger]');
    assert.dom('[data-test-oidc-client-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-client-menu-link="edit"]').exists('Edit link is rendered');
  });

  test('it renders popup menu based on permissions', async function (assert) {
    this.server.post('/sys/capabilities-self', (schema, req) => {
      const { paths } = JSON.parse(req.requestBody);
      if (paths[0] === 'identity/oidc/client/first-client') {
        return capabilitiesStub('identity/oidc/client/first-client', ['read']);
      } else {
        return capabilitiesStub('identity/oidc/client/second-client', ['deny']);
      }
    });
    await render(hbs`<Oidc::ClientList @model={{this.model}} />`);

    assert.dom('[data-test-popup-menu-trigger]').exists({ count: 1 }, 'Only one popup menu is rendered');
    await click('[data-test-popup-menu-trigger]');
    assert.dom('[data-test-oidc-client-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-client-menu-link="edit"]').doesNotExist('Edit link is not rendered');
  });
});
