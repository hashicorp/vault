/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { SecretsApiTransformListAlphabetsListEnum } from '@hashicorp/vault-client-typescript';

/**
 * @module TransformTemplateEdit
 * `TransformTemplateEdit` is a component that allows you to create/edit or view a transform template.
 *
 * @example
 * ```js
 *   <TransformTemplateEdit @form={{this.form}} @capabilities={{this.capabilities}} @mode={{this.mode}} />
 * ```
 * @param {object} form - TemplateForm instance with data and formFields.
 * @param {object} capabilities - Object with canDelete, canUpdate, canRead capabilities.
 * @param {string} mode - Is either show, create or edit.
 */

export default class TransformTemplateEditComponent extends Component {
  @service flashMessages;
  @service router;
  @service api;

  @tracked errorMessage = '';
  @tracked modelValidations;
  @tracked alphabets = [];

  constructor() {
    super(...arguments);
    this.fetchAlphabets();
  }

  async fetchAlphabets() {
    try {
      const resp = await this.api.secrets.transformListAlphabets(
        this.args.form.data.backend,
        SecretsApiTransformListAlphabetsListEnum.TRUE
      );
      this.alphabets = (resp.keys ?? []).map((key) => ({ id: key }));
    } catch {
      // swallow errors, SearchSelect will fall back to string-list
    }
  }

  get breadcrumbs() {
    // ideally this is created on the controller in the parent route but this is a generic route and adding breadcrumbs to the controller requires a larger refactor.
    const backend = this.args.form?.data?.backend;
    return [
      {
        label: backend,
        route: 'vault.cluster.secrets.backend.list-root',
        model: backend,
        query: { tab: 'template' },
      },
      { label: 'Template' },
    ];
  }

  get title() {
    if (this.args.mode === 'create') {
      return 'Create Template';
    } else if (this.args.mode === 'edit') {
      return 'Edit Template';
    } else {
      return 'Template';
    }
  }

  get subtitle() {
    if (this.args.mode === 'create' || this.args.mode === 'edit') return '';

    return this.args.form?.data?.name;
  }

  transition(route = 'show') {
    this.errorMessage = '';
    this.modelValidations = null;
    const { backend, name } = this.args.form.data;
    if (route === 'list') {
      this.router.transitionTo('vault.cluster.secrets.backend.list-root', backend, {
        queryParams: { tab: 'template' },
      });
      return;
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.show', `template/${name}`);
    }
  }

  @action async createOrUpdate(event) {
    event.preventDefault();

    const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
    this.modelValidations = isValid ? null : state;
    this.errorMessage = invalidFormMessage;
    if (!isValid) return;

    const { name, pattern, alphabet, encode_format, decode_formats, backend } = data;

    try {
      await this.api.secrets.transformWriteTemplate(name, backend, {
        type: 'regex',
        pattern,
        alphabet: Array.isArray(alphabet) ? alphabet[0] : alphabet,
        encode_format,
        decode_formats,
      });
      this.flashMessages.success('Transform template saved.');
      this.transition();
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
    }
  }

  @action async onDelete() {
    const { name, backend } = this.args.form.data;
    try {
      await this.api.secrets.transformDeleteTemplate(name, backend);
      this.flashMessages.success('Transform template deleted.');
      this.transition('list');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    }
  }
}
