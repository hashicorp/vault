import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';
import sinon from 'sinon';

/*
EXAMPLE 4

Notice that Ember gives us the ability to create a stub service by
importing and extending Service directly. This is helpful because
the real permissions service makes API calls which are unnecessary
for these tests.
*/

const Permissions = Service.extend({
  globPaths: null,
  hasNavPermission() {
    return this.globPaths ? true : false;
  },
  // we can accomplish similar behavior above by using a sinon stub
  // hasNavPermission: sinon.stub().returns(console.log('hello')),
});

module('Integration | Helper | has-permission | ember learn', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    /*
    Thanks to setupRenderingTest above, which gives us access to
    Ember's dependency injection system, we can ensure that the
    helper here uses our stub permissions service instead of the
    real one. We do this by registering the service and pointing the
    test instance of the permissions service to the one we just
    registered.
    */
    this.owner.register('service:permissions', Permissions);
    this.permissions = this.owner.lookup('service:permissions');
  });

  test('it renders', async function(assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);

    assert.equal(this.element.textContent.trim(), 'No');
    await run(() => {
      this.permissions.set('globPaths', { 'test/': { capabilities: ['update'] } });
    });
    assert.equal(this.element.textContent.trim(), 'Yes', 'the helper re-computes when globPaths changes');
  });
});
