import trimRight from 'vault/utils/trim-right';
import { module, test } from 'qunit';

module('Unit | Util | trim right');

test('it trims extension array from end of string', function(assert) {
  const trimmedName = trimRight('my-file-name-is-cool.json', ['.json', '.txt', '.hcl', '.policy']);

  assert.equal(trimmedName, 'my-file-name-is-cool');
});

test('it only trims extension array from the very end of string', function(assert) {
  const trimmedName = trimRight('my-file-name.json-is-cool.json', ['.json', '.txt', '.hcl', '.policy']);

  assert.equal(trimmedName, 'my-file-name.json-is-cool');
});

test('it returns string as is if trim array is empty', function(assert) {
  const trimmedName = trimRight('my-file-name-is-cool.json', []);

  assert.equal(trimmedName, 'my-file-name-is-cool.json');
});

test('it returns string as is if trim array is not passed in', function(assert) {
  const trimmedName = trimRight('my-file-name-is-cool.json');

  assert.equal(trimmedName, 'my-file-name-is-cool.json');
});
