/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { fetchRoleCapabilities } from 'ldap/utils/capabilities-helper';

import { ModelFrom } from 'vault/vault/route';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type CapabilitiesService from 'vault/services/capabilities';
import type { LdapStaticRole, LdapDynamicRole } from 'vault/secrets/ldap';

export type LdapRolesRoleRouteModel = ModelFrom<LdapRolesRoleRoute>;

export default class LdapRolesRoleRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async fetchCapabilities(backend: string, name: string, roleType: string) {
    const { pathFor } = this.capabilities;
    const pathType = roleType === 'static' ? 'Static' : 'Dynamic';

    const paths: { role: string; rotate?: string; creds: string } = {
      role: pathFor(`ldap${pathType}Role`, { backend, name }),
      creds: this.capabilities.pathFor(`ldap${pathType}RoleCreds`, { backend, name }),
    };
    if (roleType === 'static') {
      paths.rotate = pathFor('ldapRotateStaticRole', { backend, name });
    }

    const capabilities = await this.capabilities.fetch(Object.values(paths));
    return {
      canDelete: capabilities[paths.role]?.canDelete,
      canEdit: capabilities[paths.role]?.canUpdate,
      canReadCreds: capabilities[paths.creds]?.canRead,
      canRotateStaticCreds: paths.rotate ? capabilities[paths.rotate]?.canCreate : false,
    };
  }

  async model(params: { name: string; type: 'static' | 'dynamic' }) {
    const backend = this.secretMountPath.currentPath;
    const { name, type } = params;

    const [capabilities] = await fetchRoleCapabilities(this.capabilities, backend, [
      { name, completeRoleName: name, type },
    ]);
    const { data } =
      type === 'static'
        ? await this.api.secrets.ldapReadStaticRole(name, backend)
        : await this.api.secrets.ldapReadDynamicRole(name, backend);

    const role = {
      name,
      type,
      completeRoleName: name,
      ...(data || {}),
    } as LdapStaticRole | LdapDynamicRole;

    return { capabilities, role };
  }
}
