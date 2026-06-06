/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import {
  SecretsApiLdapLibraryListListEnum,
  SecretsApiLdapLibraryListLibraryPathListEnum,
} from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type FlashMessageService from 'ember-cli-flash/services/flash-messages';
import type CapabilitiesService from 'vault/services/capabilities';

// Base class for libraries/index and libraries/subdirectory routes
export default class LdapLibrariesRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly capabilities: CapabilitiesService;

  async fetchLibrariesAndCapabilities(path_to_library?: string) {
    try {
      const { currentPath } = this.secretMountPath;
      const { keys } = path_to_library
        ? await this.api.secrets.ldapLibraryListLibraryPath(
            path_to_library,
            currentPath,
            SecretsApiLdapLibraryListLibraryPathListEnum.TRUE
          )
        : await this.api.secrets.ldapLibraryList(currentPath, SecretsApiLdapLibraryListListEnum.TRUE);

      const libraries =
        keys?.map((name) => {
          // if path is provided combine with name for completeLibraryName
          const completeLibraryName = path_to_library ? `${path_to_library}${name}` : name;
          return { name, completeLibraryName };
        }) || [];
      // fetch capabilities for each path
      const paths = libraries.map(({ completeLibraryName: name }) =>
        this.capabilities.pathFor('ldapLibrary', { backend: currentPath, name })
      );
      const capabilities = await this.capabilities.fetch(paths);
      return { libraries, capabilities };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status !== 404) {
        throw error;
      }
      return { libraries: [], capabilities: {} };
    }
  }
}
