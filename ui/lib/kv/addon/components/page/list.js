/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { getOwner } from '@ember/application';
import { ancestorKeysForKey } from 'core/utils/key-utils';
import errorMessage from 'vault/utils/error-message';
import { pathIsDirectory } from 'kv/utils/kv-breadcrumbs';

/**
 * @module List
 * ListPage component is a component to show a list of kv/metadata secrets.
 *
 * @param {array} model - An array of models generated form kv/metadata query.
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 * @param {boolean} noMetadataListPermissions - true if the return to query metadata LIST is 403, indicating the user does not have permissions to that endpoint.
 */

export default class KvListPageComponent extends Component {
  @service flashMessages;
  @service router;

  @tracked secretPath;

  get mountPoint() {
    // mountPoint tells the LinkedBlock component where to start the transition. In this case, mountPoint will always be vault.cluster.secrets.backend.kv.
    return getOwner(this).mountPoint;
  }

  get buttonText() {
    return pathIsDirectory(this.secretPath) ? 'View directory' : 'View secret';
  }

  @action
  async onDelete(model) {
    try {
      const message = `Successfully deleted the metadata and all version data of the secret ${model.fullSecretPath}.`;
      await model.destroyRecord();
      this.flashMessages.success(message);
      // if you've deleted a secret from within a directory, transition to its parent directory.
      if (this.args.routeName === 'list-directory') {
        const ancestors = ancestorKeysForKey(model.fullSecretPath);
        const nearest = ancestors.pop();
        this.router.transitionTo(`${this.mountPoint}.list-directory`, nearest);
      }
    } catch (error) {
      const message = errorMessage(error, 'Error deleting secret. Please try again or contact support.');
      this.flashMessages.danger(message);
    }
  }

  @action
  handleSecretPathInput(value) {
    this.secretPath = value;
  }

  @action
  transitionToSecretDetail() {
    pathIsDirectory(this.secretPath)
      ? this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', this.secretPath)
      : this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details', this.secretPath);
  }
}
