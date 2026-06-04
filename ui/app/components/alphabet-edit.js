/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';

/**
 * @module AlphabetEdit
 * `AlphabetEdit` is a component that allows you to create/edit or view an alphabet.
 *
 * @example
 * ```js
 *   <AlphabetEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />
 * ```
 * @param {object} form - AlphabetForm instance with data and formFields.
 * @param {object} capabilities - Object with canDelete, canUpdate, canRead capabilities.
 * @param {string} mode - Is either show, create or edit.
 */

export default class AlphabetEditComponent extends Component {
  @service flashMessages;
  @service router;
  @service api;

  @tracked errorMessage = '';

  get breadcrumbs() {
    // ideally this is created on the controller in the parent route but this is a generic route and adding breadcrumbs to the controller requires a larger refactor.
    const backend = this.args.form?.data?.backend;
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

  get title() {
    if (this.args?.mode === 'create') {
      return 'Create Alphabet';
    } else if (this.args?.mode === 'edit') {
      return 'Edit Alphabet';
    } else {
      return this.args?.form?.data?.name;
    }
  }

  transition(route = 'show') {
    this.errorMessage = '';
    const { backend, name } = this.args.form.data;
    if (route === 'list') {
      this.router.transitionTo('vault.cluster.secrets.backend.list-root', backend, {
        queryParams: { tab: 'alphabet' },
      });
      return;
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.show', `alphabet/${name}`);
    }
  }

  @action async createOrUpdate(event) {
    event.preventDefault();

    const { name, alphabet, backend } = this.args.form.data;

    try {
      await this.api.secrets.transformWriteAlphabet(name, backend, { alphabet });
      this.flashMessages.success('Alphabet saved.');
      this.transition();
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
    }
  }

  @action async onDelete() {
    const { name, backend } = this.args.form.data;
    try {
      await this.api.secrets.transformDeleteAlphabet(name, backend);
      this.flashMessages.success('Alphabet deleted.');
      this.transition('list');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
    }
  }
}
