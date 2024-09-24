/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import validators from 'vault/utils/model-helpers/validators';

module('Unit | Util | validators', function (hooks) {
  setupTest(hooks);

  // * MODEL VALIDATORS

  test('it should validate presence', function (assert) {
    let isValid;
    const check = (value) => (isValid = validators.presence(value));
    check(null);
    assert.false(isValid, 'Invalid when value is null');
    check('');
    assert.false(isValid, 'Invalid when value is empty string');
    check(true);
    assert.true(isValid, 'Valid when value is true');
    check(0);
    assert.true(isValid, 'Valid when value is 0 as integer');
    check('0');
    assert.true(isValid, 'Valid when value is 0 as string');
  });

  test('it should validate length', function (assert) {
    let isValid;
    const options = { nullable: true, min: 3, max: 5 };
    const check = (prop) => (isValid = validators.length(prop, options));
    check(null);
    assert.true(isValid, 'Valid when nullable is true');
    options.nullable = false;
    check(null);
    assert.false(isValid, 'Invalid when nullable is false');
    check('12');
    assert.false(isValid, 'Invalid when string not min length');
    check('123456');
    assert.false(isValid, 'Invalid when string over max length');
    check('1234');
    assert.true(isValid, 'Valid when string between min and max length');
    check(12);
    assert.false(isValid, 'Invalid when integer not min length');
    check(123456);
    assert.false(isValid, 'Invalid when integer over max length');
    check(1234);
    assert.true(isValid, 'Valid when integer between min and max length');
    options.min = 1;
    check(0);
    assert.true(isValid, 'Valid when integer is 0 and min is 1');
    check('0');
    assert.true(isValid, 'Valid when string is 0 and min is 1');
  });

  test('it should validate number', function (assert) {
    let isValid;
    const options = { nullable: true };
    const check = (prop) => (isValid = validators.number(prop, options));
    check(null);
    assert.true(isValid, 'Valid when nullable is true');
    options.nullable = false;
    check(null);
    assert.false(isValid, 'Invalid when nullable is false');
    check(9);
    assert.true(isValid, 'Valid for number');
    check('9');
    assert.true(isValid, 'Valid for number as string');
    check('foo');
    assert.false(isValid, 'Invalid for string that is not a number');
    check('12foo');
    assert.false(isValid, 'Invalid for string that contains a number');
    check(0);
    assert.true(isValid, 'Valid for 0 as an integer');
    check('0');
    assert.true(isValid, 'Valid for 0 as a string');
  });

  test('it should validate whitespace', function (assert) {
    let isValid;
    const check = (prop) => (isValid = validators.containsWhiteSpace(prop));
    check('validText');
    assert.true(isValid, 'Valid when text contains no spaces');
    check('valid-text');
    assert.true(isValid, 'Valid when text contains no spaces and hyphen');
    check('some space');
    assert.false(isValid, 'Invalid when text contains single space');
    check('text with spaces');
    assert.false(isValid, 'Invalid when text contains multiple spaces');
    check(' leadingSpace');
    assert.false(isValid, 'Invalid when text has leading whitespace');
    check('trailingSpace ');
    assert.false(isValid, 'Invalid when text has trailing whitespace');
  });

  test('it should validate value ends in a slash', function (assert) {
    let isValid;
    const check = (prop) => (isValid = validators.endsInSlash(prop));
    check('validText');
    assert.true(isValid, 'Valid when text does not end in slash');
    check('valid/Text');
    assert.true(isValid, 'Valid when text only contains slash');
    check('invalid/');
    assert.false(isValid, 'Invalid when text ends in slash');
    check('also/invalid/');
    assert.false(isValid, 'Invalid when text contains and ends in slash');
  });

  // * GENERAL VALIDATORS
  test('it returns whether a value has whitespace or not', function (assert) {
    let hasWhitespace;
    const check = (value) => (hasWhitespace = validators.hasWhitespace(value));

    check('someText');
    assert.false(hasWhitespace, 'False when text contains no spaces');

    check('some-text');
    assert.false(hasWhitespace, 'False when text contains no spaces and hyphen');

    check('some space');
    assert.true(hasWhitespace, 'True when text contains single space');

    check('text with spaces');
    assert.true(hasWhitespace, 'True when text contains multiple spaces');

    check(' leadingSpace');
    assert.true(hasWhitespace, 'True when text has leading whitespace');

    check('trailingSpace ');
    assert.true(hasWhitespace, 'True when text has trailing whitespace');
  });

  test('it returns whether a string input values evaluated as non-strings', function (assert) {
    let isNonString;
    const check = (value) => (isNonString = validators.isNonString(value));
    check(' {"foo": "bar"} ');
    assert.true(isNonString, 'returns true when value contains an object');

    check(' ["a", "b", "c"] ');
    assert.true(isNonString, 'returns true when value contains an array');

    check('123');
    assert.true(isNonString, 'returns true when value is numbers');

    check('123e6');
    assert.true(isNonString, 'returns true when value is numbers with exponents');

    check('true');
    assert.true(isNonString, 'returns true when value is true');

    // falsy values that return true because JSON.parse() is successful
    check('null');
    assert.true(isNonString, 'returns true when value is null');

    check('false');
    assert.true(isNonString, 'returns true when value is false');

    check('0');
    assert.true(isNonString, 'returns true when value is "0"');

    // falsy
    check('undefined');
    assert.false(isNonString, 'returns false when value is undefined');

    check('my string');
    assert.false(isNonString, 'returns false when value is letters');
  });
});
