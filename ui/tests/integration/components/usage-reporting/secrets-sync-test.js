/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | usage-reporting/secrets-sync', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders destinations count and badge list when data is present', async function (assert) {
    this.set('destinations', { aws: 1, gcp: 2 });
    await render(hbs`
      <UsageReporting::SecretsSync
        @totalDestinations={{3}}
        @destinations={{this.destinations}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-row]')
      .exists('destinations row renders');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-row]')
      .includesText('3 destinations', 'renders correct destinations count text');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-badge]')
      .exists({ count: 2 }, 'renders a badge for each destination type');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-row]')
      .includesText('AWS: 1', 'renders AWS badge');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-row]')
      .includesText('GCP: 2', 'renders GCP badge');
  });

  test('it uses singular "Destination" when totalDestinations is 1', async function (assert) {
    this.set('destinations', { aws: 1 });
    await render(hbs`
      <UsageReporting::SecretsSync
        @totalDestinations={{1}}
        @destinations={{this.destinations}}
      />
    `);
    assert
      .dom('[data-test-vault-reporting-secrets-sync-destinations-row]')
      .includesText('1 destination', 'uses singular form for one destination');
  });

  test('it renders the default empty state when totalDestinations is 0', async function (assert) {
    this.set('destinations', {});
    await render(hbs`
      <UsageReporting::SecretsSync
        @totalDestinations={{0}}
        @destinations={{this.destinations}}
      />
    `);
    assert.dom('[data-test-vault-reporting-secrets-sync-empty-state]').exists('empty state renders');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-empty-state-description]')
      .exists('empty state description renders');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-empty-state-link]')
      .exists('empty state docs link renders');
  });

  test('it renders a named block empty state', async function (assert) {
    this.set('destinations', {});
    await render(hbs`
      <UsageReporting::SecretsSync @totalDestinations={{0}} @destinations={{this.destinations}}>
        <:empty as |A|>
          <A.Body @text="Custom empty message." />
        </:empty>
      </UsageReporting::SecretsSync>
    `);
    assert
      .dom('[data-test-vault-reporting-secrets-sync-empty-state]')
      .includesText('Custom empty message.', 'renders custom named block empty state');
    assert
      .dom('[data-test-vault-reporting-secrets-sync-empty-state-description]')
      .doesNotExist('default description is not rendered when named block is used');
  });
});
