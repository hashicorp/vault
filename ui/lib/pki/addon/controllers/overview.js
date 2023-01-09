import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

export default class PkiOverviewController extends Controller {
  @service router;

  @action
  navigateToIssueCertificate(roleName) {
    this.router.transitionTo(`vault.cluster.secrets.backend.pki.roles.role.generate`, roleName);
  }
  @action
  navigateToViewCertificate(certificateSerialNumber) {
    this.router.transitionTo(
      `vault.cluster.secrets.backend.pki.certificates.certificate.details`,
      certificateSerialNumber
    );
  }
}
