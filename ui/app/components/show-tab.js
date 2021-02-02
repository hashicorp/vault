import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

export default class ShowTab extends Component {
  @service store;
  @tracked dontShowTab;
  constructor() {
    super(...arguments);
    this.fetchCapabilities.perform();
  }

  pathQuery(backend, path) {
    return {
      id: `${backend}/${path}/`,
    };
  }

  @task(function*() {
    let peekRecordRoles = yield this.store.peekRecord('capabilities', 'database/roles/');
    let peekRecordConnections = yield this.store.peekRecord('capabilities', 'database/config/');
    // peekRecord if the capabilities store data is there for the connections (config) and roles model
    if (peekRecordRoles && this.args.path === 'roles') {
      this.dontShowTab = !peekRecordRoles.canList && !peekRecordRoles.canCreate && !peekRecordRoles.canUpdate;
      return;
    }
    if (peekRecordConnections && this.args.path === 'config') {
      this.dontShowTab =
        !peekRecordConnections.canList &&
        !peekRecordConnections.canCreate &&
        !peekRecordConnections.canUpdate;
      return;
    }
    // otherwise queryRecord and create an instance on the capabilities.
    let response = yield this.store.queryRecord('capabilities', this.pathQuery(this.args.id, this.args.path));
    this.dontShowTab = !response.canList && !response.canCreate && !response.canUpdate;
  })
  fetchCapabilities;
}
