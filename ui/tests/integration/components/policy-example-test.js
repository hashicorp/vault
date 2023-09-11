/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const SELECTORS = {
  policyText: '[data-test-modal-title]',
  policyDescription: (type) => `[data-test-example-modal-text=${type}]`,
  jsonText: '[data-test-component="code-mirror-modifier"]',
  informationLink: '[data-test-example-modal-information-link]',
};

module('Integration | Component | policy-example', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the correct paragraph for ACL policy', async function (assert) {
    await render(hbs`
    <PolicyExample
      @policyType="acl"
    />
    `);
    assert
      .dom(SELECTORS.policyDescription('acl'))
      .hasText(
        'ACL Policies are written in Hashicorp Configuration Language ( HCL ) or JSON and describe which paths in Vault a user or machine is allowed to access. Here is an example policy:'
      );
  });

  test('it renders the correct paragraph for RGP policy', async function (assert) {
    await render(hbs`
    <PolicyExample
      @policyType="rgp"
    />
    `);
    assert
      .dom(SELECTORS.policyDescription('rgp'))
      .hasText(
        'Role Governing Policies (RGPs) are tied to client tokens or identities which is similar to ACL policies . They use Sentinel as a language framework to enable fine-grained policy decisions.'
      );
  });

  test('it renders the correct paragraph for EGP policy', async function (assert) {
    await render(hbs`
    <PolicyExample
      @policyType="egp"
    />
    `);
    assert
      .dom(SELECTORS.policyDescription('egp'))
      .hasText(
        `Endpoint Governing Policies (EGPs) are tied to particular paths (e.g. aws/creds/ ) instead of tokens. They use Sentinel as a language to access properties of the incoming requests.`
      );
  });

  test('it renders the correct JSON editor text for ACL policy', async function (assert) {
    await render(hbs`
    <PolicyExample
      @policyType="acl"
    />
    `);
    assert.dom(SELECTORS.jsonText).includesText(`# Grant 'create', 'read' , 'update', and ‘list’ permission`);
  });

  test('it renders the correct JSON editor text for RGP policy', async function (assert) {
    await render(hbs`
    <PolicyExample
      @policyType="rgp"
    />
    `);
    assert
      .dom(SELECTORS.jsonText)
      .includesText(`# Import strings library that exposes common string operations`);
  });

  test('it renders the correct JSON editor text for EGP policy', async function (assert) {
    await render(hbs`
    <PolicyExample
      @policyType="egp"
    />
    `);
    assert.dom(SELECTORS.jsonText).includesText(`# Expect requests to only happen during work days (Monday`);
  });
});
