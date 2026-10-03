/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  content: '[data-test-content]',
  guidedSetup: '[data-test-guided-setup]',
  stepTitle: '[data-test-step-title]',
  welcome: '[data-test-welcome]',
  inputRow: (index) => (index ? `[data-test-input-row="${index}"]` : '[data-test-input-row]'),
};

module('Integration | Component | page/namespaces | Namespace Wizard', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.onFilterChange = sinon.spy();
    this.onDismiss = sinon.spy();
    this.onRefresh = sinon.spy();

    this.renderComponent = () => {
      return render(hbs`
        <Wizard::Namespaces::NamespaceWizard 
          @onDismiss={{this.onDismiss}}
          @onRefresh={{this.onRefresh}}
        />
      `);
    };
  });

  hooks.afterEach(async function () {
    // ensure clean state
    localStorage.clear();
  });

  test('it shows wizard when no namespaces exist', async function (assert) {
    await this.renderComponent();
    assert.dom(SELECTORS.welcome).exists('Wizard welcome is rendered');
  });

  test('it progresses through wizard steps with strict policy', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided setup'));

    // Step 1: Choose security policy
    assert.dom(GENERAL.button('Next')).isDisabled('Next button disabled with no policy choice');
    await click(GENERAL.radioByAttr('strict'));
    await click(GENERAL.button('Next'));

    // Step 2: Add namespace data
    assert.dom(SELECTORS.stepTitle).hasText('Map out your namespaces');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('global-0')}`, 'global');

    await click(GENERAL.button('Next'));

    // Step 3: Choose implementation method
    assert.dom(SELECTORS.stepTitle).hasText('Choose your implementation method');
    assert.dom(GENERAL.copyButton).exists();
  });

  test('it skips step 2 with flexible policy', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided setup'));

    // Step 1: Choose flexible policy
    await click(GENERAL.radioByAttr('flexible'));
    await click(GENERAL.button('Next'));

    // Should skip directly to step 3
    assert.dom(SELECTORS.stepTitle).hasText(`No action needed, you're all set.`);
    assert.dom(GENERAL.button('identities')).exists();
  });

  test('it shows different code snippets per creation method option', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided setup'));
    await click(GENERAL.radioByAttr('strict'));
    await click(GENERAL.button('Next'));

    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('global-0')}`, 'global');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-0')}`, 'org1');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-0')}`, 'proj1');
    await click(GENERAL.button('Next'));

    // Assert code snippet changes
    assert.dom(GENERAL.radioCardByAttr('Terraform automation')).exists('Terraform option exists');
    assert
      .dom(GENERAL.fieldByAttr('snippets'))
      .hasTextContaining(`variable "global_child_namespaces"`, 'shows terraform code snippet by default');

    await click(GENERAL.radioCardByAttr('API/CLI'));
    assert
      .dom(GENERAL.fieldByAttr('snippets'))
      .hasTextContaining(`curl`, 'shows API code snippet by default for API/CLI radio card');

    await click(GENERAL.hdsTab('CLI'));
    assert
      .dom(GENERAL.fieldByAttr('snippets'))
      .hasTextContaining(`vault namespace create`, 'shows CLI code snippet by for CLI tab');

    await click(GENERAL.radioCardByAttr('Vault UI workflow'));
    assert.dom(GENERAL.fieldByAttr('snippets')).doesNotExist('does not render a code snippet for UI flow');
  });

  test('it allows adding and removing blocks, org, and project inputs', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided setup'));
    await click(GENERAL.radioByAttr('strict'));
    await click(GENERAL.button('Next'));

    // Add a second block
    await click(GENERAL.button('add namespace'));
    assert.dom(`${SELECTORS.inputRow(1)}`).exists('Second input block exists');
    await click(`${SELECTORS.inputRow(1)} ${GENERAL.button('delete namespace')}`);
    assert.dom(`${SELECTORS.inputRow(1)}`).doesNotExist('Second input block was removed');

    // Test adding and removing project input
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('add project')}`);
    assert
      .dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-1')}`)
      .exists('project input was added');
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('delete project')}`);
    assert
      .dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-1')}`)
      .doesNotExist('project input was removed');

    // Test adding and removing org input
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('add org')}`);
    assert.dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-1')}`).exists('org input was added');
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('delete org')}`);
    assert
      .dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-1')}`)
      .doesNotExist('org input was removed');
  });
});
