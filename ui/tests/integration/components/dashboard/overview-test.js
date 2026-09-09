/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

module('Integration | Component | dashboard/overview', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.flags = this.owner.lookup('service:flags');
    this.namespace = this.owner.lookup('service:namespace');
    this.permissions = this.owner.lookup('service:permissions');
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.isRootNamespace = true;
    this.replication = {
      dr: { clusterId: '123', state: 'running' },
      performance: { clusterId: 'abc-1', state: 'running', isPrimary: true },
    };
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_f3400dee',
        path: 'kv-test/',
        type: 'kv',
      },
    });
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_f3300dee',
        path: 'kv-1/',
        type: 'kv',
      },
    });
    this.secretsEngines = this.store.peekAll('secret-engine', {});
    this.vaultConfiguration = {
      api_addr: 'http://127.0.0.1:8200',
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      listeners: [
        {
          config: {
            address: '127.0.0.1:8200',
            tls_disable: 1,
          },
          type: 'tcp',
        },
      ],
    };
    this.refreshModel = () => {};
    // Disable checklist by default so existing tests are unaffected by the new
    // checklist lifecycle rendering logic.
    this.owner.lookup('service:checklist-state').isAvailable = false;
    this.renderComponent = async () => {
      return render(
        hbs`
        <Dashboard::Overview
          @secretsEngines={{this.secretsEngines}}
          @vaultConfiguration={{this.vaultConfiguration}}
          @replication={{this.replication}}
          @version={{this.version}}
          @isRootNamespace={{this.isRootNamespace}}
          @refreshModel={{this.refreshModel}} 
          @replicationUpdatedAt={{this.replicationUpdatedAt}}
          />
      `
      );
    };
  });

  test('it should show dashboard empty states in root namespace', async function (assert) {
    this.version.version = '1.13.1';
    this.secretsEngines = null;
    this.replication = null;
    this.vaultConfiguration = null;
    await this.renderComponent();
    assert.dom(GENERAL.hdsPageHeaderTitle).exists();
    assert.dom(GENERAL.textDisplay('Secrets engines')).exists();
    assert.dom(GENERAL.emptyState('secrets-engines')).exists();
    assert.dom(GENERAL.textDisplay('Learn more')).exists();
    assert.dom(GENERAL.textDisplay('Quick actions')).exists();
    assert
      .dom(GENERAL.cardContainer('feature-spotlight'))
      .doesNotExist('feature spotlight card is not shown for community');
    assert.dom(GENERAL.textDisplay('Cluster information')).doesNotExist();
    assert.dom(GENERAL.textDisplay('Cluster replication')).doesNotExist();
    assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
  });

  test('it renders the secrets engine card', async function (assert) {
    assert.expect(3);
    await this.renderComponent();
    assert.dom(GENERAL.textDisplay('Secrets engines')).hasText('Secrets engines');
    assert.dom(SES.secretPath('kv-1/')).exists();
    assert.dom(SES.secretPath('kv-test/')).exists();
  });

  module('client count and replication card', function (hooks) {
    hooks.beforeEach(function () {
      this.version.version = '1.13.1+ent';
      this.version.type = 'enterprise';
    });

    test('it should hide cards on community in root namespace', async function (assert) {
      this.version.version = '1.13.1';
      this.version.type = 'community';
      this.server.get(
        'sys/internal/counters/activity',
        () => new Error('uh oh! a request was made to sys/internal/counters/activity')
      );
      await this.renderComponent();

      assert.dom(GENERAL.hdsPageHeaderTitle).exists();
      assert.dom(GENERAL.textDisplay('Secrets engines')).exists();
      assert.dom(GENERAL.textDisplay('Learn more')).exists();
      assert.dom(GENERAL.textDisplay('Quick actions')).exists();
      assert.dom(GENERAL.textDisplay('Cluster information')).exists();
      assert.dom(GENERAL.textDisplay('Replication')).doesNotExist();
      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
    });

    test('it should hide cards on enterprise if permission but not in root namespace', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };
      this.isRootNamespace = false;
      await this.renderComponent();
      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
      assert.dom(GENERAL.textDisplay('Replication')).doesNotExist();
    });

    test('it should show cards on enterprise if has permission and in root namespace', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };
      await this.renderComponent();
      assert.dom(GENERAL.hdsPageHeaderTitle).exists();
      assert.dom(GENERAL.textDisplay('Secrets engines')).exists();
      assert.dom(GENERAL.textDisplay('Learn more')).exists();
      assert.dom(GENERAL.textDisplay('Quick actions')).exists();
      assert.dom(GENERAL.textDisplay('Cluster information')).exists();
      assert.dom(GENERAL.textDisplay('Replication')).doesNotExist();
      assert.dom(GENERAL.textDisplay('Client count')).exists();
    });

    test('it should show client count on enterprise in admin namespace when running a managed mode', async function (assert) {
      this.permissions.exactPaths = {
        'admin/sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'admin/sys/replication/status': {
          capabilities: ['read'],
        },
      };

      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespace.path = 'admin';
      this.isRootNamespace = false;

      await this.renderComponent();

      assert.dom(GENERAL.textDisplay('Client count')).exists();
    });

    test('it should hide client count on enterprise in child namespaces called "admin" when running a managed mode', async function (assert) {
      this.permissions.exactPaths = {
        'admin/sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'admin/sys/replication/status': {
          capabilities: ['read'],
        },
      };

      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespace.path = 'ns1/admin';
      this.isRootNamespace = false;

      await this.renderComponent();

      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
    });

    test('it should hide client count on enterprise in any other namespace when running a managed mode', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };

      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      this.namespace.path = 'groceries';
      this.isRootNamespace = false;

      await this.renderComponent();

      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
    });

    test('it should hide client count on PKI-only Secrets clusters', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
      };
      this.version.features = ['PKI-only Secrets'];
      await this.renderComponent();
      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
    });

    test('it should hide client count on enterprise when Consumption Billing is enabled', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
      };
      this.version.features = ['Consumption Billing'];

      await this.renderComponent();
      assert.dom(GENERAL.widget('client count')).doesNotExist();
    });

    test('it should hide cards on enterprise in root namespace but no permission', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
      assert.dom(GENERAL.textDisplay('Replication')).doesNotExist();
    });

    test('it should hide cards on enterprise if no permission and not in root namespace', async function (assert) {
      this.isRootNamespace = false;
      await this.renderComponent();
      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
      assert.dom(GENERAL.textDisplay('Replication')).doesNotExist();
    });

    test('it should hide client count on enterprise in root namespace if no activity permission', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['deny'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };

      await this.renderComponent();

      assert.dom(GENERAL.textDisplay('Client count')).doesNotExist();
      assert.dom(GENERAL.textDisplay('Cluster replication')).exists();
    });

    test('it should hide replication on enterprise in root namespace if no replication status permission', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['deny'],
        },
      };

      await this.renderComponent();
      assert.dom(GENERAL.textDisplay('Client count')).exists();
      assert.dom(GENERAL.textDisplay('Replication')).doesNotExist();
    });

    test('it should hide replication on enterprise if has permission and in root namespace but is empty', async function (assert) {
      this.permissions.exactPaths = {
        'sys/internal/counters/activity': {
          capabilities: ['read'],
        },
        'sys/replication/status': {
          capabilities: ['read'],
        },
      };
      this.replication = {};
      await this.renderComponent();
      assert.dom(GENERAL.textDisplay('Client count')).exists();
      assert.dom(GENERAL.textDisplay('Replication')).doesNotExist();
    });
  });

  test('it shows the feature spotlight card on enterprise', async function (assert) {
    this.version.version = '1.13.1+ent';
    this.version.type = 'enterprise';
    await this.renderComponent();

    assert
      .dom(GENERAL.cardContainer('feature-spotlight'))
      .exists('feature spotlight card is visible on enterprise');
    assert.dom(GENERAL.textDisplay('New Agent Registry in Vault')).hasText('New Agent Registry in Vault');
  });

  test('it does not show the feature spotlight card on community', async function (assert) {
    this.version.version = '1.13.1';
    this.version.type = 'community';
    await this.renderComponent();

    assert
      .dom(GENERAL.cardContainer('feature-spotlight'))
      .doesNotExist('feature spotlight card is not shown for community');
  });

  test('it shows the learn more card on community', async function (assert) {
    this.version.version = '1.13.1';
    this.version.type = 'community';
    await this.renderComponent();

    assert.dom(GENERAL.textDisplay('Learn more')).hasText('Learn more');
    assert
      .dom(GENERAL.textBody('Learn more description'))
      .hasText(
        'Explore the features of Vault and learn advance practices with the following tutorials and documentation.'
      );
    assert.dom('[data-test-learn-more-links] a').exists({ count: 3 });
  });

  test('it shows the learn more card on enterprise', async function (assert) {
    this.version.type = 'enterprise';
    this.version.features = [
      'Performance Replication',
      'DR Replication',
      'Namespaces',
      'Transform Secrets Engine',
    ];
    await this.renderComponent();
    assert.dom(GENERAL.textDisplay('Learn more')).hasText('Learn more');
    assert
      .dom(GENERAL.textBody('Learn more description'))
      .hasText(
        'Explore the features of Vault and learn advance practices with the following tutorials and documentation.'
      );
    assert.dom('[data-test-learn-more-links] a').exists({ count: 4 });
  });

  module('checklist lifecycle states', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
      this.checklistState = this.owner.lookup('service:checklist-state');
      this.checklistState['_state'] = {};
      // Reset the hidden-checklists list so tests that call hideChecklist() do
      // not bleed localStorage state into subsequent tests.
      this.checklistState['_hiddenChecklists'] = [];
      // Re-enable checklist for these tests
      this.checklistState.isAvailable = true;
      // Grant the minimum permissions required for hasChecklistEntryAccess to
      // return true (sys/mounts satisfies hasPermission('sys/mounts')).
      this.permissions.exactPaths = {
        'sys/mounts': { capabilities: ['read'] },
      };
      // Stub the API-backed tasks so step-completion toggles directly update
      // in-memory state without needing a real network round-trip in these
      // component integration tests. The service's API behavior is covered
      // separately in its unit tests.
      sinon.stub(this.checklistState.updateStep, 'perform').callsFake((checklistId, stepId, completed) => {
        this.checklistState['_state'] = {
          ...this.checklistState['_state'],
          [checklistId]: { ...(this.checklistState['_state'][checklistId] ?? {}), [stepId]: completed },
        };
      });
      sinon.stub(this.checklistState.markComplete, 'perform').callsFake((checklistId, stepId) => {
        this.checklistState['_state'] = {
          ...this.checklistState['_state'],
          [checklistId]: { ...(this.checklistState['_state'][checklistId] ?? {}), [stepId]: true },
        };
      });
    });

    hooks.afterEach(function () {
      sinon.restore();
    });

    test('shows checklist widget in active lifecycle state', async function (assert) {
      // No steps complete → active
      await this.renderComponent();

      assert.dom('[data-test-widget="checklist"]').exists('Checklist widget is rendered');
      assert.dom('[data-test-widget="congrats-banner"]').doesNotExist('No congrats banner');
      assert.dom('[data-test-widget="explore-vault"]').doesNotExist('No explore vault');
    });

    test('shows congrats banner when all visible steps are complete', async function (assert) {
      // Seed all inferred steps as complete (tvp-cli is excluded from visible on CE, but we're on enterprise)
      // Use only the non-Enterprise-only inferred steps to keep test simple
      this.checklistState['_state'] = {
        'cluster-startup': { 'tvp-cli': true, policy: true, auth: true, kv: true, namespaces: true },
      };

      await this.renderComponent();

      assert.dom('[data-test-widget="congrats-banner"]').exists('Congrats banner is rendered');
      assert.dom('[data-test-widget="checklist"]').doesNotExist('No checklist widget');
      assert.dom('[data-test-widget="explore-vault"]').doesNotExist('No explore vault');
    });

    test('does not show congrats banner when there are zero visible steps', async function (assert) {
      // Simulate a token with no nav permissions → zero visible steps
      this.checklistState['_state'] = {
        'cluster-startup': { 'tvp-cli': true, policy: true, auth: true, kv: true, namespaces: true },
      };
      // Block all nav permissions so every permission-gated step is hidden
      // tvp-cli and kv have no permission check so they're always visible — set all complete but block them via the stub
      // Instead, override getVisibleSteps directly
      sinon.stub(this.checklistState, 'getVisibleSteps').returns([]);

      await this.renderComponent();

      assert.dom('[data-test-widget="congrats-banner"]').doesNotExist('No congrats when zero visible steps');
      sinon.restore();
    });

    test('shows explore vault banner in hidden lifecycle state', async function (assert) {
      this.checklistState.hideChecklist('cluster-startup');

      await this.renderComponent();

      assert.dom('[data-test-widget="explore-vault"]').exists('Explore Vault banner shown');
      assert.dom('[data-test-explore-vault-restore]').exists('Restore button present');
      assert.dom('[data-test-widget="checklist"]').doesNotExist('No checklist widget');
    });

    test('shows explore vault banner in post-completion state', async function (assert) {
      this.checklistState['_state'] = {
        'cluster-startup': { 'tvp-cli': true, policy: true, auth: true, kv: true, namespaces: true },
      };
      this.checklistState.hideChecklist('cluster-startup');

      await this.renderComponent();

      assert
        .dom('[data-test-widget="explore-vault"]')
        .exists('Explore Vault banner shown after completion + dismiss');
    });

    test('restore button in explore vault banner brings back the checklist', async function (assert) {
      this.checklistState.hideChecklist('cluster-startup');
      await this.renderComponent();

      assert.dom('[data-test-widget="explore-vault"]').exists('Starts in explore vault');

      await click('[data-test-explore-vault-restore]');

      assert.dom('[data-test-widget="checklist"]').exists('Checklist restored after clicking restore');
    });

    test('clicking hide button in checklist transitions to explore vault banner', async function (assert) {
      await this.renderComponent();

      assert.dom('[data-test-widget="checklist"]').exists('Starts in active checklist state');

      await click('[data-test-checklist-hide]');

      assert.dom('[data-test-widget="explore-vault"]').exists('Explore vault banner shown after hide');
      assert.dom('[data-test-widget="checklist"]').doesNotExist('Checklist is no longer visible');
    });

    test('after back to setup, incomplete then re-complete shows congrats again', async function (assert) {
      this.checklistState['_state'] = {
        'cluster-startup': { 'tvp-cli': true, policy: true, auth: true, kv: true, namespaces: true },
      };

      await this.renderComponent();

      assert.dom('[data-test-widget="congrats-banner"]').exists('Starts on congrats when fully complete');

      await click('[data-test-congrats-back]');
      assert.dom('[data-test-widget="checklist"]').exists('Back shows checklist with completed steps');

      await click('[data-test-checklist-step="tvp-cli"] [data-test-checklist-step-mark-complete="tvp-cli"]');
      assert
        .dom('[data-test-widget="checklist"]')
        .exists('Checklist remains visible after marking incomplete');

      await click('[data-test-checklist-step="tvp-cli"] [data-test-checklist-step-mark-complete="tvp-cli"]');
      assert.dom('[data-test-widget="congrats-banner"]').exists('Congrats shows again after re-completing');
    });
  });
});
