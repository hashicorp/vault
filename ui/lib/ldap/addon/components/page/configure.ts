import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type LdapConfigModel from 'vault/models/ldap/config';
import { Breadcrumb, ValidationMap } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: LdapConfigModel;
  breadcrumbs: Array<Breadcrumb>;
}
interface SchemaOption {
  title: string;
  icon: string;
  description: string;
  value: string;
}

export default class LdapConfigurePageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked showRotatePrompt = false;
  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormMessage = '';
  @tracked error = '';

  get schemaOptions(): Array<SchemaOption> {
    return [
      {
        title: 'OpenLDAP',
        icon: 'folder',
        description:
          'OpenLDAP is one of the most popular open source directory service developed by the OpenLDAP Project.',
        value: 'openldap',
      },
      {
        title: 'AD',
        icon: 'microsoft',
        description:
          'Active Directory is a directory service developed by Microsoft for Windows domain networks.',
        value: 'ad',
      },
      {
        title: 'RACF',
        icon: 'users',
        description:
          "For managing IBM's Resource Access Control Facility (RACF) security system, the generated passwords must be 8 characters or less.",
        value: 'racf',
      },
    ];
  }

  leave(route: string) {
    this.router.transitionTo(`vault.cluster.secrets.backend.ldap.${route}`);
  }

  validate() {
    const { isValid, state, invalidFormMessage } = this.args.model.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormMessage = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  async rotateRoot() {
    try {
      await this.args.model.rotateRoot();
    } catch (error) {
      // since config save was successful at this point we only want to show the error in a flash message
      this.flashMessages.danger(`Error rotating root password \n ${errorMessage(error)}`);
    }
  }

  @task
  @waitFor
  *save(event: Event | null, rotate: boolean) {
    if (event) {
      event.preventDefault();
    }
    const isValid = this.validate();
    // show rotate creds prompt for new models when form state is valid
    this.showRotatePrompt = isValid && this.args.model.isNew && !this.showRotatePrompt;

    if (isValid && !this.showRotatePrompt) {
      try {
        yield this.args.model.save();
        // if save was triggered from confirm action in rotate password prompt we need to make an additional request
        if (rotate) {
          yield this.rotateRoot();
        }
        this.flashMessages.success('Successfully configured LDAP engine');
        this.leave('configuration');
      } catch (error) {
        this.error = errorMessage(error, 'Error saving configuration. Please try again or contact support.');
      }
    }
  }

  @action
  cancel() {
    const { model } = this.args;
    const transitionRoute = model.isNew ? 'overview' : 'configuration';
    const cleanupMethod = model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    model[cleanupMethod]();
    this.leave(transitionRoute);
  }
}
