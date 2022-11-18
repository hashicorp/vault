import { computed } from '@ember/object';
import ClusterRoute from 'vault/mixins/cluster-route';
import UnloadModelRoute from 'vault/routes/-unload-model-route';

export default UnloadModelRoute.extend(ClusterRoute, {
  modelTypes: computed(function () {
    return ['identity/group', 'identity/group-alias', 'identity/alias', 'generated-user-userpass'];
  }),
  model() {
    return {};
  },
});
