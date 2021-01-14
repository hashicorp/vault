import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class GenerateCredentialsDatabase extends Component {
  @service store;
  // set on the component
  backendType = null;
  backendPath = null;
  roleName = null;
  @tracked model = null;

  constructor() {
    super(...arguments);
    this.fetchCredentials.perform();
  }
  @task(function*() {
    let { roleName, backendType } = this.args;
    let newModel = yield this.store.queryRecord('database/credential', {
      backend: backendType,
      secret: roleName,
    });
    this.model = newModel;
  })
  fetchCredentials;

  @action redirectPreviousPage() {
    window.history.back();
  }
}
