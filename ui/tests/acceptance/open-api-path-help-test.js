/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteAuthCmd, deleteEngineCmd, mountAuthCmd, mountEngineCmd, runCmd } from '../helpers/commands';
import expectedSecretAttrs from 'vault/tests/helpers/openapi/expected-secret-attrs';
import expectedAuthAttrs from 'vault/tests/helpers/openapi/expected-auth-attrs';
import { getHelpUrlForModel } from 'vault/utils/openapi-helpers';

/**
 * This set of tests is for ensuring that backend changes to the OpenAPI spec
 * are known by UI developers and adequately addressed in the UI. When changes
 * are detected from this set of tests, they should be updated to pass and
 * smoke tested to ensure changes to not break the GUI workflow.
 * In some cases, a ticket should be made to track updating the relevant model or form,
 * if it is not updated automatically or is a more involved feature request.
 * Marked as enterprise so it only runs periodically
 */
module(
  'Acceptance | Heads up - backend param changes! Expected OpenAPI attributes enterprise',
  function (hooks) {
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
  }
);

function secretEngineHelper(test, secretEngine) {
  const engineData = expectedSecretAttrs[secretEngine];
  if (!engineData)
    throw new Error(`No engine attributes found in secret-model-attributes for ${secretEngine}`);

  const modelNames = Object.keys(engineData);
  // A given secret engine might have multiple models that are openApi driven
  modelNames.forEach((modelName) => {
    test(`${modelName} model getProps returns correct attributes`, async function (assert) {
      const helpUrl = getHelpUrlForModel(modelName, this.backend);
      const result = await this.pathHelp.getProps(helpUrl, this.backend);
      const expected = engineData[modelName];
      // Expected values should be updated to match "actual" (result)
      assert.deepEqual(
        Object.keys(result).sort(),
        Object.keys(expected).sort(),
        `getProps returns expected attributes for ${modelName} (help url: "${helpUrl}")`
      );
      Object.keys(expected).forEach((attrName) => {
        assert.deepEqual(result[attrName], expected[attrName], `${attrName} attribute details match`);
      });
    });
  });
}

function authEngineHelper(test, authBackend) {
  const authData = expectedAuthAttrs[authBackend];
  if (!authData) throw new Error(`No auth attributes found in auth-model-attributes for ${authBackend}`);

  const itemNames = Object.keys(authData);
  itemNames.forEach((itemName) => {
    if (itemName.startsWith('auth-config/')) {
      // Config test doesn't need to instantiate a new model
      test(`${itemName} model`, async function (assert) {
        const helpUrl = getHelpUrlForModel(itemName, this.mount);
        const result = await this.pathHelp.getProps(helpUrl, this.mount);
        const expected = authData[itemName];
        assert.deepEqual(
          Object.keys(result).sort(),
          Object.keys(expected).sort(),
          `getProps returns expected attributes for ${itemName}`
        );
        Object.keys(expected).forEach((attrName) => {
          assert.propEqual(result[attrName], expected[attrName], `${attrName} attribute details match`);
        });
      });
    } else {
      test.skip(`generated-${itemName}-${authBackend} model`, async function (assert) {
        const modelName = `generated-${itemName}-${authBackend}`;
        // Generated items need to instantiate the model first via getNewModel
        await this.pathHelp.getNewModel(modelName, this.mount, `auth/${this.mount}/`, itemName);
        // Generated items don't have helpUrl method -- helpUrl is calculated in path-help.js line 101
        const helpUrl = `/v1/auth/${this.mount}?help=1`;
        const result = await this.pathHelp.getProps(helpUrl, this.mount);
        const expected = authData[modelName];
        assert.deepEqual(result, expected, `getProps returns expected attributes for ${modelName}`);
      });
    }
  });
}
