import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll, fillIn, blur } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | mount filter config list', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    this.set('config', { mode: 'whitelist', paths: [] });
    this.set('mounts', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
    await render(hbs`{{mount-filter-config-list config=config mounts=mounts}}`);

    assert.equal(findAll('#filter-userpass').length, 1);
  });

  test('it sets config.paths', async function(assert) {
    this.set('config', { mode: 'whitelist', paths: [] });
    this.set('mounts', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
    await render(hbs`{{mount-filter-config-list config=config mounts=mounts}}`);

    await click('#filter-userpass');
    assert.ok(this.get('config.paths').includes('userpass/'), 'adds to paths');

    await click('#filter-userpass');
    assert.equal(this.get('config.paths').length, 0, 'removes from paths');
  });

  test('it sets config.mode', async function(assert) {
    this.set('config', { mode: 'whitelist', paths: [] });
    await render(hbs`{{mount-filter-config-list config=config}}`);
    await fillIn('#filter-mode', 'blacklist');
    await blur('#filter-mode');
    assert.equal(this.get('config.mode'), 'blacklist');
  });
});
