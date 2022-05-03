import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module MfaSetupStepThree
 * MfaSetupStepThree components are used to...
 *
 * @example
 * ```js
 * <MfaSetupStepThree @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class MfaSetupStepThree extends Component {
  @action restartSetup() {
    console.log('restart??');
  }
}
