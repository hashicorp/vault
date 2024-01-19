/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const selectors = {
  versionDisplay: '[data-test-footer-version]',
  upgradeLink: '[data-test-footer-upgrade-link]',
  docsLink: '[data-test-footer-documentation-link]',
};

module('Integration | Component | app-footer', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.versionSvc = this.owner.lookup('service:version');
  });

  test('it renders a sane default', async function (assert) {
    await render(hbs`<AppFooter />`);
    assert.dom(selectors.versionDisplay).hasText('Vault', 'Vault without version by default');
    assert.dom(selectors.upgradeLink).hasText('Upgrade to Vault Enterprise', 'upgrade link shows');
    assert.dom(selectors.docsLink).hasText('Documentation', 'displays docs link');
  });

  test('it renders for community version', async function (assert) {
    this.versionSvc.version = '1.15.1';
    this.versionSvc.type = 'community';
    await render(hbs`<AppFooter />`);
    assert.dom(selectors.versionDisplay).hasText('Vault 1.15.1', 'Vault shows version when available');
    assert.dom(selectors.upgradeLink).hasText('Upgrade to Vault Enterprise', 'upgrade link shows');
    assert.dom(selectors.docsLink).hasText('Documentation', 'displays docs link');
  });
  test('it renders for ent version', async function (assert) {
    this.versionSvc.version = '1.15.1+hsm';
    this.versionSvc.type = 'enterprise';
    await render(hbs`<AppFooter />`);
    assert.dom(selectors.versionDisplay).hasText('Vault 1.15.1+hsm', 'shows version when available');
    assert.dom(selectors.upgradeLink).doesNotExist('upgrade link not shown');
    assert.dom(selectors.docsLink).hasText('Documentation', 'displays docs link');
  });
});
