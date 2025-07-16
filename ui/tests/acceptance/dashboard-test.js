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
import { selectChoose } from 'ember-power-select/test-support';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import clientsHandlers from 'vault/mirage/handlers/clients';
import { formatNumber } from 'core/helpers/format-number';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';
import { disableReplication } from 'vault/tests/helpers/replication';
import connectionPage from 'vault/tests/pages/secrets/backend/database/connection';
import { v4 as uuidv4 } from 'uuid';
import { runCmd, deleteEngineCmd, createNS, deleteNS } from 'vault/tests/helpers/commands';

import { DASHBOARD } from 'vault/tests/helpers/components/dashboard/dashboard-selectors';
import { CUSTOM_MESSAGES } from 'vault/tests/helpers/config-ui/message-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const authenticatedMessageResponse = {
  request_id: '664fbad0-fcd8-9023-4c5b-81a7962e9f4b',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    key_info: {
      'some-awesome-id-2': {
        authenticated: true,
        end_time: null,
        link: {
          'some link title': 'www.link.com',
        },
        message: 'aGVsbG8gd29ybGQgaGVsbG8gd29scmQ=',
        options: null,
        start_time: '2024-01-04T08:00:00Z',
        title: 'Banner title',
        type: 'banner',
      },
      'some-awesome-id-1': {
        authenticated: true,
        end_time: null,
        message: 'aGVyZSBpcyBhIGNvb2wgbWVzc2FnZQ==',
        options: null,
        start_time: '2024-01-01T08:00:00Z',
        title: 'Modal title',
        type: 'modal',
      },
    },
    keys: ['some-awesome-id-2', 'some-awesome-id-1'],
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: '',
};

