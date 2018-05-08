import Ember from 'ember';
import columnify from 'columnify';
import { capitalize } from 'vault/helpers/capitalize';
const { computed } = Ember;

export default Ember.Component.extend({
	content: null,
	columns: computed('content', function(){
		return columnify(this.get('content'), { 
			preserveNewLines: true, 
            headingTransform: function(heading) {
              	return capitalize([heading]);
            }
		});
	}),
});
