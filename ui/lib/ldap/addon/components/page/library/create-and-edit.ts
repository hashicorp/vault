import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type LdapLibraryModel from 'vault/models/ldap/library';
import { Breadcrumb, ValidationMap } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: LdapLibraryModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapCreateAndEditLibraryPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();

    const { model } = this.args;
    const { isValid, state, invalidFormMessage } = model.validate();

    this.modelValidations = isValid ? null : state;
    this.invalidFormMessage = isValid ? '' : invalidFormMessage;

    if (isValid) {
      try {
        const action = model.isNew ? 'created' : 'updated';
        yield model.save();
        this.flashMessages.success(`Successfully ${action} the library ${model.name}.`);
        this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries.library.details', model.name);
      } catch (error) {
        this.error = errorMessage(error, 'Error saving library. Please try again or contact support.');
      }
    }
  }

  @action
  cancel() {
    this.args.model.rollbackAttributes();
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries');
  }
}
