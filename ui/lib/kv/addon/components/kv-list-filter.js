/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import keys from 'core/utils/key-codes';
import { keyIsFolder, parentKeyForKey, keyWithoutParentKey } from 'core/utils/key-utils';
import escapeStringRegexp from 'escape-string-regexp';
import { tracked } from '@glimmer/tracking';

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
  @tracked filterIsFocused = false;

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

  /*
  - partialMatch returns the secret that most closely matches the pageFilter queryParam.
  - Searches pageFilter and not filterValue because if we're inside a directory we only care about the secrets listed there and not the directory. 
  - If pageFilter is empty this returns the first secret model in the list.
**/
  get partialMatch() {
    // If pageFilter is empty we replace it with an empty string because you cannot pass 'undefined' to RegEx.
    const value = !this.args.pageFilter ? '' : this.args.pageFilter;
    const reg = new RegExp('^' + escapeStringRegexp(value));
    const match = this.args.secrets.filter((path) => reg.test(path.fullSecretPath))[0];
    if (this.isFilterMatch || !match) return null;

    return match.fullSecretPath;
  }
  /*
  - isFilterMatch returns true if the filterValue matches a fullSecretPath.
**/
  get isFilterMatch() {
    return !!this.args.secrets?.findBy('fullSecretPath', this.args.filterValue);
  }
  /*
  -handleInput is triggered after the value of the input has changed. It is not triggered when input looses focus.
**/
  @action
  handleInput(event) {
    const input = event.target.value;
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
  /*
  -handleKeyDown handles: tab, enter, backspace and escape. Ignores everything else.
**/
  @action
  handleKeyDown(event) {
    const input = event.target.value;
    const parentDirectory = parentKeyForKey(input);

    if (event.keyCode === keys.BACKSPACE) {
      this.handleBackspace(input, parentDirectory);
    }

    if (event.keyCode === keys.TAB) {
      event.preventDefault();
      this.handleTab();
    }

    if (event.keyCode === keys.ENTER) {
      event.preventDefault();
      this.handleEnter(input);
    }

    if (event.keyCode === keys.ESC) {
      this.handleEscape(parentDirectory);
    }
    // ignore all other key events
    return;
  }
  // key-code specific methods
  handleBackspace(input, parentDirectory) {
    const isInputDirectory = keyIsFolder(input);
    const inputWithoutParentKey = keyWithoutParentKey(input);
    const pageFilter = isInputDirectory ? '' : inputWithoutParentKey.slice(0, -1);
    this.navigate(parentDirectory, pageFilter);
  }
  handleTab() {
    const isMatchDirectory = keyIsFolder(this.partialMatch);
    const matchParentDirectory = parentKeyForKey(this.partialMatch);
    const matchWithinDirectory = keyWithoutParentKey(this.partialMatch);

    if (isMatchDirectory) {
      // ex: beep/boop/
      this.navigate(this.partialMatch);
    } else if (!isMatchDirectory && matchParentDirectory) {
      // ex: beep/boop/my-
      this.navigate(matchParentDirectory, matchWithinDirectory);
    } else {
      // ex: my-
      this.navigate(null, this.partialMatch);
    }
  }
  handleEnter(input) {
    if (this.isFilterMatch) {
      // if secret exists send to details
      this.router.transitionTo(`${this.args.mountPoint}.secret.details`, input);
    } else {
      // if secret does not exists send to create with the path prefilled with input value.
      this.router.transitionTo(`${this.args.mountPoint}.create`, {
        queryParams: { initialKey: input },
      });
    }
  }
  handleEscape(parentDirectory) {
    // transition to the nearest parentDirectory. If no parentDirectory, then to the list route.
    !parentDirectory ? this.navigate() : this.navigate(parentDirectory);
  }

  @action
  setFilterIsFocused() {
    // tracked property used to show or hide the help-text next to the input. Not involved in focus event itself.
    this.filterIsFocused = true;
  }

  @action
  focusInput() {
    // set focus to the input when there is either a pageFilter queryParam value and/or list-directory's dynamic path-to-secret has a value.
    if (this.args.filterValue) {
      document.getElementById('secret-filter')?.focus();
    }
  }
}
