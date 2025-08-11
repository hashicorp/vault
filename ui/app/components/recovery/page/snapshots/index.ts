/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import type CapabilitiesService from 'vault/services/capabilities';
import type NamespaceService from 'vault/services/namespace';
import type PermissionsService from 'vault/services/permissions';
import type VersionService from 'vault/services/version';

enum State {
  COMMUNITY_VERSION = 'community',
  NON_ROOT_NAMESPACE = 'non-root-namespace',
  ALLOW_UPLOAD = 'default',
  CANNOT_UPLOAD = 'cannot-upload',
}

interface Args {
  model: { canLoadSnapshot: boolean; snapshots: Record<string, unknown>[] };
}

export default class Index extends Component<Args> {
  @service declare readonly version: VersionService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly permissions: PermissionsService;
  @service declare readonly capabilities: CapabilitiesService;

  viewState = State;

  emptyStateDetails = {
    [this.viewState.COMMUNITY_VERSION]: {
      title: 'Secrets Recovery is an enterprise feature',
      icon: 'sync-reverse',
      message:
        'Secrets Recovery allows you to restore accidentally deleted or lost secrets from a snapshot. The snapshots can be provided via upload or loaded from external storage.',
      buttonText: 'Learn more about upgrading',
      buttonHref: '/vault/docs/enterprise',
      buttonIcon: 'docs-link',
      buttonColor: 'tertiary',
    },
    [this.viewState.NON_ROOT_NAMESPACE]: {
      title: 'Snapshot upload is restricted',
      icon: 'sync-reverse',
      message:
        'Snapshot uploading is only available in root namespace. Please navigate to root and upload your snapshot.',
      buttonText: 'Take me to root namespace',
      buttonRoute: 'vault.cluster.dashboard',
      buttonIcon: 'arrow-right',
      buttonColor: 'tertiary',
    },
    [this.viewState.CANNOT_UPLOAD]: {
      title: 'No snapshot available',
      icon: 'skip',
      message:
        'Ready to restore secrets? Please contact your admin to either upload a snapshot or grant you uploading permissions to get started.',
      buttonText: 'Learn more about Secrets Recovery',
      buttonHref: '/vault/docs/sysadmin/snapshots/restore',
      buttonIcon: 'docs-link',
      buttonColor: 'tertiary',
    },
    [this.viewState.ALLOW_UPLOAD]: {
      title: 'Upload a snapshot to get started',
      icon: 'sync-reverse',
      message:
        'Secrets Recovery allows you to restore accidentally deleted or lost secrets from a snapshot. The snapshots can be provided via upload or loaded from external storage.',
      buttonText: 'Upload snapshot',
      buttonColor: 'primary',
    },
  };
  get loadedSnapshots() {
    const { snapshots } = this.args.model;

    return snapshots;
  }

  get state() {
    const { canLoadSnapshot } = this.args.model;

    if (this.version.isCommunity) {
      return this.viewState.COMMUNITY_VERSION;
    } else if (!this.namespace.inRootNamespace) {
      return this.viewState.NON_ROOT_NAMESPACE;
    } else if (!canLoadSnapshot) {
      return this.viewState.CANNOT_UPLOAD;
    } else {
      return this.viewState.ALLOW_UPLOAD;
    }
  }

  @action
  uploadSnapshot() {
    // TODO this will be completed in further feature work for single item recovery in the 1.21 release
  }
}
