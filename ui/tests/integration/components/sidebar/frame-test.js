/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | sidebar-frame', function (hooks) {
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

  test('it should hide and show sidebar', async function (assert) {
    this.set('showSidebar', true);
    await render(hbs`
      <Sidebar::Frame @showSidebar={{this.showSidebar}} />
    `);
    assert.dom('[data-test-sidebar-nav]').exists('Sidebar renders');

    this.set('showSidebar', false);
    assert.dom('[data-test-sidebar-nav]').doesNotExist('Sidebar is hidden');
  });

  test('it should render link status, console ui panel container and yield block for app content', async function (assert) {
    const currentCluster = this.owner.lookup('service:currentCluster');
    currentCluster.setCluster({ hcpLinkStatus: 'connected' });
    const version = this.owner.lookup('service:version');
    version.type = 'enterprise';

    await render(hbs`
      <Sidebar::Frame @showSidebar={{true}}>
        <div class="page-container">
          App goes here!
        </div>
      </Sidebar::Frame>
    `);

    assert.dom('[data-test-link-status]').exists('Link status component renders');
    assert.dom('[data-test-console-panel]').exists('Console UI panel container renders');
    assert.dom('.page-container').exists('Block yields for app content');
  });

  test('it should render logo and actions in sidebar header', async function (assert) {
    setRunOptions({
      rules: {
        'aria-prohibited-attr': { enabled: false },
        'nested-interactive': { enabled: false },
        label: { enabled: false },
      },
    });
    this.owner.lookup('service:currentCluster').setCluster({ name: 'vault' });

    await render(hbs`
      <Sidebar::Frame @showSidebar={{true}} />
    `);

    assert.dom('[data-test-sidebar-logo]').exists('Vault logo renders in sidebar header');
    assert.dom('[data-test-console-toggle]').exists('Console toggle button renders in sidebar header');
    await click('[data-test-console-toggle]');
    assert.dom('.panel-open').exists('Console ui panel opens');

    await click('[data-test-console-toggle]');
    assert.dom('.panel-open').doesNotExist('Console ui panel closes');
    assert.dom('[data-test-user-menu]').exists('User menu renders');
  });

  test('it should render namespace picker in sidebar footer', async function (assert) {
    const version = this.owner.lookup('service:version');
    version.features = ['Namespaces'];
    const auth = this.owner.lookup('service:auth');
    sinon.stub(auth, 'authData').value({});

    await render(hbs`
      <Sidebar::Frame @showSidebar={{true}} />
    `);

    assert.dom('.namespace-picker').exists('Namespace picker renders in sidebar footer');
  });
});
