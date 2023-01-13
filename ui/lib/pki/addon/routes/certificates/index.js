import PkiOverviewRoute from '../overview';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiCertificatesIndexRoute extends PkiOverviewRoute {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/certificate', this.secretMountPath.currentPath);
  }

  async fetchCertificateModel() {
    try {
      return await this.store.query('pki/certificate', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      if (e.httpStatus === 404) {
        return { parentModel: this.modelFor('certificates') };
      } else {
        throw e;
      }
    }
  }

  model() {
    return hash({
      hasConfig: this.hasConfig(),
      certificateModel: this.fetchCertificateModel(),
      parentModel: this.modelFor('certificates'),
    });
  }
}
