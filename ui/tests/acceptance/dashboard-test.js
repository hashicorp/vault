/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import authPage from 'vault/tests/pages/auth';
import SECRETS_ENGINE_SELECTORS from 'vault/tests/helpers/components/dashboard/secrets-engines-card';
import VAULT_CONFIGURATION_SELECTORS from 'vault/tests/helpers/components/dashboard/vault-configuration-details-card';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | landing page dashboard', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  // TODO LANDING PAGE: create a test that will navigate to dashboard if user opts into new dashboard ui
  test('does not navigate to dashboard on login when user has not opted into dashboard ui', async function (assert) {
    assert.strictEqual(currentURL(), '/vault/secrets');

    await visit('/vault/dashboard');

    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  test('display the version number for the title', async function (assert) {
    await visit('/vault/dashboard');
    assert.dom('[data-test-dashboard-version-header]').hasText('Vault v1.9.0');
  });

  module('secrets engines card', function () {
    test('shows a secrets engine card', async function (assert) {
      await mountSecrets.enable('pki', 'pki');
      await visit('/vault/dashboard');
      assert.dom(SECRETS_ENGINE_SELECTORS.cardTitle).hasText('Secrets engines');
      assert.dom(SECRETS_ENGINE_SELECTORS.getSecretEngineAccessor('pki')).exists();
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      await mountSecrets.enable('nomad', 'nomad');
      await visit('/vault/dashboard');
      assert.dom('[data-test-secrets-engines-row="nomad"] [data-test-view]').doesNotExist();
    });
  });

  module('learn more card', function () {
    test('shows the learn more card', async function (assert) {
      await visit('/vault/dashboard');
      assert.dom('[data-test-learn-more-title]').hasText('Learn more');
      assert
        .dom('[data-test-learn-more-subtext]')
        .hasText(
          'Explore the features of Vault and learn advance practices with the following tutorials and documentation.'
        );
      assert.dom('[data-test-learn-more-links] a').exists({ count: 4 });
    });
  });

  module('configuration details card', function () {
    test('shows the configuration details card', async function (assert) {
      this.server.get('sys/config/state/sanitized', () => ({
        data: {
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
        },
        wrap_info: null,
        warnings: null,
        auth: null,
      }));
      await visit('/vault/dashboard');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.cardTitle).hasText('Configuration details');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.apiAddr).hasText('http://127.0.0.1:8200');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.defaultLeaseTtl).hasText('0');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.maxLeaseTtl).hasText('2 days');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.tlsDisable).hasText('true');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.logFormat).hasText('None');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.logLevel).hasText('debug');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.storageType).hasText('raft');
    });
  });
});
