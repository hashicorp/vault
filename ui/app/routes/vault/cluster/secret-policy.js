import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
const SUPPORTED_BACKENDS = supportedSecretBackends();

const getBackendPolicy = (backendType, backendName) => {
  let backend = backendName || '<secret-engine-name>';
  if (backendType === 'kv-v2') {
    return `
# control over secrets engine config
path "${backend}/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# control over secret engine data
path "${backend}/metadata/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# List enabled secrets engine
path "sys/mounts/${backend}" {
  capabilities = [ "read", "list" ]
}`;
  }
  return `
# Work with transform secrets engine
path "${backend}/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# See secrets engine in list
path "sys/mounts/${backend}" {
  capabilities = [ "read", "list" ]
}`;
};

const getPolicy = (backendType, backendName, secretName) => {
  let backend = backendName || '<secret-engine-name>';
  return `
# Work with transform secrets engine
path "${backend}/${secretName}/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# Enable secrets engine
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# List enabled secrets engine
path "sys/mounts/${backendName}" {
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
    let backendType = params.backend;
    let backendName = params.backend_name;
    let secretName = params.secret_name;
    if (!SUPPORTED_BACKENDS.includes(backendType)) {
      console.log('TODO: redirect to normal policies page');
    }
    // if (!backendType) {
    //   return this.transitionTo(this.routeParent, 'policies');
    // }
    return {
      backend: backendName,
      secretName,
      policy: secretName
        ? getPolicy(backendType, backendName, secretName)
        : getBackendPolicy(backendType, backendName),
    };
  },
});
