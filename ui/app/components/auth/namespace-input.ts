/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { getRelativePath } from 'core/utils/sanitize-path';
import { restartableTask, timeout } from 'ember-concurrency';
import { service } from '@ember/service';

import type { HTMLElementEvent } from 'vault/forms';
import type ApiService from 'vault/services/api';
import type FlagsService from 'vault/services/flags';

/**
 * @module Auth::NamespaceInput
 * Renders the namespace input for the login form. As a user types, the updateNamespace callback fires in the controller to update the query param in the URL.
 * When a namespace is updated, the controller sets `shouldRefocusNamespaceInput = true` which refocuses the input after the route refreshes.
 * For HVD managed clusters the input prepends the administrative namespace: `admin/`. The input is disabled if the url has an OIDC query param: "?o=someprovider"
 *
 * @param {boolean} disabled - determines whether or not the namespace input is disabled
 * @param {function} handleNamespaceUpdate - fires updateNamespace callback in controller
 * @param {string} namespaceQueryParam - namespace query param from the url
 * @param {boolean} shouldRefocusNamespaceInput - if true, refocuses the input on `{{did-insert}}`
 * */

interface Args {
  handleNamespaceUpdate: CallableFunction;
  namespaceQueryParam: string;
  shouldRefocusNamespaceInput: boolean;
}

export default class AuthNamespaceInput extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;

  get namespaceInput() {
    const namespaceQueryParam = this.args.namespaceQueryParam;
    if (this.flags.hvdManagedNamespaceRoot) {
      // When managed, the user isn't allowed to edit the prefix `admin/`
      // so prefill just the relative path in the namespace input
      const path = getRelativePath(namespaceQueryParam, this.flags.hvdManagedNamespaceRoot);
      return path ? `/${path}` : '';
    }
    return namespaceQueryParam;
  }

  @action
  async handleInput(event: HTMLElementEvent<HTMLInputElement>) {
    // user has typed something, so input should be refocused
    const value = event.target.value;
    this.updateNamespace.perform(value);
  }

  updateNamespace = restartableTask(async (namespace) => {
    await timeout(500);
    await this.args.handleNamespaceUpdate(namespace);
  });

  @action
  maybeRefocus(element: HTMLElement) {
    if (this.args.shouldRefocusNamespaceInput) {
      element.focus();
    }
  }
}
