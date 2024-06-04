/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import layout from '../templates/components/read-more';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module ReadMore
 * @description
 * ReadMore components are used to wrap long text that we'd like to show as one line initially with the option to expand and read.
 * Text which is shorter than the surrounding div will not truncate or show the See More button.
 *
 * @example
 * <div style="width:100px;border: solid purple 2px; margin:20px">
 * <ReadMore>My <em>super</em> long text goes in here.</ReadMore>
 * </div>
 *
 * <div style="width:100px;border: solid purple 2px; margin:20px">
 * <ReadMore>This text fits!</ReadMore>
 * </div>
 */

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
  toggleOpen() {
    this.isOpen = !this.isOpen;
  }
}

export default setComponentTemplate(layout, ReadMoreComponent);
