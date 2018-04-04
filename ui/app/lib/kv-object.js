import Ember from 'ember';

export default Ember.ArrayProxy.extend({
  fromJSON(json) {
    const contents = Object.keys(json || []).map(key => {
      let obj = {
        name: key,
        value: json[key],
      };
      Ember.guidFor(obj);
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
      obj[item.name || ''] = item.value || '';
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
