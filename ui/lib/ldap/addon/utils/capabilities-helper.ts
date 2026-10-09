/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type CapabilitiesService from 'vault/services/capabilities';
import type PermissionsService from 'vault/services/permissions';
import type { LdapRole } from 'vault/vault/secrets/ldap';

export async function fetchRoleCapabilities(
  capabilities: CapabilitiesService,
  backend: string,
  roles: LdapRole[]
) {
  const { pathFor } = capabilities;

  const paths = roles.map(({ completeRoleName: name, type }) => {
    const pathType = type === 'static' ? 'Static' : 'Dynamic';

    const pathMap: { role: string; rotate?: string; creds: string } = {
      role: pathFor(`ldap${pathType}Role`, { backend, name }),
      creds: pathFor(`ldap${pathType}RoleCreds`, { backend, name }),
    };
    if (type === 'static') {
      pathMap.rotate = pathFor('ldapRotateStaticRole', { backend, name });
    }
    return pathMap;
  });

  // flatten paths array to pass into fetch method
  const allPaths = paths.map((pathMap) => Object.values(pathMap)).flat();
  const perms = await capabilities.fetch(allPaths);

  // map permissions back to array of objects
  // when used in the list view, the array order will be the same as the roles input array
  // index of each loop corresponds to the same index in capabilities array
  return paths.map((pathMap) => {
    return {
      canDelete: perms[pathMap.role]?.canDelete,
      canEdit: perms[pathMap.role]?.canUpdate,
      canReadCreds: perms[pathMap.creds]?.canRead,
      canRotateStaticCreds: pathMap.rotate ? perms[pathMap.rotate]?.canUpdate : false,
    };
  });
}

// Creating a role is a create on <mount>/static-role/:name or <mount>/role/:name. The name is not chosen
// yet, so check for create anywhere beneath either path (policies may be scoped, e.g. team-*).
// Self-managed mounts only support static roles.
export function canCreateRole(permissions: PermissionsService, backend: string, isSelfManaged = false) {
  const paths = isSelfManaged ? ['static-role'] : ['static-role', 'role'];
  return paths.some((path) => permissions.hasPermissionBeneath(`${backend}/${path}`, ['create']));
}
