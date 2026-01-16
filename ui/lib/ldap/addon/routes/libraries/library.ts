/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type CapabilitiesService from 'vault/services/capabilities';
import type { LdapLibrary } from 'vault/vault/secrets/ldap';

export type LdapLibraryRouteModel = ModelFrom<LdapLibraryRoute>;

interface LdapLibraryRouteParams {
  name?: string;
}

export default class LdapLibraryRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async model(params: LdapLibraryRouteParams) {
    const backend = this.secretMountPath.currentPath;
    const { name } = params;
    // Decode URL-encoded hierarchical paths (e.g., "service-account1%2Fsa1" -> "service-account1/sa1")
    const decodedName = decodeURIComponent(name || '');
    const { data } = await this.api.secrets.ldapLibraryRead(decodedName, backend);

    // If the decoded name contains a slash, it's hierarchical
    let libraryName = decodedName;
    if (decodedName.includes('/')) {
      const lastSlashIndex = decodedName.lastIndexOf('/');
      libraryName = decodedName.substring(lastSlashIndex + 1);
    }
    const library = {
      ...(data as object),
      name: libraryName,
      completeLibraryName: decodedName,
    } as LdapLibrary;

    // fetch capabilities for this library
    const { pathFor } = this.capabilities;
    const paths = [
      pathFor('ldapLibrary', { backend, name: library.completeLibraryName }),
      pathFor('ldapLibraryCheckOut', { backend, name: library.completeLibraryName }),
      pathFor('ldapLibraryCheckIn', { backend, name: library.completeLibraryName }),
    ];
    const capabilities = await this.capabilities.fetch(paths);

    return { library, capabilities };
  }
}
