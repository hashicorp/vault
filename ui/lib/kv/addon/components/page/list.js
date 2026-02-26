/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { getOwner } from '@ember/owner';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { ancestorKeysForKey } from 'core/utils/key-utils';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';

/**
 * @module List
 * ListPage component is a component to show a list of secrets.
 *
 * @param {array} secrets - An array of secrets
 * @param {string} backend - The name of the kv secret engine.
 * @param {string} pathToSecret - The directory name that the secret belongs to ex: beep/boop/
 * @param {string} filterValue - The concatenation of the pathToSecret and pageFilter ex: beep/boop/my-
 * @param {boolean} failedDirectoryQuery - true if the query was a 403 and the search was for a directory. Used to display inline alert message on the overview card.
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 * @param {object} capabilities - capabilities for metadata path
 */

export default class KvListPageComponent extends Component {
  @service flashMessages;
  @service('app-router') router;
  @service api;

  @tracked secretPath;
  @tracked metadataToDelete = null; // set to the metadata intended to delete

  // used for KV list and list-directory view
  // ex: beep/
  isDirectory = (path) => pathIsDirectory(path);
  fullSecretPath = (secret) => `${this.args.pathToSecret}${secret}`;

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
  async onDelete(secretPath) {
    try {
      const fullSecretPath = this.fullSecretPath(secretPath);
      await this.api.secrets.kvV2DeleteMetadataAndAllVersions(fullSecretPath, this.args.backend);
      const message = `Successfully deleted the metadata and all version data of the secret ${fullSecretPath}.`;
      this.flashMessages.success(message);
      // if you've deleted a secret from within a directory, transition to its parent directory.
      if (this.router.currentRoute.localName === 'list-directory') {
        const ancestors = ancestorKeysForKey(fullSecretPath);
        const nearest = ancestors.pop();
        this.router.transitionTo(`${this.mountPoint}.list-directory`, nearest);
      } else {
        // Transition to refresh the model
        this.router.transitionTo(`${this.mountPoint}.list`);
      }
    } catch (error) {
      const { message } = await this.api.parseError(
        error,
        'Error deleting secret. Please try again or contact support.'
      );
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
