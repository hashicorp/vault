/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import keys from 'core/utils/key-codes';
import { keyIsFolder, parentKeyForKey, keyWithoutParentKey } from 'core/utils/key-utils';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';

/**
 * @module KvListFilter
 * `KvListFilter` filters through the KV metadata LIST response. It allows users to search through the current list, navigate into directories, and use keyboard functions to: autocomplete, view a secret, create a new secret, or clear the input field.
 *  *
 * <KvListFilter
 *  @secrets={{this.model.secrets}}
 *  @mountPoint={{this.model.mountPoint}}
 *  @filterValue="beep/my-"
 *  @pageFilter="my-"
 * />
 * @param {array} secrets - An array of secret models.
 * @param {string} mountPoint - Where in the router files we're located. For this component it will always be vault.cluster.secrets.backend.kv
 * @param {string} filterValue - A concatenation between the list-directory's dynamic path "path-to-secret" and the queryParam "pageFilter". For example, if we're inside the beep/ directory searching for any secret that starts with "my-" this value will equal "beep/my-".
 * @param {string} pageFilter - The queryParam value, does not include pathToSecret ex: my-.
 */

export default class KvListFilterComponent extends Component {
  @service router;
  @tracked query;
  @tracked filterIsFocused = false;

  constructor() {
    super(...arguments);
    this.query = this.args.filterValue;
  }

  navigate(pathToSecret, pageFilter) {
    const route = pathToSecret ? `${this.args.mountPoint}.list-directory` : `${this.args.mountPoint}.list`;
    const args = [route];
    if (pathToSecret) {
      args.push(pathToSecret);
    }
    args.push({
      queryParams: {
        pageFilter: pageFilter ? pageFilter : null,
      },
    });
    this.router.transitionTo(...args);
  }

  @action
  handleKeyDown(event) {
    if (event.keyCode === keys.ESC) {
      // On escape, transition to the nearest parentDirectory.
      // If no parentDirectory, then to the list route.
      const input = event.target.value;
      const parentDirectory = parentKeyForKey(input);
      !parentDirectory ? this.navigate() : this.navigate(parentDirectory);
    }
    // ignore all other key events
    return;
  }

  @action handleInput(evt) {
    this.query = evt.target.value;
  }

  @task
  *handleSearch(evt) {
    evt.preventDefault();
    // shows loader to indicate that the search was executed
    yield timeout(250);
    const input = this.query;
    const isDirectory = keyIsFolder(input);
    const parentDirectory = parentKeyForKey(input);
    const secretWithinDirectory = keyWithoutParentKey(input);
    if (isDirectory) {
      this.navigate(input);
    } else if (parentDirectory) {
      this.navigate(parentDirectory, secretWithinDirectory);
    } else {
      this.navigate(null, input);
    }
  }
}
