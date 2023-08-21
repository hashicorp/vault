/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, currentURL, settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import authPage from 'vault/tests/pages/auth';
import SECRETS_ENGINE_SELECTORS from 'vault/tests/helpers/components/dashboard/secrets-engines-card';
// import VAULT_CONFIGURATION_SELECTORS from 'vault/tests/helpers/components/dashboard/vault-configuration-details-card';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { deleteEngineCmd } from 'vault/tests/helpers/commands';
import { create } from 'ember-cli-page-object';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { selectChoose } from 'ember-power-select/test-support/helpers';

const consoleComponent = create(consoleClass);

module('Acceptance | landing page dashboard', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('navigate to dashboard on login', async function (assert) {
    await authPage.login();
    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  test('display the version number for the title', async function (assert) {
    await authPage.login();
    await visit('/vault/dashboard');
    assert.dom('[data-test-dashboard-version-header]').hasText('Vault v1.9.0 root');
  });

  module('secrets engines card', function (hooks) {
    hooks.beforeEach(function () {
      return authPage.login();
    });

    test('shows a secrets engine card', async function (assert) {
      await mountSecrets.enable('pki', 'pki');
      await settled();
      await visit('/vault/dashboard');
      assert.dom(SECRETS_ENGINE_SELECTORS.cardTitle).hasText('Secrets engines');
      assert.dom(SECRETS_ENGINE_SELECTORS.getSecretEngineAccessor('pki')).exists();
      // cleanup engine
      await consoleComponent.runCommands(deleteEngineCmd('pki'));
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      await mountSecrets.enable('nomad', 'nomad');
      await settled();
      await visit('/vault/dashboard');
      assert.dom('[data-test-secrets-engines-row="nomad"] [data-test-view]').doesNotExist();
      // cleanup engine
      await consoleComponent.runCommands(deleteEngineCmd('nomad'));
    });
  });

  module('learn more card', function (hooks) {
    hooks.beforeEach(function () {
      return authPage.login();
    });
    test('shows the learn more card', async function (assert) {
      await visit('/vault/dashboard');
      assert.dom('[data-test-learn-more-title]').hasText('Learn more');
      assert
        .dom('[data-test-learn-more-subtext]')
        .hasText(
          'Explore the features of Vault and learn advance practices with the following tutorials and documentation.'
        );
      assert.dom('[data-test-learn-more-links] a').exists({ count: 4 });
      assert
        .dom('[data-test-feedback-form]')
        .hasText("Don't see what you're looking for on this page? Let us know via our feedback form. ");
    });
  });

  module('quick actions card', function (hooks) {
    hooks.beforeEach(function () {
      return authPage.login();
    });

    test('shows the default state of the quick actions card', async function (assert) {
      await mountSecrets.enable('pki', 'pki');
      await mountSecrets.enable('database', 'database');
      await mountSecrets.enable('kv', 'kv');
      await visit('/vault/dashboard');
      assert.dom('[data-test-no-mount-selected-empty]').exists();
      await selectChoose('.search-select', 'pki-0-test');
      // await this.pauseTest();
    });
  });
});
