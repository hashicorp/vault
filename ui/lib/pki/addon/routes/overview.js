import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class PkiOverviewRoute extends Route {
  @service secretMountPath;
  @service auth;
  @service store;

  get win() {
    return this.window || window;
  }

  hasConfig() {
    // When the engine is configured, it creates a default issuer.
    // If the issuers list is empty, we know it hasn't been configured
    const endpoint = `${this.win.origin}/v1/${this.secretMountPath.currentPath}/issuers?list=true`;
    return this.auth
      .ajax(endpoint, 'GET', {})
      .then(() => true)
      .catch(() => false);
  }

  async model() {
    return hash({
      hasConfig: this.hasConfig(),
      engine: this.store
        .query('secret-engine', {
          path: this.secretMountPath.currentPath,
        })
        .then((model) => {
          if (model) {
            return model.get('firstObject');
          }
        }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.secretMountPath.currentPath || 'pki';
    controller.breadcrumbs = [{ label: 'secrets', route: 'secrets', linkExternal: true }, { label: backend }];
  }
}
