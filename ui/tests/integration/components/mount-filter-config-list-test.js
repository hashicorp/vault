import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('mount-filter-config-list', 'Integration | Component | mount filter config list', {
  integration: true,
});

test('it renders', function(assert) {
  this.set('config', { mode: 'whitelist', paths: [] });
  this.set('mounts', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
  this.render(hbs`{{mount-filter-config-list config=config mounts=mounts}}`);

  assert.equal(this.$('#filter-userpass').length, 1);
});

test('it sets config.paths', function(assert) {
  this.set('config', { mode: 'whitelist', paths: [] });
  this.set('mounts', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
  this.render(hbs`{{mount-filter-config-list config=config mounts=mounts}}`);

  this.$('#filter-userpass').click();
  assert.ok(this.get('config.paths').includes('userpass/'), 'adds to paths');

  this.$('#filter-userpass').click();
  assert.equal(this.get('config.paths').length, 0, 'removes from paths');
});

test('it sets config.mode', function(assert) {
  this.set('config', { mode: 'whitelist', paths: [] });
  this.render(hbs`{{mount-filter-config-list config=config}}`);
  this.$('#filter-mode').val('blacklist').change();
  assert.equal(this.get('config.mode'), 'blacklist');
});
