/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { click } from '@ember/test-helpers';

module('Integration | Component | pagination-controls', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders correct number of pages', async function (assert) {
    const totals = [
      [10, 1],
      [40, 3],
      [100, 5],
    ];
    for (const [total, count] of totals) {
      this.total = total;
      await render(hbs`<PaginationControls @total={{this.total}} />`);
      assert
        .dom('[data-test-page]')
        .exists({ count }, `Correct page count of ${count} renders for ${total} total items`);
      assert.dom('[data-test-more-pages]')[count === 5 ? 'exists' : 'doesNotExist']();
    }
  });

  test('it changes pages', async function (assert) {
    assert.expect(10);

    let expectedPage = 2;
    this.onChange = (page) => {
      assert.strictEqual(page, expectedPage, 'onChange callback is fired with correct page number');
    };

    await render(hbs`<PaginationControls @total={{75}} @onChange={{this.onChange}} />`);

    const isActive = (page) => {
      return this.element
        .querySelector(`[data-test-page="${page}"]`)
        .classList.value.includes('is-primary is-underlined is-active');
    };

    assert.ok(isActive(1), 'Page 1 is active by default');
    assert.dom('[data-test-previous-page]').isDisabled('Previous page button is disabled on page 1');

    await click('[data-test-next-page]');
    assert.ok(isActive(2), 'Page 2 is active');
    assert.dom('[data-test-previous-page]').isNotDisabled('Previous page button is disabled on page 1');

    expectedPage = 5;
    await click('[data-test-page="5"]');
    assert.ok(isActive(5), 'Page 5 is active');
    assert.dom('[data-test-next-page]').isDisabled('Next page button is disabled on last page');

    expectedPage = 4;
    await click('[data-test-previous-page]');
    assert.ok(isActive(4), 'Page 4 is active');
  });

  test('it renders correct display info', async function (assert) {
    this.onChange = () => {};
    await render(hbs`<PaginationControls @total={{68}} @onChange={{this.onChange}} />`);

    const ranges = ['1-15', '16-30', '31-45', '46-60', '61-68'];
    for (const [i, range] of ranges.entries()) {
      assert
        .dom('[data-test-page-display-info]')
        .hasText(`${range} of 68`, `Correct display info renders for page ${i + 1}`);

      if (i < 4) {
        await click(`[data-test-next-page]`);
      }
    }
  });
});
