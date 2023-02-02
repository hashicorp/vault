import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { lastDayOfMonth, subMonths, format, fromUnixTime, addMonths } from 'date-fns';
import { parseAPITimestamp, ARRAY_OF_MONTHS } from 'core/utils/date-formatters';

module('Unit | Adapter | clients activity', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'clients/activity';
    this.startDate = subMonths(new Date(), 6);
    this.endDate = new Date();
    this.readableUnix = (unix) => parseAPITimestamp(fromUnixTime(unix).toISOString(), 'MMMM dd yyyy');
  });

  test('it does not format if both params are timestamp strings', async function (assert) {
    assert.expect(1);
    const queryParams = {
      start_time: { timestamp: this.startDate.toISOString() },
      end_time: { timestamp: this.endDate.toISOString() },
    };
    this.server.get('sys/internal/counters/activity', (schema, req) => {
      assert.propEqual(req.queryParams, {
        start_time: this.startDate.toISOString(),
        end_time: this.endDate.toISOString(),
      });
    });

    this.store.queryRecord(this.modelName, queryParams);
  });

  test('it formats start_time if only end_time is a timestamp string', async function (assert) {
    assert.expect(2);
    const twoMonthsAhead = addMonths(this.startDate, 2);
    const month = twoMonthsAhead.getMonth();
    const year = twoMonthsAhead.getFullYear();
    const queryParams = {
      start_time: {
        monthIdx: month,
        year,
      },
      end_time: {
        timestamp: this.endDate.toISOString(),
      },
    };

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      const { start_time, end_time } = req.queryParams;
      const readableStart = this.readableUnix(start_time);
      assert.strictEqual(
        readableStart,
        `${ARRAY_OF_MONTHS[month]} 01 ${year}`,
        `formatted unix start time is the first of the month: ${readableStart}`
      );
      assert.strictEqual(end_time, this.endDate.toISOString(), 'end time is a timestamp string');
    });

    this.store.queryRecord(this.modelName, queryParams);
  });

  test('it formats end_time only if only start_time is a timestamp string', async function (assert) {
    assert.expect(2);
    const twoMothsAgo = subMonths(this.endDate, 2);
    const month = twoMothsAgo.getMonth() - 2;
    const year = twoMothsAgo.getFullYear();
    const dayOfMonth = format(lastDayOfMonth(new Date(year, month, 10)), 'dd');
    const queryParams = {
      start_time: {
        timestamp: this.startDate.toISOString(),
      },
      end_time: {
        monthIdx: month,
        year,
      },
    };

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      const { start_time, end_time } = req.queryParams;
      const readableEnd = this.readableUnix(end_time);
      assert.strictEqual(start_time, this.startDate.toISOString(), 'start time is a timestamp string');
      assert.strictEqual(
        readableEnd,
        `${ARRAY_OF_MONTHS[month]} ${dayOfMonth} ${year}`,
        `formatted unix end time is the last day of the month: ${readableEnd}`
      );
    });

    this.store.queryRecord(this.modelName, queryParams);
  });

  test('it formats both params if neither are a timestamp', async function (assert) {
    assert.expect(2);
    const startDate = subMonths(this.startDate, 2);
    const endDate = addMonths(this.endDate, 2);
    const startMonth = startDate.getMonth() + 2;
    const startYear = startDate.getFullYear();
    const endMonth = endDate.getMonth() - 2;
    const endYear = endDate.getFullYear();
    const endDay = format(lastDayOfMonth(new Date(endYear, endMonth, 10)), 'dd');
    const queryParams = {
      start_time: {
        monthIdx: startMonth,
        year: startYear,
      },
      end_time: {
        monthIdx: endMonth,
        year: endYear,
      },
    };

    this.server.get('sys/internal/counters/activity', (schema, req) => {
      const { start_time, end_time } = req.queryParams;
      const readableEnd = this.readableUnix(end_time);
      const readableStart = this.readableUnix(start_time);
      assert.strictEqual(
        readableStart,
        `${ARRAY_OF_MONTHS[startMonth]} 01 ${startYear}`,
        `formatted unix start time is the first of the month: ${readableStart}`
      );
      assert.strictEqual(
        readableEnd,
        `${ARRAY_OF_MONTHS[endMonth]} ${endDay} ${endYear}`,
        `formatted unix end time is the last day of the month: ${readableEnd}`
      );
    });

    this.store.queryRecord(this.modelName, queryParams);
  });
});
