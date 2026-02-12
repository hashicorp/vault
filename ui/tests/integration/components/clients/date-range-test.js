/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';
import timestamp from 'core/utils/timestamp';
import { format, subYears } from 'date-fns';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

const DATE_RANGE = CLIENT_COUNT.dateRange;
module('Integration | Component | clients/date-range', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    Sinon.replace(timestamp, 'now', Sinon.fake.returns(new Date('2018-04-03T14:15:30')));
    this.now = timestamp.now();
    this.startTimestamp = new Date('2018-01-01T14:15:30');
    this.endTimestamp = new Date('2019-01-31T14:15:30');
    this.billingStartTime = '';
    this.retentionMonths = 48;
    this.onChange = Sinon.stub();
    this.setEditModalVisible = Sinon.stub().callsFake((visible) => {
      this.set('showEditModal', visible);
    });
    this.showEditModal = false;
    this.renderComponent = async () => {
      await render(
        hbs`<Clients::DateRange
          @startTimestamp={{this.startTimestamp}}
          @endTimestamp={{this.endTimestamp}}
          @onChange={{this.onChange}}
          @billingStartTime={{this.billingStartTime}}
          @retentionMonths={{this.retentionMonths}}
          @setEditModalVisible={{this.setEditModalVisible}}
          @showEditModal={{this.showEditModal}}
        />`
      );
    };
  });

  test('it does not render if a start and end timestamp are already provided', async function (assert) {
    await this.renderComponent();
    assert.dom(DATE_RANGE.edit).doesNotExist('it does not render if timestamps are provided');
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
    this.version.type = 'community';
    this.endTimestamp = undefined;
    const currentMonth = format(timestamp.now(), 'yyyy-MM');

    await this.renderComponent();
    await click(DATE_RANGE.edit);
    await fillIn(DATE_RANGE.editDate('start'), currentMonth);
    await fillIn(DATE_RANGE.editDate('end'), currentMonth);

    assert.dom(DATE_RANGE.validation).hasText('You cannot select the current month or beyond.');
    await click(GENERAL.submitButton);
    assert.false(this.onChange.called, 'it does not call @onChange callback');

    //  This tests validation when the end date is the current month and start is valid.
    //  If start is current month and end is a valid prior selection, it will run into the validation error of start being after end date
    //  which is covered by prior tests.
    await fillIn(DATE_RANGE.editDate('start'), '2018-01');
    await fillIn(DATE_RANGE.editDate('end'), currentMonth);
    await click(GENERAL.submitButton);
    assert.false(this.onChange.called, 'it does not call @onChange callback');
  });

  test('it allows the current month to be selected if enterprise and there is not a @billingStartTime', async function (assert) {
    this.version.type = 'enterprise';
    this.endTimestamp = undefined;
    const currentMonth = format(timestamp.now(), 'yyyy-MM');

    await this.renderComponent();
    await click(DATE_RANGE.edit);
    await fillIn(DATE_RANGE.editDate('start'), currentMonth);
    await fillIn(DATE_RANGE.editDate('end'), currentMonth);

    assert.dom(DATE_RANGE.validation).doesNotExist();
    await click(GENERAL.submitButton);
    assert.true(this.onChange.called, 'it calls @onChange callback');
  });

  module('enterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version = this.owner.lookup('service:version');
      this.version.type = 'enterprise';
      this.billingStartTime = new Date('2018-01-01T14:15:30');
    });

    test('it renders billing start date dropdown for enterprise', async function (assert) {
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

    test('it renders date range modal if there are no timestamps provided', async function (assert) {
      this.billingStartTime = '';
      this.startTimestamp = '';
      this.endTimestamp = '';
      await this.renderComponent();
      assert
        .dom(DATE_RANGE.edit)
        .exists('it renders button to open date range modal')
        .hasText('Set date range');
      await click(DATE_RANGE.edit);
      assert.dom(DATE_RANGE.editModal).exists();
    });

    test('it updates toggle text when a new date is selected', async function (assert) {
      this.onChange.callsFake(({ start_time }) => this.set('startTimestamp', new Date(start_time)));

      await this.renderComponent();
      assert.dom(DATE_RANGE.edit).hasText('January 2018').hasAttribute('aria-expanded', 'false');
      await click(DATE_RANGE.edit);
      assert.dom(DATE_RANGE.edit).hasAttribute('aria-expanded', 'true');
      await click(DATE_RANGE.dropdownOption(1));
      assert
        .dom(DATE_RANGE.edit)
        .hasText('January 2017')
        .hasAttribute('aria-expanded', 'false', 'it closes dropdown after selection');
    });

    test('it renders billing period text', async function (assert) {
      await this.renderComponent();
      assert
        .dom(this.element)
        .hasText('Change billing period January 2018', 'it renders billing related text');
    });

    test('it renders data period text for HVD managed clusters', async function (assert) {
      this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
      await this.renderComponent();
      assert.dom(this.element).hasText('Change data period January 2018');
    });

    test('it should send an empty string for start_time when selecting current period', async function (assert) {
      await this.renderComponent();

      await click(DATE_RANGE.edit);
      await click(DATE_RANGE.dropdownOption(1));
      assert.true(
        this.onChange.calledWith({
          start_time: subYears(this.billingStartTime, 1).toISOString(),
          end_time: '',
        }),
        'correct start_time sent on change for prior period'
      );

      await click(DATE_RANGE.edit);
      await click(DATE_RANGE.dropdownOption(0));
      assert.true(
        this.onChange.calledWith({ start_time: '', end_time: '' }),
        'start_time is empty string on current period change'
      );
    });
  });
});
