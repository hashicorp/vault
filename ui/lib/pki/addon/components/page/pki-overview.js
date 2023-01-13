import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';

export default class PkiOverview extends Component {
  @service router;
  @service store;

  @tracked rolesValue = '';
  @tracked certificateValue = '';

  get totalRoles() {
    return this.args.issuers === 404 ? 0 : this.args.issuers.length;
  }
  get totalIssuers() {
    return this.args.issuers === 404 ? 0 : this.args.issuers.length;
  }

  @action
  transitionToViewCertificates(event) {
    event.preventDefault();
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.certificates.certificate.details',
      this.certificateValue
    );
  }
  @action
  transitionToIssueCertificates(event) {
    event.preventDefault();
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.generate', this.rolesValue);
  }

  @action
  handleRolesInput(roles) {
    if (Array.isArray(roles)) {
      this.rolesValue = roles[0];
    } else {
      this.rolesValue = roles;
    }
  }

  @action
  handleCertificateInput(certificate) {
    if (Array.isArray(certificate)) {
      this.certificateValue = certificate[0];
    } else {
      this.certificateValue = certificate;
    }
  }
}
