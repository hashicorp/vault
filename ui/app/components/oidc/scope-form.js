/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';

/**
 * @module OidcScopeForm
 * Oidc scope form components are used to create and edit oidc scopes
 *
 * @example
 * ```js
 * <Oidc::ScopeForm @form={{this.model}} />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} form - oidc scope form
 * @param {onCancel} onCancel - callback triggered when cancel button is clicked
 * @param {onSave} onSave - callback triggered on save success
 */

export default class OidcScopeFormComponent extends Component {
  @service api;
  @service flashMessages;

  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked modelValidations;
  // formatting here is purposeful so that whitespace renders correctly in JsonEditor
  exampleTemplate = `{
  "username": {{identity.entity.aliases.$MOUNT_ACCESSOR.name}},
  "contact": {
    "email": {{identity.entity.metadata.email}},
    "phone_number": {{identity.entity.metadata.phone_number}}
  },
  "groups": {{identity.entity.groups.names}}
}`;

  get breadcrumbs() {
    const crumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'OIDC provider: Scopes', route: 'vault.cluster.access.oidc.scopes' },
    ];

    if (!this.args.form.isNew) {
      crumbs.push({
        label: this.args.form.data.name,
        route: 'vault.cluster.access.oidc.scopes.scope.details',
        model: this.args.form.data.name,
      });
    }

    crumbs.push({ label: this.args.form.isNew ? 'Create scope' : 'Edit scope' });
    return crumbs;
  }

  save = task(
    waitFor(async (event) => {
      event.preventDefault();
      try {
        const { isNew } = this.args.form;
        const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;

        if (isValid) {
          const { name, ...payload } = data;
          await this.api.identity.oidcWriteScope(name, payload);
          this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the scope ${name}.`);
          this.args.onSave();
        }
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorBanner = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );
}
