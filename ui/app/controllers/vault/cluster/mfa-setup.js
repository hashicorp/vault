import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class VaultClusterMfaSetupController extends Controller {
  @service auth;
  @tracked onStep = 1;

  get entityId() {
    // ARG TODO if root this will return empty string.
    return this.auth.authData.entity_id;
  }

  @action isUUIDVerified(response) {
    if (response) {
      this.onStep = 2;
    } else {
      this.isError = 'UUID was not verified';
      // ARG TODO work with Ivana on error message.
      // try and figure out API response.
    }
  }
  @action isAuthenticationCodeVerified(response) {
    if (response) {
      this.onStep = 3;
    } else {
      this.isError = 'Authentication code not verified';
      // ARG TODO work with Ivana on error message.
      // try and figure out API response.
    }
  }
}
