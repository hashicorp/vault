/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type CapabilitiesService from 'vault/services/capabilities';

export async function fetchRoleCapabilities(
  capabilities: CapabilitiesService,
  backend: string,
  roles: Array<{ name: string; type: string }>
) {
  const { pathFor } = capabilities;

  const paths = roles.map(({ name, type }) => {
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
      canRotateStaticCreds: pathMap.rotate ? perms[pathMap.rotate]?.canCreate : false,
    };
  });
}
