import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { changelogUrlFor } from '../../../helpers/changelog-url-for';

const CHANGELOG_URL = 'https://www.github.com/hashicorp/vault/blob/master/CHANGELOG.md#';

module('Integration | Helper | changelog-url-for', function(hooks) {
  setupRenderingTest(hooks);

  test('it builds an enterprise URL', function(assert) {
    const result = changelogUrlFor(['1.5.0+prem']);
    assert.equal(result, CHANGELOG_URL.concat('150'));
  });

  test('it builds an OSS URL', function(assert) {
    const result = changelogUrlFor(['1.4.3']);
    assert.equal(result, CHANGELOG_URL.concat('143'));
  });

  test('it returns the base changelog URL if the version is less than 1.4.3', function(assert) {
    const result = changelogUrlFor(['1.4.0']);
    assert.equal(result, CHANGELOG_URL);
  });

  test('it returns the base changelog URL if version cannot be found', function(assert) {
    const result = changelogUrlFor(['']);
    assert.equal(result, CHANGELOG_URL);
  });
});
