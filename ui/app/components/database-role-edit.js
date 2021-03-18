import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class DatabaseRoleEdit extends Component {
  @service router;
  @service flashMessages;
  @service wizard;

  constructor() {
    super(...arguments);
    if (
      this.wizard.featureState === 'displayConnection' ||
      this.wizard.featureState === 'displayRoleDatabase'
    ) {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', 'database');
    }
    if (this.args.initialKey) {
      this.args.model.database = [this.args.initialKey];
    }
  }

  @tracked loading = false;

  get warningMessages() {
    let warnings = {};
    if (this.args.model.canUpdateDb === false) {
      warnings.database = `You don’t have permissions to update this database connection, so this role cannot be created.`;
    }
    if (
      (this.args.model.type === 'dynamic' && this.args.model.canCreateDynamic === false) ||
      (this.args.model.type === 'static' && this.args.model.canCreateStatic === false)
    ) {
      warnings.type = `You don't have permissions to create this type of role.`;
    }
    return warnings;
  }

  get databaseType() {
    if (this.args.model?.database) {
      // TODO: Calculate this
      return 'mongodb-database-plugin';
    }
    return null;
  }

  @action
  generateCreds(roleId) {
    this.router.transitionTo('vault.cluster.secrets.backend.credentials', roleId);
  }

  @action
  delete() {
    const secret = this.args.model;
    const backend = secret.backend;
    secret
      .destroyRecord()
      .then(() => {
        try {
          this.router.transitionTo(LIST_ROOT_ROUTE, backend, { queryParams: { tab: 'role' } });
        } catch (e) {
          console.debug(e);
        }
      })
      .catch(e => {
        this.flashMessages.danger(e.errors?.join('. '));
      });
  }

  @action
  handleCreateEditRole(evt) {
    evt.preventDefault();
    this.loading = true;

    const mode = this.args.mode;
    let roleSecret = this.args.model;
    let secretId = roleSecret.name;
    if (mode === 'create') {
      roleSecret.set('id', secretId);
      let path = roleSecret.type === 'static' ? 'static-roles' : 'roles';
      roleSecret.set('path', path);
    }
    roleSecret
      .save()
      .then(() => {
        try {
          this.router.transitionTo(SHOW_ROUTE, `role/${secretId}`);
        } catch (e) {
          console.debug(e);
        }
      })
      .catch(e => {
        const errorMessage = e.errors?.join('. ') || e.message;
        this.flashMessages.danger(
          errorMessage || 'Could not save the role. Please check Vault logs for more information.'
        );
        this.loading = false;
      });
  }
}
