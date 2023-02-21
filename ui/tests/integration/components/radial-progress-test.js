import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import radialProgress from 'vault/tests/pages/components/radial-progress';

const component = create(radialProgress);

module('Integration | Component | radial progress', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    // We have to manually round the circumference, strokeDash, and strokeDashOffset because
    // ie11 truncates decimals differently than other browsers.
    const circumference = ((19 / 2) * Math.PI * 2).toFixed(2);
    await render(hbs`{{radial-progress progressDecimal=0.5}}`);

    assert.strictEqual(component.viewBox, '0 0 20 20');
    assert.strictEqual(component.height, '20');
    assert.strictEqual(component.width, '20');
    assert.strictEqual(component.strokeWidth, '1');
    assert.strictEqual(component.r, (19 / 2).toString());
    assert.strictEqual(component.cx, '10');
    assert.strictEqual(component.cy, '10');
    assert.strictEqual(Number(component.strokeDash).toFixed(2), circumference);
    assert.strictEqual(Number(component.strokeDashOffset).toFixed(3), (circumference * 0.5).toFixed(3));
  });
});
