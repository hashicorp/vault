import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class OverviewPageComponent extends Component {
  @service router;

  @tracked selectedRole = null;
  @tracked roleOptions = [];

  constructor() {
    super(...arguments);
    this.roleOptions = this.args.model.roles.map((role) => {
      return { name: role.name, id: role.name };
    });
  }

  @action
  selectRole([roleName]) {
    this.selectedRole = roleName;
  }

  @action
  generateCredential() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.kubernetes.roles.role.credentials',
      this.selectedRole
    );
  }
}
