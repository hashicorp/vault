import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import validators from 'vault/utils/validators';

module('Unit | Util | validators', function (hooks) {
  setupTest(hooks);

  test('it should validate presence', function (assert) {
    let isValid = validators.presence(null);
    assert.false(isValid);
    isValid = validators.presence(true);
    assert.true(isValid);
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
  });

  test('it should validate number', function (assert) {
    let isValid;
    const options = { nullable: true, asString: false };
    const check = (prop) => (isValid = validators.number(prop, options));
    check(null);
    assert.true(isValid, 'Valid when nullable is true');
    options.nullable = false;
    check(null);
    assert.false(isValid, 'Invalid when nullable is false');
    check('9');
    assert.false(isValid, 'Invalid for string when asString is false');
    check(9);
    assert.true(isValid, 'Valid for number');
    options.asString = true;
    check('9');
    assert.true(isValid, 'Valid for number as string');
    check('foo');
    assert.false(isValid, 'Invalid for string that is not a number');
    check(0);
    assert.true(isValid, 'Valid for 0 as an integer');
    check('0');
    assert.true(isValid, 'Valid for 0 as a string');
  });
});
