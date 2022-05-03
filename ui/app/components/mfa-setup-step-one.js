import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module MfaSetupStepOne
 * MfaSetupStepOne components are used to...
 *
 * @example
 * ```js
 * <MfaSetupStepOne @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} isUUIDVerified - action from parent that returns a boolean from the child regarding whether or not the UUID was verified. If true then proceed to step 2.
 */

export default class MfaSetupStepOne extends Component {
  @action
  verifyUUID(evt) {
    evt.preventDefault();
    // let UUID = this.UUID;
    // ARG TODO verify the UUID;
    // if verified send confirm boolean to the parent?
    this.args.isUUIDVerified(true);
  }
  @action
  redirectPreviousPage() {
    window.history.back();
  }
}
