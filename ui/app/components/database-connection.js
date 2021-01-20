import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

import { action } from '@ember/object';

const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class DatabaseConnectionEdit extends Component {
  store = service();
  router = service();

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
    // this.args.sendMessage(this.body);
    // this.body = '';
    // this.store.createRecord(secretId, {
    //   title: 'Rails is Omakase',
    //   body: 'Lorem ipsum'
    // });
    // this.store
    //   .adapterFor('secret-v2-version')
    //   .deleteRecord()
  }
}
