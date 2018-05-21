import Ember from 'ember';
import columnify from 'columnify';
import { capitalize } from 'vault/helpers/capitalize';
const { computed } = Ember;

export default Ember.Component.extend({
    content: null,
    columns: computed('content', function(){
        let data = this.get('content');
        Object.keys(data).forEach((item) => {
          let val = data[item];
          if (typeof val !== 'string') {
            val = JSON.stringify(val);
          }
          data[item] = val;
        });

        return columnify(data, {
            preserveNewLines: true,
            headingTransform: function(heading) {
                return capitalize([heading]);
            }
        });
    }),
});
