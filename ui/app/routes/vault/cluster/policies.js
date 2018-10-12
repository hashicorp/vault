import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

const ALLOWED_TYPES = ['acl', 'egp', 'rgp'];

export default Route.extend(ClusterRoute, {
  version: service(),

  beforeModel() {
    return this.get('version')
      .fetchFeatures()
      .then(() => {
        return this._super(...arguments);
      });
  },

  model(params) {
    let policyType = params.type;
    if (!ALLOWED_TYPES.includes(policyType)) {
      return this.transitionTo(this.routeName, ALLOWED_TYPES[0]);
    }
    return {};
  },
});
