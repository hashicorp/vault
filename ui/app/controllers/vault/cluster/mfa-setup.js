import Controller from '@ember/controller';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class VaultClusterMfaSetupController extends Controller {
  @tracked onStep = 1;

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
