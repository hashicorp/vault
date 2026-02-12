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
      <Sidebar::Nav::ResilienceAndRecovery />
    </Sidebar::Frame>
  `);
};

module('Integration | Component | sidebar-nav-resilience-and-recovery', function (hooks) {
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

  test('it should hide links user does not have access to other than secrets recovery', async function (assert) {
    await renderComponent();
    stubFeaturesAndPermissions(this.owner);
    assert
      .dom(GENERAL.navLink())
      .exists({ count: 2 }, 'Nav links are hidden other than back and recovery link');
  });

  test('it should render nav headings and links', async function (assert) {
    const links = [
      'Back to main navigation',
      'Secrets recovery',
      'Seal Vault',
      'Overview',
      'Performance replication',
      'Disaster recovery',
    ];
    stubFeaturesAndPermissions(this.owner, true);
    await renderComponent();

    assert.dom(GENERAL.navHeading()).exists({ count: 2 }, 'Correct number of headings render');
    assert
      .dom(GENERAL.navHeading('Resilience and recovery'))
      .hasText('Resilience and recovery', 'Resilience and recovery heading renders');
    assert.dom(GENERAL.navHeading('Replication')).hasText('Replication', 'Replication heading renders');

    assert.dom(GENERAL.navLink()).exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      const name = capitalize(link);
      assert.dom(GENERAL.navLink(name)).hasText(name, `${name} link renders`);
    });
  });

  test('it shows Seal Vault when user is enterprise and in root namespace and has nav permissions', async function (assert) {
    stubFeaturesAndPermissions(this.owner, true, true, [], true);
    await renderComponent();
    assert.dom(GENERAL.navLink('Seal Vault')).exists();
  });

  test('it does NOT show snapshots when user is in HVD admin namespace', async function (assert) {
    const stubs = stubFeaturesAndPermissions(this.owner, true, false, [], false);
    stubs.hasNavPermission.callsFake((route) => route === 'snapshots');

    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    const namespace = this.owner.lookup('service:namespace');
    namespace.setNamespace('admin');

    await renderComponent();

    assert.dom(GENERAL.navLink('Secrets recovery')).doesNotExist();
  });
});
