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
      <Sidebar::Nav::Cluster />
    </Sidebar::Frame>
  `);
};

module('Integration | Component | sidebar-nav-cluster', function (hooks) {
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

  test('it should render nav headings', async function (assert) {
    const headings = ['Vault', 'Monitoring'];
    stubFeaturesAndPermissions(this.owner, true, true);
    await renderComponent();

    assert.dom(GENERAL.navHeading()).exists({ count: headings.length }, 'Correct number of headings render');
    headings.forEach((heading) => {
      assert.dom(GENERAL.navHeading(heading)).hasText(heading, `${heading} heading renders`);
    });
  });

  test('it should hide links and headings user does not have access to', async function (assert) {
    await renderComponent();

    assert
      .dom(GENERAL.navLink())
      .exists({ count: 3 }, 'Nav links are hidden other than secrets, recovery and dashboard');
    assert.dom(GENERAL.navHeading()).exists({ count: 1 }, 'Headings are hidden other than Vault');
  });

  test('it should render nav links', async function (assert) {
    const links = [
      'Dashboard',
      'Secrets Engines',
      'Secrets Sync',
      'Secrets Recovery',
      'Access',
      'Operational tools',
      'Replication',
      'Raft Storage',
      'Client Count',
      'Vault Usage',
      'License',
      'Seal Vault',
    ];
    // do not add PKI-only Secrets feature as it hides Client Count nav link
    const features = allFeatures().filter((feat) => feat !== 'PKI-only Secrets');
    stubFeaturesAndPermissions(this.owner, true, true, features);
    await renderComponent();

    assert.dom(GENERAL.navLink()).exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      assert.dom(GENERAL.navLink(link)).hasText(link, `${link} link renders`);
    });
  });

  test('it should hide enterprise related links in child namespace', async function (assert) {
    const links = [
      'Disaster Recovery',
      'Performance',
      'Replication',
      'Raft Storage',
      'License',
      'Seal Vault',
    ];
    this.owner.lookup('service:namespace').set('path', 'foo');
    const stubs = stubFeaturesAndPermissions(this.owner, true, true);
    stubs.hasNavPermission.callsFake((route) => route !== 'clients');

    await renderComponent();

    assert
      .dom(GENERAL.navHeading('Monitoring'))
      .doesNotExist(
        'Monitoring heading is hidden in child namespace when user does not have access to Client Count'
      );

    links.forEach((link) => {
      assert.dom(GENERAL.navLink(link)).doesNotExist(`${link} is hidden in child namespace`);
    });
  });

  test('it should hide client counts link in chroot namespace', async function (assert) {
    this.owner.lookup('service:permissions').setPaths({
      data: {
        chroot_namespace: 'admin',
        root: true,
      },
    });
    this.owner.lookup('service:currentCluster').setCluster({
      id: 'foo',
      anyReplicationEnabled: true,
      usingRaft: true,
      hasChrootNamespace: true,
    });
    const links = ['Client Counts', 'Replication', 'Raft Storage', 'License', 'Seal Vault'];

    await renderComponent();
    assert
      .dom(GENERAL.navHeading('Monitoring'))
      .doesNotExist('Monitoring heading is hidden in chroot namespace');
    links.forEach((link) => {
      assert.dom(GENERAL.navLink(link)).doesNotExist(`${link} is hidden in chroot namespace`);
    });
  });

  test('it should hide client counts link in PKI-only Secrets clusters', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false);
    await renderComponent();
    assert.dom(GENERAL.navHeading('Client Counts')).doesNotExist('Client count link is hidden.');
  });

  test('it should render badge for promotional links on managed clusters', async function (assert) {
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    const promotionalLinks = ['Secrets Sync'];
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
    assert.dom(GENERAL.navLink('Secrets Sync')).doesNotExist();
  });

  test('ent but feature is not on license: it hides Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, []);
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets Sync')).doesNotExist();
  });

  test('ent (on license), activated and permissions: it shows Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync']);
    this.flags.activatedFlags = ['secrets-sync'];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets Sync')).exists();
  });

  test('ent (on license), activated and no permissions: it hides Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync'], false);
    this.flags.activatedFlags = ['secrets-sync'];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets Sync')).doesNotExist();
  });

  test('ent (on license), not activated and permissions: it shows Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync']);
    this.flags.activatedFlags = [];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets Sync')).exists();
  });

  test('ent (on license), not activated and no permissions: it shows Secrets Sync nav link', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, ['Secrets Sync'], false);
    this.flags.activatedFlags = [];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets Sync')).exists();
  });

  test('hvd managed: it shows Secrets Sync nav link regardless of activation status or permissions', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, false, [], false);
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    this.flags.activatedFlags = [];
    await renderComponent();
    assert.dom(GENERAL.navLink('Secrets Sync')).exists();
  });

  test('it shows Vault Usage when user is enterprise and in root namespace', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true);
    await renderComponent();
    assert.dom(GENERAL.navLink('Vault Usage')).exists();
  });

  test('it does NOT show Vault Usage when user is user is on CE || OSS || community', async function (assert) {
    stubFeaturesAndPermissions(this.owner, false);
    await renderComponent();
    assert.dom(GENERAL.navLink('Vault Usage')).doesNotExist();
  });

  test('it does NOT show Vault Usage when user is enterprise but not in root namespace', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true);

    this.owner.lookup('service:namespace').set('path', 'foo');

    await renderComponent();
    assert.dom(GENERAL.navLink('Vault Usage')).doesNotExist();
  });

  test('it does NOT show Vault Usage when user lacks the necessary permission', async function (assert) {
    // no permissions
    stubFeaturesAndPermissions(this.owner, true, false, [], false);

    await renderComponent();
    assert.dom(GENERAL.navLink('Vault Usage')).doesNotExist();
  });

  test('it does NOT Vault Usage if the user has the necessary permission but user is on CE || OSS || community', async function (assert) {
    // no permissions
    const stubs = stubFeaturesAndPermissions(this.owner, false, false, [], false);

    // allow the route
    stubs.hasNavPermission.callsFake((route) => route === 'monitoring');

    await renderComponent();

    assert.dom(GENERAL.navLink('Vault Usage')).doesNotExist();
  });

  test('it shows Vault Usage when user is in HVD admin namespace', async function (assert) {
    const stubs = stubFeaturesAndPermissions(this.owner, true, false, [], false);
    stubs.hasNavPermission.callsFake((route) => route === 'monitoring');

    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    const namespace = this.owner.lookup('service:namespace');
    namespace.setNamespace('admin');

    await renderComponent();

    assert.dom(GENERAL.navLink('Vault Usage')).exists();
  });

  test('it does NOT show Secrets Recovery when user is in HVD admin namespace', async function (assert) {
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    const namespace = this.owner.lookup('service:namespace');
    namespace.setNamespace('admin');

    await renderComponent();

    assert.dom(GENERAL.navLink('Secrets Recovery')).doesNotExist();
  });
});
