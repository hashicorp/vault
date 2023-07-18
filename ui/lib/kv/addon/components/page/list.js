/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import { ancestorKeysForKey } from 'core/utils/key-utils';
import errorMessage from 'vault/utils/error-message';

/**
 * @module List
 * ListPage component is a component to show a list of kv/metadata secrets.
 *
 * @param {array} model - An array of models generated form kv/metadata query.
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 * @param {string} filterValue - The input on the Filter secrets Navigate input or the current secret directory.
 */

export default class KvListPageComponent extends Component {
  @service flashMessages;
  @service router;

  get mountPoint() {
    // mountPoint tells the LinkedBlock component where to start the transition. In this case, mountPoint will always be vault.cluster.secrets.backend.kv.
    return getOwner(this).mountPoint;
  }

  @action
  async onDelete(model) {
    try {
      const message = `Successfully deleted secret ${model.fullSecretPath}.`;
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
}
