import ArrayProxy from '@ember/array/proxy';
import { typeOf } from '@ember/utils';
import { guidFor } from '@ember/object/internals';

export default ArrayProxy.extend({
  fromJSON(json) {
    if (json && typeOf(json) !== 'object') {
      throw new Error('Vault expects data to be formatted as an JSON object.');
    }
    let contents = Object.keys(json || []).map(key => {
      let obj = {
        name: key,
        value: json[key],
      };
      guidFor(obj);
      return obj;
    });
    this.setObjects(
      contents.sort((a, b) => {
        if (a.name === '') {
          return 1;
        }
        if (b.name === '') {
          return -1;
        }
        return a.name.localeCompare(b.name);
      })
    );
    return this;
  },

  fromJSONString(jsonString) {
    return this.fromJSON(JSON.parse(jsonString));
  },

  toJSON(includeBlanks = false) {
    return this.reduce((obj, item) => {
      if (!includeBlanks && item.value === '' && item.name === '') {
        return obj;
      }
      let val = typeof item.value === 'undefined' ? '' : item.value;
      obj[item.name || ''] = val;
      return obj;
    }, {});
  },

  toJSONString(includeBlanks) {
    return JSON.stringify(this.toJSON(includeBlanks), null, 2);
  },

  isAdvanced() {
    return this.any(item => typeof item.value !== 'string');
  },
});
