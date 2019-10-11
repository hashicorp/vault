import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll, fillIn, blur } from '@ember/test-helpers';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import Service from '@ember/service';
import sinon from 'sinon';
import { Promise } from 'rsvp';
import ss from 'vault/tests/pages/components/search-select';

const searchSelect = create(ss);

const resolver = engineResolverFor('replication');
const NAMESPACE_MOUNTS_RESPONSE = {
  data: {
    secret: {},
    auth: {},
  },
};

module('Integration | Component | path filter config list', function(hooks) {
  setupRenderingTest(hooks, { resolver });
  beforeEach(function() {
    let ajaxStub = sinon
      .stub()
      .usingPromise(Promise)
      .resolves(NAMESPACE_MOUNTS_RESPONSE);
    this.set('ajaxStub', ajaxStub);
    const namespaceServiceStub = Service.extend({
      init() {
        this.set('accessibleNamespaces', ['ns1']);
      },
    });

    const storeServiceStub = Service.extend({
      adapterFor() {
        return {
          ajax: ajaxStub,
        };
      },
    });
    this.register('service:namespace', namespaceServiceStub);
    this.register('service:store', storeServiceStub);
  });

  test('it renders', async function(assert) {
    this.set('config', { mode: null, paths: [] });
    this.set('paths', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
    await render(hbs`<PathFilterConfigList @config={{config}} @paths={{paths}} />`);

    assert.equal(findAll('#filter-userpass').length, 1);
  });

  test('it sets config.paths', async function(assert) {
    this.set('config', { mode: 'whitelist', paths: [] });
    this.set('paths', [{ path: 'userpass/', type: 'userpass', accessor: 'userpass' }]);
    await render(hbs`<PathFilterConfigList @config={{config}} @paths={{paths}} />`);

    await searchSelect.selectOption();
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
