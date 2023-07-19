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
 * KvFilterList component replaces the older NavigateInput for the KV located in the ember engine.
 *
 * @param {object} model - The model object which contains model.secrets, and array of model objects that make up the list view. The route controls what is return in this object based on the pageFilter query param and the list-directory route path-to-secret param. Together these are combined to form the filterValue param.
 * @param {boolean} mountPoint - Tells us where in the router files we're located. For this component it will always be vault.cluster.secrets.backend.kv
 */

export default class KvFilterListComponent extends Component {
  @service router;

  @tracked filterIsFocused = false;

  kvRoute(route) {
    return `${this.args.mountPoint}.${route}`;
  }
  /*
  partialMatch returns the fullSecretPath of the secret that most closely matches the pageFilter queryParam.
  If pageFilter is 'b' and the list of secrets '[ae, be, ce]', then the match will 'be' secret. 
  We return the fullSecretPath in case we're inside a directory.
  If pageFilter is empty this returns the first secret model in the list.
**/
  get partialMatch() {
    // you can't pass undefined to RegExp so we replace pageFilter with an empty string.
    const value = !this.args.pageFilter ? '' : this.args.pageFilter;
    const reg = new RegExp('^' + escapeStringRegexp(value));
    const match = this.args.model.secrets.filter((path) => reg.test(path.fullSecretPath))[0];
    if (this.filterMatchesASecretPath || !match) return null;

    return match.fullSecretPath;
  }
  /*
  filterMatchesASecretPath returns true if the `path-to-secret + pageFilter` matches a fullSecretPath
  within the list of models.
  Ex: path-to-secret: `beep/boop/` + pageFilter: `bop` === fullSecretPath`beep/boop/bop` 
**/
  get filterMatchesASecretPath() {
    return !!(
      this.args.model.secrets &&
      this.args.model.secrets.length &&
      this.args.model.secrets.findBy('fullSecretPath', this.args.filterValue)
    );
  }
  /*
  Handles onInput event. Trigger occurs after the value of the input has changed. Is not triggered when input looses focus.
**/
  @action
  handleInput(event) {
    const input = event.target.value;
    const isDirectory = keyIsFolder(input);
    const parentDirectory = parentKeyForKey(input);
    const secretWithinDirectory = keyWithoutParentKey(input);
    // TODO ideally when it's not a directory we could just filter through the current models and remove pageFilter refresh on the list route.
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
  Handles specific key events: tab, enter, backspace and escape. Ignores everything else.
**/
  @action
  handleKeyDown(event) {
    const inputValue = event.target.value;
    const parentDirectory = parentKeyForKey(inputValue);
    const isInputDirectory = keyIsFolder(inputValue);
    const inputWithoutParentKey = keyWithoutParentKey(inputValue);
    if (event.keyCode === keys.BACKSPACE && parentDirectory) {
      const pageFilter = isInputDirectory ? '' : inputWithoutParentKey.slice(0, -1);
      this.router.transitionTo(this.kvRoute('list-directory'), parentDirectory, {
        queryParams: {
          pageFilter,
        },
      });
    }

    if (event.keyCode === keys.TAB) {
      event.preventDefault();
      const isMatchDirectory = keyIsFolder(this.partialMatch);
      const matchParentDirectory = parentKeyForKey(this.partialMatch);
      const matchMinusTheParentDirectory = keyWithoutParentKey(this.partialMatch);

      if (isMatchDirectory) {
        // beep/boop/
        this.router.transitionTo(this.kvRoute('list-directory'), this.partialMatch);
      } else if (!isMatchDirectory && matchParentDirectory) {
        // beep/boop/my-
        this.router.transitionTo(this.kvRoute('list-directory'), matchParentDirectory, {
          queryParams: { pageFilter: matchMinusTheParentDirectory },
        });
      } else {
        this.router.transitionTo(this.kvRoute('list'), {
          queryParams: { pageFilter: this.partialMatch },
        });
      }
    }
    if (event.keyCode === keys.ENTER) {
      event.preventDefault();
      // if secret exists navigate to the details page. Otherwise, send to create with path prefilled.
      // TODO inputValue queryParam on details and create pages.
      if (this.filterMatchesASecretPath) {
        this.router.transitionTo(this.kvRoute('secret.details'), inputValue);
      } else {
        this.router.transitionTo(this.kvRoute('create'), {
          queryParams: { initialKey: inputValue },
        });
      }
      return;
    }
    if (event.keyCode === keys.ESC) {
      // transition to the nearest parentDirectory or to the list route.
      // clear pageFilter queryParam each time
      return !parentDirectory
        ? this.router.transitionTo(this.kvRoute('list'), {
            queryParams: { pageFilter: '' },
          })
        : this.router.transitionTo(this.kvRoute('list-directory'), parentDirectory, {
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
    if (this.args.filterValue) {
      document.getElementById('secret-filter')?.focus();
    }
  }
}
