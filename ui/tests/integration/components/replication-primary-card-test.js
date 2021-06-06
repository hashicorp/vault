import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import { CLUSTER_STATES } from 'core/helpers/cluster-states';
import hbs from 'htmlbars-inline-precompile';
const resolver = engineResolverFor('replication');

module('Integration | Component | replication-primary-card', function(hooks) {
  setupRenderingTest(hooks, { resolver });

  test('it renders', async function(assert) {
    const title = 'Last WAL';
    const description = 'WALL-E';
    const metric = '3000';

    this.set('title', title);
    this.set('description', description);
    this.set('metric', metric);

    await render(hbs`
      <ReplicationPrimaryCard
        @title={{title}}
        @description={{description}}
        @metric='3000' />`);

    assert.dom('[data-test-hasError]').doesNotExist('shows no error for non-State cards');

    assert.dom('.last-wal').includesText(title);
    assert.dom('[data-test-description]').includesText(description);
    assert.dom('[data-test-metric]').includesText(metric);
  });

  Object.keys(CLUSTER_STATES).forEach(state => {
    test(`it renders a card when cluster has the ${state} state`, async function(assert) {
      this.set('glyph', CLUSTER_STATES[state].glyph);
      this.set('state', state);

      await render(hbs`
      <ReplicationPrimaryCard
        @title='State'
        @description='Updated every ten seconds.'
        @glyph={{glyph}}
        @metric={{state}} />`);

      if (CLUSTER_STATES[state].isOk) {
        assert.dom('[data-test-hasError]').doesNotExist();
        assert.dom('[data-test-icon]').exists('shows an icon if state is ok');
      } else {
        assert.dom('[data-test-hasError]').exists('shows an error if the cluster state is not ok');
        assert.dom('[data-test-icon]').doesNotExist('does not show an icon if state is not ok');
      }
    });
  });
});
