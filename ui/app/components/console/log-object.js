import Ember from 'ember';
import columnify from 'columnify';
const { computed } = Ember;

export function stringifyObjectValues(data) {
  Object.keys(data).forEach(item => {
    let val = data[item];
    if (typeof val !== 'string') {
      val = JSON.stringify(val);
    }
    data[item] = val;
  });
}

export default Ember.Component.extend({
  content: null,
  columns: computed('content', function() {
    let data = this.get('content');
    stringifyObjectValues(data);

    return columnify(data, {
      preserveNewLines: true,
      headingTransform: function(heading) {
        return Ember.String.capitalize(heading);
      },
    });
  }),
});
