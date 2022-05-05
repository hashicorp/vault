import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class VaultClusterMfaSetupController extends Controller {
  @service auth;
  @tracked onStep = 1;
  @tracked warning = '';
  @tracked uuid = '';

  get entityId() {
    // ARG TODO if root this will return empty string.
    return this.auth.authData.entity_id;
  }

  @action isUUIDVerified(verified) {
    if (verified) {
      this.onStep = 2;
    } else {
      this.onStep = 1;
    }
  }

  @action
  goToReset(warning) {
    this.warning = warning;
    this.onStep = 3;
  }

  @action isAuthenticationCodeVerified(response) {
    if (response) {
      this.onStep = 3;
    } else {
      // ARG TODO work with Ivana on error message.
      // try and figure out API response.
    }
  }

  @action
  saveUUID(uuid) {
    this.uuid = uuid;
    console.log(this.uuid, 'Save UUID should be called on verify');
  }
}
