/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { assert } from '@ember/debug';

/**
 * @module KvDeleteModal displays a button for a delete type and launches a modal (undelete is the only mode that does not launch the modal).
 *  All destroyRecord network requests are handled inside the component.
 *
 * <KvDeleteModal
 *  @mode="destroy"
 *  @path="my-secret"
 *  @secret={{this.model.secret}}
 *  @metadata={{this.model.metadata}}
 * />
 *
 * @param {string} mode - matches key in serializer: delete, destroy-version or delete-metadata
 * @param {object} secret - The kv/data model.
 * @param {object} [metadata] - The kv/metadata model. It is only required when mode is "delete" or "metadata-delete".
 */

export default class KvDeleteModal extends Component {
  @service flashMessages;
  @service router;
  @service store;
  @tracked deleteType = null;
  @tracked modalOpen = false;

  get modalDisplay() {
    switch (this.args.mode) {
      // does not match serializer key directly because a delete type must be selected
      case 'delete':
        return {
          title: 'Delete version?',
          type: 'warning',
          intro:
            'There are two ways to delete a version of a secret. Both delete actions can be un-deleted later. How would you like to proceed?',
        };
      case 'destroy-version':
        return {
          title: 'Destroy version?',
          type: 'danger',
          intro: `This action will permanently destroy Version ${this.args.secret.version} of the secret, and the secret data cannot be read or recovered later.`,
        };
      case 'delete-metadata':
        return {
          title: 'Delete metadata?',
          type: 'danger',
          intro:
            'This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone.',
        };
      default:
        return assert('mode must be one of delete, destroy-version, delete-metadata.');
    }
  }

  get deleteOptions() {
    const { secret, metadata } = this.args;
    const isDeactivated = secret.canReadMetadata ? metadata?.currentSecret.isDeactivated : false;
    return [
      {
        key: 'delete-version',
        label: 'Delete this version',
        description: `This deletes Version ${secret.version} of the secret.`,
        disabled: !secret.canDeleteVersion,
        tooltipMessage: 'You do not have permission to delete a specific version.',
      },
      {
        key: 'delete-latest-version',
        label: 'Delete latest version',
        description: 'This deletes the most recent version of the secret.',
        disabled: !secret.canDeleteLatestVersion || isDeactivated,
        tooltipMessage: isDeactivated
          ? `The latest version of the secret is already ${metadata.secretState.state}.`
          : 'You do not have permission to delete the latest version of this secret.',
      },
    ];
  }

  @action
  onDelete() {
    const type = this.args.mode === 'delete' ? this.deleteType : this.args.mode;
    this.args.onDelete(type);
  }
}
