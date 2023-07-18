/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';
import { ancestorKeysForKey, keyIsFolder, parentKeyForKey, keyWithoutParentKey } from 'core/utils/key-utils';
import keys from 'core/utils/key-codes';
import { tracked } from '@glimmer/tracking';
import escapeStringRegexp from 'escape-string-regexp';

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

  @tracked filterIsFocused = false;

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

  // filter operations and getters
  get filterMatchesASecretPath() {
    return !!(
      this.args.model.secrets &&
      this.args.model.secrets.length &&
      this.args.model.secrets.findBy('fullSecretPath', this.args.model.filterValue)
    );
  }

  get partialMatch() {
    // you can't pass undefined to RegExp so replacing with empty string if there is no pageFilter value
    const value = !this.pageFilter ? '' : this.pageFilter;
    const reg = new RegExp('^' + escapeStringRegexp(value));
    const match = this.args.model.secrets.filter((path) => reg.test(path.fullSecretPath))[0];

    if (this.filterMatchesASecretPath || !match) return null;
    // TODO not doing the shared prefix?
    return match.fullSecretPath;
  }

  @action
  handleInput(event) {
    // handling typing
    const input = event.target.value;
    const isDirectory = keyIsFolder(input);
    const parentDirectory = parentKeyForKey(input);
    const secretWithinDirectory = keyWithoutParentKey(input);

    if (isDirectory) {
      this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', input);
    } else if (parentDirectory) {
      this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', parentDirectory, {
        queryParams: { pageFilter: secretWithinDirectory },
      });
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.kv.list', {
        queryParams: { pageFilter: input },
      });
    }
  }
  @action
  handleKeyDown(event) {
    // handle keyboard events: tab, enter, escape
    const inputValue = event.target.value;
    const parentDirectory = parentKeyForKey(inputValue);

    if (event.keyCode === keys.TAB) {
      event.preventDefault();
      const isMatchDirectory = keyIsFolder(this.partialMatch);
      const parentDirectoryFromMatch = parentKeyForKey(this.partialMatch);
      const withoutDirectory = keyWithoutParentKey(this.partialMatch);

      if (isMatchDirectory) {
        // beep/boop/
        this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', this.partialMatch);
      } else if (!isMatchDirectory && parentDirectoryFromMatch) {
        // beep/boop/my-
        this.router.transitionTo(
          'vault.cluster.secrets.backend.kv.list-directory',
          parentDirectoryFromMatch,
          {
            queryParams: { pageFilter: withoutDirectory },
          }
        );
      } else {
        this.router.transitionTo('vault.cluster.secrets.backend.kv.list', {
          queryParams: { pageFilter: this.partialMatch },
        });
      }
    }
    if (event.keyCode === keys.ENTER) {
      event.preventDefault();
      if (this.filterMatchesASecretPath) {
        // check if secret exists if it does, navigate to the details page
        this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details', inputValue);
      } else {
        this.router.transitionTo('vault.cluster.secrets.backend.kv.create', {
          queryParams: { initialKey: inputValue },
        });
      }
      return;
    }
    if (event.keyCode === keys.ESC) {
      // transition to the nearest directory or to the list route.
      // clear pageFilter queryParam each time
      return !parentDirectory
        ? this.router.transitionTo('vault.cluster.secrets.backend.kv.list', {
            queryParams: { pageFilter: '' },
          })
        : this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', parentDirectory, {
            queryParams: { pageFilter: '' },
          });
    }
    return;
  }

  @action
  setFilterIsFocused() {
    this.filterIsFocused = true;
  }
  @action
  focusInput() {
    if (this.args.model.filterValue) {
      document.getElementById('secret-filter')?.focus();
    }
  }
}
