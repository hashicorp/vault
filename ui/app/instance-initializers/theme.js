/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

export function initialize(appInstance) {
  // Eagerly instantiate the theme service so _applyTheme() runs on page load
  // and sets the correct data-theme attribute before any component renders.
  appInstance.lookup('service:theme');
}

export default {
  name: 'theme',
  initialize,
};
