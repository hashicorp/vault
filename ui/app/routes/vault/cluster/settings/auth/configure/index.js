import Route from '@ember/routing/route';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';

export default Route.extend({
  beforeModel() {
    const model = this.modelFor('vault.cluster.settings.auth.configure');
    const section = tabsForAuthSection([model]).firstObject.routeParams.lastObject;
    return this.transitionTo('vault.cluster.settings.auth.configure.section', section);
  },
});
