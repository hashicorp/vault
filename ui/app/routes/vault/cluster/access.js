import Ember from 'ember';
import ClusterRoute from 'vault/mixins/cluster-route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

export default Ember.Route.extend(ModelBoundaryRoute, ClusterRoute, {
  modelTypes: ['capabilities', 'control-group', 'identity/group', 'identity/group-alias', 'identity/alias'],
  model() {
    return {};
  },
});
