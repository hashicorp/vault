/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CLUSTER_STARTUP_CHECKLIST } from 'vault/utils/constants/checklist';

module('Integration | Component | dashboard/widgets/checklist', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.checklist = CLUSTER_STARTUP_CHECKLIST;
    this.visibleSteps = [...this.checklist.order];

    this.checklistState = this.owner.lookup('service:checklist-state');
    this.checklistState['_state'] = {};
    this.checklistState.isAvailable = true;
    this.onHide = () => {};

    this.renderComponent = () =>
      render(hbs`
        <Dashboard::Widgets::Checklist
          @checklist={{this.checklist}}
          @checklistId={{this.checklist.id}}
          @visibleSteps={{this.visibleSteps}}
          @onHide={{this.onHide}}
        />
      `);
  });

  test('clicking the toggle on an already-expanded step collapses it', async function (assert) {
    await this.renderComponent();

    const [firstStepId, secondStepId] = this.visibleSteps;

    assert.dom(GENERAL.accordionContent(firstStepId)).exists('First step is open by default');

    await click(GENERAL.accordionButton(secondStepId));
    assert.dom(GENERAL.accordionContent(secondStepId)).exists('Clicking a different step opens it');
    assert.dom(GENERAL.accordionContent(firstStepId)).doesNotExist('Previously open step collapses');

    await click(GENERAL.accordionButton(secondStepId));
    assert
      .dom(GENERAL.accordionContent(secondStepId))
      .doesNotExist('Clicking the open step again collapses it');
    assert.dom(GENERAL.accordionContent(firstStepId)).doesNotExist('No other step re-opens automatically');
  });
});
