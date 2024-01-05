/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from 'vault/app';
import config from 'vault/config/environment';
import * as QUnit from 'qunit';
import { setApplication } from '@ember/test-helpers';
import { setup } from 'qunit-dom';
import { start } from 'ember-qunit';
import './helpers/flash-message';
import preloadAssets from 'ember-asset-loader/test-support/preload-assets';
import { setupGlobalA11yHooks, setRunOptions } from 'ember-a11y-testing/test-support';
import manifest from 'vault/config/asset-manifest';

preloadAssets(manifest).then(() => {
  setup(QUnit.assert);
  setApplication(Application.create(config.APP));
  setupGlobalA11yHooks(() => true, {
    helpers: ['render'],
  });
  setRunOptions({
    runOnly: {
      type: 'tag',
      values: ['wcag2a'],
    },
  });

  start({
    setupTestIsolationValidation: true,
  });
});
