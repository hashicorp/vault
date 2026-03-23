/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import MatchesCurrentUrl, { matchesCurrentUrl } from 'core/helpers/matches-current-url';
import sinon from 'sinon';

module('Unit | Helper | matches-current-url', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.helper = new MatchesCurrentUrl(this.owner);
    this.recognizeStub = sinon.stub(this.router, 'recognize').returns({ name: 'vault.secrets.secret.show' });
    this.currentURL = sinon.stub(this.router, 'currentURL').value('/vault/secrets/secret/show');
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('it looks up router based on context', function (assert) {
    assert.strictEqual(this.helper.router, this.router, 'helper looks up correct router');
  });

  test('it returns false when the route info cannot be recognized', function (assert) {
    this.router.recognize = () => null;
    const result = this.helper.compute(['vault.secrets']);
    assert.false(result, 'returns false when route info is null');
  });

  test('it returns an error if routeName is an empty string', function (assert) {
    assert.throws(
      () => {
        this.helper.compute(['']);
      },
      /Error: Assertion Failed: routeName is required, you passed an empty string/,
      'it throws when routeName is an empty string'
    );
  });

  test('it returns an error for undefined args', function (assert) {
    assert.throws(
      () => {
        this.helper.compute([]);
      },
      /Error: Assertion Failed: routeName is required, you passed undefined/,
      'it throws when args are undefined'
    );
  });

  test('it handles failed route lookup args', function (assert) {
    sinon.stub(this.helper, 'router').value(undefined);
    const result = this.helper.compute(['vault.secrets.secret.show']);
    assert.false(result, 'returns false when no route info is available');
  });

  test('it computes', function (assert) {
    const result = this.helper.compute(['vault.secrets']);
    assert.true(result, 'returns true when current route name includes the passed route name');
    const [url] = this.recognizeStub.lastCall.args;
    // rootURL is `/ui/`
    assert.strictEqual(url, '/ui/vault/secrets/secret/show', 'recognize is called with rootURL + currentURL');
  });

  test('it returns true when isExactMatch is true and names are equal', function (assert) {
    const result = this.helper.compute(['vault.secrets.secret.show'], { isExactMatch: true });
    assert.true(result, 'returns true when route name matches exactly');
  });

  test('matchesCurrentUrl: it calls recognize with rootURL + currentURL', function (assert) {
    matchesCurrentUrl(this.router, 'vault.secrets');
    const [url] = this.recognizeStub.lastCall.args;
    assert.strictEqual(url, '/ui/vault/secrets/secret/show', 'recognize is called with rootURL + currentURL');
  });

  test('matchesCurrentUrl: it returns true for substring match', function (assert) {
    const result = matchesCurrentUrl(this.router, 'vault.secrets');
    assert.true(result, 'returns true when current route name includes passed route name');
  });

  test('matchesCurrentUrl: it returns false when substring does not match', function (assert) {
    const result = matchesCurrentUrl(this.router, 'vault.access');
    assert.false(result, 'returns false when passed route name is not in current route name');
  });

  test('matchesCurrentUrl: it returns true when names are equal and isExactMatch is false', function (assert) {
    const result = matchesCurrentUrl(this.router, 'vault.secrets.secret.show');
    assert.true(result, 'returns true when route names match exactly');
  });

  test('matchesCurrentUrl: it returns true when names are equal and isExactMatch is true', function (assert) {
    const result = matchesCurrentUrl(this.router, 'vault.secrets.secret.show', { isExactMatch: true });
    assert.true(result, 'returns true when route names match exactly');
  });

  test('matchesCurrentUrl: it returns false when isExactMatch is true and routes are different', function (assert) {
    const result = matchesCurrentUrl(this.router, 'vault.secrets', { isExactMatch: true });
    assert.false(result);
  });

  test('matchesCurrentUrl: it returns false when router is undefined', function (assert) {
    const result = matchesCurrentUrl(undefined, 'vault.secrets');
    assert.false(result);
  });

  test('matchesCurrentUrl: it returns an error when route is an empty string', function (assert) {
    assert.throws(
      () => {
        matchesCurrentUrl(undefined, '');
      },
      /Error: Assertion Failed: routeName is required, you passed an empty string/,
      'it throws when routeName is an empty string'
    );
  });

  test('matchesCurrentUrl: it returns an error when args are undefined', function (assert) {
    assert.throws(
      () => {
        matchesCurrentUrl();
      },
      /Error: Assertion Failed: routeName is required, you passed undefined/,
      'it throws when args are undefined'
    );
  });
});