module('Acceptance | landing page dashboard', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('navigate to dashboard on login', async function (assert) {
    assert.expect(1);
    await login();
    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  test('display the version number for the title', async function (assert) {
    assert.expect(1);
    await login();
    await visit('/vault/dashboard');
    const version = this.owner.lookup('service:version');
    // Since we're using mirage, version is mocked static value
    const versionText = version.isEnterprise
      ? `Vault ${version.versionDisplay} root`
      : `Vault ${version.versionDisplay}`;

    assert.dom(DASHBOARD.cardHeader('Vault version')).hasText(versionText);
  });

  module('secrets engines card', function (hooks) {
    hooks.beforeEach(async function () {
      await login();
    });

    test('shows a secrets engine card', async function (assert) {
      assert.expect(1);
      await mountSecrets.enable('pki', 'pki');
      await settled();
      await visit('/vault/dashboard');
      assert.dom(DASHBOARD.cardHeader('Secrets engines')).hasText('Secrets engines');
      // cleanup engine mount
      await runCmd(deleteEngineCmd('pki'));
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      assert.expect(1);
      await mountSecrets.enable('nomad', 'nomad');
      await settled();
      await visit('/vault/dashboard');
      assert.dom('[data-test-secrets-engines-row="nomad"] [data-test-view]').doesNotExist();
      // cleanup engine mount
      await runCmd(deleteEngineCmd('nomad'));
    });
  });

  module('configuration details card', function (hooks) {
    hooks.beforeEach(async function () {
      this.data = {
        apiAddr: 'http://127.0.0.1:8200',
        cacheSize: 0,
        clusterAddr: 'https://127.0.0.1:8201',
        clusterCipherSuites: '',
        clusterName: '',
        defaultLeaseTtl: 0,
        defaultMaxRequestDuration: 0,
        detectDeadlocks: '',
        disableCache: false,
        disableClustering: false,
        disableIndexing: false,
        disableMlock: true,
        disablePerformanceStandby: false,
        disablePrintableCheck: false,
        disableSealwrap: false,
        disableSentinelTrace: false,
        enableResponseHeaderHostname: false,
        enableResponseHeaderRaftNodeId: false,
        enableUi: true,
        experiments: null,
        introspectionEndpoint: false,
        listeners: [
          {
            config: {
              address: '0.0.0.0:8200',
              clusterAddress: '0.0.0.0:8201',
              tlsDisable: true,
            },
            type: 'tcp',
          },
        ],
        logFormat: '',
        logLevel: 'debug',
        logRequestsLevel: '',
        maxLeaseTtl: '48h',
        pidFile: '',
        pluginDirectory: '',
        pluginFilePermissions: 0,
        pluginFileUid: 0,
        rawStorageEndpoint: true,
        seals: [
          {
            disabled: false,
            type: 'shamir',
          },
        ],
        storage: {
          clusterAddr: 'https://127.0.0.1:8201',
          disableClustering: false,
          raft: {
            maxEntrySize: '',
          },
          redirectAddr: 'http://127.0.0.1:8200',
          type: 'raft',
        },
        telemetry: {
          addLeaseMetricsNamespaceLabels: false,
          circonusApiApp: '',
          circonusApiToken: '',
          circonusApiUrl: '',
          circonusBrokerId: '',
          circonusBrokerSelectTag: '',
          circonusCheckDisplayName: '',
          circonusCheckForceMetricActivation: '',
          circonusCheckId: '',
          circonusCheckInstanceId: '',
          circonusCheckSearchTag: '',
          circonusCheckTags: '',
          circonusSubmissionInterval: '',
          circonusSubmissionUrl: '',
          disableHostname: true,
          dogstatsdAddr: '',
          dogstatsdTags: null,
          leaseMetricsEpsilon: 3600000000000,
          maximumGaugeCardinality: 500,
          metricsPrefix: '',
          numLeaseMetricsBuckets: 168,
          prometheusRetentionTime: 86400000000000,
          stackdriverDebugLogs: false,
          stackdriverLocation: '',
          stackdriverNamespace: '',
          stackdriverProjectId: '',
          statsdAddress: '',
          statsiteAddress: '',
          usageGaugePeriod: 5000000000,
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
      assert.expect(3);
      await login();
      await visit('/vault/dashboard');
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'vault is enterprise');
      assert.dom(DASHBOARD.cardName('configuration-details')).exists();
      await runCmd(createNS('world'), false);
      await visit('/vault/dashboard?namespace=world');
      assert.dom(DASHBOARD.cardName('configuration-details')).doesNotExist();

      // navigate to "root" before deleting
      await visit('vault/dashboard');
      // clean up namespace pollution
      await runCmd(deleteNS('world'));
    });

    test('shows the configuration details card', async function (assert) {
      assert.expect(8);
      await login();
      await visit('/vault/dashboard');
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
      assert.expect(1);
      this.data.listeners[0].config.tlsDisable = false;
      this.data.listeners[0].config.tlsCertFile = './cert.pem';
      this.data.listeners[0].config.tlsKeyFile = './key.pem';

      await login();
      await visit('/vault/dashboard');
      assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('tls')).hasText('Enabled');
    });

    test('it should show tls as enabled if only cert and key exist in config', async function (assert) {
      assert.expect(1);
      delete this.data.listeners[0].config.tlsDisable;
      this.data.listeners[0].config.tlsCertFile = './cert.pem';
      this.data.listeners[0].config.tlsKeyFile = './key.pem';
      await login();
      await visit('/vault/dashboard');
      assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('tls')).hasText('Enabled');
    });

    test('it should show tls as disabled if there is no tls information in the config', async function (assert) {
      assert.expect(1);
      this.data.listeners = [];
      await login();
      await visit('/vault/dashboard');
      assert.dom(DASHBOARD.vaultConfigurationCard.configDetailsField('tls')).hasText('Disabled');
    });
  });

  module('quick actions card', function (hooks) {
    hooks.beforeEach(async function () {
      await login();
    });

    test('shows the default state of the quick actions card', async function (assert) {
      assert.expect(1);
      assert.dom(DASHBOARD.emptyState('no-mount-selected')).exists();
    });

    test('shows the correct actions and links associated with pki', async function (assert) {
      assert.expect(12);
      const backend = 'pki-dashboard';
      await mountSecrets.enable('pki', backend);
      await runCmd([
        `write ${backend}/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
      ]);
      await runCmd([`write ${backend}/root/generate/internal issuer_name="Hashicorp" common_name="Hello"`]);
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(DASHBOARD.searchSelect('secrets-engines'), backend);
      await fillIn(DASHBOARD.selectEl, 'Issue certificate');
      assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
      assert.dom(DASHBOARD.subtitle('param')).hasText('Role to use');

      await selectChoose(DASHBOARD.searchSelect('params'), 'some-role');
      assert.dom(DASHBOARD.actionButton('Issue leaf certificate')).exists({ count: 1 });
      await click(DASHBOARD.actionButton('Issue leaf certificate'));
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.roles.role.generate');

      await visit('/vault/dashboard');

      await selectChoose(DASHBOARD.searchSelect('secrets-engines'), backend);
      await fillIn(DASHBOARD.selectEl, 'View certificate');
      assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
      assert.dom(DASHBOARD.subtitle('param')).hasText('Certificate serial number');
      assert.dom(DASHBOARD.actionButton('View certificate')).exists({ count: 1 });
      await selectChoose(DASHBOARD.searchSelect('params'), '.ember-power-select-option', 0);
      await click(DASHBOARD.actionButton('View certificate'));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.secrets.backend.pki.certificates.certificate.details'
      );

      await visit('/vault/dashboard');

      await selectChoose(DASHBOARD.searchSelect('secrets-engines'), backend);
      await fillIn(DASHBOARD.selectEl, 'View issuer');
      assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
      assert.dom(DASHBOARD.subtitle('param')).hasText('Issuer');
      assert.dom(DASHBOARD.actionButton('View issuer')).exists({ count: 1 });
      await selectChoose(DASHBOARD.searchSelect('params'), '.ember-power-select-option', 0);
      await click(DASHBOARD.actionButton('View issuer'));
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.pki.issuers.issuer.details');

      // cleanup engine mount
      await runCmd(deleteEngineCmd(backend));
    });

    const newConnection = async (backend, plugin = 'mongodb-database-plugin') => {
      const name = `connection-${Date.now()}`;
      await connectionPage.visitCreate({ backend });
      await connectionPage.dbPlugin(plugin);
      await connectionPage.name(name);
      await connectionPage.connectionUrl(`mongodb://127.0.0.1:4321/${name}`);
      await connectionPage.toggleVerify();
      await click(GENERAL.submitButton);
      await connectionPage.enable();
      return name;
    };

    test('shows the correct actions and links associated with database', async function (assert) {
      assert.expect(4);
      const databaseBackend = `database-${uuidv4()}`;
      await mountSecrets.enable('database', databaseBackend);
      await newConnection(databaseBackend);
      await runCmd([
        `write ${databaseBackend}/roles/my-role \
        db_name=mongodb-database-plugin \
        creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "foo"}] }' \
        default_ttl="1h" \
        max_ttl="24h`,
      ]);
      await settled();
      await visit('/vault/dashboard');
      await selectChoose(DASHBOARD.searchSelect('secrets-engines'), databaseBackend);
      await fillIn(DASHBOARD.selectEl, 'Generate credentials for database');
      assert.dom(DASHBOARD.emptyState('quick-actions')).doesNotExist();
      assert.dom(DASHBOARD.subtitle('param')).hasText('Role to use');
      assert.dom(DASHBOARD.actionButton('Generate credentials')).exists({ count: 1 });
      await selectChoose(DASHBOARD.searchSelect('params'), '.ember-power-select-option', 0);
      await click(DASHBOARD.actionButton('Generate credentials'));
      assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.backend.credentials');
      await runCmd(deleteEngineCmd(databaseBackend));
    });

    test('does not show kv1 mounts', async function (assert) {
      assert.expect(1);
      // delete before in case you are rerunning the test and it fails without deleting
      await runCmd(deleteEngineCmd('kv1'));
      await runCmd([`write sys/mounts/kv1 type=kv`]);
      await settled();
      await visit('/vault/dashboard');
      await clickTrigger('#type-to-select-a-mount');
      assert
        .dom('.ember-power-select-option')
        .doesNotHaveTextContaining('kv1', 'dropdown does not show kv1 mount');
      await runCmd(deleteEngineCmd('kv1'));
    });
  });

  module('client counts card enterprise', function (hooks) {
    hooks.beforeEach(async function () {
      clientsHandlers(this.server);
      this.store = this.owner.lookup('service:store');

      await login();
    });

    test('shows the client count card for enterprise', async function (assert) {
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'version is enterprise');
      assert.strictEqual(currentURL(), '/vault/dashboard');
      assert.dom(DASHBOARD.cardName('client-count')).exists();
      const response = await this.store.peekRecord('clients/activity', 'clients/activity');
      assert.dom('[data-test-client-count-title]').hasText('Client count');
      assert.dom('[data-test-stat-text="Total"] .stat-label').hasText('Total');
      assert.dom('[data-test-stat-text="Total"] .stat-value').hasText(formatNumber([response.total.clients]));
      assert.dom('[data-test-stat-text="New"] .stat-label').hasText('New');
      assert
        .dom('[data-test-stat-text="New"] .stat-text')
        .hasText('The number of clients new to Vault in the current month.');
      assert
        .dom('[data-test-stat-text="New"] .stat-value')
        .hasText(formatNumber([response.byMonth.lastObject.new_clients.clients]));
      assert
        .dom(`${GENERAL.flashMessage}.is-info`)
        .doesNotExist('Does not show warning about client count estimations');
    });
  });

  module('replication card enterprise', function (hooks) {
    hooks.beforeEach(async function () {
      await login();
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
      assert.dom(DASHBOARD.emptyState('replication')).exists();
      assert.dom(DASHBOARD.emptyStateTitle('replication')).hasText('Replication not set up');
      assert
        .dom(DASHBOARD.emptyStateMessage('replication'))
        .hasText('Data will be listed here. Enable a primary replication cluster to get started.');
      assert.dom(DASHBOARD.emptyStateActions('replication')).hasText('Enable replication');
    });

    test('hides the replication card on a non-root namespace enterprise version', async function (assert) {
      await visit('/vault/dashboard');
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'vault is enterprise');
      assert.dom(DASHBOARD.cardName('replication')).exists();
      await runCmd(createNS('blah'), false);
      await visit('/vault/dashboard?namespace=blah');
      assert.dom(DASHBOARD.cardName('replication')).doesNotExist();

      // navigate to "root" before deleting
      await visit('vault/dashboard');
      // clean up namespace pollution
      await runCmd(deleteNS('blah'));
    });

    test('it should show replication status if both dr and performance replication are enabled as features in enterprise', async function (assert) {
      const version = this.owner.lookup('service:version');
      assert.true(version.isEnterprise, 'vault is enterprise');
      await visit('/vault/replication');
      assert.strictEqual(currentURL(), '/vault/replication');
      await click('[data-test-replication-type-select="performance"]');
      await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
      await click(GENERAL.submitButton);
      await pollCluster(this.owner);
      assert.ok(
        await waitUntil(() => find('[data-test-replication-dashboard]')),
        'details dashboard is shown'
      );
      await visit('/vault/dashboard');
      assert.dom(DASHBOARD.title('DR primary')).hasText('DR primary');
      assert.dom(DASHBOARD.tooltipTitle('DR primary')).hasText('not set up');
      assert.dom(DASHBOARD.tooltipIcon('dr-perf', 'DR primary', 'x-circle')).exists();
      assert.dom(DASHBOARD.title('Performance primary')).hasText('Performance primary');
      assert.dom(DASHBOARD.tooltipTitle('Performance primary')).hasText('running');
      assert.dom(DASHBOARD.tooltipIcon('dr-perf', 'Performance primary', 'check-circle')).exists();
    });
  });

  module('custom messages auth tests', function (hooks) {
    hooks.beforeEach(function () {
      return this.server.get('/sys/internal/ui/mounts', () => ({}));
    });

    test('it shows the alert banner and modal message', async function (assert) {
      this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
        return authenticatedMessageResponse;
      });
      await visit('/vault/dashboard');
      const modalId = 'some-awesome-id-1';
      const alertId = 'some-awesome-id-2';
      assert.dom(CUSTOM_MESSAGES.modal(modalId)).exists();
      assert.dom(CUSTOM_MESSAGES.modalTitle(modalId)).hasText('Modal title');
      assert.dom(CUSTOM_MESSAGES.modalBody(modalId)).exists();
      assert.dom(CUSTOM_MESSAGES.modalBody(modalId)).hasText('here is a cool message');
      assert.dom(CUSTOM_MESSAGES.alertTitle(alertId)).hasText('Banner title');
      assert.dom(CUSTOM_MESSAGES.alertDescription(alertId)).hasText('hello world hello wolrd');
      assert.dom(CUSTOM_MESSAGES.alertAction('link')).hasText('some link title');
    });

    test('it shows the multiple modal messages', async function (assert) {
      const modalIdOne = 'some-awesome-id-2';
      const modalIdTwo = 'some-awesome-id-1';

      this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
        authenticatedMessageResponse.data.key_info[modalIdOne].type = 'modal';
        authenticatedMessageResponse.data.key_info[modalIdOne].title = 'Modal title 1';
        authenticatedMessageResponse.data.key_info[modalIdTwo].type = 'modal';
        authenticatedMessageResponse.data.key_info[modalIdTwo].title = 'Modal title 2';
        return authenticatedMessageResponse;
      });
      await visit('/vault/dashboard');
      assert.dom(CUSTOM_MESSAGES.modal(modalIdOne)).exists();
      assert.dom(CUSTOM_MESSAGES.modalTitle(modalIdOne)).hasText('Modal title 1');
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdOne)).exists();
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdOne)).hasText('hello world hello wolrd some link title');
      assert.dom(CUSTOM_MESSAGES.modal(modalIdTwo)).exists();
      assert.dom(CUSTOM_MESSAGES.modalTitle(modalIdTwo)).hasText('Modal title 2');
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdTwo)).exists();
      assert.dom(CUSTOM_MESSAGES.modalBody(modalIdTwo)).hasText('here is a cool message');
    });

    test('it shows the multiple banner messages', async function (assert) {
      const bannerIdOne = 'some-awesome-id-2';
      const bannerIdTwo = 'some-awesome-id-1';

      this.server.get('/sys/internal/ui/unauthenticated-messages', function () {
        authenticatedMessageResponse.data.key_info[bannerIdOne].type = 'banner';
        authenticatedMessageResponse.data.key_info[bannerIdOne].title = 'Banner title 1';
        authenticatedMessageResponse.data.key_info[bannerIdTwo].type = 'banner';
        authenticatedMessageResponse.data.key_info[bannerIdTwo].title = 'Banner title 2';
        return authenticatedMessageResponse;
      });
      await visit('/vault/dashboard');
      assert.dom(CUSTOM_MESSAGES.alertTitle(bannerIdOne)).hasText('Banner title 1');
      assert.dom(CUSTOM_MESSAGES.alertDescription(bannerIdOne)).hasText('hello world hello wolrd');
      assert.dom(CUSTOM_MESSAGES.alertAction('link')).hasText('some link title');
      assert.dom(CUSTOM_MESSAGES.alertTitle(bannerIdTwo)).hasText('Banner title 2');
      assert.dom(CUSTOM_MESSAGES.alertDescription(bannerIdTwo)).hasText('here is a cool message');
    });
  });
});
