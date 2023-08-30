import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type LdapRoleModel from 'vault/models/ldap/role';
import { Breadcrumb } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: LdapRoleModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapRoleDetailsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @action
  async delete() {
    try {
      await this.args.model.destroyRecord();
      this.flashMessages.success('Role deleted successfully.');
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles');
    } catch (error) {
      const message = errorMessage(error, 'Unable to delete role. Please try again or contact support.');
      this.flashMessages.danger(message);
    }
  }

  @task
  @waitFor
  *rotateCredentials() {
    try {
      yield this.args.model.rotateStaticPassword();
      this.flashMessages.success('Credentials successfully rotated.');
    } catch (error) {
      this.flashMessages.danger(`Error rotating credentials \n ${errorMessage(error)}`);
    }
  }
}
