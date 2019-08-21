import Route from '@ember/routing/route';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';
export default Route.extend({
  beforeModel() {
    let { methodType, paths } = this.modelFor('vault.cluster.access.method');
    let navigationPaths = paths.paths.filter(path => path.navigation);
    const activeTab = tabsForAuthSection([methodType, 'authConfig', navigationPaths])[0].routeParams;
    return this.transitionTo(...activeTab);
  },
});
