import PkiKeysIndexRoute from '.';

export default class PkiKeyEditRoute extends PkiKeysIndexRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push(
      { label: this.controller.model.id, route: 'keys.key.details' },
      { label: 'edit' }
    );
  }
}
