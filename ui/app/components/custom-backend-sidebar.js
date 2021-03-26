/**
 * @module CustomBackendSidebar
 * CustomBackendSidebar components are used to...
 *
 * @example
 * ```js
 * <CustomBackendSidebar @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class CustomBackendSidebar extends Component {
  @tracked
  param = '';

  @action
  updateParam(evt) {
    if (evt.target.name === 'clear') {
      this.param = '';
    } else {
      this.param = evt.target.value.trim();
    }
  }
}
