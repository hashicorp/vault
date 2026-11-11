/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import {
  click,
  currentRouteName,
  currentURL,
  fillIn,
  settled,
  visit,
  waitUntil,
  waitFor,
} from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { ALL_FEATURES, mockReplication } from 'vault/tests/helpers/replication-mocks';

// These tests characterize the replication routes: URLs, full route names, redirects and the
// pages they render. Replication endpoints are mocked, so they do not depend on the backend edition.
// They are frozen while the replication engine is inlined into the main app (VAULT-51483).

const ROUTE_PREFIX = 'vault.cluster.replication';
const s = {
  tab: (name) => `[data-test-replication-link="${name}"]`,
  enableForm: '[data-test-replication-enable-form]',
  dashboard: '[data-test-replication-dashboard]',
  summaryCard: '[data-test-replication-summary-card]',
  secondaryAdd: '[data-test-secondary-add]',
  secondaryName: (name) => `[data-test-secondary-name="${name}"]`,
  emptyTitle: '[data-test-empty-state-title]',
};

const assertRoute = (assert, url, routeName, message) => {
  assert.strictEqual(currentURL(), url, `${message}: url`);
  assert.strictEqual(currentRouteName(), routeName, `${message}: route name`);
};

module('Acceptance | replication routes', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.mockReplication = async (options) => {
      this.state = mockReplication(this.server, options);
      await login();
      // the path filter form expects the accessible namespaces, which the namespace picker usually loads
      this.owner.lookup('service:namespace').accessibleNamespaces = [];
    };
  });

  module('url and route name mapping', function () {
    test('every replication url resolves to its route name', async function (assert) {
      await this.mockReplication({ dr: 'primary', performance: 'primary' });

      const routes = [
        ['/vault/replication', `${ROUTE_PREFIX}.index`],
        ['/vault/replication/dr', `${ROUTE_PREFIX}.mode.index`],
        ['/vault/replication/performance', `${ROUTE_PREFIX}.mode.index`],
        ['/vault/replication/dr/manage', `${ROUTE_PREFIX}.mode.manage`],
        ['/vault/replication/performance/manage', `${ROUTE_PREFIX}.mode.manage`],
        ['/vault/replication/dr/secondaries', `${ROUTE_PREFIX}.mode.secondaries.index`],
        ['/vault/replication/performance/secondaries', `${ROUTE_PREFIX}.mode.secondaries.index`],
        ['/vault/replication/dr/secondaries/add', `${ROUTE_PREFIX}.mode.secondaries.add`],
        ['/vault/replication/performance/secondaries/add', `${ROUTE_PREFIX}.mode.secondaries.add`],
        ['/vault/replication/dr/secondaries/revoke', `${ROUTE_PREFIX}.mode.secondaries.revoke`],
        ['/vault/replication/performance/secondaries/revoke', `${ROUTE_PREFIX}.mode.secondaries.revoke`],
        [
          '/vault/replication/performance/secondaries/config/show/filtered-secondary',
          `${ROUTE_PREFIX}.mode.secondaries.config-show`,
        ],
        [
          '/vault/replication/performance/secondaries/config/edit/filtered-secondary',
          `${ROUTE_PREFIX}.mode.secondaries.config-edit`,
        ],
      ];

      for (const [url, routeName] of routes) {
        await visit(url);
        assertRoute(assert, url, routeName, url);
      }
    });

    test('route names generate the expected urls', async function (assert) {
      await this.mockReplication({ dr: 'primary', performance: 'primary' });
      // engine routes are only known to the router once the engine has loaded
      await visit('/vault/replication');
      const router = this.owner.lookup('service:router');
      const urlFor = (...args) => router.urlFor(...args).replace(/^\/ui/, '');

      assert.strictEqual(urlFor(`${ROUTE_PREFIX}.index`), '/vault/replication');
      assert.strictEqual(urlFor(`${ROUTE_PREFIX}.mode.index`, 'dr'), '/vault/replication/dr');
      assert.strictEqual(urlFor(`${ROUTE_PREFIX}.mode.manage`, 'dr'), '/vault/replication/dr/manage');
      assert.strictEqual(
        urlFor(`${ROUTE_PREFIX}.mode.secondaries.index`, 'performance'),
        '/vault/replication/performance/secondaries'
      );
      assert.strictEqual(
        urlFor(`${ROUTE_PREFIX}.mode.secondaries.add`, 'performance'),
        '/vault/replication/performance/secondaries/add'
      );
      assert.strictEqual(
        urlFor(`${ROUTE_PREFIX}.mode.secondaries.revoke`, 'dr'),
        '/vault/replication/dr/secondaries/revoke'
      );
      assert.strictEqual(
        urlFor(`${ROUTE_PREFIX}.mode.secondaries.config-show`, 'performance', 'my-secondary'),
        '/vault/replication/performance/secondaries/config/show/my-secondary'
      );
      assert.strictEqual(
        urlFor(`${ROUTE_PREFIX}.mode.secondaries.config-edit`, 'performance', 'my-secondary'),
        '/vault/replication/performance/secondaries/config/edit/my-secondary'
      );
      assert.strictEqual(
        urlFor(`${ROUTE_PREFIX}.mode.secondaries.config-create`, 'performance', 'my-secondary'),
        '/vault/replication/performance/secondaries/config/create/my-secondary'
      );
    });

    test('the replication tabs link to the mode routes', async function (assert) {
      await this.mockReplication({ dr: 'primary', performance: 'primary' });
      await visit('/vault/replication/dr');

      await click(s.tab('manage'));
      assertRoute(assert, '/vault/replication/dr/manage', `${ROUTE_PREFIX}.mode.manage`, 'manage tab');
      await click(s.tab('secondaries'));
      assertRoute(
        assert,
        '/vault/replication/dr/secondaries',
        `${ROUTE_PREFIX}.mode.secondaries.index`,
        'secondaries tab'
      );
      await click(s.tab('details'));
      assertRoute(assert, '/vault/replication/dr', `${ROUTE_PREFIX}.mode.index`, 'details tab');
    });

    test('the sidebar links to the replication routes', async function (assert) {
      await this.mockReplication({ dr: 'primary', performance: 'primary' });
      await visit('/vault/replication');

      await click(GENERAL.navLink('Disaster recovery'));
      assertRoute(assert, '/vault/replication/dr', `${ROUTE_PREFIX}.mode.index`, 'disaster recovery link');
      await click(GENERAL.navLink('Performance replication'));
      assertRoute(
        assert,
        '/vault/replication/performance',
        `${ROUTE_PREFIX}.mode.index`,
        'performance replication link'
      );
      await click(GENERAL.navLink('Overview'));
      assertRoute(assert, '/vault/replication', `${ROUTE_PREFIX}.index`, 'overview link');
    });
  });

  module('landing pages', function () {
    test('index shows the enable view when replication is disabled', async function (assert) {
      await this.mockReplication();
      await visit('/vault/replication');

      assertRoute(assert, '/vault/replication', `${ROUTE_PREFIX}.index`, 'index');
      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Enable replication');
      assert.dom(s.enableForm).exists();
    });

    test('index shows the unsupported state without replication support', async function (assert) {
      await this.mockReplication({ dr: 'unsupported' });
      await visit('/vault/replication');

      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Replication unsupported');
      assert.dom(s.emptyTitle).hasText('The current cluster configuration does not support replication');
      assert.dom(GENERAL.navLink('Disaster recovery')).doesNotExist();
    });

    test('index shows the overview when both modes are primary', async function (assert) {
      await this.mockReplication({ dr: 'primary', performance: 'primary' });
      await visit('/vault/replication');

      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Disaster recovery and performance');
      assert.dom(s.summaryCard).exists({ count: 2 });
    });

    test('mode index shows the enable form for a disabled mode and the dashboard for an enabled one', async function (assert) {
      await this.mockReplication({ dr: 'primary' });

      await visit('/vault/replication/performance');
      assert.dom(s.enableForm).exists('disabled performance mode shows the enable form');
      assert.dom(s.tab('manage')).doesNotExist('tabs are hidden until replication is enabled');

      await visit('/vault/replication/dr');
      await waitFor(s.dashboard);
      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Disaster recovery');
      assert.dom(GENERAL.badge('primary')).exists();
      assert.dom(s.tab('details')).exists();
      assert.dom(s.tab('manage')).exists();
      assert.dom(s.tab('secondaries')).exists('primary clusters have a secondaries tab');
    });

    test('secondary clusters do not get a secondaries tab', async function (assert) {
      await this.mockReplication({ performance: 'secondary' });
      await visit('/vault/replication/performance');

      assert.dom(GENERAL.badge('secondary')).exists();
      assert.dom(s.tab('manage')).exists();
      assert.dom(s.tab('secondaries')).doesNotExist();
    });
  });

  module('redirects', function () {
    test('an unsupported replication mode redirects to the index', async function (assert) {
      await this.mockReplication({ dr: 'primary', performance: 'primary' });
      await visit('/vault/replication/not-a-mode');

      assertRoute(assert, '/vault/replication', `${ROUTE_PREFIX}.index`, 'unknown mode');
    });

    test('manage and secondaries redirect to the mode details when the mode is disabled', async function (assert) {
      await this.mockReplication({ dr: 'primary' });

      await visit('/vault/replication/performance/manage');
      assertRoute(
        assert,
        '/vault/replication/performance',
        `${ROUTE_PREFIX}.mode.index`,
        'manage on a disabled mode'
      );

      await visit('/vault/replication/performance/secondaries/add');
      assertRoute(
        assert,
        '/vault/replication/performance',
        `${ROUTE_PREFIX}.mode.index`,
        'add on a disabled mode'
      );
    });

    test('add secondary redirects to the mode details on a secondary cluster', async function (assert) {
      await this.mockReplication({ performance: 'secondary' });
      await visit('/vault/replication/performance/secondaries/add');

      assertRoute(assert, '/vault/replication/performance', `${ROUTE_PREFIX}.mode.index`, 'add secondary');
    });

    test('revoke secondary redirects away on a secondary cluster', async function (assert) {
      await this.mockReplication({ performance: 'secondary' });
      await visit('/vault/replication/performance/secondaries/revoke');

      assertRoute(assert, '/vault/replication/performance', `${ROUTE_PREFIX}.mode.index`, 'revoke secondary');
    });

    test('add and revoke secondary redirect without the capability to run them', async function (assert) {
      await this.mockReplication({
        dr: 'primary',
        deniedCapabilityPaths: ['primary/secondary-token', 'primary/revoke-secondary'],
      });

      await visit('/vault/replication/dr/secondaries/add');
      assertRoute(assert, '/vault/replication/dr', `${ROUTE_PREFIX}.mode.index`, 'add without capability');

      await visit('/vault/replication/dr/secondaries/revoke');
      assertRoute(assert, '/vault/replication', `${ROUTE_PREFIX}.index`, 'revoke without capability');
    });

    test('path filter config edit redirects without the capability to add secondaries', async function (assert) {
      await this.mockReplication({
        performance: 'primary',
        deniedCapabilityPaths: ['primary/secondary-token'],
      });

      await visit('/vault/replication/performance/secondaries/config/edit/filtered-secondary');
      assertRoute(assert, '/vault/replication/performance', `${ROUTE_PREFIX}.mode.index`, 'config edit');
    });

    test('path filter config edit redirects for dr', async function (assert) {
      await this.mockReplication({ dr: 'primary', performance: 'primary' });

      await visit('/vault/replication/dr/secondaries/config/edit/filtered-secondary');
      assertRoute(assert, '/vault/replication/dr', `${ROUTE_PREFIX}.mode.index`, 'dr config edit');
    });

    test('performance routes redirect to the mode details without the performance license feature', async function (assert) {
      await this.mockReplication({ performance: 'primary', features: ['DR Replication'] });

      await visit('/vault/replication/performance/secondaries/config/edit/filtered-secondary');
      assertRoute(assert, '/vault/replication/performance', `${ROUTE_PREFIX}.mode.index`, 'config edit');
      assert
        .dom(GENERAL.navLink('Performance replication'))
        .doesNotExist('the sidebar hides performance replication');
      assert.dom(GENERAL.navLink('Disaster recovery')).exists();
    });

    test('replication redirects to the cluster when the endpoints are redacted', async function (assert) {
      await this.mockReplication({ redacted: true });
      await visit('/vault/replication');

      assert.strictEqual(currentURL(), '/vault/dashboard', 'it lands on the dashboard');
      assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard');
    });
  });

  module('workflows', function () {
    test('enabling a mode as primary posts to the enable endpoint and shows the dashboard', async function (assert) {
      await this.mockReplication();
      await visit('/vault/replication/dr');
      await click(s.enableForm + ' ' + GENERAL.submitButton);
      await waitFor(s.dashboard);
      await settled();

      const [request] = this.state.requests;
      assert.strictEqual(request.method, 'POST');
      assert.strictEqual(request.url, 'sys/replication/dr/primary/enable');
      assert.dom(GENERAL.badge('primary')).exists('it shows the primary badge once enabled');
      assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.mode.index`);
    });

    test('disabling a mode from the manage page posts to the disable endpoint and returns to the index', async function (assert) {
      await this.mockReplication({ dr: 'primary' });
      await visit('/vault/replication/dr/manage');
      await click(GENERAL.button('disable'));
      await fillIn(`[data-test-confirmation-modal-input="Disable Replication?"]`, 'Disaster recovery');
      await click(`[data-test-confirm-button="Disable Replication?"]`);

      const request = this.state.requests.find((r) => r.url.endsWith('/disable'));
      assert.strictEqual(request.method, 'POST');
      assert.strictEqual(request.url, 'sys/replication/dr/primary/disable');
      assertRoute(assert, '/vault/replication', `${ROUTE_PREFIX}.index`, 'after disabling');
    });

    test('the secondaries list shows known secondaries and links to add and revoke', async function (assert) {
      await this.mockReplication({ dr: 'primary' });
      await visit('/vault/replication/dr/secondaries');

      assert.dom(s.secondaryName('4')).exists('it lists the known secondary');
      await click(s.secondaryAdd);
      assertRoute(
        assert,
        '/vault/replication/dr/secondaries/add',
        `${ROUTE_PREFIX}.mode.secondaries.add`,
        'add link'
      );
    });

    test('adding a secondary generates a token and returns to the secondaries list', async function (assert) {
      await this.mockReplication({ dr: 'primary' });
      await visit('/vault/replication/dr/secondaries/add');
      await fillIn('[data-test-input="Secondary ID"]', 'new-secondary');
      await click(GENERAL.submitButton);

      const request = this.state.requests.find((r) => r.url.endsWith('secondary-token'));
      assert.strictEqual(request.url, 'sys/replication/dr/primary/secondary-token');
      assert.strictEqual(request.body.id, 'new-secondary', 'it sends the secondary id');
      await waitFor('#replication-copy-token-modal');
      assert.dom('#replication-copy-token-modal').exists('it shows the token modal');

      await click(GENERAL.button('Copy token'));
      await click(GENERAL.cancelButton);
      assertRoute(
        assert,
        '/vault/replication/dr/secondaries',
        `${ROUTE_PREFIX}.mode.secondaries.index`,
        'after closing the token modal'
      );
    });

    test('the revoke page asks for the secondary id and a confirmation', async function (assert) {
      await this.mockReplication({ dr: 'primary' });
      await visit('/vault/replication/dr/secondaries/revoke');

      assertRoute(
        assert,
        '/vault/replication/dr/secondaries/revoke',
        `${ROUTE_PREFIX}.mode.secondaries.revoke`,
        'revoke'
      );
      await fillIn('[data-test-input="Secondary ID"]', '4');
      await click('[data-test-confirm-action-trigger]');
      assert.dom(GENERAL.confirmButton).exists('it asks for confirmation before revoking');
      await click(GENERAL.cancelButton);
      assert.strictEqual(this.state.requests.length, 0, 'cancelling does not call the api');
    });

    test('path filter config show renders the config and links to edit', async function (assert) {
      await this.mockReplication({ performance: 'primary', features: ALL_FEATURES });

      await visit('/vault/replication/performance/secondaries/config/show/filtered-secondary');
      assert.dom('[data-test-mount-config-mode]').includesText('allow');
      assert.dom('[data-test-mount-config-paths]').includesText('kv-one/');
      await click('[data-test-replication-link="edit-mount-config"]');
      assertRoute(
        assert,
        '/vault/replication/performance/secondaries/config/edit/filtered-secondary',
        `${ROUTE_PREFIX}.mode.secondaries.config-edit`,
        'edit link'
      );
    });

    test('path filter config show prompts to create a config when there is none', async function (assert) {
      await this.mockReplication({ performance: 'primary', features: ALL_FEATURES });

      await visit('/vault/replication/performance/secondaries/config/show/unfiltered-secondary');
      assert.dom(s.emptyTitle).hasText('No config yet');
      assert
        .dom('[data-test-replication-link="create-mount-config"]')
        .hasAttribute(
          'href',
          /\/vault\/replication\/performance\/secondaries\/config\/create\/unfiltered-secondary$/
        );
    });

    test('removing the path filter config deletes it and returns to the secondaries list', async function (assert) {
      await this.mockReplication({ performance: 'primary' });
      await visit('/vault/replication/performance/secondaries/config/edit/filtered-secondary');
      await click('[data-test-filter-mode="none"]');
      await click('[data-test-config-save]');
      await waitUntil(() => currentURL().endsWith('/secondaries'));

      const request = this.state.requests.find((r) => r.method === 'DELETE');
      assert.strictEqual(request.url, 'sys/replication/performance/primary/paths-filter/filtered-secondary');
      assertRoute(
        assert,
        '/vault/replication/performance/secondaries',
        `${ROUTE_PREFIX}.mode.secondaries.index`,
        'after deleting the config'
      );
    });
  });
});
