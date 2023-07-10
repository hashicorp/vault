/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

// To-Do: add to tests
const SELECTORS = {
  policyText: '[data-test-modal-title]',
  aclPolicyParagraph: '[data-test-example-modal-text-acl]',
  rgpPolicyParagraph: '[data-test-example-modal-text-rgp]',
  jsonText: '[data-test-example-modal-json-text]',
  informationLink: '[data-test-example-modal-information-link]',
};

module('Integration | Component | policy-example', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it renders the correct paragraph for ACL policy', async function (assert) {
    this.model = this.store.createRecord('policy/acl');

    await render(hbs`
    <PolicyExample
      @policyType={{this.model.policyType}}
    />
    `);
    assert
      .dom(SELECTORS.aclPolicyParagraph)
      .hasText(
        'ACL Policies are written in Hashicorp Configuration Language ( HCL ) or JSON and describe which paths in Vault a user or machine is allowed to access. Here is an example policy:'
      );
  });

  test('it renders the correct paragraph for RGP policy', async function (assert) {
    this.model = this.store.createRecord('policy/rgp');

    await render(hbs`
    <PolicyExample
      @policyType={{this.model.policyType}}
    />
    `);
    assert
      .dom(SELECTORS.rgpPolicyParagraph)
      .hasText(
        'Role Governing Policies (RGPs) are tied to client tokens or identities which is similar to ACL policies . They use Sentinel as a language framework to enable fine-grained policy decisions.'
      );
  });

  test('it renders the correct JSON editor text for ACL policy', async function (assert) {
    this.model = this.store.createRecord('policy/acl');

    await render(hbs`
    <PolicyExample
      @policyType={{this.model.policyType}}
    />
    `);
    //await this.pauseTest();
    assert.dom(SELECTORS.jsonText).includesText(`# Grant 'create', 'read' , 'update', and ‘list’ permission`);
  });

  test('it renders the correct JSON editor text for RGP policy', async function (assert) {
    this.model = this.store.createRecord('policy/rgp');

    await render(hbs`
    <PolicyExample
      @policyType={{this.model.policyType}}
    />
    `);
    assert
      .dom(SELECTORS.jsonText)
      .includesText(`# Import strings library that exposes common string operations`);
  });
});
