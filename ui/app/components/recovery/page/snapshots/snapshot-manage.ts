/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { restartableTask } from 'ember-concurrency';
import { sanitizePath } from 'core/utils/sanitize-path';
import SecretsEngineResource, { RecoverySupportedEngines } from 'vault/resources/secrets/engine';
import { SupportedSecretBackendsEnum } from 'vault/helpers/supported-secret-backends';
import { ROOT_NAMESPACE } from 'vault/services/namespace';
import {
  createPollingTask,
  getSnapshotStatusBadge,
} from 'vault/components/recovery/page/snapshots/snapshot-utils';

import type ApiService from 'vault/services/api';
import type AuthService from 'vault/vault/services/auth';
import type NamespaceService from 'vault/services/namespace';
import type { SnapshotManageModel } from 'vault/routes/vault/cluster/recovery/snapshots/snapshot/manage';
import { ResponseError } from '@hashicorp/vault-client-typescript';

interface Args {
  model: SnapshotManageModel;
}

type SecretData = { [key: string]: unknown };

type RecoveryData = {
  models: string[];
  query?: { [key: string]: string };
};

type MountOption = { type: RecoverySupportedEngines; path: string };
type GroupedOption = { groupName: string; options: MountOption[] };

export default class SnapshotManage extends Component<Args> {
  // TODO remove once endpoint is updated to accepted `read_snapshot_id`
  @service declare readonly auth: AuthService;
  @service declare readonly api: ApiService;
  @service declare readonly currentCluster: any;
  @service declare readonly namespace: NamespaceService;

  @tracked selectedNamespace: string;
  @tracked selectedMount?: MountOption;
  @tracked resourcePath = '';
  @tracked copyPath = '';

  @tracked mountOptions: GroupedOption[] = [];
  @tracked secretData: SecretData | undefined;

  @tracked mountError = '';
  @tracked resourcePathError = '';
  @tracked copyPathError = '';
  @tracked bannerError = '';

  @tracked showReadModal = false;
  @tracked showJson = false;
  @tracked recoveryData?: RecoveryData;
  @tracked recoverMethod = 'original';

  @tracked snapshotStatus: string | null = null;

  private pollingController: { start: () => Promise<void>; cancel: () => void } | null = null;

  recoverySupportedEngines = [
    { display: 'Database', value: SupportedSecretBackendsEnum.DATABASE },
    { display: 'Cubbyhole', value: SupportedSecretBackendsEnum.CUBBYHOLE },
    { display: 'KV v1', value: SupportedSecretBackendsEnum.KV },
  ];

  recoverMethods = [
    { label: 'Recover to original path (recover in place)', value: 'original' },
    { label: 'Recover to a new path (copy)', value: 'copy' },
  ];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.selectedNamespace = this.namespace.inRootNamespace ? 'root' : this.namespace.path;
    this.fetchMounts.perform();

