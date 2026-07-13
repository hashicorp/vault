/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  egpPaths: '[data-test-egp-paths]',
};

const aclModel = {
  name: 'my-acl-policy',
  policy: 'path "secret/*" { capabilities = ["read"] }',
  policyType: 'acl',
  format: 'hcl',
  enforcement_level: null,
  paths: null,
  capabilities: { canUpdate: true, canDelete: true },
};

const rgpModel = {
  name: 'my-rgp-policy',
  policy: 'import "strings"\nmain = rule { true }',
  policyType: 'rgp',
  format: 'sentinel',
  enforcement_level: 'advisory',
  paths: ['/sys/mounts/*'],
  capabilities: { canUpdate: true, canDelete: true },
};

const renderComponent = async () => await render(hbs`<Page::PolicyShow @model={{this.model}} />`);

module('Integration | Component | page/policy-show', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('model', aclModel);
  });

  module('ACL policy', function () {
    test('it renders the policy name in the header', async function (assert) {
      await renderComponent();
      assert.ok(GENERAL.hdsPageHeaderTitle, 'my-acl-policy');
    });

    test('it renders the policy content in a code block', async function (assert) {
      await renderComponent();
      assert
        .dom('.hds-code-block')
        .includesText(
          'path "secret/*" { capabilities = ["read"] }',
          'Policy content is rendered in a code block'
        );
    });

    test('it shows the automation snippets accordion', async function (assert) {
      await renderComponent();
      assert
        .dom(GENERAL.accordionButton('Automation snippets'))
        .exists('Automation snippets accordion is rendered');
    });

    test('it does not render EGP paths', async function (assert) {
      await renderComponent();
      assert.dom(SELECTORS.egpPaths).doesNotExist('paths list is not rendered for ACL policies');
    });
  });

  module('non-ACL policy (RGP/EGP)', function (hooks) {
    hooks.beforeEach(function () {
      this.set('model', rgpModel);
    });

    test('it renders the enforcement level badge', async function (assert) {
      await renderComponent();
      assert.dom('[aria-label="Enforcement level: advisory"]').exists('enforcement level badge is rendered');
    });

    test('it does not show the automation snippets accordion', async function (assert) {
      await renderComponent();
      assert
        .dom(GENERAL.accordionButton('Automation snippets'))
        .doesNotExist('Automation snippets accordion is not rendered for non-ACL policies');
    });

    test('it renders policy paths when present', async function (assert) {
      await renderComponent();
      assert.dom(SELECTORS.egpPaths).exists('paths list is rendered for non-ACL policies');
      assert.dom(`${SELECTORS.egpPaths} li`).hasText('/sys/mounts/*', 'path value is displayed');
    });
  });
});
