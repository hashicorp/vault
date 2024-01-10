/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteAuthCmd, deleteEngineCmd, mountAuthCmd, mountEngineCmd, runCmd } from '../helpers/commands';
import { authEngineHelper, secretEngineHelper } from '../helpers/openapi/test-helpers';

/**
 * This set of tests is for ensuring that backend changes to the OpenAPI spec
 * are known by UI developers and adequately addressed in the UI. When changes
 * are detected from this set of tests, they should be updated to pass and
 * smoke tested to ensure changes to not break the GUI workflow.
 * Marked as enterprise so it only runs periodically
 */
module('Acceptance | OpenAPI provides expected attributes enterprise', function (hooks) {
  setupApplicationTest(hooks);
  hooks.beforeEach(function () {
    this.pathHelp = this.owner.lookup('service:pathHelp');
    this.store = this.owner.lookup('service:store');
    return authPage.login();
  });

  // Secret engines that use OpenAPI
  ['ssh', 'kmip', 'pki'].forEach(function (testCase) {
    return module(`${testCase} engine`, function (hooks) {
      hooks.beforeEach(async function () {
        this.backend = `${testCase}-openapi`;
        await runCmd(mountEngineCmd(testCase, this.backend), false);
      });
      hooks.afterEach(async function () {
        await runCmd(deleteEngineCmd(this.backend), false);
      });

      secretEngineHelper(test, testCase);
    });
  });

  // All auth backends use OpenAPI except aws
  ['azure', 'userpass', 'cert', 'gcp', 'github', 'jwt', 'kubernetes', 'ldap', 'okta', 'radius'].forEach(
    function (testCase) {
      return module(`${testCase} auth`, function (hooks) {
        hooks.beforeEach(async function () {
          this.mount = `${testCase}-openapi`;
          await runCmd(mountAuthCmd(testCase, this.mount), false);
        });
        hooks.afterEach(async function () {
          await runCmd(deleteAuthCmd(this.backend), false);
        });

        authEngineHelper(test, testCase);
      });
    }
  );
});
