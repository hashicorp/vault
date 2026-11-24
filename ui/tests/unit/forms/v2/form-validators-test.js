/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { validators } from 'vault/forms/v2/form-validators';

module('Unit | forms/v2 | form-validators', function () {
  module('required validator', function () {
    test('it rejects null', function (assert) {
      assert.notOk(validators.required(null), 'null is invalid');
    });

    test('it rejects undefined', function (assert) {
      assert.notOk(validators.required(undefined), 'undefined is invalid');
    });

    test('it rejects empty string', function (assert) {
      assert.notOk(validators.required(''), 'empty string is invalid');
    });

    test('it rejects whitespace-only string', function (assert) {
      assert.notOk(validators.required('   '), 'whitespace-only string is invalid');
    });

    test('it rejects empty array', function (assert) {
      assert.notOk(validators.required([]), 'empty array is invalid');
    });

    test('it rejects empty object', function (assert) {
      assert.notOk(validators.required({}), 'empty object is invalid');
    });

    test('it accepts non-empty string', function (assert) {
      assert.ok(validators.required('value'), 'non-empty string is valid');
    });

    test('it accepts non-empty array', function (assert) {
      assert.ok(validators.required(['item']), 'non-empty array is valid');
    });

    test('it accepts non-empty object', function (assert) {
      assert.ok(validators.required({ key: 'value' }), 'non-empty object is valid');
    });

    test('it accepts number zero', function (assert) {
      assert.ok(validators.required(0), 'zero is valid');
    });

    test('it accepts boolean false', function (assert) {
      assert.ok(validators.required(false), 'false is valid');
    });
  });

  module('email validator', function () {
    test('it accepts valid email addresses', function (assert) {
      assert.ok(validators.email('user@example.com'), 'simple email is valid');
      assert.ok(validators.email('user.name@example.com'), 'email with dot is valid');
      assert.ok(validators.email('user+tag@example.co.uk'), 'email with plus and subdomain is valid');
      assert.ok(
        validators.email('user_name@example-domain.com'),
        'email with underscore and hyphen is valid'
      );
    });

    test('it rejects invalid email addresses', function (assert) {
      assert.notOk(validators.email('invalid'), 'missing @ is invalid');
      assert.notOk(validators.email('invalid@'), 'missing domain is invalid');
      assert.notOk(validators.email('@example.com'), 'missing local part is invalid');
      assert.notOk(validators.email('invalid@domain'), 'missing TLD is invalid');
      assert.notOk(validators.email('invalid @example.com'), 'space in email is invalid');
    });

    test('it accepts empty value (use with required for mandatory)', function (assert) {
      assert.ok(validators.email(''), 'empty string is valid');
      assert.ok(validators.email(null), 'null is valid');
      assert.ok(validators.email(undefined), 'undefined is valid');
    });
  });

  module('url validator', function () {
    test('it accepts valid URLs', function (assert) {
      assert.ok(validators.url('https://example.com'), 'https URL is valid');
      assert.ok(validators.url('http://example.com'), 'http URL is valid');
      assert.ok(validators.url('https://example.com/path'), 'URL with path is valid');
      assert.ok(validators.url('https://example.com:8080'), 'URL with port is valid');
      assert.ok(validators.url('https://example.com?query=value'), 'URL with query is valid');
      assert.ok(validators.url('https://sub.example.com'), 'URL with subdomain is valid');
    });

    test('it rejects invalid URLs', function (assert) {
      assert.notOk(validators.url('not-a-url'), 'plain text is invalid');
      assert.notOk(validators.url('example.com'), 'missing protocol is invalid');
      assert.notOk(validators.url('//example.com'), 'protocol-relative URL is invalid');
    });

    test('it accepts empty value', function (assert) {
      assert.ok(validators.url(''), 'empty string is valid');
      assert.ok(validators.url(null), 'null is valid');
      assert.ok(validators.url(undefined), 'undefined is valid');
    });
  });

  module('pattern validator', function () {
    test('it validates against string pattern', function (assert) {
      const options = { pattern: '^[A-Z]{3}$' };
      assert.ok(validators.pattern('ABC', options), 'matching pattern is valid');
      assert.notOk(validators.pattern('abc', options), 'non-matching pattern is invalid');
      assert.notOk(validators.pattern('ABCD', options), 'too long is invalid');
    });

    test('it validates against RegExp pattern', function (assert) {
      const options = { pattern: /^[A-Z]{3}$/ };
      assert.ok(validators.pattern('ABC', options), 'matching RegExp is valid');
      assert.notOk(validators.pattern('abc', options), 'non-matching RegExp is invalid');
    });

    test('it supports regex flags', function (assert) {
      const options = { pattern: '^[A-Z]{3}$', flags: 'i' };
      assert.ok(validators.pattern('abc', options), 'case-insensitive match is valid');
      assert.ok(validators.pattern('ABC', options), 'uppercase match is valid');
    });

    test('it accepts empty value', function (assert) {
      const options = { pattern: '^[A-Z]{3}$' };
      assert.ok(validators.pattern('', options), 'empty string is valid');
      assert.ok(validators.pattern(null, options), 'null is valid');
    });

    test('it returns true when no pattern provided', function (assert) {
      assert.ok(validators.pattern('anything', {}), 'no pattern always valid');
    });
  });

  module('minLength validator', function () {
    test('it validates minimum string length', function (assert) {
      const options = { minLength: 3 };
      assert.ok(validators.minLength('abc', options), 'exact length is valid');
      assert.ok(validators.minLength('abcd', options), 'longer is valid');
      assert.notOk(validators.minLength('ab', options), 'shorter is invalid');
    });

    test('it rejects empty value when minLength is set', function (assert) {
      const options = { minLength: 3 };
      assert.notOk(validators.minLength('', options), 'empty string is invalid');
      assert.notOk(validators.minLength(null, options), 'null is invalid');
    });

    test('it returns true when no minLength provided', function (assert) {
      assert.ok(validators.minLength('ab', {}), 'no minLength always valid');
    });
  });

  module('maxLength validator', function () {
    test('it validates maximum string length', function (assert) {
      const options = { maxLength: 5 };
      assert.ok(validators.maxLength('abc', options), 'shorter is valid');
      assert.ok(validators.maxLength('abcde', options), 'exact length is valid');
      assert.notOk(validators.maxLength('abcdef', options), 'longer is invalid');
    });

    test('it accepts empty value', function (assert) {
      const options = { maxLength: 5 };
      assert.ok(validators.maxLength('', options), 'empty string is valid');
      assert.ok(validators.maxLength(null, options), 'null is valid');
    });

    test('it returns true when no maxLength provided', function (assert) {
      assert.ok(validators.maxLength('very long string', {}), 'no maxLength always valid');
    });
  });

  module('min validator', function () {
    test('it validates minimum numeric value', function (assert) {
      const options = { min: 10 };
      assert.ok(validators.min(10, options), 'exact value is valid');
      assert.ok(validators.min(15, options), 'greater value is valid');
      assert.ok(validators.min('15', options), 'string number is valid');
      assert.notOk(validators.min(5, options), 'lesser value is invalid');
      assert.notOk(validators.min('5', options), 'string lesser value is invalid');
    });

    test('it rejects non-numeric values', function (assert) {
      const options = { min: 10 };
      assert.notOk(validators.min('abc', options), 'non-numeric string is invalid');
      assert.notOk(validators.min(NaN, options), 'NaN is invalid');
    });

    test('it accepts empty value', function (assert) {
      const options = { min: 10 };
      assert.ok(validators.min('', options), 'empty string is valid');
      assert.ok(validators.min(null, options), 'null is valid');
      assert.ok(validators.min(undefined, options), 'undefined is valid');
    });

    test('it returns true when no min provided', function (assert) {
      assert.ok(validators.min(5, {}), 'no min always valid');
    });

    test('it handles zero correctly', function (assert) {
      const options = { min: 0 };
      assert.ok(validators.min(0, options), 'zero is valid when min is zero');
      assert.ok(validators.min(5, options), 'positive is valid when min is zero');
      assert.notOk(validators.min(-5, options), 'negative is invalid when min is zero');
    });

    test('it handles negative numbers', function (assert) {
      const options = { min: -10 };
      assert.ok(validators.min(-5, options), 'greater negative is valid');
      assert.ok(validators.min(0, options), 'zero is valid');
      assert.notOk(validators.min(-15, options), 'lesser negative is invalid');
    });
  });

  module('max validator', function () {
    test('it validates maximum numeric value', function (assert) {
      const options = { max: 10 };
      assert.ok(validators.max(10, options), 'exact value is valid');
      assert.ok(validators.max(5, options), 'lesser value is valid');
      assert.ok(validators.max('5', options), 'string number is valid');
      assert.notOk(validators.max(15, options), 'greater value is invalid');
      assert.notOk(validators.max('15', options), 'string greater value is invalid');
    });

    test('it rejects non-numeric values', function (assert) {
      const options = { max: 10 };
      assert.notOk(validators.max('abc', options), 'non-numeric string is invalid');
      assert.notOk(validators.max(NaN, options), 'NaN is invalid');
    });

    test('it accepts empty value', function (assert) {
      const options = { max: 10 };
      assert.ok(validators.max('', options), 'empty string is valid');
      assert.ok(validators.max(null, options), 'null is valid');
      assert.ok(validators.max(undefined, options), 'undefined is valid');
    });

    test('it returns true when no max provided', function (assert) {
      assert.ok(validators.max(100, {}), 'no max always valid');
    });

    test('it handles zero correctly', function (assert) {
      const options = { max: 0 };
      assert.ok(validators.max(0, options), 'zero is valid when max is zero');
      assert.ok(validators.max(-5, options), 'negative is valid when max is zero');
      assert.notOk(validators.max(5, options), 'positive is invalid when max is zero');
    });

    test('it handles negative numbers', function (assert) {
      const options = { max: -10 };
      assert.ok(validators.max(-15, options), 'lesser negative is valid');
      assert.ok(validators.max(-10, options), 'exact negative is valid');
      assert.notOk(validators.max(-5, options), 'greater negative is invalid');
      assert.notOk(validators.max(0, options), 'zero is invalid');
    });
  });
});
