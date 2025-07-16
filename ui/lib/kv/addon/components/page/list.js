/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/owner';
import { ancestorKeysForKey } from 'core/utils/key-utils';
import errorMessage from 'vault/utils/error-message';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';

/**
 * @module List
 * ListPage component is a component to show a list of kv/metadata secrets.
 *
 * @param {array} secrets - An array of models generated form kv/metadata query.
 * @param {string} backend - The name of the kv secret engine.
 * @param {string} pathToSecret - The directory name that the secret belongs to ex: beep/boop/
 * @param {string} filterValue - The concatenation of the pathToSecret and pageFilter ex: beep/boop/my-
 * @param {boolean} failedDirectoryQuery - true if the query was a 403 and the search was for a directory. Used to display inline alert message on the overview card.
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 */

export default class KvListPageComponent extends Component {
  @service flashMessages;
  @service('app-router') router;
  @service pagination;

  @tracked secretPath;
  @tracked metadataToDelete = null; // set to the metadata intended to delete

  get mountPoint() {
    // mountPoint tells transition where to start. In this case, mountPoint will always be vault.cluster.secrets.backend.kv.
    return getOwner(this).mountPoint;
  }

  get buttonText() {
    // if secretPath is an empty string it could be because the user hit a permissions error.
    const path = this.secretPath || this.args.pathToSecret;
    return pathIsDirectory(path) ? 'View list' : 'View secret';
  }

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  }

  @action
  async onDelete(model) {
    try {
      // The model passed in is a kv/metadata model
      await model.destroyRecord();
      this.pagination.clearDataset('kv/metadata'); // Clear out the pagination cache so that the metadata/list view is updated.
      const message = `Successfully deleted the metadata and all version data of the secret ${model.fullSecretPath}.`;
      this.flashMessages.success(message);
      // if you've deleted a secret from within a directory, transition to its parent directory.
      if (this.router.currentRoute.localName === 'list-directory') {
        const ancestors = ancestorKeysForKey(model.fullSecretPath);
        const nearest = ancestors.pop();
        this.router.transitionTo(`${this.mountPoint}.list-directory`, nearest);
      } else {
        // Transition to refresh the model
        this.router.transitionTo(`${this.mountPoint}.list`);
      }
    } catch (error) {
      const message = errorMessage(error, 'Error deleting secret. Please try again or contact support.');
      this.flashMessages.danger(message);
    } finally {
      this.metadataToDelete = null;
    }
  }

  @action
  handleSecretPathInput(value) {
    this.secretPath = value;
  }

  @action
  transitionToSecretDetail(evt) {
    evt.preventDefault();
    pathIsDirectory(this.secretPath)
      ? this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', this.secretPath)
      : this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index', this.secretPath);
  }
}
