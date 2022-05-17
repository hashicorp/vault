import { formatNumbers, formatTooltipNumber } from 'vault/utils/chart-helpers';
import { module, test } from 'qunit';

const SMALL_NUMBERS = [0, 7, 27, 103, 999];
const LARGE_NUMBERS = {
  1001: '1k',
  33777: '34k',
  532543: '530k',
  2100100: '2.1M',
  54500200100: '55B',
};

module('Unit | Utility | chart-helpers', function () {
  test('formatNumbers renders number correctly', function (assert) {
    assert.expect(11);
    const method = formatNumbers();
    assert.ok(method);
    SMALL_NUMBERS.forEach(function (num) {
      assert.equal(formatNumbers(num), num, `Does not format small number ${num}`);
    });
    Object.keys(LARGE_NUMBERS).forEach(function (num) {
      const expected = LARGE_NUMBERS[num];
      assert.equal(formatNumbers(num), expected, `Formats ${num} as ${expected}`);
    });
  });

  test('formatTooltipNumber renders number correctly', function (assert) {
    const formatted = formatTooltipNumber(120300200100);
    assert.equal(formatted.length, 15, 'adds punctuation at proper place for large numbers');
  });
});
