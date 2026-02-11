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
      .exists(
        { count: 3 },
        'Nav links are hidden other than secrets, recovery (nested in Resilience and recovery nav link) and dashboard'
      );
    assert.dom(GENERAL.navHeading()).exists({ count: 1 }, 'Headings are hidden other than Vault');
  });

  test('it should render nav links on community version', async function (assert) {
    const links = [
      'Dashboard',
      'Secrets',
      'Resilience and recovery',
      'Access control',
      'Operational tools',
      'Raft storage',
      'Client count',
    ];

    const features = allFeatures().filter((feat) => feat !== 'PKI-only Secrets');
    stubFeaturesAndPermissions(this.owner, false, true, features);
    await renderComponent();

    assert.dom(GENERAL.navLink()).exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      assert.dom(GENERAL.navLink(link)).hasText(link, `${link} link renders`);
    });
  });

  test('it should render nav links', async function (assert) {
    const links = [
      'Dashboard',
      'Secrets',
      'Access control',
      'Operational tools',
      'Resilience and recovery',
      'Reporting',
      'Raft storage',
      'Client count',
    ];
    // do not add PKI-only Secrets feature as it hides Client count nav link
    const features = allFeatures().filter((feat) => feat !== 'PKI-only Secrets');
    stubFeaturesAndPermissions(this.owner, true, true, features);
    await renderComponent();

    assert.dom(GENERAL.navLink()).exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      assert.dom(GENERAL.navLink(link)).hasText(link, `${link} link renders`);
    });
  });

  test('it should hide enterprise related links in child namespace', async function (assert) {
    const links = ['Raft storage', 'License'];
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
    const links = ['Client Counts', 'Replication', 'Raft storage', 'License', 'Seal Vault'];

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

  test('it does NOT show Secrets Recovery when user is in HVD admin namespace', async function (assert) {
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    const namespace = this.owner.lookup('service:namespace');
    namespace.setNamespace('admin');

    await renderComponent();

    assert.dom(GENERAL.navLink('Secrets Recovery')).doesNotExist();
  });
});
