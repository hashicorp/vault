import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-config';
import { hash } from 'rsvp';

@withConfig()
export default class PkiRolesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  async fetchRoles() {
    try {
      return await this.store.query('pki/role', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      if (e.httpStatus === 404) {
        return { parentModel: this.modelFor('roles') };
      } else {
        throw e;
      }
    }
  }

  model() {
    return hash({
      hasConfig: this.shouldPromptConfig,
      roles: this.fetchRoles(),
      parentModel: this.modelFor('roles'),
    });
  }
}
