/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { clickTrigger, typeInSearch } from 'ember-power-select/test-support/helpers';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import ss from 'vault/tests/pages/components/search-select';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const component = create(ss);

module('Integration | Component | search select with modal', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  hooks.beforeEach(function () {
    this.set('onChange', sinon.spy());
    this.server.get('sys/policies/acl', () => {
      return {
        request_id: 'acl-policy-list',
        data: {
          keys: ['default', 'root', 'acl-test'],
        },
      };
    });
    this.server.get('sys/policies/rgp', () => {
      return {
        request_id: 'rgp-policy-list',
        data: {
          keys: ['rgp-test'],
        },
      };
    });
    this.server.get('/sys/policies/acl/acl-test', () => {
      return {
        request_id: 'policy-acl',
        data: {
          name: 'acl-test',
          policy:
            '\n# Grant \'create\', \'read\' , \'update\', and ‘list’ permission\n# to paths prefixed by \'secret/*\'\npath "secret/*" {\n  capabilities = [ "create", "read", "update", "list" ]\n}\n\n# Even though we allowed secret/*, this line explicitly denies\n# secret/super-secret. This takes precedence.\npath "secret/super-secret" {\n  capabilities = ["deny"]\n}\n',
        },
      };
    });
    this.server.get('/sys/policies/rgp/rgp-test', () => {
      return {
        request_id: 'policy-rgp',
        data: {
          name: 'rgp-test',
          enforcement_level: 'hard-mandatory',
          policy:
            '\n# Import strings library that exposes common string operations\nimport "strings"\n\n# Conditional rule (precond) checks the incoming request endpoint\n# targeted to sys/policies/acl/admin\nprecond = rule {\n    strings.has_prefix(request.path, "sys/policies/admin")\n}\n\n# Vault checks to see if the request was made by an entity\n# named James Thomas or Team Lead role defined as its metadata\nmain = rule when precond {\n    identity.entity.metadata.role is "Team Lead" or\n      identity.entity.name is "James Thomas"\n}\n',
        },
      };
    });
    setRunOptions({
      rules: {
        // TODO: Fix this component
        'color-contrast': { enabled: false },
        label: { enabled: false },
        'aria-input-field-name': { enabled: false },
        'aria-required-attr': { enabled: false },
        'aria-valid-attr-value': { enabled: false },
      },
    });
  });

  test('it renders passed in models', async function (assert) {
    await render(hbs`
    <SearchSelectWithModal
      @id="policies"
      @label="Policies"
      @labelClass="title is-4"
      @models={{array "policy/acl" "policy/rgp"}}
      @inputValue={{this.policies}}
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormTemplate="modal-form/policy-template"
      @excludeOptions={{array "root"}}
      @subText="Some modal subtext"
    />
      `);
    assert.dom('[data-test-search-select-with-modal]').exists('the component renders');
    assert.dom('[data-test-modal-subtext]').hasText('Some modal subtext', 'renders modal text');
    assert.strictEqual(component.labelText, 'Policies', 'label text is correct');
    assert.ok(component.hasTrigger, 'it renders the power select trigger');
    assert.strictEqual(component.selectedOptions.length, 0, 'there are no selected options');

    await clickTrigger();
    const dropdownOptions = findAll('[data-option-index]').map((o) => o.innerText);
    assert.notOk(dropdownOptions.includes('root'), 'root policy is not listed as option');
    assert.strictEqual(component.options.length, 3, 'dropdown renders passed in models as options');
    assert.ok(this.onChange.notCalled, 'onChange is not called');
  });

  test('it renders input value', async function (assert) {
    this.policies = ['acl-test'];
    await render(hbs`
    <SearchSelectWithModal
      @id="policies"
      @label="Policies"
      @labelClass="title is-4"
      @models={{array "policy/acl" "policy/rgp"}}
      @inputValue={{this.policies}}
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormTemplate="modal-form/policy-template"
      @subText="Some modal subtext"
    />
  `);
    assert.strictEqual(component.selectedOptions.length, 1, 'there is one selected option');
    assert.strictEqual(component.selectedOptions.objectAt(0).text, 'acl-test', 'renders inputted policies');

    await clickTrigger();
    assert.strictEqual(component.options.length, 3, 'does not render all options returned from query');
    const dropdownOptions = findAll('[data-option-index]').map((o) => o.innerText);
    assert.notOk(dropdownOptions.includes('acl-test'), 'selected option is not included in the dropdown');
    assert.ok(this.onChange.notCalled, 'onChange is not called');
  });

  test('it filters options, shows option to create new item and opens modal on select', async function (assert) {
    assert.expect(7);
    await render(hbs`
    <SearchSelectWithModal
      @id="policies"
      @label="Policies"
      @labelClass="title is-4"
      @models={{array "policy/acl" "policy/rgp"}}
      @inputValue={{this.policies}}
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormTemplate="modal-form/policy-template"
    />
      `);

    await clickTrigger();
    assert.strictEqual(component.options.length, 4, 'dropdown renders all options');

    await typeInSearch('a');
    assert.strictEqual(component.options.length, 3, 'dropdown renders all matching options plus add option');
    await typeInSearch('acl-test');
    assert.strictEqual(component.options[0].text, 'acl-test', 'dropdown renders only matching option');

    await typeInSearch('acl-test-new');
    assert.strictEqual(
      component.options[0].text,
      'No results found for "acl-test-new". Click here to create it.',
      'dropdown gives option to create new option'
    );
    await component.selectOption();

    assert.dom('#search-select-modal').exists('modal is active');
    assert.dom('[data-test-empty-state-title]').hasText('No policy type selected');
    assert.ok(this.onChange.notCalled, 'onChange is not called');
  });

  test('it renders policy template and selects policy type', async function (assert) {
    assert.expect(9);
    this.server.put('/sys/policies/acl/acl-test-new', async (schema, req) => {
      const requestBody = JSON.parse(req.requestBody);
      assert.propEqual(
        requestBody,
        {
          name: 'acl-test-new',
          policy: 'path "secret/super-secret" { capabilities = ["deny"] }',
        },
        'onSave sends request to endpoint with correct policy attributes'
      );
    });
    await render(hbs`
    <SearchSelectWithModal
      @id="policies"
      @label="Policies"
      @labelClass="title is-4"
      @models={{array "policy/acl" "policy/rgp"}}
      @inputValue={{this.policies}}
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormTemplate="modal-form/policy-template"
    />
      `);
    await clickTrigger();
    await typeInSearch('acl-test-new');
    assert.strictEqual(
      component.options[0].text,
      'No results found for "acl-test-new". Click here to create it.',
      'dropdown gives option to create new option'
    );
    await component.selectOption();
    assert.dom('[data-test-empty-state-title]').hasText('No policy type selected');
    await fillIn('[data-test-select="policyType"]', 'acl');
    assert.dom('[data-test-policy-form]').exists('policy form renders after type is selected');
    await click('[data-test-tab-example-policy] button');
    assert.dom('[data-test-tab-example-policy] button').hasAttribute('aria-selected', 'true');
    await click('[data-test-tab-your-policy] button');
    assert.dom('[data-test-tab-your-policy] button').hasAttribute('aria-selected', 'true');
    await fillIn(
      '[data-test-component="code-mirror-modifier"] textarea',
      'path "secret/super-secret" { capabilities = ["deny"] }'
    );
    await click('[data-test-policy-save]');
    assert.dom('[data-test-modal-div]').doesNotExist('modal closes after save');
    assert
      .dom('[data-test-selected-option="0"]')
      .hasText('acl-test-new', 'adds newly created policy to selected options');
    assert.ok(
      this.onChange.calledWithExactly(['acl-test-new']),
      'onChange is called only after item is created'
    );
  });

  test('it still renders search select if only second model returns 403', async function (assert) {
    assert.expect(4);
    this.server.get('sys/policies/rgp', () => {
      return new Response(
        403,
        { 'Content-Type': 'application/json' },
        JSON.stringify({ errors: ['permission denied'] })
      );
    });

    await render(hbs`
    <SearchSelectWithModal
      @id="policies"
      @label="Policies"
      @labelClass="title is-4"
      @models={{array "policy/acl" "policy/rgp"}}
      @inputValue={{this.policies}}
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormTemplate="modal-form/policy-template"
    />
      `);

    assert.dom('[data-test-search-select-with-modal]').exists('the component renders');
    assert.dom('[data-test-component="string-list"]').doesNotExist('does not render fallback component');
    await clickTrigger();
    assert.strictEqual(component.options.length, 3, 'only options from successful query render');
    assert.ok(this.onChange.notCalled, 'onChange is not called');
  });

  test('it renders fallback component if both models return 403', async function (assert) {
    assert.expect(7);
    this.server.get('sys/policies/acl', () => {
      return new Response(
        403,
        { 'Content-Type': 'application/json' },
        JSON.stringify({ errors: ['permission denied'] })
      );
    });
    this.server.get('sys/policies/rgp', () => {
      return new Response(
        403,
        { 'Content-Type': 'application/json' },
        JSON.stringify({ errors: ['permission denied'] })
      );
    });

    await render(hbs`
    <SearchSelectWithModal
      @id="policies"
      @label="Policies"
      @labelClass="title is-4"
      @models={{array "policy/acl" "policy/rgp"}}
      @inputValue={{this.policies}}
      @onChange={{this.onChange}}
      @fallbackComponent="string-list"
      @modalFormTemplate="modal-form/policy-template"
    />
      `);
    assert.dom('[data-test-component="string-list"]').exists('renders fallback component');
    assert.false(component.hasTrigger, 'does not render power select trigger');
    await fillIn('[data-test-string-list-input="0"]', 'string-list-policy');
    await click('[data-test-string-list-button="add"]');
    assert
      .dom('[data-test-string-list-input="0"]')
      .hasValue('string-list-policy', 'first row renders inputted string');
    assert
      .dom('[data-test-string-list-row="0"] [data-test-string-list-button="delete"]')
      .exists('first row renders delete icon');
    assert.dom('[data-test-string-list-row="1"]').exists('renders second input row');
    assert
      .dom('[data-test-string-list-row="1"] [data-test-string-list-button="add"]')
      .exists('second row renders add icon');
    assert.ok(
      this.onChange.calledWithExactly(['string-list-policy']),
      'onChange is called only after item is created'
    );
  });
});
