import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'ember-qunit';
import Pretender from 'pretender';
import authPage from 'vault/tests/pages/auth';
import { addMonths, format, formatRFC3339, startOfMonth, subMonths } from 'date-fns';

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
function generateActivityResponse(nsCount = 1, startDate, endDate) {
  if (nsCount === 0) {
    return {
      id: 'some-activity-id',
      data: {
        start_time: formatRFC3339(startDate),
        end_time: formatRFC3339(endDate),
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
      start_time: formatRFC3339(startDate),
      end_time: formatRFC3339(endDate),
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
function generateLicenseResponse(startDate, endDate) {
  return {
    request_id: 'my-license-request-id',
    data: {
      autoloaded: {
        license_id: 'my-license-id',
        start_time: formatRFC3339(startDate),
        expiration_time: formatRFC3339(endDate),
      },
    },
  };
}
function send(data, httpStatus = 200) {
  return [httpStatus, { 'Content-Type': 'application/json' }, JSON.stringify(data)];
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

  test('visiting /clients with zero data and config on', async function (assert) {
    const licenseStart = startOfMonth(subMonths(new Date(), 6));
    const licenseEnd = addMonths(new Date(), 6);
    const config = generateConfigResponse();
    const monthly = generateMonthlyResponse();
    const activity = generateActivityResponse(0, licenseStart, licenseEnd);
    const license = generateLicenseResponse(licenseStart, licenseEnd);
    this.server = new Pretender(function () {
      this.get('/v1/sys/health', this.passthrough);
      this.get('/v1/sys/seal-status', this.passthrough);
      this.post('/v1/sys/capabilities-self', this.passthrough);
      this.get('/v1/sys/license/status', () => send(license));
      this.get('/v1/sys/internal/counters/activity', () => send(activity));
      this.get('/v1/sys/internal/counters/activity/monthly', () => send(monthly));
      this.get('/v1/sys/internal/counters/config', () => send(config));
      this.get('/v1/sys/version-history', () => send({ keys: [] }));
    });
    await visit('/vault/clients');
    assert.equal(currentURL(), '/vault/clients');
    assert.dom('.nav-tab-link.is-active').hasText('Current month', 'current month tab is active');
    assert.dom('[data-test-usage-stats]').exists('usage stats block exists');
    assert.dom('[data-test-stat-text-container]').exists({ count: 3 }, '3 stat texts exist');
    assert.dom('[data-test-clients-attribution]').doesNotExist('Does not show attribution with no data');
    // History Tab, zero data
    await visit('/vault/clients/history');
    assert.equal(currentURL(), '/vault/clients/history');
    assert.dom('.nav-tab-link.is-active').hasText('History', 'history tab is active');
    assert
      .dom('[data-test-date-display]')
      .hasText(format(licenseStart, 'MMMM yyyy'), 'billing start month is correctly parsed from license');
    assert.dom('[data-test-clients-attribution]').exists('Attribution block is shown');
    // TODO: Export attribution data not shown
    assert.dom('[data-test-clients-attribution] [data-test-empty-state-title]').hasText('No data received');
    // TODO: Attribution empty state text is correct
    // TODO: Filters correct
  });
});
