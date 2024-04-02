/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from 'vault/app';
import config from 'vault/config/environment';
import * as QUnit from 'qunit';
import { setApplication } from '@ember/test-helpers';
import { setup } from 'qunit-dom';
import start from 'ember-exam/test-support/start';
import './helpers/flash-message';
import { setupGlobalA11yHooks, setRunOptions } from 'ember-a11y-testing/test-support';

setApplication(Application.create(config.APP));

setup(QUnit.assert);
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
