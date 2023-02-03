import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiIssuerSignRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/sign-intermediate', this.secretMountPath.currentPath);
  }

  model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    return this.store.createRecord('pki/sign-intermediate', { issuerRef: issuer_ref });
  }
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: resolvedModel.issuerRef, route: 'issuers.issuer.details' },
      { label: 'sign intermediate' },
    ];
  }
}
