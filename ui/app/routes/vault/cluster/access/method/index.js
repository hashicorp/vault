import Route from '@ember/routing/route';
import { tabsForAuthSection } from 'vault/helpers/tabs-for-auth-section';
export default Route.extend({
  beforeModel() {
    let { methodType, paths } = this.modelFor('vault.cluster.access.method');
    paths = paths ? paths.navPaths.reduce((acc, cur) => acc.concat(cur.path), []) : null;
    const activeTab = tabsForAuthSection([methodType, 'authConfig', paths])[0].routeParams;
    return this.transitionTo(...activeTab);
  },
});
