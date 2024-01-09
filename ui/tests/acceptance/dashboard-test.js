/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  visit,
  currentURL,
  settled,
  fillIn,
  click,
  waitUntil,
  find,
  currentRouteName,
} from '@ember/test-helpers';
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
import { pollCluster } from 'vault/tests/helpers/poll-cluster';
import { disableReplication } from 'vault/tests/helpers/replication';
import connectionPage from 'vault/tests/pages/secrets/backend/database/connection';

import { SELECTORS } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';

const consoleComponent = create(consoleClass);

const createNS = async (name) => consoleComponent.runCommands(`write sys/namespaces/${name} -force`);

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
    const version = this.owner.lookup('service:version');
    const versionName = version.version;
    const versionText = version.isEnterprise
      ? `Vault v${versionName.slice(0, versionName.indexOf('+'))} root`
      : `Vault v${versionName}`;

    assert.dom(SELECTORS.cardHeader('Vault version')).hasText(versionText);
  });

  module('secrets engines card', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
    });

    test('shows a secrets engine card', async function (assert) {
      await mountSecrets.enable('pki', 'pki');
      await settled();
      await visit('/vault/dashboard');
      assert.dom(SELECTORS.cardHeader('Secrets engines')).hasText('Secrets engines');
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

      this.server.get('sys/config/state/sanitized', () => ({
        data: this.data,
        wrap_info: null,
        warnings: null,
        auth: null,
      }));
    });

    test('hides the configuration details card on a non-root namespace enterprise version', async function (assert) {
      await authPage.login();
      await visit('/vault/dashboard');
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'vault is enterprise');
      assert.dom(SELECTORS.cardName('configuration-details')).exists();
      createNS('world');
      await visit('/vault/dashboard?namespace=world');
      assert.dom(SELECTORS.cardName('configuration-details')).doesNotExist();
    });

    test('shows the configuration details card', async function (assert) {
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(SELECTORS.cardHeader('configuration')).hasText('Configuration details');
      assert
        .dom(SELECTORS.vaultConfigurationCard.configDetailsField('api_addr'))
        .hasText('http://127.0.0.1:8200');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('default_lease_ttl')).hasText('0');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('max_lease_ttl')).hasText('2 days');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('tls')).hasText('Disabled'); // tls_disable=true
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('log_format')).hasText('None');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('log_level')).hasText('debug');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('type')).hasText('raft');
    });

    test('it should show tls as enabled if tls_disable, tls_cert_file and tls_key_file are in the config', async function (assert) {
      this.data.listeners[0].config.tls_disable = false;
      this.data.listeners[0].config.tls_cert_file = './cert.pem';
      this.data.listeners[0].config.tls_key_file = './key.pem';

      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('tls')).hasText('Enabled');
    });

    test('it should show tls as enabled if only cert and key exist in config', async function (assert) {
      delete this.data.listeners[0].config.tls_disable;
      this.data.listeners[0].config.tls_cert_file = './cert.pem';
      this.data.listeners[0].config.tls_key_file = './key.pem';
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('tls')).hasText('Enabled');
    });

    test('it should show tls as disabled if there is no tls information in the config', async function (assert) {
      this.data.listeners = [];
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(SELECTORS.vaultConfigurationCard.configDetailsField('tls')).hasText('Disabled');
    });
  });

  module('quick actions card', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
    });

    test('shows the default state of the quick actions card', async function (assert) {
      assert.dom(SELECTORS.emptyState('no-mount-selected')).exists();
    });

    test('shows the correct actions and links associated with pki', async function (assert) {
      const backend = 'pki-dashboard';
      await mountSecrets.enable('pki', backend);
      await runCommands([
        `write ${backend}/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
      ]);
      await runCommands([
        `write ${backend}/root/generate/internal issuer_name="Hashicorp" common_name="Hello"`,
      ]);
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(SELECTORS.searchSelect('secrets-engines'), backend);
      await fillIn(SELECTORS.selectEl, 'Issue certificate');
      assert.dom(SELECTORS.emptyState('quick-actions')).doesNotExist();
      assert.dom(SELECTORS.subtitle('param')).hasText('Role to use');

      await selectChoose(SELECTORS.searchSelect('params'), 'some-role');
      assert.dom(SELECTORS.actionButton('Issue leaf certificate')).exists({ count: 1 });
      await click(SELECTORS.actionButton('Issue leaf certificate'));
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.roles.role.generate');

      await visit('/vault/dashboard');

      await selectChoose(SELECTORS.searchSelect('secrets-engines'), backend);
      await fillIn(SELECTORS.selectEl, 'View certificate');
      assert.dom(SELECTORS.emptyState('quick-actions')).doesNotExist();
      assert.dom(SELECTORS.subtitle('param')).hasText('Certificate serial number');
      assert.dom(SELECTORS.actionButton('View certificate')).exists({ count: 1 });
      await selectChoose(SELECTORS.searchSelect('params'), '.ember-power-select-option', 0);
      await click(SELECTORS.actionButton('View certificate'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.pki.certificates.certificate.details'
      );

      await visit('/vault/dashboard');

      await selectChoose(SELECTORS.searchSelect('secrets-engines'), backend);
      await fillIn(SELECTORS.selectEl, 'View issuer');
      assert.dom(SELECTORS.emptyState('quick-actions')).doesNotExist();
      assert.dom(SELECTORS.subtitle('param')).hasText('Issuer');
      assert.dom(SELECTORS.actionButton('View issuer')).exists({ count: 1 });
      await selectChoose(SELECTORS.searchSelect('params'), '.ember-power-select-option', 0);
      await click(SELECTORS.actionButton('View issuer'));
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.issuers.issuer.details');

      // cleanup engine mount
      await consoleComponent.runCommands(deleteEngineCmd(backend));
    });

    const newConnection = async (backend, plugin = 'mongodb-database-plugin') => {
      const name = `connection-${Date.now()}`;
      await connectionPage.visitCreate({ backend });
      await connectionPage.dbPlugin(plugin);
      await connectionPage.name(name);
      await connectionPage.connectionUrl(`mongodb://127.0.0.1:4321/${name}`);
      await connectionPage.toggleVerify();
      await connectionPage.save();
      await connectionPage.enable();
      return name;
    };

    test('shows the correct actions and links associated with database', async function (assert) {
      await mountSecrets.enable('database', 'database');
      await newConnection('database');
      await runCommands([
        `write database/roles/my-role \
        db_name=mongodb-database-plugin \
        creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "foo"}] }' \
        default_ttl="1h" \
        max_ttl="24h`,
      ]);
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(SELECTORS.searchSelect('secrets-engines'), 'database');
      await fillIn(SELECTORS.selectEl, 'Generate credentials for database');
      assert.dom(SELECTORS.emptyState('quick-actions')).doesNotExist();
      assert.dom(SELECTORS.subtitle('param')).hasText('Role to use');
      assert.dom(SELECTORS.actionButton('Generate credentials')).exists({ count: 1 });
      await selectChoose(SELECTORS.searchSelect('params'), '.ember-power-select-option', 0);
      await click(SELECTORS.actionButton('Generate credentials'));
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.credentials');
      await consoleComponent.runCommands(deleteEngineCmd('database'));
    });

    test('shows the correct actions and links associated with kv v1', async function (assert) {
      await runCommands(['write sys/mounts/kv type=kv', 'write kv/foo bar=baz']);
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(SELECTORS.searchSelect('secrets-engines'), 'kv');
      await fillIn(SELECTORS.selectEl, 'Find KV secrets');
      assert.dom(SELECTORS.emptyState('quick-actions')).doesNotExist();
      assert.dom(SELECTORS.subtitle('param')).hasText('Secret path');
      assert.dom(SELECTORS.actionButton('Read secrets')).exists({ count: 1 });
      await consoleComponent.runCommands(deleteEngineCmd('kv'));
    });
  });

  module('client counts card enterprise', function (hooks) {
    hooks.before(async function () {
      ENV['ember-cli-mirage'].handler = 'clients';
    });

    hooks.beforeEach(async function () {
      this.store = this.owner.lookup('service:store');

      await authPage.login();
    });

    hooks.after(function () {
      ENV['ember-cli-mirage'].handler = null;
    });

    test('shows the client count card for enterprise', async function (assert) {
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'version is enterprise');
      assert.strictEqual(currentURL(), '/vault/dashboard');
      assert.dom(SELECTORS.cardName('client-count')).exists();
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
  });

  module('replication card enterprise', function (hooks) {
    hooks.beforeEach(async function () {
      await authPage.login();
      await settled();
      await disableReplication('dr');
      await settled();
      await disableReplication('performance');
      await settled();
    });

    test('shows the replication card empty state in enterprise version', async function (assert) {
      await visit('/vault/dashboard');
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'vault is enterprise');
      assert.dom(SELECTORS.emptyState('replication')).exists();
      assert.dom(SELECTORS.emptyStateTitle('replication')).hasText('Replication not set up');
      assert
        .dom(SELECTORS.emptyStateMessage('replication'))
        .hasText('Data will be listed here. Enable a primary replication cluster to get started.');
      assert.dom(SELECTORS.emptyStateActions('replication')).hasText('Enable replication');
    });

    test('hides the replication card on a non-root namespace enterprise version', async function (assert) {
      await visit('/vault/dashboard');
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'vault is enterprise');
      assert.dom(SELECTORS.cardName('replication')).exists();
      createNS('blah');
      await visit('/vault/dashboard?namespace=blah');
      assert.dom(SELECTORS.cardName('replication')).doesNotExist();
    });

    test('it should show replication status if both dr and performance replication are enabled as features in enterprise', async function (assert) {
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'vault is enterprise');
      await visit('/vault/replication');
      assert.strictEqual(currentURL(), '/vault/replication');
      await click('[data-test-replication-type-select="performance"]');
      await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
      await click('[data-test-replication-enable]');
      await pollCluster(this.owner);
      assert.ok(
        await waitUntil(() => find('[data-test-replication-dashboard]')),
        'details dashboard is shown'
      );
      await visit('/vault/dashboard');
      assert.dom(SELECTORS.title('DR primary')).hasText('DR primary');
      assert.dom(SELECTORS.tooltipTitle('DR primary')).hasText('not set up');
      assert.dom(SELECTORS.tooltipIcon('dr-perf', 'DR primary', 'x-circle')).exists();
      assert.dom(SELECTORS.title('Performance primary')).hasText('Performance primary');
      assert.dom(SELECTORS.tooltipTitle('Performance primary')).hasText('running');
      assert.dom(SELECTORS.tooltipIcon('dr-perf', 'Performance primary', 'check-circle')).exists();
    });
  });
});
