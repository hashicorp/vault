import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import authPage from 'vault/tests/pages/auth';

function generateNamespaceBlock(idx = 0, skipMounts = false) {
  let mountCount = 1;
  const nsBlock = {
    namespace_id: `${idx}UUID`,
    namespace_path: `my-namespace-${idx}/`,
    counts: {
      entity_clients: mountCount * 15,
      non_entity_clients: mountCount * 10,
      clients: mountCount * 5,
    },
  };
  if (!skipMounts) {
    mountCount = Math.floor((Math.random() + idx) * 20);
    let mounts = [];
    if (!skipMounts) {
      Array.from(Array(mountCount)).forEach((v, index) => {
        mounts.push({
          id: index,
          path: `auth/method/authid${index}`,
          counts: {
            clients: 15,
            entity_clients: 10,
            non_entity_clients: 5,
          },
        });
      });
    }
    nsBlock.mounts = mounts;
  }
  return nsBlock;
}

function generateConfigResponse() {
  return {
    request_id: 'some-config-id',
    data: {
      default_report_months: 12,
      enabled: 'default-enable',
      queries_available: true,
      retention_months: 24,
    },
  };
}
function generateActivityResponse(nsCount = 1) {
  if (nsCount === 0) {
    return {
      id: 'some-activity-id',
      data: {
        start_time: '2021-03-17T00:00:00Z',
        end_time: '2021-12-31T23:59:59Z',
        total: {
          clients: 0,
          entity_clients: 0,
          non_entity_clients: 0,
        },
        by_namespace: [
          {
            namespace_id: `root`,
            namespace_path: '',
            counts: {
              entity_clients: 0,
              non_entity_clients: 0,
              clients: 0,
            },
          },
        ],
        // months: [],
      },
    };
  }
  let namespaces = Array.from(Array(nsCount)).map((v, idx) => {
    return generateNamespaceBlock(idx);
  });
  console.log({ namespaces });
  return {
    id: 'some-activity-id',
    data: {
      start_time: '2021-03-17T00:00:00Z',
      end_time: '2021-12-31T23:59:59Z',
      total: {
        clients: 3637,
        entity_clients: 1643,
        non_entity_clients: 1994,
      },
      by_namespace: namespaces,
      // months: [],
    },
  };
}
function generateMonthlyResponse(namespaces) {
  if (!namespaces) {
    return {
      request_id: 'monthly-response-id',
      data: {
        by_namespace: [],
        clients: 0,
        entity_clients: 0,
        non_entity_clients: 0,
      },
    };
  }
  // generate by_namespace data
  const by_namespace = namespaces.map((ns, idx) => generateNamespaceBlock(idx));
  const counts = by_namespace.reduce(
    (prev, curr) => ({
      clients: prev.clients + curr.clients,
      entity_clients: prev.entity_clients + curr.entity_clients,
      non_entity_clients: prev.non_entity_clients + curr.non_entity_clients,
    }),
    { clients: 0, entity_clients: 0, non_entity_clients: 0 }
  );
  return {
    request_id: 'monthly-response-id',
    data: {
      by_namespace,
      ...counts,
    },
  };
}
// const RES = {
//   permissionDenied: [403, {}, JSON.stringify({ errors: ['permission denied'] })],
//   empty: [204, {}],
//   empty200: [200, {}, JSON.stringify({ id: 'empty', data: {} })],
// };
module('Acceptance | clients route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  test('visiting /clients with zero data', async function (assert) {
    const configResponse = generateConfigResponse();
    const activity = generateActivityResponse(0);
    const monthly = generateMonthlyResponse();
    this.server = new Pretender(function () {
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/license/status', this.passthrough);
      this.get('/v1/sys/internal/counters/activity', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify(activity)];
      });
      this.get('/v1/sys/internal/counters/activity/monthly', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify(monthly)];
      });
      this.get('/v1/sys/internal/counters/config', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify(configResponse)];
      });
      this.get('/v1/sys/version-history', () => {
        return [200, { 'Content-Type': 'application/json' }, JSON.stringify({ keys: [] })];
      });
    });
    await visit('/vault/clients');
    assert.equal(currentURL(), '/vault/clients');
    assert.dom('.nav-tab-link.is-active').hasText('Current month', 'current month tab is active');
    assert.dom('[data-test-usage-stats]').exists('usage stats block exists');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    assert.dom('[data-test-clients-attribution]').doesNotExist('Does not show attribution with no data');
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
  });
});
