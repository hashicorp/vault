import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import columnify from 'columnify';
import { capitalize } from 'vault/helpers/capitalize';


moduleForComponent('console/log-object', 'Integration | Component | console/log object', {
  integration: true
});

test('it renders', function(assert) {

  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });
  const objectContent = {'one':'two', 'three':'four', 'seven': {'five': 'six'}, 'eight': [5, 6]};
  const data = {'one':'two', 'three':'four', 'seven': {'five': 'six'}, 'eight': [5, 6]};
  Object.keys(data).forEach((item) => {
    data[item] = JSON.stringify(data[item]);
  });
  const expectedText = columnify(data, { 
    preserveNewLines: true, 
    headingTransform: function(heading) {
        return capitalize([heading]);
    }
  });

  this.set('content', objectContent);

  this.render(hbs`{{console/log-object content=content}}`);

  assert.dom('pre').includesText(`${expectedText}`);


});
