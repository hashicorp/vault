import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

const TOTAL = 15;
const CARD_TITLE = 'Tokens';

module('Integration | Component meep', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('total', TOTAL);
    this.set('cardTitle', CARD_TITLE);
  });

  test('the total number renders', async function(assert) {
    await render(hbs`<SelectableCard @total={{total}} @cardTitle={{cardTitle}}/>`);
    let titleNumber = this.element.querySelector('.title-number').innerText;

    assert.equal(titleNumber, 15);
  });

  test('if total is 1, return non-plural version of card title', async function(assert) {
    await render(hbs`<SelectableCard @total={{1}} @cardTitle={{cardTitle}}/>`);
    let titleText = this.element.querySelector('.title').innerText;

    assert.equal(titleText, 'Token');
  });
});
