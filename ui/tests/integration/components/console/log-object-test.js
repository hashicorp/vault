import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import columnify from 'columnify';
import { capitalize } from 'vault/helpers/capitalize';
import { stringifyObjectValues } from 'vault/components/console/log-object';

module('Integration | Component | console/log object', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    const objectContent = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
    const data = { one: 'two', three: 'four', seven: { five: 'six' }, eight: [5, 6] };
    stringifyObjectValues(data);
    const expectedText = columnify(data, {
      preserveNewLines: true,
      headingTransform: function(heading) {
        return capitalize([heading]);
      },
    });

    this.set('content', objectContent);

    await render(hbs`{{console/log-object content=content}}`);

    assert.dom('pre').includesText(`${expectedText}`);
  });
});
