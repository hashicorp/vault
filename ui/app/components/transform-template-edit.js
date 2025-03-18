/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { ROUTES } from 'vault/utils/routes';

/**
 * @module TransformTemplateEdit
 * `TransformTemplateEdit` is a component that allows you to create/edit or view a transform template.
 *
 * @example
 * ```js
 *   <TransformTemplateEdit }} />
 * ```
 * @param {object} model - This is the transform template model.
 * @param {string} mode - Is either show, create or edit.
 */

export default class TransformTemplateEditComponent extends Component {
  @service flashMessages;
  @service router;

  @tracked errorMessage = '';

  get breadcrumbs() {
    // ideally this is created on the controller in the parent route but this is a generic route and adding breadcrumbs to the controller requires a larger refactor.
    const { backend } = this.args.model;
    return [
      {
        label: backend,
        route: ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT,
        model: backend,
        query: { tab: 'template' },
      },
      { label: 'Template' },
    ];
  }

  transition(route = 'show') {
    this.errorMessage = '';
    const { backend, id } = this.args.model;
    if (route === 'list') {
      this.router.transitionTo(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT, backend, {
        queryParams: { tab: 'template' },
      });
      return;
    } else {
      this.router.transitionTo(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_SHOW, `template/${id}`);
    }
  }

  @action async createOrUpdate(event) {
    event.preventDefault();

    if (!this.args.model.hasDirtyAttributes) {
      this.flashMessages.info('No changes detected.');
      this.transition();
      return;
    }

    try {
      await this.args.model.save();
      this.flashMessages.success('Transform template saved.');
      this.transition();
    } catch (e) {
      this.errorMessage = errorMessage(e);
    }
  }

  @action async onDelete() {
    try {
      await this.args.model.destroyRecord();
      this.flashMessages.success('Transform template deleted.');
      this.transition('list');
    } catch (e) {
      this.errorMessage = errorMessage(e);
    }
  }
}
