import Route from '@ember/routing/route';
export default class ConfigRoute extends Route {
  model() {
    return this.store.queryRecord('clients/config', {});
  }
}
