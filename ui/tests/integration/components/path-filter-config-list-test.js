import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { typeInSearch, clickTrigger } from 'ember-power-select/test-support/helpers';
import hbs from 'htmlbars-inline-precompile';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import Service from '@ember/service';
import sinon from 'sinon';
import { Promise } from 'rsvp';
import { create } from 'ember-cli-page-object';
import ss from 'vault/tests/pages/components/search-select';

const searchSelect = create(ss);

const resolver = engineResolverFor('replication');

const MOUNTS_RESPONSE = {
  data: {
    secret: {},
    auth: {
      'userpass/': { type: 'userpass', accessor: 'userpass' },
    },
  },
};
const NAMESPACE_MOUNTS_RESPONSE = {
  data: {
    secret: {
      'namespace-kv/': { type: 'kv', accessor: 'kv' },
    },
    auth: {},
  },
};

module('Integration | Component | path filter config list', function(hooks) {
  setupRenderingTest(hooks, { resolver });
  hooks.beforeEach(function() {
    let ajaxStub = sinon.stub().usingPromise(Promise);
    ajaxStub.withArgs('/v1/sys/internal/ui/mounts', 'GET').resolves(MOUNTS_RESPONSE);
    ajaxStub
      .withArgs('/v1/sys/internal/ui/mounts', 'GET', { namespace: 'ns1' })
      .resolves(NAMESPACE_MOUNTS_RESPONSE);
    this.set('ajaxStub', ajaxStub);
    const namespaceServiceStub = Service.extend({
      init() {
        this._super(...arguments);
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
    this.owner.register('service:namespace', namespaceServiceStub);
    this.owner.register('service:store', storeServiceStub);
  });

  test('it renders', async function(assert) {
    this.set('config', { mode: null, paths: [] });
    await render(hbs`<PathFilterConfigList @config={{config}} @paths={{paths}} />`);

    assert.dom('[data-test-component=path-filter-config]').exists();
  });

  test('it sets config.paths', async function(assert) {
    this.set('config', { mode: 'allow', paths: [] });
    this.set('paths', []);
    await render(hbs`<PathFilterConfigList @config={{config}} @paths={{paths}} />`);

    await clickTrigger();
    await typeInSearch('auth');
    await searchSelect.options.objectAt(1).click();
    assert.ok(this.config.paths.includes('auth/userpass/'), 'adds to paths');

    await clickTrigger();
    await assert.equal(searchSelect.options.length, 1, 'has one option left');

    await searchSelect.deleteButtons.objectAt(0).click();
    assert.equal(this.config.paths.length, 0, 'removes from paths');
    await clickTrigger();
    await assert.equal(searchSelect.options.length, 2, 'has both options');
  });

  test('it sets config.mode', async function(assert) {
    this.set('config', { mode: 'allow', paths: [] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} />`);
    await click('#deny');
    assert.equal(this.config.mode, 'deny');
    await click('#no-filtering');
    assert.equal(this.config.mode, null);
  });

  test('it shows a warning when going from a mode to allow all', async function(assert) {
    this.set('config', { mode: 'allow', paths: [] });
    await render(hbs`<PathFilterConfigList @config={{this.config}} />`);
    await click('#no-filtering');
    assert.dom('[data-test-remove-warning]').exists('shows removal warning');
  });

  test('it fetches mounts from a namespace when namespace name is entered', async function(assert) {
    this.set('config', { mode: 'allow', paths: [] });
    this.set('paths', []);
    await render(hbs`<PathFilterConfigList @config={{config}} @paths={{paths}} />`);

    await clickTrigger();
    assert.equal(searchSelect.options.length, 2, 'shows userpass and namespace as an option');
    // type the namespace to trigger an ajax request
    await typeInSearch('ns1');
    assert.equal(searchSelect.options.length, 2, 'has ns and ns mount in the list');
    await searchSelect.options.objectAt(1).click();
    assert.ok(this.config.paths.includes('ns1/namespace-kv/'), 'adds namespace mount to paths');
  });
});
