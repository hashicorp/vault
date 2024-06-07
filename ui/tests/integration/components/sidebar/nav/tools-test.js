/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { stubFeaturesAndPermissions } from 'vault/tests/helpers/components/sidebar-nav';
import { toolsActions } from 'vault/helpers/tools-actions';
import { capitalize } from '@ember/string';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const renderComponent = () => {
  return render(hbs`
    <Sidebar::Frame @isVisible={{true}}>
            <Sidebar::Nav::Tools />
    </Sidebar::Frame>
  `);
};

module('Integration | Component | sidebar-nav-tools', function (hooks) {
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

  test('it should hide links user does not have access too', async function (assert) {
    await renderComponent();
    assert
      .dom('[data-test-sidebar-nav-link]')
      .exists({ count: 2 }, 'Nav links are hidden other than back link and api explorer.');
  });

  test('it should render nav headings and links', async function (assert) {
    const links = ['Back to main navigation', ...toolsActions(), 'API Explorer'];
    stubFeaturesAndPermissions(this.owner);
    await renderComponent();

    assert.dom('[data-test-sidebar-nav-heading]').exists({ count: 1 }, 'Correct number of headings render');
    assert.dom('[data-test-sidebar-nav-heading="Tools"]').hasText('Tools', 'Tools heading renders');

    assert
      .dom('[data-test-sidebar-nav-link]')
      .exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      const name = capitalize(link);
      assert.dom(`[data-test-sidebar-nav-link="${name}"]`).hasText(name, `${name} link renders`);
    });
  });
});
