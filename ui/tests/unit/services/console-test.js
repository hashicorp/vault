import { moduleFor, test } from 'ember-qunit';
import { sanitizePath, ensureTrailingSlash } from 'vault/services/console';
import Pretender from 'pretender';

moduleFor('service:console', 'Unit | Service | console', {
  needs: ['service:auth'],
  beforeEach() {
    this.server = new Pretender();
  },
  afterEach() {
    this.server.shutdown();
  },
});

test('#sanitizePath', function(assert) {
  assert.equal(sanitizePath(' /foo/bar/baz/ '), 'foo/bar/baz', 'removes spaces and slashs on either side');
  assert.equal(sanitizePath('//foo/bar/baz/'), 'foo/bar/baz', 'removes more than one slash');
});

test('#ensureTrailingSlash', function(assert) {
  assert.equal(ensureTrailingSlash('foo/bar'), 'foo/bar/', 'adds trailing slash');
  assert.equal(ensureTrailingSlash('baz/'), 'baz/', 'keeps trailing slash if there is one');
});
