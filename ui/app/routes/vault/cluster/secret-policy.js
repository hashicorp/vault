import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
const SUPPORTED_BACKENDS = supportedSecretBackends();

const getTransformPolicy = secret => {
  let secretName = secret || '<secret-name-here>';
  return `
# Work with transform secrets engine
path "transform/${secretName}/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Enable secrets engine
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# List enabled secrets engine
path "sys/mounts" {
  capabilities = [ "read", "list" ]
}`;
};

export default Route.extend(ClusterRoute, {
  version: service(),

  backendType() {
    return this.modelFor('vault.cluster.secrets.backend').get('engineType');
  },

  beforeModel() {
    console.log('BEFORE_MODEL');
    return this.get('version')
      .fetchFeatures()
      .then(() => {
        return this._super(...arguments);
      });
  },

  model(params) {
    console.log('---- MODEL');
    console.log(params);
    let backend = params.backend;
    let secretName = params.secret_name;
    console.log(backend, secretName);
    // let secretName = params.secret_name;
    // let policyType = params.type;
    if (!SUPPORTED_BACKENDS.includes(backend)) {
      console.log('TODO: redirect to normal policies page');
    }
    // if (secretName) {
    // return this.transitionTo(this.routeParent, 'policies');
    // }
    return {
      backend,
      secretName,
      policy: getTransformPolicy(secretName),
    };
  },
});