    // Create and start polling task
    this.pollingController = createPollingTask(
      this.args.model.snapshot.snapshot_id,
      this.api,
      this.onPollSuccess,
      this.onPollError
    );
    this.pollingController.start();
  }

  willDestroy() {
    super.willDestroy();
    if (this.pollingController) {
      this.pollingController.cancel();
    }
  }

  get hasValidationErrors() {
    return !!(this.mountError || this.resourcePathError);
  }

  // Form secret data to display in accordance with <SecretFormShow/> expectations
  get modelForData() {
    return {
      secretData: this.secretData,
      secretKeyAndValue: this.secretKeyAndValue,
    };
  }

  get namespaceOptions() {
    const { namespaces } = this.args.model;
    // Remove slashes as query param does not expect them
    const sanitized = namespaces.map((ns) => sanitizePath(ns));
    // Add the root namespace because `sys/internal/ui/namespaces` does not include it.
    return ['root', ...sanitized];
  }

  get recoverToCopy() {
    return this.recoverMethod === 'copy';
  }

  get secretKeyAndValue() {
    if (!this.secretData || typeof this.secretData !== 'object') {
      return [];
    }

    return Object.entries(this.secretData).map(([key, value]) => ({
      key,
      value: typeof value === 'string' ? value : JSON.stringify(value),
    }));
  }

  get badge() {
    // Use polled status if available, otherwise fall back to initial model status
    const status = this.snapshotStatus || this.args.model.snapshot?.status;
    return getSnapshotStatusBadge(status);
  }

  onPollError = async (e: unknown) => {
    const error = await this.api.parseError(e);
    this.bannerError = `Snapshot load error: ${error.message}`;
    this.snapshotStatus = 'error';
    // Update the model directly for reactive updates across routes
    this.args.model.snapshot.status = 'error';
  };

  onPollSuccess = async (status: string) => {
    this.snapshotStatus = status;
    // Update the model directly for reactive updates across routes
    this.args.model.snapshot.status = status as 'ready' | 'error' | 'loading';
  };

  fetchMounts = restartableTask(async () => {
    try {
      const namespace = this.selectedNamespace === 'root' ? ROOT_NAMESPACE : this.selectedNamespace;
      const headers = this.api.buildHeaders({ namespace });
      const { secret } = await this.api.sys.internalUiListEnabledVisibleMounts(headers);

      const mounts = this.api.responseObjectToArray(secret, 'path').flatMap((engine) => {
        const eng = new SecretsEngineResource(engine);

        // performance secondaries cannot perform snapshot operations on shared paths
        // operations can still be performed on local mounts
        if (this.currentCluster?.performance?.isSecondary && !eng.local) return [];

        // Use `engineType` as it is the normalized version of `type`
        return eng.supportsRecovery
          ? [
              {
                path: sanitizePath(eng.path),
                type: eng.engineType,
              } as MountOption,
            ]
          : [];
      });

      const databases: MountOption[] = mounts.filter((m) => m.type === SupportedSecretBackendsEnum.DATABASE);
      const secretEngines: MountOption[] = mounts.filter(
        (m) => m.type !== SupportedSecretBackendsEnum.DATABASE
      );

      this.mountOptions = [
        ...(databases.length ? [{ groupName: 'Databases', options: databases }] : []),
        { groupName: 'Secret Engines', options: secretEngines },
      ];
    } catch {
      this.mountOptions = [];
    }
  });

  @action
  clearSelections() {
    this.selectedNamespace = this.namespace.inRootNamespace ? 'root' : this.namespace.path;
    this.selectedMount = undefined;
    this.resourcePath = '';
    this.copyPath = '';
    this.mountError = '';
    this.resourcePathError = '';
    this.secretData = undefined;

    // Refetch mounts after clearing
    if (this.namespace.inRootNamespace) {
      this.fetchMounts.perform();
    }
  }

  @action
  handleSelectNamespace(selection: string) {
    this.selectedMount = undefined;
    this.selectedNamespace = selection as string;
    this.fetchMounts.perform();
  }

  @action
  handleSelectMount(selection: MountOption) {
    this.mountError = '';
    this.selectedMount = selection;

    // Cubbyhole path is always the same, set for manual path selection when user does not have LIST permissions
    if (this.selectedMount.type === SupportedSecretBackendsEnum.CUBBYHOLE) {
      this.selectedMount.path = 'cubbyhole';
    }
  }

  @action
  handleSelectRecoverMethod({ target }: { target: HTMLInputElement }) {
    this.recoverMethod = target.value;
    if (this.recoverToCopy && this.resourcePath) {
      this.copyPath = this.resourcePath + '-copy';
    }
  }

  @action
  updateCopyPath({ target }: { target: HTMLInputElement }) {
    this.copyPath = target.value.trim();
  }

  @action
  updateResourcePath({ target }: { target: HTMLInputElement }) {
    this.resourcePath = target.value.trim();
    this.resourcePathError = '';
  }

  @action
  async readResource() {
    const isValid = this.validateFields();
    if (!isValid) {
      return;
    }
    try {
      this.bannerError = '';
      this.recoveryData = undefined;

      const { snapshot_id } = this.args.model.snapshot as { snapshot_id: string };
      const mountType = this.selectedMount?.type;
      const mountPath = this.selectedMount?.path as string;
      const namespace = this.selectedNamespace === 'root' ? ROOT_NAMESPACE : this.selectedNamespace;
      const headers = this.api.buildHeaders({ namespace });

      switch (mountType) {
        case SupportedSecretBackendsEnum.KV: {
          const { data } = await this.api.secrets.kvV1Read(
            this.resourcePath,
            mountPath,
            snapshot_id,
            headers
          );
          this.secretData = data as SecretData;
          break;
        }
        case SupportedSecretBackendsEnum.CUBBYHOLE: {
          const { data } = await this.api.secrets.cubbyholeRead(this.resourcePath, snapshot_id, headers);
          this.secretData = data as SecretData;
          break;
        }
        case SupportedSecretBackendsEnum.DATABASE: {
          // TODO remove once endpoint is updated to accept `read_snapshot_id`
          const { currentToken } = this.auth;

          const resp = await fetch(
            `/v1/${mountPath}/static-roles/${this.resourcePath}?read_snapshot_id=${snapshot_id}`,
            {
              method: 'GET',
              headers: {
                'X-Vault-Namespace': namespace,
                'X-Vault-Token': currentToken,
              },
            }
          );

          if (!resp.ok) {
            throw new ResponseError(resp);
          }

          const { data } = await resp.json();
          this.secretData = data as SecretData;
          break;
        }
        default: {
          // This should never be reached, but just in case
          throw new Error('Unsupported recovery engine');
        }
      }

      this.showReadModal = true;
    } catch (e) {
      const error = await this.api.parseError(e);
      this.bannerError = `Snapshot read error: ${error.message}`;
    }
  }

  @action
  async recover() {
    const isValid = this.validateFields(true);
    if (!isValid) {
      return;
    }
    try {
      this.bannerError = '';
      const { snapshot_id } = this.args.model.snapshot as { snapshot_id: string };
      const namespace = this.selectedNamespace === 'root' ? ROOT_NAMESPACE : this.selectedNamespace;
      const mountType = this.selectedMount?.type;
      const mountPath = this.selectedMount?.path as string;
      // if recovering to a copy, the new path (the copy input) becomes the target path
      // and the original path (the resource path input) becomes the source path header
      const path = this.recoverToCopy ? this.copyPath : this.resourcePath;
      const recoverSourcePath =
        mountType === SupportedSecretBackendsEnum.DATABASE
          ? mountPath + '/static-roles/' + this.resourcePath
          : mountPath + '/' + this.resourcePath;

      const headers = this.api.buildHeaders({
        namespace,
        recoverSnapshotId: snapshot_id,
        ...(this.recoverToCopy && { recoverSourcePath }),
      });

      // this query is used to build the recovered resource link in the success message
      let query: { [key: string]: string } = {};
      if (namespace && namespace !== this.namespace.path) {
        query = { namespace };
      }

      // Certain backends have a prefix which is needed for the recovery link we show to the user
      let modelPrefix = '';
      switch (mountType) {
        case SupportedSecretBackendsEnum.KV: {
          await this.api.secrets.kvV1Write(path, mountPath, {}, snapshot_id, undefined, headers);
          break;
        }
        case SupportedSecretBackendsEnum.CUBBYHOLE: {
          await this.api.secrets.cubbyholeWrite(path, {}, snapshot_id, undefined, headers);
          break;
        }
        case SupportedSecretBackendsEnum.DATABASE: {
          await this.api.secrets.databaseWriteStaticRole(path, mountPath, {}, headers);
          modelPrefix = 'role/';
          query = {
            ...query,
            type: 'static',
          };
          break;
        }
        default: {
          // This should never be reached, but just in case
          throw new Error('Unsupported recovery engine');
        }
      }

      this.recoveryData = {
        models: [mountPath, modelPrefix + path],
        query,
      };
    } catch (e) {
      const error = await this.api.parseError(e);
      this.bannerError = `Snapshot recovery error: ${error.message}`;
      this.recoveryData = undefined;
    }
  }

  @action
  closeReadModal() {
    this.showReadModal = false;
  }

  @action
  toggleJson(event: { target: { checked: boolean } }) {
    this.showJson = event.target.checked;
  }

  @action
  validateFields(isRecover = false): boolean {
    this.mountError = '';
    this.resourcePathError = '';
    this.copyPathError = '';

    let hasErrors = false;

    if (!this.selectedMount) {
      this.mountError = 'Please select a secret mount';
      hasErrors = true;
    }

    if (!this.resourcePath) {
      this.resourcePathError = 'Please enter a resource path';
      hasErrors = true;
    }

    if (isRecover && this.recoverToCopy && !this.copyPath) {
      this.copyPathError = 'Please enter a copy path';
      hasErrors = true;
    }

    return !hasErrors;
  }
}
