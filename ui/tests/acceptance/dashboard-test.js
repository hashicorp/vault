/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, currentURL, settled, fillIn, click } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { create } from 'ember-cli-page-object';
import { selectChoose } from 'ember-power-select/test-support/helpers';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';
import { deleteEngineCmd } from 'vault/tests/helpers/commands';
import authPage from 'vault/tests/pages/auth';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import ENV from 'vault/config/environment';
import { formatNumber } from 'core/helpers/format-number';
// import { pollCluster } from 'vault/tests/helpers/poll-cluster';
import { disableReplication } from 'vault/tests/helpers/replication';

// selectors
import SECRETS_ENGINE_SELECTORS from 'vault/tests/helpers/components/dashboard/secrets-engines-card';
import VAULT_CONFIGURATION_SELECTORS from 'vault/tests/helpers/components/dashboard/vault-configuration-details-card';
import QUICK_ACTION_SELECTORS from 'vault/tests/helpers/components/dashboard/quick-actions-card';
import REPLICATION_CARD_SELECTORS from 'vault/tests/helpers/components/dashboard/replication-card';

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
    hooks.beforeEach(async function () {
      await authPage.login();
    });

    test('shows a secrets engine card', async function (assert) {
      await mountSecrets.enable('pki', 'pki');
      await settled();
      await visit('/vault/dashboard');
      assert.dom(SECRETS_ENGINE_SELECTORS.cardTitle).hasText('Secrets engines');
      assert.dom(SECRETS_ENGINE_SELECTORS.getSecretEngineAccessor('pki')).exists();
      // cleanup engine mount
      await consoleComponent.runCommands(deleteEngineCmd('pki'));
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      await mountSecrets.enable('nomad', 'nomad');
      await settled();
      await visit('/vault/dashboard');
      assert.dom('[data-test-secrets-engines-row="nomad"] [data-test-view]').doesNotExist();
      // cleanup engine mount
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
        .hasText("Don't see what you're looking for on this page? Let us know via our feedback form .");
    });
  });

  module('configuration details card', function (hooks) {
    hooks.beforeEach(async function () {
      this.data = {
        api_addr: 'http://127.0.0.1:8200',
        cache_size: 0,
        cluster_addr: 'https://127.0.0.1:8201',
        cluster_cipher_suites: '',
        cluster_name: '',
        default_lease_ttl: 0,
        default_max_request_duration: 0,
        detect_deadlocks: '',
        disable_cache: false,
        disable_clustering: false,
        disable_indexing: false,
        disable_mlock: true,
        disable_performance_standby: false,
        disable_printable_check: false,
        disable_sealwrap: false,
        disable_sentinel_trace: false,
        enable_response_header_hostname: false,
        enable_response_header_raft_node_id: false,
        enable_ui: true,
        experiments: null,
        introspection_endpoint: false,
        listeners: [
          {
            config: {
              address: '0.0.0.0:8200',
              cluster_address: '0.0.0.0:8201',
              tls_disable: true,
            },
            type: 'tcp',
          },
        ],
        log_format: '',
        log_level: 'debug',
        log_requests_level: '',
        max_lease_ttl: '48h',
        pid_file: '',
        plugin_directory: '',
        plugin_file_permissions: 0,
        plugin_file_uid: 0,
        raw_storage_endpoint: true,
        seals: [
          {
            disabled: false,
            type: 'shamir',
          },
        ],
        storage: {
          cluster_addr: 'https://127.0.0.1:8201',
          disable_clustering: false,
          raft: {
            max_entry_size: '',
          },
          redirect_addr: 'http://127.0.0.1:8200',
          type: 'raft',
        },
        telemetry: {
          add_lease_metrics_namespace_labels: false,
          circonus_api_app: '',
          circonus_api_token: '',
          circonus_api_url: '',
          circonus_broker_id: '',
          circonus_broker_select_tag: '',
          circonus_check_display_name: '',
          circonus_check_force_metric_activation: '',
          circonus_check_id: '',
          circonus_check_instance_id: '',
          circonus_check_search_tag: '',
          circonus_check_tags: '',
          circonus_submission_interval: '',
          circonus_submission_url: '',
          disable_hostname: true,
          dogstatsd_addr: '',
          dogstatsd_tags: null,
          lease_metrics_epsilon: 3600000000000,
          maximum_gauge_cardinality: 500,
          metrics_prefix: '',
          num_lease_metrics_buckets: 168,
          prometheus_retention_time: 86400000000000,
          stackdriver_debug_logs: false,
          stackdriver_location: '',
          stackdriver_namespace: '',
          stackdriver_project_id: '',
          statsd_address: '',
          statsite_address: '',
          usage_gauge_period: 5000000000,
        },
      };
      await authPage.login();
    });

    test('shows the configuration details card', async function (assert) {
      this.server.get('sys/config/state/sanitized', () => ({
        data: this.data,
        wrap_info: null,
        warnings: null,
        auth: null,
      }));
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.cardTitle).hasText('Configuration details');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.apiAddr).hasText('http://127.0.0.1:8200');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.defaultLeaseTtl).hasText('0');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.maxLeaseTtl).hasText('2 days');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.tlsDisable).hasText('Enabled');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.logFormat).hasText('None');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.logLevel).hasText('debug');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.storageType).hasText('raft');
    });
    test('shows the tls disabled if it is disabled', async function (assert) {
      this.server.get('sys/config/state/sanitized', () => {
        this.data.listeners[0].config.tls_disable = false;
        return {
          data: this.data,
          wrap_info: null,
          warnings: null,
          auth: null,
        };
      });
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.tlsDisable).hasText('Disabled');
    });
    test('shows the tls disabled if there is no tlsDisabled returned from server', async function (assert) {
      this.server.get('sys/config/state/sanitized', () => {
        this.data.listeners = [];

        return {
          data: this.data,
          wrap_info: null,
          warnings: null,
          auth: null,
        };
      });
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.tlsDisable).hasText('Disabled');
    });
  });

  module('quick actions card', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
    });

    test('shows the default state of the quick actions card', async function (assert) {
      assert.dom(QUICK_ACTION_SELECTORS.emptyState).exists();
    });

    test('shows the correct actions and links associated with pki', async function (assert) {
      await mountSecrets.enable('pki', 'pki');
      // generate role, issuer
      await runCommands([`write pki/root/generate/internal common_name="Hashicorp Test"`]);
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(QUICK_ACTION_SELECTORS.secretsEnginesSelect, 'pki');
      await fillIn(QUICK_ACTION_SELECTORS.actionSelect, 'Issue certificate');
      assert.dom(QUICK_ACTION_SELECTORS.emptyState).doesNotExist();
      assert.dom(QUICK_ACTION_SELECTORS.paramsTitle).hasText('Role to use');
      assert.dom(QUICK_ACTION_SELECTORS.getActionButton('Issue leaf certificate')).exists({ count: 1 });
      // select param input, click action, check route (issue cert)
      // await selectChoose(assert.dom(QUICK_ACTION_SELECTORS.paramSelect), 'pki');
      // await click(QUICK_ACTION_SELECTORS.getActionButton('Issue leaf certificate'));

      await fillIn(QUICK_ACTION_SELECTORS.actionSelect, 'View certificate');
      assert.dom(QUICK_ACTION_SELECTORS.emptyState).doesNotExist();
      assert.dom(QUICK_ACTION_SELECTORS.paramsTitle).hasText('Certificate serial number');
      assert.dom(QUICK_ACTION_SELECTORS.getActionButton('View certificate')).exists({ count: 1 });
      // select param input, click action, check route (view cert)

      await fillIn(QUICK_ACTION_SELECTORS.actionSelect, 'View issuer');
      assert.dom(QUICK_ACTION_SELECTORS.emptyState).doesNotExist();
      assert.dom(QUICK_ACTION_SELECTORS.paramsTitle).hasText('Issuer');
      assert.dom(QUICK_ACTION_SELECTORS.getActionButton('View issuer')).exists({ count: 1 });
      // select param input, click action, check route (view issuer)

      // cleanup engine mount
      await consoleComponent.runCommands(deleteEngineCmd('pki'));
    });

    test('shows the correct actions and links associated with database', async function (assert) {
      await mountSecrets.enable('database', 'database');
      // create a role
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(QUICK_ACTION_SELECTORS.secretsEnginesSelect, 'database');
      await fillIn(QUICK_ACTION_SELECTORS.actionSelect, 'Generate credentials for database');
      assert.dom(QUICK_ACTION_SELECTORS.emptyState).doesNotExist();
      assert.dom(QUICK_ACTION_SELECTORS.paramsTitle).hasText('Role to use');
      assert.dom(QUICK_ACTION_SELECTORS.getActionButton('Generate credentials')).exists({ count: 1 });
      // select param input, click action, check route
      // cleanup engine mount
      await consoleComponent.runCommands(deleteEngineCmd('database'));
    });

    test('shows the correct actions and links associated with kv', async function (assert) {
      await mountSecrets.enable('kv', 'kv');
      // create a secret
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(QUICK_ACTION_SELECTORS.secretsEnginesSelect, 'kv');
      await fillIn(QUICK_ACTION_SELECTORS.actionSelect, 'Find KV secrets');
      assert.dom(QUICK_ACTION_SELECTORS.emptyState).doesNotExist();
      assert.dom(QUICK_ACTION_SELECTORS.paramsTitle).hasText('Secret path');
      assert.dom(QUICK_ACTION_SELECTORS.getActionButton('Read secrets')).exists({ count: 1 });
      // select param input, click action, check route
      // cleanup engine mount
      await consoleComponent.runCommands(deleteEngineCmd('database'));
    });
  });
  module('client counts card', function (hooks) {
    hooks.before(async function () {
      ENV['ember-cli-mirage'].handler = 'clients';
    });

    hooks.beforeEach(async function () {
      this.store = this.owner.lookup('service:store');
      this.version = this.owner.lookup('service:version');
      this.server.get('sys/seal-status', function () {
        return {
          type: 'shamir',
          initialized: true,
          sealed: false,
          t: 1,
          n: 1,
          progress: 0,
          nonce: '',
          version: '1.14.0+ent',
          build_date: '2023-06-19T16:20:54Z',
          migration: false,
          cluster_name: 'vault-cluster-107c3b64',
          cluster_id: '7883553f-51ea-bea4-6135-2deb4f147fd1',
          recovery_seal: false,
          storage_type: 'inmem_transactional_ha',
        };
      });
      await authPage.login();
    });

    hooks.after(function () {
      ENV['ember-cli-mirage'].handler = null;
    });

    test('shows the client count card', async function (assert) {
      assert.true(this.version.isEnterprise);
      assert.strictEqual(this.version.version, '1.14.0+ent');
      assert.strictEqual(currentURL(), '/vault/dashboard');
      assert.dom('[data-test-client-count-card]').exists();
      const response = await this.store.peekRecord('clients/activity', 'some-activity-id');
      assert.dom('[data-test-client-count-title]').hasText('Client count');
      assert.dom('[data-test-stat-text="total-clients"] .stat-label').hasText('Total');
      assert
        .dom('[data-test-stat-text="total-clients"] .stat-value')
        .hasText(formatNumber([response.total.clients]));
      assert.dom('[data-test-stat-text="new-clients"] .stat-label').hasText('New');
      assert
        .dom('[data-test-stat-text="new-clients"] .stat-text')
        .hasText('The number of clients new to Vault in the current month.');
      assert
        .dom('[data-test-stat-text="new-clients"] .stat-value')
        .hasText(formatNumber([response.byMonth.lastObject.new_clients.clients]));
    });
    test('hides the client count card', async function (assert) {
      await authPage.logout();
      this.server.get(
        'sys/license/status',
        () => new Response(403, {}, { errors: ['API returns this error'] })
      );
      await authPage.login();
      assert.dom('[data-test-client-count-card]').doesNotExist();
    });
  });
  module('replication card', function (hooks) {
    hooks.beforeEach(async function () {
      this.version = this.owner.lookup('service:version');
      this.server.get('sys/seal-status', function () {
        return {
          type: 'shamir',
          initialized: true,
          sealed: false,
          t: 1,
          n: 1,
          progress: 0,
          nonce: '',
          version: '1.14.0+ent',
          build_date: '2023-06-19T16:20:54Z',
          migration: false,
          cluster_name: 'vault-cluster-107c3b64',
          cluster_id: '7883553f-51ea-bea4-6135-2deb4f147fd1',
          recovery_seal: false,
          storage_type: 'inmem_transactional_ha',
        };
      });
      await authPage.login();
      await settled();
      await disableReplication('dr');
      await settled();
      await disableReplication('performance');
      await settled();
    });

    test('shows the replication card empty state', async function (assert) {
      assert.true(this.version.isEnterprise);
      assert.strictEqual(this.version.version, '1.14.0+ent');
      await visit('/vault/dashboard');
      assert.dom(REPLICATION_CARD_SELECTORS.replicationEmptyState).exists();
      assert.dom(REPLICATION_CARD_SELECTORS.replicationEmptyStateTitle).hasText('Replication not set up');
      assert
        .dom(REPLICATION_CARD_SELECTORS.replicationEmptyStateMessage)
        .hasText('Data will be listed here. Enable a primary replication cluster to get started.');
      assert.dom(REPLICATION_CARD_SELECTORS.replicationEmptyStateActions).hasText('Enable replication');
    });

    test('it should show replication status if both dr and performance replication are enabled as features in version', async function (assert) {
      assert.true(this.version.isEnterprise);
      assert.strictEqual(this.version.version, '1.14.0+ent');
      await visit('/vault/dashboard');
      assert.strictEqual(currentURL(), '/vault/clients/dashboard');
      await click(REPLICATION_CARD_SELECTORS.replicationEmptyStateActionsLink);

      await click('[data-test-replication-enable]');
      await pollCluster(this.owner);
      await settled();
      await visit('/vault/dashboard');
      assert
        .dom(REPLICATION_CARD_SELECTORS.getReplicationTitle('dr-perf', 'DR primary'))
        .hasText('DR primary');
      assert.dom(REPLICATION_CARD_SELECTORS.getStateTooltipTitle('dr-perf', 'DR primary')).hasText('running');
      assert
        .dom(REPLICATION_CARD_SELECTORS.getStateTooltipIcon('dr-perf', 'DR primary', 'check-circle'))
        .exists();
      assert
        .dom(REPLICATION_CARD_SELECTORS.getReplicationTitle('dr-perf', 'Perf primary'))
        .hasText('Perf primary');
      assert
        .dom(REPLICATION_CARD_SELECTORS.getStateTooltipTitle('dr-perf', 'Perf primary'))
        .hasText('not set up');
      assert
        .dom(REPLICATION_CARD_SELECTORS.getStateTooltipIcon('dr-perf', 'Perf primary', 'x-circle'))
        .exists();
    });
  });
});
