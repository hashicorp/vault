import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

export default class CredentialsPageComponent extends Component {
  @service store;
  @service router;

  @tracked ttl;
  @tracked kubernetesNamespace;
  @tracked error;
  @tracked clusterRoleBinding = false;

  @tracked serviceAcctName;
  @tracked serviceAcctNamespace;
  @tracked serviceAcctToken;
  @tracked leaseDuration;
  @tracked leaseId;

  @tracked showCredentialDetails;

  constructor() {
    super(...arguments);
  }

  get leaseExpiry() {
    let date = new Date();
    date.setSeconds(date.getSeconds() + this.leaseDuration);
    date = new Date(date);
    return date;
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
        role: this.args.model.roleModel.name,
        kubernetes_namespace: this.kubernetesNamespace,
        cluster_role_binding: this.clusterRoleBinding,
        ttl: this.ttl,
      };

      const credentials = yield this.store
        .adapterFor('kubernetes/role')
        .generateCredentials(this.args.model.kubernetesBackend, payload);

      const {
        lease_duration,
        lease_id,
        data: { service_account_token, service_account_name, service_account_namespace },
      } = credentials;

      this.showCredentialDetails = true;
      this.serviceAcctName = service_account_name;
      this.serviceAcctNamespace = service_account_namespace;
      this.serviceAcctToken = service_account_token;
      this.leaseDuration = lease_duration;
      this.leaseId = lease_id;
    } catch (error) {
      this.error = errorMessage(error);
    }
  }
}
