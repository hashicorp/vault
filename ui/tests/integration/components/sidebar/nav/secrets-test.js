/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { stubFeaturesAndPermissions } from 'vault/tests/helpers/components/sidebar-nav';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { allFeatures } from 'core/utils/all-features';

const renderComponent = () => {
  return render(hbs`
    <Sidebar::Frame @isVisible={{true}}>
      <Sidebar::Nav::Secrets />
    </Sidebar::Frame>
  `);
};

module('Integration | Component | sidebar-nav-secrets', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.flags = this.owner.lookup('service:flags');
    setRunOptions({
      rules: {
        // This is an issue with Hds::AppHeader::HomeLink
        'aria-prohibited-attr': { enabled: false },
      },
    });
  });

  test('it should hide links and headings user does not have access to', async function (assert) {
    await renderComponent();
    assert.dom(GENERAL.navLink()).exists({ count: 2 }, 'Nav links are hidden other than secrets engines');
    assert.dom(GENERAL.navHeading()).exists({ count: 1 }, 'Headings are hidden other than Secrets engines');
  });

  test('it should render nav links', async function (assert) {
    const links = ['Secrets engines', 'Secrets sync'];
    // do not add PKI-only Secrets feature as it hides Client count nav link
    const features = allFeatures().filter((feat) => feat !== 'PKI-only Secrets');
    stubFeaturesAndPermissions(this.owner, true, true, features);
    await renderComponent();

    assert.dom(GENERAL.navLink()).exists({ count: links.length + 1 }, 'Correct number of links render');
    links.forEach((link) => {
      assert.dom(GENERAL.navLink(link)).hasText(link, `${link} link renders`);
    });
  });

  test('it should render badge for promotional links on managed clusters', async function (assert) {
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    const promotionalLinks = ['Secrets sync'];
    stubFeaturesAndPermissions(this.owner, true, true);
    await renderComponent();

    promotionalLinks.forEach((link) => {
      assert.dom(GENERAL.navLink(link)).hasText(`${link} Plus`, `${link} link renders Plus badge`);
    });
  });

  // Secrets Sync side nav link has multiple combinations of three variables to test:
  // 1. cluster type: enterprise (on and off license), HVD managed or community
  // 2. activation status: activated or not
  // 3. permissions: policy access to sys/sync routes or not

  test('community: it hides Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, false, false);
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets sync')).doesNotExist();
  });

  test('ent but feature is not on license: it hides Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, []);
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets sync')).doesNotExist();
  });

  test('ent (on license), activated and permissions: it shows Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync']);
    this.flags.activatedFlags = ['secrets-sync'];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets sync')).exists();
  });

  test('ent (on license), activated and no permissions: it hides Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync'], false);
    this.flags.activatedFlags = ['secrets-sync'];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets sync')).doesNotExist();
  });

  test('ent (on license), not activated and permissions: it shows Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync']);
    this.flags.activatedFlags = [];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets sync')).exists();
  });

  test('ent (on license), not activated and no permissions: it shows Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync'], false);
    this.flags.activatedFlags = [];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets sync')).exists();
  });

  test('hvd managed: it shows Secrets Sync nav link regardless of activation status or permissions', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, [], false);
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    this.flags.activatedFlags = [];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets sync')).exists();
  });
});
