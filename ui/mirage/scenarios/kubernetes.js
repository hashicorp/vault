/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server, shouldConfigureRoles = true) {
  server.create('kubernetes-config', { path: 'kubernetes' });
  if (shouldConfigureRoles) {
    server.create('kubernetes-role');
    server.create('kubernetes-role', 'withRoleName');
    server.create('kubernetes-role', 'withRoleRules');
  }
}
