/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | oidc/provider-list', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = (isLimited) => {
      const providers = [
        { name: 'first-provider', issuer: 'foobar' },
        { name: 'second-provider', issuer: 'foobar' },
      ];

      const capabilities = providers.reduce((capabilities, provider) => {
        const path = this.owner.lookup('service:capabilities').pathFor('oidcProvider', provider);
        const isFirstProvider = provider.name === 'first-provider';
        const canRead = isLimited ? isFirstProvider : true;
        const canUpdate = !isLimited;
        capabilities[path] = { canRead, canUpdate };
        return capabilities;
      }, {});

      this.model = {
        providers,
        capabilities,
      };
      return render(hbs`<Oidc::ProviderList @model={{this.model}} />`);
    };
  });

  test('it renders list of providers', async function (assert) {
    await this.renderComponent(false);

    assert.dom('[data-test-oidc-provider-linked-block]').exists({ count: 2 }, 'Two providers are rendered');
    assert
      .dom('[data-test-oidc-provider-linked-block="first-provider"]')
      .exists('First provider is rendered');
    assert
      .dom('[data-test-oidc-provider-linked-block="second-provider"]')
      .exists('Second provider is rendered');

    await click('[data-test-oidc-provider-linked-block="first-provider"] [data-test-popup-menu-trigger]');
    assert.dom('[data-test-oidc-provider-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-provider-menu-link="edit"]').exists('Edit link is rendered');
  });

  test('it renders popup menu based on permissions', async function (assert) {
    await this.renderComponent(true);
    assert.dom('[data-test-popup-menu-trigger]').exists({ count: 1 }, 'Only one popup menu is rendered');
    await click(GENERAL.menuTrigger);
    assert.dom('[data-test-oidc-provider-menu-link="details"]').exists('Details link is rendered');
    assert.dom('[data-test-oidc-provider-menu-link="edit"]').doesNotExist('Edit link is not rendered');
  });
});
