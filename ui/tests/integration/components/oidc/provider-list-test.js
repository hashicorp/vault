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

module('Integration | Component | oidc/provider-list', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.createRecord('oidc/provider', { name: 'first-provider', issuer: 'foobar' });
    this.store.createRecord('oidc/provider', { name: 'second-provider', issuer: 'foobar' });
    this.model = this.store.peekAll('oidc/provider');
  });

  test('it renders list of providers', async function (assert) {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub(['read', 'update']));
    await render(hbs`<Oidc::ProviderList @model={{this.model}} />`);

    assert.dom('[data-test-oidc-provider-linked-block]').exists({ count: 2 }, 'Two clients are rendered');
    assert.dom('[data-test-oidc-provider-linked-block="first-provider"]').exists('First client is rendered');
    assert
      .dom('[data-test-oidc-provider-linked-block="second-provider"]')
      .exists('Second client is rendered');

    await click('[data-test-oidc-provider-linked-block="first-provider"] [data-test-popup-menu-trigger]');
    assert.dom('[data-test-oidc-provider-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-provider-menu-link="edit"]').exists('Edit link is rendered');
  });

  test('it renders popup menu based on permissions', async function (assert) {
    this.server.post('/sys/capabilities-self', (schema, req) => {
      const { paths } = JSON.parse(req.requestBody);
      if (paths[0] === 'identity/oidc/provider/first-provider') {
        return capabilitiesStub('identity/oidc/provider/first-provider', ['read']);
      } else {
        return capabilitiesStub('identity/oidc/provider/second-provider', ['deny']);
      }
    });
    await render(hbs`<Oidc::ProviderList @model={{this.model}} />`);
    assert.dom('[data-test-popup-menu-trigger]').exists({ count: 1 }, 'Only one popup menu is rendered');
    await click('[data-test-popup-menu-trigger]');
    assert.dom('[data-test-oidc-provider-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-provider-menu-link="edit"]').doesNotExist('Edit link is not rendered');
  });
});
