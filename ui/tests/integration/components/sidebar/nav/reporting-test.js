/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { stubFeaturesAndPermissions } from 'vault/tests/helpers/components/sidebar-nav';
import { capitalize } from '@ember/string';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const renderComponent = () => {
  return render(hbs`
    <Sidebar::Frame @isVisible={{true}}>
            <Sidebar::Nav::Reporting />
    </Sidebar::Frame>
  `);
};

module('Integration | Component | sidebar-nav-reporting', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.flags = this.owner.lookup('service:flags');

    setRunOptions({
      rules: {
        // This is an issue with Hds::AppHeader::HomeLink
        'aria-prohibited-attr': { enabled: false },
        // TODO: fix use Dropdown on user-menu
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should hide links user does not have access to', async function (assert) {
    await renderComponent();
    stubFeaturesAndPermissions(this.owner);
    assert
      .dom(GENERAL.navLink())
      .exists({ count: 1 }, 'Nav links are hidden other than back link and license');
  });

  test('it should render nav headings and links', async function (assert) {
    const links = ['Back to main navigation', 'Vault usage', 'License'];
    stubFeaturesAndPermissions(this.owner, true);
    await renderComponent();

    assert.dom(GENERAL.navHeading()).exists({ count: 1 }, 'Correct number of headings render');
    assert.dom(GENERAL.navHeading('Reporting')).hasText('Reporting', 'Reporting heading renders');

    assert.dom(GENERAL.navLink()).exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      const name = capitalize(link);
      assert.dom(GENERAL.navLink(name)).hasText(name, `${name} link renders`);
    });
  });

  test('it shows Vault Usage when user is enterprise and in root namespace', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true);
    await renderComponent();
    assert.dom(GENERAL.navLink('Vault usage')).exists();
  });

  test('it does NOT show Vault Usage when user is user is on CE || OSS || community', async function (assert) {
    stubFeaturesAndPermissions(this.owner, false);
    await renderComponent();
    assert.dom(GENERAL.navLink('Vault usage')).doesNotExist();
  });

  test('it does NOT show Vault Usage when user is enterprise but not in root namespace', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true);

    this.owner.lookup('service:namespace').set('path', 'foo');

    await renderComponent();
    assert.dom(GENERAL.navLink('Vault usage')).doesNotExist();
  });

  test('it does NOT show Vault Usage when user lacks the necessary permission', async function (assert) {
    // no permissions
    stubFeaturesAndPermissions(this.owner, true, false, [], false);

    await renderComponent();
    assert.dom(GENERAL.navLink('Vault usage')).doesNotExist();
  });

  test('it does NOT Vault Usage if the user has the necessary permission but user is on CE || OSS || community', async function (assert) {
    // no permissions
    const stubs = stubFeaturesAndPermissions(this.owner, false, false, [], false);

    // allow the route
    stubs.hasNavPermission.callsFake((route) => route === 'monitoring');

    await renderComponent();

    assert.dom(GENERAL.navLink('Vault usage')).doesNotExist();
  });

  test('it shows Vault Usage when user is in HVD admin namespace', async function (assert) {
    const stubs = stubFeaturesAndPermissions(this.owner, true, false, [], false);
    stubs.hasNavPermission.callsFake((route) => route === 'monitoring');

    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    const namespace = this.owner.lookup('service:namespace');
    namespace.setNamespace('admin');

    await renderComponent();

    assert.dom(GENERAL.navLink('Vault usage')).exists();
  });
});
