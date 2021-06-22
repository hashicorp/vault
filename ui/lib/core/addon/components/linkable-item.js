/**
 * @module LinkableItem
 * LinkableItem components are used to...
 *
 * @example
 * ```js
 * <LinkableItem @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import layout from '../templates/components/linkable-item';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

class LinkableItemComponent extends Component {
  //track height
  //open vs close text

  @tracked
  isOpen = false;

  @action
  toggleOpen() {
    this.isOpen = !this.isOpen;
  }
}

export default setComponentTemplate(layout, LinkableItemComponent);
