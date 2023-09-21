import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

import type LdapLibraryModel from 'vault/models/ldap/library';
import { Breadcrumb } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: LdapLibraryModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class LdapLibraryDetailsPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @action
  async delete() {
    try {
      await this.args.model.destroyRecord();
      this.flashMessages.success('Library deleted successfully.');
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.libraries');
    } catch (error) {
      const message = errorMessage(error, 'Unable to delete library. Please try again or contact support.');
      this.flashMessages.danger(message);
    }
  }
}
