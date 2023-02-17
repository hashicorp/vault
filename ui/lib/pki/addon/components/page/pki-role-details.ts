import { action } from '@ember/object';
import RouterService from '@ember/routing/router-service';
import Component from '@glimmer/component';
import FlashMessageService from 'vault/services/flash-messages';
import SecretMountPath from 'vault/services/secret-mount-path';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import PkiRoleModel from 'vault/models/pki/role';

interface Args {
  role: PkiRoleModel;
}

export default class DetailsPage extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPath;

  get breadcrumbs() {
    return [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'roles', route: 'roles.index' },
      { label: this.args.role.id },
    ];
  }

  get arrayAttrs() {
    return ['keyUsage', 'extKeyUsage', 'extKeyUsageOids'];
  }

  @action
  async deleteRole() {
    try {
      await this.args.role.destroyRecord();
      this.flashMessages.success('Role deleted successfully');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.index');
    } catch (error) {
      this.args.role.rollbackAttributes();
      this.flashMessages.danger(errorMessage(error));
    }
  }
}
