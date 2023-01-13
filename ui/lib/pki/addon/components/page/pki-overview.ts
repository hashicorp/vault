import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';
// TYPES
import Store from '@ember-data/store';
import RouterService from '@ember/routing/router-service';
import PkiIssuerModel from 'vault/models/pki/issuer';
import PkiRoleModel from 'vault/models/pki/role';
import { waitFor } from '@ember/test-waiters';

interface Args {
  issuers: number | PkiIssuerModel[];
  roles: number | PkiRoleModel[];
  engine: string;
}

export default class PkiOverview extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly store: Store;

  @tracked rolesValue = '';
  @tracked commonNameValue = '';
  @tracked certificateValue = '';

  @action
  transitionToViewCertificates(event: Event) {
    event.preventDefault();
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.certificates.certificate.details',
      this.certificateValue
    );
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    const role = this.rolesValue;
    const commonName = this.commonNameValue;

    const certificateModel = this.store.createRecord('pki/certificate/generate', {
      role,
      commonName,
    });

    try {
      yield certificateModel.save();

      this.router.transitionTo(
        'vault.cluster.secrets.backend.pki.certificates.certificate.details',
        certificateModel.serialNumber
      );
    } catch (err) {
      errorMessage(err);
    }
  }

  @action
  handleRolesInput(roles: string) {
    if (Array.isArray(roles)) {
      this.rolesValue = roles[0];
    } else {
      this.rolesValue = roles;
    }
  }
  @action
  handleCommonNameInput(commonName: string) {
    if (Array.isArray(commonName)) {
      this.commonNameValue = commonName[0];
    } else {
      this.commonNameValue = commonName;
    }
  }

  @action
  handleCertificateInput(certificate: string) {
    if (Array.isArray(certificate)) {
      this.certificateValue = certificate[0];
    } else {
      this.certificateValue = certificate;
    }
  }
}
