import Controller from '@ember/controller';
import { task } from 'ember-concurrency';

export default Controller.extend({
  showRoute: 'vault.cluster.access.identity.show',
  showTab: 'details',
  navAfterSave: task(function* ({ saveType, model }) {
    const isDelete = saveType === 'delete';
    const type = model.get('identityType');
    const listRoutes = {
      'entity-alias': 'vault.cluster.access.identity.aliases.index',
      'group-alias': 'vault.cluster.access.identity.aliases.index',
      group: 'vault.cluster.access.identity.index',
      entity: 'vault.cluster.access.identity.index',
    };
    const routeName = listRoutes[type];
    if (!isDelete) {
      yield this.transitionToRoute(this.showRoute, model.id, this.showTab);
      return;
    }
    yield this.transitionToRoute(routeName);
  }),
});
