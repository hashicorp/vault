/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | usage-reporting/base/title-row', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the title', async function (assert) {
    await render(hbs`<UsageReporting::Base::TitleRow @title="Secret engines" />`);
    assert
      .dom('[data-test-vault-reporting-dashboard-card-title]')
      .hasText('Secret engines', 'renders title text');
  });

  test('it renders description when provided', async function (assert) {
    await render(hbs`
      <UsageReporting::Base::TitleRow
        @title="Secret engines"
        @description="Enabled secret engines for this cluster."
      />
    `);
    assert
      .dom('[data-test-vault-reporting-dashboard-card-description]')
      .hasText('Enabled secret engines for this cluster.', 'renders description');
  });

  test('it does not render description when not provided', async function (assert) {
    await render(hbs`<UsageReporting::Base::TitleRow @title="Secret engines" />`);
    assert.dom('[data-test-vault-reporting-dashboard-card-description]').doesNotExist();
  });

  test('it renders external link when @linkUrl is provided', async function (assert) {
    await render(hbs`
      <UsageReporting::Base::TitleRow
        @title="Global lease count quota"
        @linkUrl="https://developer.hashicorp.com/vault/docs"
        @linkText="Documentation"
        @linkIcon="docs-link"
        @linkTarget="_blank"
      />
    `);
    assert
      .dom('[data-test-vault-reporting-dashboard-card-title-link]')
      .exists('renders the link element')
      .hasText('Documentation', 'link has correct text');
  });

  test('it renders internal link when @linkRoute is provided', async function (assert) {
    await render(hbs`
      <UsageReporting::Base::TitleRow
        @title="Authentication methods"
        @linkRoute="vault.cluster.access"
      />
    `);
    assert.dom('[data-test-vault-reporting-dashboard-card-title-link]').exists('renders internal route link');
  });

  test('it does not render a link when neither @linkUrl nor @linkRoute is provided', async function (assert) {
    await render(hbs`<UsageReporting::Base::TitleRow @title="No link" />`);
    assert.dom('[data-test-vault-reporting-dashboard-card-title-link]').doesNotExist();
  });

  test('it uses "View all" as the default link text', async function (assert) {
    await render(hbs`
      <UsageReporting::Base::TitleRow @title="Secret engines" @linkUrl="https://example.com" />
    `);
    assert
      .dom('[data-test-vault-reporting-dashboard-card-title-link]')
      .hasText('View all', 'defaults to "View all" link text');
  });
});
