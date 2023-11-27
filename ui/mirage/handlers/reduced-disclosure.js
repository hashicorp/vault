/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import modifyPassthroughResponse from '../helpers/modify-passthrough-response';

export default function (server) {
  server.get('/sys/health', (schema, req) =>
    modifyPassthroughResponse(req, { version: '', cluster_name: '' })
  );
  server.get('/sys/seal-status', (schema, req) =>
    modifyPassthroughResponse(req, { version: '', cluster_name: '', build_date: '' })
  );
  server.get('sys/replication/status', () => new Response(404));
  server.get('sys/replication/dr/status', () => new Response(404));
  server.get('sys/replication/performance/status', () => new Response(404));
}
