import Route from '@ember/routing/route';
import { get } from '@ember/object';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';

export default Route.extend({
  beforeModel() {
    const type = this.modelFor('vault.cluster.settings.auth.configure').get('type');
    const section = get(tabsForAuthSection([type]), 'firstObject.routeParams.lastObject');
    return this.transitionTo('vault.cluster.settings.auth.configure.section', section);
  },
});
