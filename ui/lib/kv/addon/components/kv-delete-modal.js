/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { assert } from '@ember/debug';

/**
 * @module KvDeleteModal displays a button for a delete type and launches a modal. Undelete is the only mode that does not launch the modal and is not handled in this component.
 *
 * <KvDeleteModal
 *  @mode="destroy"
 *  @secret={{this.model.secret}}
 *  @metadata={{this.model.metadata}}
 *  @onDelete={{this.handleDestruction}}
 * />
 *
 * @param {string} mode - delete, delete-metadata, or destroy.
 * @param {object} secret - The kv/data model.
 * @param {object} [metadata] - The kv/metadata model. It is only required when mode is "delete" or "metadata-delete".
 * @param {callback} onDelete - callback function fired to handle delete event.
 */

export default class KvDeleteModal extends Component {
  @tracked deleteType = null; // Either delete-version or delete-current-version.
  @tracked modalOpen = false;

  get modalDisplay() {
    switch (this.args.mode) {
      // Does not match adapter key directly because a delete type must be selected.
      case 'delete':
        return {
          title: 'Delete version?',
          type: 'warning',
          intro:
            'There are two ways to delete a version of a secret. Both delete actions can be undeleted later. How would you like to proceed?',
        };
      case 'destroy':
        return {
          title: 'Destroy version?',
          type: 'danger',
          intro: `This action will permanently destroy Version ${this.args.version} of the secret, and the secret data cannot be read or recovered later.`,
        };
      case 'delete-metadata':
        return {
          title: 'Delete metadata?',
          type: 'danger',
          intro:
            'This will permanently delete the metadata and versions of the secret. All version history will be removed. This cannot be undone.',
        };
      default:
        return assert('mode must be one of delete, destroy, or delete-metadata.');
    }
  }

  get deleteOptions() {
    const { secret, metadata, version } = this.args;
    const isDeactivated = secret.canReadMetadata ? metadata?.currentSecret?.isDeactivated : false;
    return [
      {
        key: 'delete-version',
        label: 'Delete this version',
        description: `This deletes ${version ? `Version ${version}` : `a specific version`} of the secret.`,
        disabled: !secret.canDeleteVersion,
        tooltipMessage: `Deleting a specific version requires "update" capabilities to ${secret.backend}/delete/${secret.path}.`,
      },
      {
        key: 'delete-latest-version',
        label: 'Delete latest version',
        description: 'This deletes the most recent version of the secret.',
        disabled: !secret.canDeleteLatestVersion || isDeactivated,
        tooltipMessage: isDeactivated
          ? `The latest version of the secret is already ${metadata.currentSecret.state}.`
          : `Deleting the latest version of this secret requires "delete" capabilities to ${secret.backend}/data/${secret.path}.`,
      },
    ];
  }

  @action
  onDelete() {
    const type = this.args.mode === 'delete' ? this.deleteType : this.args.mode;
    this.args.onDelete(type);
    this.modalOpen = false;
  }
}
