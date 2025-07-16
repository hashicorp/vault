import { modifier } from 'ember-modifier';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Applies and updates a custom CSS property based on updates to a tracked/computed ember value.
 */
var cssCustomProperty = modifier(function customProperty(element, [prop, value]) {
  if (!value) {
    element.style.removeProperty(prop);
    return;
  }
  // Allow only basic unit type values
  const safeValuePattern = /^[-a-zA-Z0-9\s#(),.%]+$/;
  if (!safeValuePattern.test(value)) {
    console.warn(`Blocked potentially unsafe variable value: ${value}`);
    return;
  }
  element.style.setProperty(prop, value);
});

export { cssCustomProperty as default };
//# sourceMappingURL=css-custom-property.js.map
