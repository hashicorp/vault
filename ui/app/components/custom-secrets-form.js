/**
 * @module CustomSecretsForm
 * CustomSecretsForm components are used to...
 *
 * @example
 * ```js
 * <CustomSecretsForm @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class CustomSecretsForm extends Component {
  @tracked
  model = {
    set: (path, val) => {
      console.log('TODO: set', path, val);
    },
  };

  @action
  handleSubmit(evt) {
    console.log('submitted');
    console.log(evt);
  }
}
