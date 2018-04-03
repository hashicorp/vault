import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import Ember from 'ember';
import wait from 'ember-test-helpers/wait';

let file;
const fileEvent = () => {
  const data = { some: 'content' };
  file = new File([JSON.stringify(data, null, 2)], 'file.json', { type: 'application/json' });
  return Ember.$.Event('change', {
    target: {
      files: [file],
    },
  });
};

moduleForComponent('pgp-file', 'Integration | Component | pgp file', {
  integration: true,
  beforeEach: function() {
    file = null;
    this.lastOnChangeCall = null;
    this.set('change', (index, key) => {
      this.lastOnChangeCall = [index, key];
      this.set('key', key);
    });
  },
});

test('it renders', function(assert) {
  this.set('key', { value: '' });
  this.set('index', 0);

  this.render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);

  assert.equal(this.$('[data-test-pgp-label]').text().trim(), 'PGP KEY 1');
  assert.equal(this.$('[data-test-pgp-file-input-label]').text().trim(), 'Choose a file…');
});

test('it accepts files', function(assert) {
  const key = { value: '' };
  const event = fileEvent();
  this.set('key', key);
  this.set('index', 0);

  this.render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);
  this.$('[data-test-pgp-file-input]').trigger(event);

  return wait().then(() => {
    // FileReader is async, but then we need extra run loop wait to re-render
    Ember.run.next(() => {
      assert.equal(
        this.$('[data-test-pgp-file-input-label]').text().trim(),
        file.name,
        'the file input shows the file name'
      );
      assert.notDeepEqual(this.lastOnChangeCall[1].value, key.value, 'onChange was called with the new key');
      assert.equal(this.lastOnChangeCall[0], 0, 'onChange is called with the index value');
      this.$('[data-test-pgp-clear]').click();
    });
    return wait().then(() => {
      assert.equal(this.lastOnChangeCall[1].value, key.value, 'the key gets reset when the input is cleared');
    });
  });
});

test('it allows for text entry', function(assert) {
  const key = { value: '' };
  const text = 'a really long pgp key';
  this.set('key', key);
  this.set('index', 0);

  this.render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);
  this.$('[data-test-text-toggle]').click();
  assert.equal(this.$('[data-test-pgp-file-textarea]').length, 1, 'renders the textarea on toggle');

  this.$('[data-test-pgp-file-textarea]').text(text).trigger('input');
  assert.equal(this.lastOnChangeCall[1].value, text, 'the key value is passed to onChange');
});

test('toggling back and forth', function(assert) {
  const key = { value: '' };
  const event = fileEvent();
  this.set('key', key);
  this.set('index', 0);

  this.render(hbs`{{pgp-file index=index key=key onChange=(action change)}}`);
  this.$('[data-test-pgp-file-input]').trigger(event);
  return wait().then(() => {
    Ember.run.next(() => {
      this.$('[data-test-text-toggle]').click();
      wait().then(() => {
        assert.equal(this.$('[data-test-pgp-file-textarea]').length, 1, 'renders the textarea on toggle');
        assert.equal(
          this.$('[data-test-pgp-file-textarea]').text().trim(),
          this.lastOnChangeCall[1].value,
          'textarea shows the value of the base64d key'
        );
      });
    });
  });
});
