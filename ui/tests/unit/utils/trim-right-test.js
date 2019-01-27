import trimRight from 'vault/utils/trim-right';
import { module, test } from 'qunit';

module('Unit | Util | trim right', function() {
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

  test('it allows the last extension to also be part of the file name', function(assert) {
    const trimmedName = trimRight('my-policy.hcl', ['.json', '.txt', '.hcl', '.policy']);

    assert.equal(trimmedName, 'my-policy');
  });

  test('it allows the last extension to also be part of the file name and the extenstion', function(assert) {
    const trimmedName = trimRight('my-policy.policy', ['.json', '.txt', '.hcl', '.policy']);

    assert.equal(trimmedName, 'my-policy');
  });

  test('it passes endings into the regex unescaped when passing false as the third arg', function(assert) {
    const trimmedName = trimRight('my-policypolicy', ['.json', '.txt', '.hcl', '.policy'], false);

    // the . gets interpreted as regex wildcard so it also trims the y character
    assert.equal(trimmedName, 'my-polic');
  });
});
