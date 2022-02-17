import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Helper | format-ttl', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the input if no match found', async function (assert) {
    this.set('inputValue', '1234');

    await render(hbs`{{format-ttl inputValue}}`);

    assert.equal(this.element.textContent.trim(), '1234');
  });

  test('it parses hours correctly', async function (assert) {
    this.set('inputValue', '12h');

    await render(hbs`{{format-ttl inputValue}}`);
    assert.equal(this.element.textContent.trim(), '12 hours');

    this.set('inputValue', '1h');
    assert.equal(this.element.textContent.trim(), '1 hour');
  });

  test('it parses minutes correctly', async function (assert) {
    this.set('inputValue', '30m');

    await render(hbs`{{format-ttl inputValue}}`);
    assert.equal(this.element.textContent.trim(), '30 minutes');

    this.set('inputValue', '1m');
    assert.equal(this.element.textContent.trim(), '1 minute');
  });

  test('it parses seconds correctly', async function (assert) {
    this.set('inputValue', '45s');

    await render(hbs`{{format-ttl inputValue}}`);
    assert.equal(this.element.textContent.trim(), '45 seconds');

    this.set('inputValue', '1s');
    assert.equal(this.element.textContent.trim(), '1 second');
  });

  test('it parses multiple matches correctly', async function (assert) {
    this.set('inputValue', '1h30m0s');

    await render(hbs`{{format-ttl inputValue}}`);
    assert.equal(this.element.textContent.trim(), '1 hour 30 minutes 0 seconds');
  });

  test('it removes 0 values if removeZero true', async function (assert) {
    this.set('inputValue', '1h30m0s');

    await render(hbs`{{format-ttl inputValue removeZero=true}}`);
    assert.equal(this.element.textContent.trim(), '1 hour 30 minutes');
  });

  test('returns empty string if all values 0 and removeZero true', async function (assert) {
    this.set('inputValue', '0h0m0s');

    await render(hbs`{{format-ttl inputValue removeZero=true}}`);
    assert.equal(this.element.textContent.trim(), '');
  });
});
