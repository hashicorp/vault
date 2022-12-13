import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { add } from 'date-fns';
import errorMessage from 'vault/utils/error-message';

export default class CredentialsPageComponent extends Component {
  @service store;
  @service router;

  @tracked ttl = '';
  @tracked clusterRoleBinding = false;
  @tracked kubernetesNamespace;
  @tracked error;

  @tracked credentials;

  get leaseExpiry() {
    return add(new Date(), { seconds: this.credentials.lease_duration });
  }

  @action
  cancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles.role.details');
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
        role: this.args.role.name,
        kubernetes_namespace: this.kubernetesNamespace,
        cluster_role_binding: this.clusterRoleBinding,
        ttl: this.ttl,
      };

      this.credentials = yield this.store
        .adapterFor('kubernetes/role')
        .generateCredentials(this.args.role.backend, payload);
    } catch (error) {
      this.error = errorMessage(error);
    }
  }
}
