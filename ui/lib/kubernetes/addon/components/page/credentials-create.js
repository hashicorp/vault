import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

export default class CredentialsCreatePageComponent extends Component {
  @service store;

  @tracked ttl = '';
  @tracked kubernetesNamespace = '';
  @tracked clusterRoleBinding = false;
  @tracked error;

  @tracked credentials;

  constructor() {
    super(...arguments);
  }

  @action
  setKubernetesNamespace({ target }) {
    this.kubernetesNamespace = target.value;
  }

  @action
  updateTtl({ goSafeTimeString }) {
    this.ttl = goSafeTimeString;
  }

  @task
  @waitFor
  *fetchCredentials() {
    try {
      const payload = {
        role: this.args.model.roleName,
        kubernetes_namespace: this.kubernetesNamespace,
        cluster_role_binding: this.clusterRoleBinding,
        ttl: this.ttl,
      };
      const credentials = yield this.store
        .adapterFor('kubernetes/role')
        .generateCredentials(this.args.model.kubernetesBackend, payload);

      this.credentials = credentials;
    } catch (error) {
      this.error = errorMessage(error);
    }
  }
}
