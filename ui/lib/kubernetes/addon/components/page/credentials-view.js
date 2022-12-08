import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
// import { action } from '@ember/object';
// import { task } from 'ember-concurrency';
// import { waitFor } from '@ember/test-waiters';
// import errorMessage from 'vault/utils/error-message';

export default class CredentialsCreatePageComponent extends Component {
  @service store;

  @tracked serviceAcctNamespace = null;
  @tracked serviceAcctToken = null;
  @tracked serviceAcctName = '';
  @tracked leaseDuration = null;
  @tracked leaseId = null;

  constructor() {
    super(...arguments);

    const {
      lease_duration,
      lease_id,
      data: { service_account_token, service_account_name, service_account_namespace },
    } = this.args.credentials;

    this.leaseId = lease_id;
    this.leaseDuration = lease_duration;
    this.serviceAcctToken = service_account_token;
    this.serviceAcctName = service_account_name;
    this.serviceAcctNamespace = service_account_namespace;
  }

  get leaseExpiry() {
    let date = new Date();
    date.setMilliseconds(date.getMilliseconds() + this.leaseDuration);
    date = new Date(date);
    return date;
  }
}
