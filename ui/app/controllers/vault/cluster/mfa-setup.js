import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class VaultClusterMfaSetupController extends Controller {
  @service auth;
  @tracked onStep = 1;
  @tracked warning = '';
  @tracked uuid = '';
  @tracked qrCode = '';

  get entityId() {
    // ARG TODO if root this will return empty string.
    return this.auth.authData.entity_id;
  }

  @action isUUIDVerified(verified) {
    if (verified) {
      this.onStep = 2;
    } else {
      this.restartFlow();
    }
  }

  @action
  restartFlow() {
    // they have hit the restart screen and hit the cancel button.
    window.history.back();
    this.onStep = 1;
  }

  @action
  saveUUIDandQrCode(uuid, qrCode) {
    // qrCode could be an empty string if the admin-generate was not successful
    this.uuid = uuid;
    this.qrCode = qrCode;
  }

  @action
  showWarning(warning) {
    this.warning = warning;
    this.onStep = 2;
  }
}
