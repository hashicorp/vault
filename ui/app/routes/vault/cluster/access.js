import { computed } from '@ember/object';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

export default Route.extend(ModelBoundaryRoute, ClusterRoute, {
  modelTypes: computed(function() {
    return ['capabilities', 'control-group', 'identity/group', 'identity/group-alias', 'identity/alias'];
  }),
  model() {
    return {};
  },
});
