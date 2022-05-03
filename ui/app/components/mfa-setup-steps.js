import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module MfaSetupSteps
 * MfaSetupSteps components are used to...
 *
 * @example
 * ```js
 * <MfaSetupSteps @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class MfaSetupSteps extends Component {
  @tracked onStep = 1;

  @action isUUIDVerified(response) {
    console.log('here');
    if (response) {
      this.onStep = 2;
    } else {
      this.isError = 'UUID was not verified';
      // ARG TODO work with Ivana on error message.
      // try and figure out API response.
    }
  }
  @action isAuthenticationCodeVerified(response) {
    console.log('here');
    if (response) {
      this.onStep = 3;
    } else {
      this.isError = 'Authentication code not verified';
      // ARG TODO work with Ivana on error message.
      // try and figure out API response.
    }
  }
}
