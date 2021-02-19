import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class DatabaseRoleEdit extends Component {
  @service router;
  @service flashMessages;

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
  handleCreateRole(evt) {
    evt.preventDefault();
    let roleSecret = this.args.model;
    let secretId = roleSecret.name;
    roleSecret.set('id', secretId);
    let path = roleSecret.type === 'static' ? 'static-roles' : 'roles';
    roleSecret.set('path', path);
    roleSecret.save().then(() => {
      try {
        this.router.transitionTo(SHOW_ROUTE, `role/${secretId}`);
      } catch (e) {
        console.debug(e);
      }
    });
  }

  @action
  handleCreateEditRole(evt) {
    evt.preventDefault();
    const mode = this.args.mode;
    let roleSecret = this.args.model;
    let secretId = roleSecret.name;
    if (mode === 'create') {
      roleSecret.set('id', secretId);
      let path = roleSecret.type === 'static' ? 'static-roles' : 'roles';
      roleSecret.set('path', path);
    }
    roleSecret.save().then(() => {
      try {
        this.router.transitionTo(SHOW_ROUTE, `role/${secretId}`);
      } catch (e) {
        console.debug(e);
      }
    });
  }
}
