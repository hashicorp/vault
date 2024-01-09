/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import modifyPassthroughResponse from '../helpers/modify-passthrough-response';

export const statuses = [
  'connected',
  'disconnected since 2022-09-21T11:25:02.196835-07:00; error: unable to establish a connection with HCP',
  'connecting since 2022-09-21T11:25:02.196835-07:00; error: unable to establish a connection with HCP',
  'connecting since 2022-09-21T11:25:02.196835-07:00; error: principal does not have the permission to register as a provider',
  'connecting since 2022-09-21T11:25:02.196835-07:00; error: could not obtain a token with the supplied credentials',
];
let index = null;

export default function (server) {
  server.get('sys/seal-status', (schema, req) => {
    // return next status from statuses array
    if (index === null || index === statuses.length - 1) {
      index = 0;
    } else {
      index++;
    }
    return modifyPassthroughResponse(req, { hcp_link_status: statuses[index] });
  });
  // enterprise only feature initially
  server.get('sys/health', (schema, req) => modifyPassthroughResponse(req, { version: '1.12.0-dev1+ent' }));
}
