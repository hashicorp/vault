/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { DASHBOARD } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';

module('Integration | Component | dashboard/vault-configuration-details-card', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
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

    this.renderComponent = () => {
      return render(hbs`<Dashboard::VaultConfigurationDetailsCard @vaultConfiguration={{this.data}} />`);
    };
  });

  test('it renders configuration details', async function (assert) {
    await this.renderComponent();
    assert.dom(DASHBOARD.cardHeader('configuration')).hasText('Configuration details');
    assert
      .dom(DASHBOARD.vaultConfigurationCard.configDetailsField('api_addr'))
      .hasText('http://127.0.0.1:8200');
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('default_lease_ttl')).hasText('0');
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('max_lease_ttl')).hasText('2 days');
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('tls')).hasText('Disabled'); // tls_disable=true
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('log_format')).hasText('None');
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('log_level')).hasText('debug');
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('type')).hasText('raft');
  });

  test('it should show tls as enabled if tls_disable, tls_cert_file and tls_key_file are in the config', async function (assert) {
    this.data.listeners[0].config.tls_disable = false;
    this.data.listeners[0].config.tls_cert_file = './cert.pem';
    this.data.listeners[0].config.tls_key_file = './key.pem';

    await this.renderComponent();
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('tls')).hasText('Enabled');
  });

  test('it should show tls as enabled if only cert and key exist in config', async function (assert) {
    delete this.data.listeners[0].config.tls_disable;
    this.data.listeners[0].config.tls_cert_file = './cert.pem';
    this.data.listeners[0].config.tls_key_file = './key.pem';

    await this.renderComponent();
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('tls')).hasText('Enabled');
  });

  test('it should show tls as disabled if there is no tls information in the config', async function (assert) {
    this.data.listeners = [];
    await this.renderComponent();
    assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('tls')).hasText('Disabled');
  });
});
