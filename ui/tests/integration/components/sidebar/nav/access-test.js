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
      <Sidebar::Nav::Access />
    </Sidebar::Frame>
  `);
};

module('Integration | Component | sidebar-nav-access', function (hooks) {
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
    const headings = ['Authentication', 'Access Control', 'Organization', 'Administration'];

    stubFeaturesAndPermissions(this.owner);
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
      .exists({ count: 1 }, 'Nav links are hidden other than back link');
    assert
      .dom('[data-test-sidebar-nav-heading]')
      .doesNotExist('Headings are hidden when user does not have access to links');
  });

  test('it should render nav links', async function (assert) {
    const links = [
      'Back to main navigation',
      'Authentication Methods',
      'Multi-Factor Authentication',
      'OIDC Provider',
      'Control Groups',
      'Namespaces',
      'Groups',
      'Entities',
      'Leases',
    ];
    stubFeaturesAndPermissions(this.owner);
    await renderComponent();

    assert
      .dom('[data-test-sidebar-nav-link]')
      .exists({ count: links.length }, 'Correct number of links render');
    links.forEach((link) => {
      assert.dom(`[data-test-sidebar-nav-link="${link}"]`).hasText(link, `${link} link renders`);
    });
  });
});
