import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, findAll, find, triggerEvent, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

let file;
const fileEvent = () => {
  const data = { some: 'content' };
  file = new File([JSON.stringify(data, null, 2)], 'file.json', { type: 'application/json' });
  return ['change', [file]];
};

module('Integration | Component | pgp file', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    file = null;
    this.lastOnChangeCall = null;
    this.set('change', (index, key) => {
      this.lastOnChangeCall = [index, key];
      this.set('key', key);
    });
  });

  test('it renders', async function(assert) {
    this.set('key', { value: '' });
    this.set('index', 0);

    await render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);

    assert.equal(find('[data-test-pgp-label]').textContent.trim(), 'PGP KEY 1');
    assert.equal(find('[data-test-pgp-file-input-label]').textContent.trim(), 'Choose a file…');
  });

  test('it accepts files', async function(assert) {
    const key = { value: '' };
    const event = fileEvent();
    this.set('key', key);
    this.set('index', 0);

    await render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);
    triggerEvent('[data-test-pgp-file-input]', ...event);

    // FileReader is async, but then we need extra run loop wait to re-render
    await waitUntil(() => {
      return !!this.lastOnChangeCall;
    });
    assert.equal(
      find('[data-test-pgp-file-input-label]').textContent.trim(),
      file.name,
      'the file input shows the file name'
    );
    assert.notDeepEqual(this.lastOnChangeCall[1].value, key.value, 'onChange was called with the new key');
    assert.equal(this.lastOnChangeCall[0], 0, 'onChange is called with the index value');
    await click('[data-test-pgp-clear]');
    assert.equal(this.lastOnChangeCall[1].value, key.value, 'the key gets reset when the input is cleared');
  });

  test('it allows for text entry', async function(assert) {
    const key = { value: '' };
    const text = 'a really long pgp key';
    this.set('key', key);
    this.set('index', 0);

    await render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);
    await click('[data-test-text-toggle]');
    assert.equal(findAll('[data-test-pgp-file-textarea]').length, 1, 'renders the textarea on toggle');

    fillIn('[data-test-pgp-file-textarea]', text);
    await waitUntil(() => {
      return !!this.lastOnChangeCall;
    });
    assert.equal(this.lastOnChangeCall[1].value, text, 'the key value is passed to onChange');
  });

  test('toggling back and forth', async function(assert) {
    const key = { value: '' };
    const event = fileEvent();
    this.set('key', key);
    this.set('index', 0);

    await render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);
    await triggerEvent('[data-test-pgp-file-input]', ...event);
    await click('[data-test-text-toggle]');
    assert.equal(findAll('[data-test-pgp-file-textarea]').length, 1, 'renders the textarea on toggle');
    assert.equal(
      find('[data-test-pgp-file-textarea]').textContent.trim(),
      this.lastOnChangeCall[1].value,
      'textarea shows the value of the base64d key'
    );
  });
});
