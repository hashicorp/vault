/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { stubFeaturesAndPermissions } from 'vault/tests/helpers/components/sidebar-nav';
import { setRunOptions } from 'ember-a11y-testing/test-support';

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
    setRunOptions({
      rules: {
        // This is an issue with Hds::SideNav::Header::HomeLink
        'aria-prohibited-attr': { enabled: false },
        // TODO: fix use Dropdown on user-menu
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should render nav headings', async function (assert) {
    const headings = ['Vault', 'Monitoring', 'Settings'];
    stubFeaturesAndPermissions(this.owner, true, true);
    await renderComponent();

    assert
      .dom('[data-test-sidebar-nav-heading]')
      .exists({ count: headings.length }, 'Correct number of headings render');
    headings.forEach((heading) => {
      assert
        .dom(`[data-test-sidebar-nav-heading="${heading}"]`)
        .hasText(heading, `${heading} heading renders`);
    });
  });

  test('it should hide links and headings user does not have access too', async function (assert) {
    await renderComponent();

    assert
      .dom('[data-test-sidebar-nav-link]')
      .exists({ count: 2 }, 'Nav links are hidden other than secrets and dashboard');
    assert
      .dom('[data-test-sidebar-nav-heading]')
      .exists({ count: 1 }, 'Headings are hidden other than Vault');
  });

  test('it should render nav links', async function (assert) {
    const links = [
      'Dashboard',
      'Secrets Engines',
      'Secrets Sync',
      'Access',
      'Policies',
      'Tools',
      'Replication',
      'Raft Storage',
      'Client Count',
      'License',
      'Seal Vault',
      'Custom Messages',
    ];
    stubFeaturesAndPermissions(this.owner, true, true);
    await renderComponent();

    assert
      .dom('[data-test-sidebar-nav-link]')
      .exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      assert.dom(`[data-test-sidebar-nav-link="${link}"]`).hasText(link, `${link} link renders`);
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
      .dom('[data-test-sidebar-nav-heading="Monitoring"]')
      .doesNotExist(
        'Monitoring heading is hidden in child namespace when user does not have access to Client Count'
      );

    links.forEach((link) => {
      assert
        .dom(`[data-test-sidebar-nav-link="${link}"]`)
        .doesNotExist(`${link} is hidden in child namespace`);
    });
  });

  test('it should render badge for promotional links on managed clusters', async function (assert) {
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    const promotionalLinks = ['Secrets Sync'];
    stubFeaturesAndPermissions(this.owner, true, true);
    await renderComponent();

    promotionalLinks.forEach((link) => {
      assert
        .dom(`[data-test-sidebar-nav-link="${link}"]`)
        .hasText(`${link} Plus`, `${link} link renders Plus badge`);
    });
  });
});
