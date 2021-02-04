// import Component from '@ember/component';
import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

// import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
// import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class DatabaseRoleEdit extends Component {
  @service router;

  get buttonDisabled() {
    if (this.args.model.canUpdateDb === false) {
      return true;
    }
    return true;
  }

  @action
  generateCreds(roleId) {
    this.router.transitionTo('vault.cluster.secrets.backend.credentials', roleId);
  }

  @action
  delete(evt) {
    evt.preventDefault();
    // const adapter = this.store.adapterFor('cluster');
    const secret = this.args.model;
    const backend = secret.backend;
    secret.destroyRecord().then(() => {
      // TODO: Update database allowed roles
      this.router.transitionTo(LIST_ROOT_ROUTE, backend);
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
      // TODO: Update database with allowed roles + id
      this.router.transitionTo(SHOW_ROUTE, `role/${secretId}`);
    });
  }
}
