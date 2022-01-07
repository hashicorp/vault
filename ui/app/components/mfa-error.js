import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { TOTP_NOT_CONFIGURED } from 'vault/services/auth';

const TOTP_NA_MSG =
  'Multi-factor authentication is required, but you have not set it up. In order to do so, please contact your administrator.';
const MFA_ERROR_MSG =
  'Multi-factor authentication is required, but failed. Go back and try again, or contact your administrator.';

export { TOTP_NA_MSG, MFA_ERROR_MSG };

/**
 * @module MfaError
 * MfaError components are used to display mfa errors
 *
 * @example
 * ```js
 * <MfaError />
 * ```
 */

export default class MfaError extends Component {
  @service auth;

  get isTotp() {
    return this.auth.mfaErrors.includes(TOTP_NOT_CONFIGURED);
  }
  get title() {
    return this.isTotp ? 'TOTP not set up' : 'Unauthorized';
  }
  get description() {
    return this.isTotp ? TOTP_NA_MSG : MFA_ERROR_MSG;
  }

  @action
  onClose() {
    this.auth.set('mfaErrors', null);
    if (this.args.onClose) {
      this.args.onClose();
    }
  }
}
