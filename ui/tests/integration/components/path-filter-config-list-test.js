import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll, fillIn, blur } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
const resolver = engineResolverFor('replication');

module('Integration | Component | path filter config list', function(hooks) {
  setupRenderingTest(hooks, { resolver });

  test('it renders', async function(assert) {
    this.set('config', { mode: 'whitelist', paths: [] });
    this.set('paths', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
    await render(hbs`{{path-filter-config-list config=config paths=paths}}`);

    assert.equal(findAll('#filter-userpass').length, 1);
  });

  test('it sets config.paths', async function(assert) {
    this.set('config', { mode: 'whitelist', paths: [] });
    this.set('paths', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
    await render(hbs`{{path-filter-config-list config=config paths=paths}}`);

    await click('#filter-userpass');
    assert.ok(this.get('config.paths').includes('userpass/'), 'adds to paths');

    await click('#filter-userpass');
    assert.equal(this.get('config.paths').length, 0, 'removes from paths');
  });

  test('it sets config.mode', async function(assert) {
    this.set('config', { mode: 'whitelist', paths: [] });
    await render(hbs`{{path-filter-config-list config=config}}`);
    await fillIn('#filter-mode', 'blacklist');
    await blur('#filter-mode');
    assert.equal(this.get('config.mode'), 'blacklist');
  });
});
