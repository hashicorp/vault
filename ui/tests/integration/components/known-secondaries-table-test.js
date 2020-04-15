import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import hbs from 'htmlbars-inline-precompile';
const resolver = engineResolverFor('replication');

const REPLICATION_ATTRS = {
  knownSecondaries: ['firstSecondary', 'secondary-2'],
};

module('Integration | Component | replication known-secondaries-table', function(hooks) {
  setupRenderingTest(hooks, { resolver });

  hooks.beforeEach(function() {
    this.set('replicationAttrs', REPLICATION_ATTRS);
  });

  test('it renders with a table of known secondaries', async function(assert) {
    await render(hbs`<KnownSecondariesCard @replicationAttrs={{replicationAttrs}} />`);

    assert.dom('[data-test-known-secondaries-table]').exists();
    REPLICATION_ATTRS.knownSecondaries.forEach(secondaryId => {
      assert
        .dom(`[data-test-secondaries-row=${secondaryId}]`)
        .exists('shows a table row for each known secondary');
    });
  });

  test('shows unknown if url or connection are unknown', async function(assert) {
    // TODO: update this test  once we know what the shape of knownSecondaries is & how we will access the url and connection
    const secondaryDetailsMissing = {
      knownSecondaries: ['firstSecondary', { id: 'secondary-2', url: 'http://stuff.com/' }],
    };
    this.set('replicationAttrs', secondaryDetailsMissing);
    await render(hbs`<KnownSecondariesCard @replicationAttrs={{replicationAttrs}} />`);

    this.element.querySelectorAll('[data-test-url]').forEach((td, i) => {
      let expectedUrl = secondaryDetailsMissing.knownSecondaries[i].url || 'unknown';
      return assert.equal(td.textContent.trim(), expectedUrl);
    });
  });
});
