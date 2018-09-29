import { capitalize } from '@ember/string';
import Component from '@ember/component';
import { computed } from '@ember/object';
import columnify from 'columnify';

export function stringifyObjectValues(data) {
  Object.keys(data).forEach(item => {
    let val = data[item];
    if (typeof val !== 'string') {
      val = JSON.stringify(val);
    }
    data[item] = val;
  });
}

export default Component.extend({
  content: null,
  columns: computed('content', function() {
    let data = this.get('content');
    stringifyObjectValues(data);

    return columnify(data, {
      preserveNewLines: true,
      headingTransform: function(heading) {
        return capitalize(heading);
      },
    });
  }),
});
