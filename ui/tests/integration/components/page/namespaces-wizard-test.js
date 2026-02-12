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
  guidedStart: '[data-test-guided-start]',
  stepTitle: '[data-test-step-title]',
  tree: '[data-test-tree]',
  intro: '[data-test-intro]',
  inputRow: (index) => (index ? `[data-test-input-row="${index}"]` : '[data-test-input-row]'),
};

module('Integration | Component | page/namespaces | Namespace Wizard', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.refreshSpy = sinon.spy();
    this.wizardService = this.owner.lookup('service:wizard');

    this.renderComponent = () => {
      return render(hbs`
        <Wizard::Namespaces::NamespaceWizard
          @onRefresh={{this.refreshSpy}}
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
    assert.dom(SELECTORS.intro).exists('Wizard intro is rendered');
  });

  test('it progresses through wizard steps with strict policy', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided start'));

    // Step 1: Choose security policy
    assert.dom(GENERAL.button('Next')).isDisabled('Next button disabled with no policy choice');
    await click(GENERAL.radioByAttr('strict'));
    assert.dom(GENERAL.button('Next')).isNotDisabled('Next button enabled after policy selection');
    await click(GENERAL.button('Next'));

    // Step 2: Add namespace data
    assert.dom(SELECTORS.stepTitle).hasText('Map out your namespaces');
    assert.dom(GENERAL.button('Next')).isDisabled('Next button disabled with no namespace data');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('global-0')}`, 'global');
    await click(GENERAL.button('Next'));

    // Step 3: Choose implementation method
    assert.dom(SELECTORS.stepTitle).hasText('Choose your implementation method');
    assert.dom(GENERAL.copyButton).exists('Copy button exists for code snippets');
  });

  test('it skips step 2 with flexible policy', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided start'));

    // Step 1: Choose flexible policy
    await click(GENERAL.radioByAttr('flexible'));
    await click(GENERAL.button('Next'));

    // Should skip directly to step 3 (final step)
    assert.dom(SELECTORS.stepTitle).hasText(`No action needed, you're all set.`);
    assert.dom(GENERAL.button('identities')).exists('Link to identities exists');
    assert.dom(GENERAL.button('Done')).exists('Done button exists for flexible policy');
  });

  test('it shows different code snippets per creation method option', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided start'));
    await click(GENERAL.radioByAttr('strict'));
    await click(GENERAL.button('Next'));

    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('global-0')}`, 'global');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-0')}`, 'org1');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-0')}`, 'proj1');
    await click(GENERAL.button('Next'));

    // Assert code snippet changes
    assert.dom(GENERAL.radioCardByAttr('Terraform automation')).exists('Terraform option exists');
    assert
      .dom(GENERAL.fieldByAttr('terraform'))
      .hasTextContaining(`variable "global_child_namespaces"`, 'shows terraform code snippet by default');

    await click(GENERAL.radioCardByAttr('API/CLI'));
    assert
      .dom(GENERAL.fieldByAttr('api'))
      .hasTextContaining(`curl`, 'shows API code snippet by default for API/CLI radio card');
    await click(GENERAL.hdsTab('cli'));
    assert
      .dom(GENERAL.fieldByAttr('cli'))
      .hasTextContaining(`vault namespace create`, 'shows CLI code snippet by for CLI tab');

    await click(GENERAL.radioCardByAttr('Vault UI workflow'));
    assert.dom(GENERAL.fieldByAttr('terraform')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('api')).doesNotExist();
    assert.dom(GENERAL.fieldByAttr('cli')).doesNotExist();
  });

  test('it allows adding and removing blocks, org, and project inputs', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided start'));
    await click(GENERAL.radioByAttr('strict'));
    await click(GENERAL.button('Next'));

    // Test adding and removing a second namespace block
    await click(GENERAL.button('add namespace'));
    assert.dom(`${SELECTORS.inputRow(1)}`).exists('Second input block exists');
    await click(`${SELECTORS.inputRow(1)} ${GENERAL.button('delete namespace')}`);
    assert.dom(`${SELECTORS.inputRow(1)}`).doesNotExist('Second input block was removed');

    // Test adding and removing project input
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('add project')}`);
    assert
      .dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-1')}`)
      .exists('Second project input was added');
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('delete project')}`);
    assert
      .dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-1')}`)
      .doesNotExist('Second project input was removed');

    // Test adding and removing org input
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('add org')}`);
    assert
      .dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-1')}`)
      .exists('Second org input was added');
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('delete org')}`);
    assert
      .dom(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-1')}`)
      .doesNotExist('Second org input was removed');
  });

  test('it dismisses from the intro page', async function (assert) {
    await this.renderComponent();

    assert.false(this.wizardService.isDismissed('namespace'), 'Wizard is not dismissed initially');

    await click(GENERAL.button('Skip'));

    assert.true(this.wizardService.isDismissed('namespace'), 'Wizard was marked as dismissed in service');
    assert.true(this.refreshSpy.calledOnce, 'onRefresh callback was called');
  });

  test('it dismisses from the Guided start', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided start'));

    assert.false(this.wizardService.isDismissed('namespace'), 'Wizard is not dismissed initially');

    await click(GENERAL.button('Exit'));

    assert.true(this.wizardService.isDismissed('namespace'), 'Wizard was marked as dismissed in service');
    assert.true(this.refreshSpy.calledOnce, 'onRefresh callback was called');
  });

  test('it dismisses after completing flexible policy flow', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided start'));
    await click(GENERAL.radioByAttr('flexible'));
    await click(GENERAL.button('Next'));

    assert.false(this.wizardService.isDismissed('namespace'), 'Wizard not dismissed before clicking Done');

    await click(GENERAL.button('Done'));

    assert.true(this.wizardService.isDismissed('namespace'), 'Wizard was marked as dismissed after Done');
    assert.true(this.refreshSpy.calledOnce, 'onRefresh callback was called');
  });

  test('it shows tree chart only when there are multiple globals, orgs, or projects', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.button('Guided start'));
    await click(GENERAL.radioByAttr('strict'));
    await click(GENERAL.button('Next'));

    // Initially with only one global and one org/project, tree should not show
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('global-0')}`, 'global1');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-0')}`, 'org1');
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-0')}`, 'proj1');
    assert.dom(SELECTORS.tree).doesNotExist('Tree chart hidden with single global, org, and project');

    // Add a second global namespace - tree should now show
    await click(GENERAL.button('add namespace'));
    await fillIn(`${SELECTORS.inputRow(1)} ${GENERAL.inputByAttr('global-1')}`, 'global2');
    assert.dom(SELECTORS.tree).exists('Tree chart shows with multiple globals');

    // Remove second global - tree is hidden again
    await click(`${SELECTORS.inputRow(1)} ${GENERAL.button('delete namespace')}`);
    assert.dom(SELECTORS.tree).doesNotExist('Tree chart hidden after removing second global');

    // Add a second org - tree should show
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('add org')}`);
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('org-1')}`, 'org2');
    assert.dom(SELECTORS.tree).exists('Tree chart shows with multiple orgs');

    // Remove second org - tree is hidden again
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('delete org')}`);
    assert.dom(SELECTORS.tree).doesNotExist('Tree chart hidden after removing second org');

    // Add a second project - tree should show
    await click(`${SELECTORS.inputRow(0)} ${GENERAL.button('add project')}`);
    await fillIn(`${SELECTORS.inputRow(0)} ${GENERAL.inputByAttr('project-1')}`, 'project2');
    assert.dom(SELECTORS.tree).exists('Tree chart shows with multiple projects');
  });
});
