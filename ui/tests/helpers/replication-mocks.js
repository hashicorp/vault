/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';
import { mockReplicationBlock } from 'vault/tests/helpers/replication';

export const ALL_FEATURES = ['DR Replication', 'Performance Replication'];

// A JWT-shaped activation token. The UI only decodes the payload to read the expiry.
export const buildActivationToken = () => {
  const b64 = (obj) => btoa(JSON.stringify(obj)).replace(/=+$/, '');
  return `${b64({ alg: 'ES512' })}.${b64({ exp: 4102444800, addr: 'http://127.0.0.1:8200' })}.signature`;
};

/**
 * Mocks everything the replication pages need so the tests behave the same against a
 * Community or an Enterprise backend. The mocked health response makes the UI treat the cluster as
 * Enterprise and the license features decide which replication modes are available.
 *
 * @param {object} server mirage server
 * @param {object} options
 * @param {string} options.dr dr mode: disabled | primary | secondary | unsupported
 * @param {string} options.performance performance mode: disabled | primary | secondary
 * @param {string[]} options.features license features
 * @param {boolean} options.redacted respond to the replication status like a cluster with redacted endpoints
 * @param {string[]} options.deniedCapabilityPaths replication paths the token cannot update
 * @returns {object} state, with `requests` recording the replication calls made by the UI
 */
export function mockReplication(server, options = {}) {
  const { dr = 'disabled', performance = 'disabled', features = ALL_FEATURES } = options;
  const { deniedCapabilityPaths = [] } = options;
  const state = {
    unsupported: dr === 'unsupported',
    dr: mockReplicationBlock(dr),
    performance: mockReplicationBlock(performance),
    pathFilters: { 'filtered-secondary': { mode: 'allow', paths: ['kv-one/', 'kv-two/'] } },
    requests: [],
  };
  const record = (request) => {
    const body = request.requestBody ? JSON.parse(request.requestBody) : null;
    state.requests.push({ method: request.method, url: request.url.replace(/^\/v1\//, ''), body });
    return body;
  };

  server.get('/sys/health', () => ({
    initialized: true,
    sealed: false,
    standby: false,
    performance_standby: false,
    replication_performance_mode: state.performance.mode,
    replication_dr_mode: state.dr.mode,
    server_time_utc: 1754948244,
    version: '1.21.0+ent',
    enterprise: true,
    cluster_name: 'vault-cluster-replication',
    cluster_id: '113a6c47-077f-bea7-0e8e-70a91821e85a',
    last_wal: 82,
    license: { state: 'autoloaded', expiry_time: '2099-01-27T00:00:00Z', terminated: false },
    echo_duration_ms: 0,
    clock_skew_ms: 0,
    replication_primary_canary_age_ms: 0,
    removed_from_cluster: false,
  }));
  // the path filter form reads the accessible namespaces
  server.get('/sys/internal/ui/namespaces', () => ({ data: { keys: [] } }));
  server.get('/sys/license/features', () => ({ features }));
  server.get('/sys/replication/status', () => {
    // vault redacts the replication endpoints on some clusters, eg. HVD
    if (options.redacted) return new Response(403, {}, { errors: ['disabled path'] });
    return {
      data: state.unsupported ? { mode: 'unsupported' } : { dr: state.dr, performance: state.performance },
    };
  });

  server.post('/sys/capabilities-self', (schema, request) => {
    const { paths } = JSON.parse(request.requestBody);
    const caps = (path) => (deniedCapabilityPaths.some((p) => path.includes(p)) ? ['deny'] : ['root']);
    const data = paths.reduce((acc, path) => ({ ...acc, [path]: caps(path) }), {});
    return { request_id: 'capabilities', data: { ...data, capabilities: caps(paths[0]) } };
  });

  ['dr', 'performance'].forEach((mode) => {
    server.post(`/sys/replication/${mode}/primary/enable`, (schema, request) => {
      record(request);
      state[mode] = mockReplicationBlock('primary');
      return {};
    });
    server.post(`/sys/replication/${mode}/primary/disable`, (schema, request) => {
      record(request);
      state[mode] = mockReplicationBlock('disabled');
      return {};
    });
    server.post(`/sys/replication/${mode}/primary/secondary-token`, (schema, request) => {
      record(request);
      return {
        wrap_info: {
          token: buildActivationToken(),
          ttl: 1800,
          creation_time: '2026-01-01T00:00:00Z',
          creation_path: `sys/replication/${mode}/primary/secondary-token`,
        },
      };
    });
    server.post(`/sys/replication/${mode}/primary/revoke-secondary`, (schema, request) => {
      record(request);
      return {};
    });
  });

  server.get('/sys/replication/performance/primary/paths-filter/:id', (schema, request) => {
    const config = state.pathFilters[request.params.id];
    if (!config) return new Response(404, {}, { errors: [] });
    return { data: config };
  });
  server.put('/sys/replication/performance/primary/paths-filter/:id', (schema, request) => {
    const body = record(request);
    state.pathFilters[request.params.id] = { mode: body.mode, paths: body.paths };
    return {};
  });
  server.delete('/sys/replication/performance/primary/paths-filter/:id', (schema, request) => {
    record(request);
    delete state.pathFilters[request.params.id];
    return new Response(204);
  });

  return state;
}
