/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';

import authPage from 'vault/tests/pages/auth';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

module('Acceptance | aws | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashDangerSpy = spy(flash, 'danger');

    this.uid = uuidv4();
    return authPage.login();
  });

  test('it should prompt configuration after mounting the engine', async function (assert) {
    const path = `aws-${this.uid}`;
    await visit('/vault/settings/mount-secret-backend');
    await click(SES.mountType('aws'));
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(SES.mountSubmit);
    await click(SES.configTab);

    assert.dom(GENERAL.emptyStateTitle).hasText('AWS not configured');
    assert.dom(GENERAL.emptyStateActions).hasText('Configure AWS');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  // 'it should transition to configure page on Configure click from toolbar',
  // ARG TODO stopped here.
  // 'it should prompt configuration after mounting the engine', async
  // 'it should show configured details and still allow you to edit configuration', async
});
