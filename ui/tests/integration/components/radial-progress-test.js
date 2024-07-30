/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { RADIAL_PROGRESS } from 'vault/tests/helpers/components/shamir-selectors';

module('Integration | Component | radial progress', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // We have to manually round the circumference, strokeDash, and strokeDashOffset because
    // ie11 truncates decimals differently than other browsers.
    const circumference = ((19 / 2) * Math.PI * 2).toFixed(2);
    await render(hbs`<RadialProgress @progressDecimal={{0.5}}/>`);
    const svg = document.querySelector(RADIAL_PROGRESS.svg);
    const path = document.querySelector(RADIAL_PROGRESS.path);
    const progress = document.querySelector(RADIAL_PROGRESS.progress);

    assert.strictEqual(svg.getAttribute('viewBox'), '0 0 20 20');
    assert.strictEqual(svg.getAttribute('height'), '20');
    assert.strictEqual(svg.getAttribute('width'), '20');
    // path
    assert.strictEqual(path.getAttribute('stroke-width'), '1');
    assert.strictEqual(path.getAttribute('r'), (19 / 2).toString());
    assert.strictEqual(path.getAttribute('cx'), '10');
    assert.strictEqual(path.getAttribute('cy'), '10');
    assert.strictEqual(Number(progress.getAttribute('strokeDash').toFixed(2)), circumference);
    assert.strictEqual(
      Number(progress.getAttribute('strokeDashOffset').toFixed(3)),
      (circumference * 0.5).toFixed(3)
    );
  });
});
