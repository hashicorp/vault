import Ember from 'ember';
import { task } from 'ember-concurrency';

export default Ember.Controller.extend({
  showRoute: 'vault.cluster.access.identity.show',
  showTab: 'details',
  navAfterSave: task(function*({saveType, model}) {
    let isDelete = saveType === 'delete';
    let type = model && model.get('identityType');
    let listRoutes= {
      'entity-alias': 'vault.cluster.access.identity.aliases',
      'group-alias': 'vault.cluster.access.identity.aliases',
      'group': 'vault.cluster.access.identity.index',
      'entity': 'vault.cluster.access.identity.index',
    };
    let routeName = listRoutes[type]
    if (!isDelete) {
      yield this.transitionToRoute(
        this.get('showRoute'),
        model.id,
        this.get('showTab')
      );
      return;
    }
    yield this.transitionToRoute(
      'vault.cluster.access.identity.index'
    );
  }),
});
