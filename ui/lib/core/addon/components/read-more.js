/**
 * @module ReadMore
 * ReadMore components are used to...
 *
 * @example
 * ```js
 * <ReadMore @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import layout from '../templates/components/read-more';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

class ReadMoreComponent extends Component {
  @action
  calculateOverflow(e) {
    const spanText = e.querySelector('.description-block');
    if (spanText.offsetWidth > e.offsetWidth) {
      this.hasOverflow = true;
    }
  }

  @tracked
  isOpen = false;

  @tracked
  hasOverflow = false;

  @action
  toggleOpen(e) {
    console.log(e);
    this.isOpen = !this.isOpen;
  }
}

export default setComponentTemplate(layout, ReadMoreComponent);
