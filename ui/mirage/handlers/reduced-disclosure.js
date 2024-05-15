/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import modifyPassthroughResponse from '../helpers/modify-passthrough-response';
import { Response } from 'miragejs';

export default function (server) {
  server.get('/sys/health', (schema, req) =>
    modifyPassthroughResponse(req, { version: '', cluster_name: '' })
  );
  server.get('/sys/seal-status', (schema, req) => {
    // When reduced disclosure is active, the version is only returned when a valid token is used
    const overrides = req.requestHeaders['X-Vault-Token']
      ? { cluster_name: '', build_date: '' }
      : { version: '', cluster_name: '', build_date: '' };
    return modifyPassthroughResponse(req, overrides);
  });
  server.get('sys/replication/status', () => new Response(404, {}, { errors: ['disabled path'] }));
  server.get('sys/replication/dr/status', () => new Response(404, {}, { errors: ['disabled path'] }));
  server.get(
    'sys/replication/performance/status',
    () => new Response(404, {}, { errors: ['disabled path'] })
  );
}
