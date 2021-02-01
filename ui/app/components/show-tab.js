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
    let response = yield this.store.queryRecord('capabilities', this.pathQuery(this.args.id, this.args.path));
    this.dontShowTab = !response.canList && !response.canCreate && !response.canUpdate;
  })
  fetchCapabilities;
}
