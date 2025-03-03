/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

/** // ARG TODO
 * @module AlphabetEdit
 * `AlphabetEdit` is a component that will allow you to pick a file from the local file system. Once
 * loaded, this file will be emitted as a JS ArrayBuffer to the passed `onChange` callback.
 *
 * @example
 * ```js
 *   <AlphabetEdit }} />
 * ```
 * @param {object} model - This is the transform template model.
 * @param {string} mode - Determines if create or edit.

 *
 */
export default class AlphabetEditComponent extends Component {
  @service flashMessages;
  @service router;

  @tracked errorMessage = '';

  get breadcrumbs() {
    // ideally this is created on the controller in the parent route but this is a generic route and requires a significant refactor.
    const { backend } = this.args.model;
    return [
      {
        label: backend,
        route: 'vault.cluster.secrets.backend.list-root',
        model: backend,
        query: { tab: 'alphabet' },
      },
      { label: 'Alphabet' },
    ];
  }

  transition(route = 'show') {
    this.displayErrors = '';
    const { backend, id } = this.args.model;
    if (route === 'list') {
      this.router.transitionTo('vault.cluster.secrets.backend.list-root', backend, {
        queryParams: { tab: 'alphabet' },
      });
      return;
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.show', `alphabet/${id}`);
    }
  }

  @action async createOrUpdate(event) {
    event.preventDefault();
    const { id, name } = this.args.model;
    const modelId = id || name;

    if (!modelId) return; // TODO should solve with modelValidations instead

    if (!this.args.model?.hasDirtyAttributes) {
      this.flashMessages.info('No changes detected.');
      this.transition();
      return;
    }

    try {
      await this.args.model.save();
      this.flashMessages.success('Alphabet saved.');
      this.transition();
    } catch (e) {
      this.errorMessage = errorMessage(e);
    }
  }

  @action async onDelete() {
    try {
      await this.args.model.destroyRecord();
      this.flashMessages.success('Alphabet deleted.');
      this.transition('list');
    } catch (e) {
      this.errorMessage = errorMessage(e);
    }
  }
}
