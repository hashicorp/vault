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
 * `KvListFilter` component is used to filter through the KV metadata LIST response. It allows users to search through the current list, navigate into directories, and use keyboard functions to: autocomplete, view a secret, create a new secret, or clear the input field.
 *
 * @param {object} model - The adapter model object which contains an array of secret models.
 * @param {string} mountPoint - Where in the router files we're located. For this component it will always be vault.cluster.secrets.backend.kv
 * @param {string} filterValue - A concatenation between the list-directory's dynamic path "path-to-secret" and the queryParam "pageFilter". For example, if we're inside the beep/ directory searching for any secret that starts with "my-" this value will equal "beep/my-".
 * @param {string} pageFilter - The queryParam value.
 */

export default class KvListFilterComponent extends Component {
  @service router;
  @tracked filterIsFocused = false;

  kvRoute(route) {
    return `${this.args.mountPoint}.${route}`;
  }
  /*
  -partialMatch returns the secret that most closely matches the pageFilter queryParam.
  -We're focused on pageFilter and not filterValue because if we're inside a directory we only care about the secrets listed there and not the directory. 
  -If pageFilter is empty this returns the first secret model in the list.
**/
  get partialMatch() {
    // you can't pass undefined to RegExp so if pageFilter is empty we replace it with an empty string.
    const value = !this.args.pageFilter ? '' : this.args.pageFilter;
    const reg = new RegExp('^' + escapeStringRegexp(value));
    const match = this.args.model.secrets.filter((path) => reg.test(path.fullSecretPath))[0];
    if (this.filterMatchesASecretPath || !match) return null;

    return match.fullSecretPath;
  }
  /*
  -filterMatchesASecretPath returns true if the filterValue matches a fullSecretPath
  within the list of models.
**/
  get filterMatchesASecretPath() {
    return !!this.args.model.secrets?.findBy('fullSecretPath', this.args.filterValue);
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
    // TODO kv engine cleanup: ideally when it's not a directory we could filter through the current models and remove pageFilter refresh on the list route.
    if (isDirectory) {
      this.router.transitionTo(this.kvRoute('list-directory'), input);
    } else if (parentDirectory) {
      this.router.transitionTo(this.kvRoute('list-directory'), parentDirectory, {
        queryParams: { pageFilter: secretWithinDirectory },
      });
    } else {
      this.router.transitionTo(this.kvRoute('list'), {
        queryParams: { pageFilter: input },
      });
    }
  }
  /*
  -handleKeyDown handles: tab, enter, backspace and escape. Ignores everything else.
**/
  @action
  handleKeyDown(event) {
    const input = event.target.value;
    const parentDirectory = parentKeyForKey(input);

    if (event.keyCode === keys.BACKSPACE && parentDirectory) {
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
    this.router.transitionTo(this.kvRoute('list-directory'), parentDirectory, {
      queryParams: {
        pageFilter,
      },
    });
  }
  handleTab() {
    const matchParentDirectory = parentKeyForKey(this.partialMatch);
    const isMatchDirectory = keyIsFolder(this.partialMatch);
    const matchWithoutParentDirectory = keyWithoutParentKey(this.partialMatch);

    if (isMatchDirectory) {
      // ex: beep/boop/
      this.router.transitionTo(this.kvRoute('list-directory'), this.partialMatch);
    } else if (!isMatchDirectory && matchParentDirectory) {
      // ex: beep/boop/my-
      this.router.transitionTo(this.kvRoute('list-directory'), matchParentDirectory, {
        queryParams: { pageFilter: matchWithoutParentDirectory },
      });
    } else {
      // ex: my-
      this.router.transitionTo(this.kvRoute('list'), {
        queryParams: { pageFilter: this.partialMatch },
      });
    }
  }
  handleEnter(input) {
    // TODO input queryParam on details and create pages.
    if (this.filterMatchesASecretPath) {
      // if secret exists send to details
      this.router.transitionTo(this.kvRoute('secret.details'), input);
    } else {
      // if secret does not exists send to create with the path prefilled with input value.
      this.router.transitionTo(this.kvRoute('create'), {
        queryParams: { initialKey: input },
      });
    }
  }
  handleEscape(parentDirectory) {
    // transition to the nearest parentDirectory. If no parentDirectory, then to the list route.
    !parentDirectory
      ? this.router.transitionTo(this.kvRoute('list'), {
          queryParams: { pageFilter: '' },
        })
      : this.router.transitionTo(this.kvRoute('list-directory'), parentDirectory, {
          queryParams: { pageFilter: '' },
        });
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
