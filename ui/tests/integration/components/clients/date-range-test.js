/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';
import timestamp from 'core/utils/timestamp';
import { format } from 'date-fns';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

const DATE_RANGE = CLIENT_COUNT.dateRange;
module('Integration | Component | clients/date-range', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    Sinon.replace(timestamp, 'now', Sinon.fake.returns(new Date('2018-04-03T14:15:30')));
    this.now = timestamp.now();
    this.startTimestamp = '2018-01-01T14:15:30';
    this.endTimestamp = '2019-01-31T14:15:30';
    this.billingStartTime = '';
    this.retentionMonths = 48;
    this.onChange = Sinon.spy();
    this.setEditModalVisible = Sinon.stub().callsFake((visible) => {
      this.set('showEditModal', visible);
    });
    this.showEditModal = false;
    this.renderComponent = async () => {
      await render(
        hbs`<Clients::DateRange @startTimestamp={{this.startTimestamp}} @endTimestamp={{this.endTimestamp}} @onChange={{this.onChange}} @billingStartTime={{this.billingStartTime}} @retentionMonths={{this.retentionMonths}} @setEditModalVisible={{this.setEditModalVisible}} @showEditModal={{this.showEditModal}}/>`
      );
    };
  });

  test('it formats modal inputs to ISO string timestamps', async function (assert) {
    this.startTimestamp = undefined;
    await this.renderComponent();

    assert.dom(DATE_RANGE.dateDisplay('start')).doesNotExist();
    assert.dom(DATE_RANGE.dateDisplay('end')).doesNotExist();
    assert.dom(DATE_RANGE.edit).hasText('Set date range');

    await click(DATE_RANGE.edit);
    assert.dom(DATE_RANGE.editModal).exists();
    assert.dom(DATE_RANGE.editDate('start')).hasValue('');
    await fillIn(DATE_RANGE.editDate('start'), '2018-01');
    await fillIn(DATE_RANGE.editDate('end'), '2018-03');
    await click(GENERAL.submitButton);
    const { start_time, end_time } = this.onChange.lastCall.args[0];
    assert.strictEqual(start_time, '2018-01-01T00:00:00Z', 'it formats start_time param');
    assert.strictEqual(end_time, '2018-03-31T23:59:59Z', 'it formats end_time param');
    assert.dom(DATE_RANGE.editModal).doesNotExist('closes modal');
  });

  test('it does not trigger onChange if dates are invalid', async function (assert) {
    this.owner.lookup('service:version').type = 'community';
    this.endTimestamp = undefined;
    await this.renderComponent();
    await click(DATE_RANGE.edit);
    await fillIn(DATE_RANGE.editDate('end'), '');
    assert.dom(DATE_RANGE.validation).hasText('You must supply both start and end dates.');
    await click(GENERAL.submitButton);
    assert.false(this.onChange.called);

    await fillIn(DATE_RANGE.editDate('start'), '2018-01');
    await fillIn(DATE_RANGE.editDate('end'), '2018-04');
    assert.dom(DATE_RANGE.validation).hasText('You cannot select the current month or beyond.');
    await fillIn(DATE_RANGE.editDate('end'), '2017-05');
    assert.dom(DATE_RANGE.validation).hasText('Start date must be before end date.');
    await click(GENERAL.submitButton);
    assert.false(this.onChange.called);

    await click(GENERAL.cancelButton);
    assert.false(this.onChange.called);
    assert.dom(DATE_RANGE.editModal).doesNotExist();
  });

  test('it does not allow the current month to be selected as a start date or as an end date', async function (assert) {
    this.owner.lookup('service:version').type = 'community';
    this.endTimestamp = undefined;
    const currentMonth = format(timestamp.now(), 'yyyy-MM');

    await this.renderComponent();
    await click(DATE_RANGE.edit);
    await fillIn(DATE_RANGE.editDate('start'), currentMonth);
    await fillIn(DATE_RANGE.editDate('end'), currentMonth);

    assert.dom(DATE_RANGE.validation).hasText('You cannot select the current month or beyond.');
    await click(GENERAL.submitButton);
    assert.false(this.onChange.called);

    //  This tests validation when the end date is the current month and start is valid.
    //  If start is current month and end is a valid prior selection, it will run into the validation error of start being after end date
    //  which is covered by prior tests.
    await fillIn(DATE_RANGE.editDate('start'), '2018-01');
    await fillIn(DATE_RANGE.editDate('end'), currentMonth);
    await click(GENERAL.submitButton);
    assert.false(this.onChange.called);
  });

  module('enterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version = this.owner.lookup('service:version');
      this.version.type = 'enterprise';
      this.billingStartTime = '2018-01-01T14:15:30';
    });

    test('it billing start date dropdown for enterprise', async function (assert) {
      await this.renderComponent();
      await click(DATE_RANGE.edit);
      const expectedPeriods = [
        'January 2018',
        'January 2017',
        'January 2016',
        'January 2015',
        'January 2014',
      ];
      const dropdownList = findAll(DATE_RANGE.dropdownOption(null));
      dropdownList.forEach((item, idx) => {
        const month = expectedPeriods[idx];
        assert.dom(item).hasText(month, `dropdown index: ${idx} renders ${month}`);
      });
    });
  });
});
