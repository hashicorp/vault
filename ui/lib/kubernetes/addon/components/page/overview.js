import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class OverviewPageComponent extends Component {
  @service router;

  @tracked roles = [];
  @tracked selectedRole = null;
  @tracked roleOptions = this.roles.map((role) => {
    return { name: role.name, id: role.name };
  });

  constructor() {
    super(...arguments);
    this.roles = this.args.model.roles.toArray();
  }

  get isDisabled() {
    return !this.selectedRole;
  }

  @action
  selectRole([roleName]) {
    this.selectedRole = roleName;
  }

  @action
  generateCredential() {
    const { selectedRole } = this;

    this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles.role.credentials', selectedRole);
  }
}
