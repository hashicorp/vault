import { modifier } from 'ember-modifier';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * By default the axis elements are outside of the bounds of the svg and rely on the containing element having
 * enough padding to compensate. A fixed padding is not flexible to varied width of axis labels.
 *
 * This modifier is used to pad compensate for the width of the axis element. It also returns a value
 * that can be used to set the chart width based on the offset amount.
 */
var axisOffset = modifier((element, [onOffset, additionalPadding = 0]) => {
  const axis = element.querySelector('g.axis');
  if (!axis) {
    return;
  }
  const axisOffset = axis.getBoundingClientRect().width;
  element.style.transform = `translateX(${axisOffset + additionalPadding}px)`;
  if (onOffset) {
    onOffset(axisOffset + additionalPadding);
  }
});

export { axisOffset as default };
//# sourceMappingURL=axis-offset.js.map
