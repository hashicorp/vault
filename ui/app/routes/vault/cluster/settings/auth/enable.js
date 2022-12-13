import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class VaultClusterSettingsAuthEnableRoute extends Route {
  @service store;

  beforeModel() {
    // Unload to prevent naming collisions when we mount a new engine
    this.store.unloadAll('auth-method');
  }

  model() {
    const authMethod = this.store.createRecord('auth-method');
    authMethod.set('config', this.store.createRecord('mount-config'));
    return authMethod;
  }
}
