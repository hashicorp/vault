/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { keyIsFolder } from 'core/utils/key-utils';
import { SECRET_TYPE_CONFIGS, getSecretTypeFromMount } from 'sync/utils/secret-type-config';

import type { Destination, SecretType } from 'vault/sync';
import type { MountOption } from 'sync/utils/secret-type-config';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  destination: Destination;
}

export default class DestinationSyncPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.fetchMounts();
  }

  @tracked allSupportedMounts: MountOption[] = [];
  @tracked mountPath = '';
  @tracked secretPath = '';
  @tracked error = '';
  @tracked syncedSecret = '';
  @tracked syncedMount = '';

  get selectedMountData(): MountOption | null {
    if (!this.mountPath) return null;
    return this.allSupportedMounts.find((mount) => mount.id === this.mountPath) || null;
  }

  get detectedSecretType(): SecretType | null {
    if (!this.selectedMountData) return null;
    return getSecretTypeFromMount(this.selectedMountData.engineType, this.selectedMountData.version);
  }

  get currentSecretTypeConfig() {
    // KV is the default when mount type cannot be determined
    const secretType = this.detectedSecretType || 'kv';
    return SECRET_TYPE_CONFIGS[secretType];
  }

  get syncedSecretTypeConfig() {
    const syncedMountData = this.allSupportedMounts.find((mount) => mount.id === this.syncedMount);
    const secretType = syncedMountData
      ? getSecretTypeFromMount(syncedMountData.engineType, syncedMountData.version)
      : null;
    const config = SECRET_TYPE_CONFIGS[secretType || 'kv'];
    return secretType ? config : { ...config, supportsExternalLink: false };
  }

  get hasSupportedMounts() {
    return this.allSupportedMounts.length > 0;
  }

  get isSecretDirectory() {
    return this.secretPath && keyIsFolder(this.secretPath);
  }

  get isKvSecret() {
    return this.detectedSecretType === 'kv';
  }

  get isSubmitDisabled() {
    if (!this.mountPath || !this.secretPath || this.setAssociation.isRunning) {
      return true;
    }

    // For KV (detected or default), also check that it's not a directory
    if ((this.detectedSecretType === 'kv' || !this.detectedSecretType) && this.isSecretDirectory) {
      return true;
    }

    return false;
  }

  async fetchMounts() {
    const supportedMounts: MountOption[] = [];
    try {
      const { secret } = await this.api.sys.internalUiListEnabledVisibleMounts();
      if (secret) {
        for (const path in secret) {
          const { type, options } = secret[path as keyof typeof secret];
          const version = options?.['version'] ? Number(options['version']) : undefined;

          const secretType = getSecretTypeFromMount(type, version);
          if (secretType) {
            supportedMounts.push({
              name: path,
              id: path,
              engineType: type,
              version,
            });
          }
        }
      }
      this.allSupportedMounts = supportedMounts;
    } catch (error) {
      /*
       * If mount list cannot be fetched (e.g., due to permissions),
       * user can still manually enter mount path via FilterInput
       */
    }
  }

  @action
  handleMountSelection(mount: MountOption | null) {
    this.mountPath = mount?.id || '';
    this.secretPath = '';
    this.syncedSecret = '';
    this.error = '';
  }

  setAssociation = task({}, async (event: Event) => {
    event.preventDefault();
    this.error = '';
    try {
      this.syncedSecret = '';
      const { name, type } = this.args.destination;
      const mount = keyIsFolder(this.mountPath) ? this.mountPath.slice(0, -1) : this.mountPath;

      const payload = { mount, secret_name: this.secretPath };

      await this.api.sys.systemWriteSyncDestinationsTypeNameAssociationsSet(name, type, payload);
      this.syncedSecret = this.secretPath;
      this.syncedMount = this.mountPath;
      this.mountPath = '';
      this.secretPath = '';
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.error = `Sync operation error: \n ${message}`;
    }
  });
}
