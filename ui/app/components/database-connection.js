import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

import { action } from '@ember/object';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

const getErrorMessage = errors => {
  let errorMessage = 'Something went wrong. Check the Vault logs for more information.';
  if (errors?.join(' ').indexOf('failed to verify')) {
    errorMessage =
      'There was a verification error for this connection. Check the Vault logs for more information.';
  }
  return errorMessage;
};

export default class DatabaseConnectionEdit extends Component {
  @service store;
  @service router;
  @service flashMessages;

  transitionToRoute() {
    return this.router.transitionTo(...arguments);
  }

  @action
  async handleCreateConnection(evt) {
    evt.preventDefault();
    let secret = this.args.model;
    let secretId = secret.name;
    secret.set('id', secretId);
    secret.save().then(() => {
      this.transitionToRoute(SHOW_ROUTE, secretId);
    });
  }

  @action
  delete(evt) {
    evt.preventDefault();
    // const adapter = this.store.adapterFor('cluster');
    const secret = this.args.model;
    const backend = secret.backend;
    secret.destroyRecord().then(() => {
      this.transitionToRoute(LIST_ROOT_ROUTE, backend);
    });
  }

  @action
  reset() {
    const { name, backend } = this.args.model;
    let adapter = this.store.adapterFor('database/connection');
    adapter
      .resetConnection(backend, name)
      .then(() => {
        // TODO: Why isn't the confirmAction closing?
        // this.args.onRefresh();
        this.flashMessages.success('Successfully reset connection');
      })
      .catch(e => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }

  @action
  rotate() {
    const { name, backend } = this.args.model;
    let adapter = this.store.adapterFor('database/connection');
    adapter
      .rotateRootCredentials(backend, name)
      .then(() => {
        // TODO: Why isn't the confirmAction closing?
        this.flashMessages.success('Successfully rotated credentials');
      })
      .catch(e => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }
}
