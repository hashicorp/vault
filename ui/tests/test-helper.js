/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from 'vault/app';
import config from 'vault/config/environment';
import * as QUnit from 'qunit';
import { setApplication } from '@ember/test-helpers';
import { setup } from 'qunit-dom';
import start from 'ember-exam/test-support/start';
import './helpers/flash-message';
import preloadAssets from 'ember-asset-loader/test-support/preload-assets';
import { setupGlobalA11yHooks, setRunOptions } from 'ember-a11y-testing/test-support';
import manifest from 'vault/config/asset-manifest';
import setupSinon from 'ember-sinon-qunit';
import { DISMISSED_WIZARD_KEY, WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { setupDarkThemeAudit } from './theme-helper';

preloadAssets(manifest).then(() => {
  setup(QUnit.assert);
  setApplication(Application.create(config.APP));
  setupSinon();
  // conditional a11y testing allows for theme specific testing
  // dark theme audit focusses specifically on the color-contrast rule
  if (config.APP.DARK_THEME_AUDIT) {
    setupDarkThemeAudit();
  } else {
    setupGlobalA11yHooks(() => true, {
      helpers: ['render'],
    });
    setRunOptions({
      runOnly: {
        type: 'tag',
        values: ['wcag2a'],
      },
    });
  }
  // dismiss all wizards before each test to have a consistent state
  // this can be overridden in individual tests when needed
  QUnit.hooks.beforeEach(function () {
    window.localStorage.setItem(DISMISSED_WIZARD_KEY, JSON.stringify(Object.values(WIZARD_ID_MAP)));
  });
  start({
    setupTestIsolationValidation: true,
  });
});
