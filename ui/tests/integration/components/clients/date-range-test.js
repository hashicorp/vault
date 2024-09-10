/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Sinon from 'sinon';
import timestamp from 'core/utils/timestamp';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

const DATE_RANGE = CLIENT_COUNT.dateRange;
module('Integration | Component | clients/date-range', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    Sinon.replace(timestamp, 'now', Sinon.fake.returns(new Date('2018-04-03T14:15:30')));
    this.now = timestamp.now();
    this.startTime = '2018-01-01T14:15:30';
    this.endTime = '2019-01-31T14:15:30';
    this.onChange = Sinon.spy();
    this.renderComponent = async () => {
      await render(
        hbs`<Clients::DateRange @startTime={{this.startTime}} @endTime={{this.endTime}} @onChange={{this.onChange}} />`
      );
    };
  });

  test('it renders prompt to set dates if no start time', async function (assert) {
    this.startTime = undefined;
    await this.renderComponent();

    assert.dom(DATE_RANGE.dateDisplay('start')).doesNotExist();
    assert.dom(DATE_RANGE.dateDisplay('end')).doesNotExist();
    assert.dom(DATE_RANGE.edit).hasText('Set date range');

    await click(DATE_RANGE.edit);
    assert.dom(DATE_RANGE.editModal).exists();
    assert.dom(DATE_RANGE.editDate('start')).hasValue('');
    await fillIn(DATE_RANGE.editDate('start'), '2018-01');
    await fillIn(DATE_RANGE.editDate('end'), '2019-01');
    await click(GENERAL.saveButton);
    assert.deepEqual(this.onChange.args[0], [
      {
        end_time: 1548892800,
        start_time: 1514764800,
      },
    ]);
    assert.dom(DATE_RANGE.editModal).doesNotExist('closes modal');
  });

  test('it renders the date range passed and can reset it (ent)', async function (assert) {
    this.owner.lookup('service:version').type = 'enterprise';
    await this.renderComponent();

    assert.dom(DATE_RANGE.dateDisplay('start')).hasText('January 2018');
    assert.dom(DATE_RANGE.dateDisplay('end')).hasText('January 2019');
    assert.dom(DATE_RANGE.edit).hasText('Edit');

    await click(DATE_RANGE.edit);
    assert.dom(DATE_RANGE.editModal).exists();
    assert.dom(DATE_RANGE.editDate('start')).hasValue('2018-01');
    assert.dom(DATE_RANGE.editDate('end')).hasValue('2019-01');
    assert.dom(DATE_RANGE.defaultRangeAlert).doesNotExist();

    await click(DATE_RANGE.editDate('reset'));
    assert.dom(DATE_RANGE.editDate('start')).hasValue('');
    assert.dom(DATE_RANGE.editDate('end')).hasValue('');
    assert.dom(DATE_RANGE.defaultRangeAlert).exists();
    await click(GENERAL.saveButton);
    assert.deepEqual(this.onChange.args[0], [{ start_time: undefined, end_time: undefined }]);
  });

  test('it renders the date range passed and cannot reset it when community', async function (assert) {
    this.owner.lookup('service:version').type = 'community';
    await this.renderComponent();

    assert.dom(DATE_RANGE.dateDisplay('start')).hasText('January 2018');
    assert.dom(DATE_RANGE.dateDisplay('end')).hasText('January 2019');
    assert.dom(DATE_RANGE.edit).hasText('Edit');

    await click(DATE_RANGE.edit);
    assert.dom(DATE_RANGE.editModal).exists();
    assert.dom(DATE_RANGE.editDate('reset')).doesNotExist();
    assert.dom(DATE_RANGE.editDate('start')).hasValue('2018-01');
    assert.dom(DATE_RANGE.editDate('end')).hasValue('2019-01');
    assert.dom(DATE_RANGE.defaultRangeAlert).doesNotExist();

    await fillIn(DATE_RANGE.editDate('start'), '');
    assert.dom(DATE_RANGE.validation).hasText('You must supply both start and end dates.');
    await click(GENERAL.saveButton);
    assert.false(this.onChange.called);
  });

  test('it does not trigger onChange if date range invalid', async function (assert) {
    this.owner.lookup('service:version').type = 'enterprise';
    await this.renderComponent();

    await click(DATE_RANGE.edit);
    await click(DATE_RANGE.editDate('reset'));
    await fillIn(DATE_RANGE.editDate('end'), '2017-05');
    assert.dom(DATE_RANGE.validation).hasText('You must supply both start and end dates.');
    await click(GENERAL.saveButton);
    assert.false(this.onChange.called);

    await fillIn(DATE_RANGE.editDate('start'), '2018-01');
    assert.dom(DATE_RANGE.validation).hasText('Start date must be before end date.');
    await click(GENERAL.saveButton);
    assert.false(this.onChange.called);

    await click(GENERAL.cancelButton);
    assert.false(this.onChange.called);
    assert.dom(DATE_RANGE.editModal).doesNotExist();
  });

  test('it resets the tracked values on close', async function (assert) {
    await this.renderComponent();

    await click(DATE_RANGE.edit);
    await fillIn(DATE_RANGE.editDate('start'), '2017-04');
    await fillIn(DATE_RANGE.editDate('end'), '2018-05');
    await click(GENERAL.cancelButton);

    await click(DATE_RANGE.edit);
    assert.dom(DATE_RANGE.editDate('start')).hasValue('2018-01');
    assert.dom(DATE_RANGE.editDate('end')).hasValue('2019-01');
  });
});
