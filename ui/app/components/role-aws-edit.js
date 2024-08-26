/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { isBlank } from '@ember/utils';
import errorMessage from 'vault/utils/error-message';

/**
 * @module ConfigureAwsComponent is generated from the SecretEditLayout route template. It is used to create, edit or view the role-aws model.
 *
 * @example
 * ```js
 * {{component
  (get (options-for-backend this.backendType this.model.idPrefix) "editComponent")
  model=this.model
  mode=this.mode
}}
 * ```
 *
 * @param {string} mode - "create" or "edit" or "show"
 * @param {object} model - role-aws model
 */

export default class RoleAwsEdit extends Component {
  @tracked errorMessage = null;
  @tracked invalidFormAlert = null;
  @tracked lintingErrors = null;
  @service router;

  get breadcrumbs() {
    const { backend } = this.args.model;
    return [
      { label: 'Secrets', route: 'vault.cluster.secrets.backends' },
      {
        label: backend,
        route: 'vault.cluster.secrets.backend',
        model: backend,
      },
      { label: this.args.mode },
    ];
  }

  get haveAttributesChanged() {
    const { model } = this.args;
    return Object.keys(model.changedAttributes()).some((key) => key === 'backend') ? false : true;
  }

  @action
  async createOrUpdate(event) {
    event.preventDefault();
    // all of the attributes with fieldValue:'id' are called `name`
    const { model } = this.args;
    const { credential_type, backend } = model;

    if (this.args.mode === 'create' && isBlank(backend)) {
      return;
    }
    // only save if there are changes
    const attrChanged =
      Object.keys(model.changedAttributes()).filter((item) => item !== 'backend').length > 0;
    if (!attrChanged) return; // todo flash message and transition

    if (credential_type === 'iam_user') {
      model.role_arns = [];
    }
    if (credential_type === 'assumed_role') {
      model.policy_arns = [];
    }
    if (credential_type === 'federation_token') {
      model.role_arns = model.policy_arns = [];
      [];
    }
    if (model.policy_document === '{}') {
      model.policy_document = '';
    }

    try {
      await model.save();
      this.router.transitionTo('vault.cluster.secrets.backend.credentials', backend);
    } catch (e) {
      this.errorMessage = errorMessage(e);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  codemirrorUpdated(attr, val, codemirror) {
    codemirror.performLint();
    this.lintingErrors = codemirror.state.lint.marked.length > 0;
    if (!this.lintingErrors) {
      this.args.model[attr] = val;
    }
  }

  @action
  delete() {
    this.args.model.destroyRecord().then(() => {
      this.router.transitionTo('vault.cluster.secrets.backend.list-root');
    });
  }
}
