import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module MfaSetupStepTwo
 * MfaSetupStepTwo components are used to...
 *
 * @example
 * ```js
 * <MfaSetupStepTwo @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class MfaSetupStepTwo extends Component {
  @action verifyAuthenticationCode(evt) {
    evt.preventDefault();
    // let authenticationCode = this.authenticationCode;
    // ARG TODO verify the UUID;
    // if verified send confirm boolean to the parent?
    this.args.isAuthenticationCodeVerified(true);
  }
}
