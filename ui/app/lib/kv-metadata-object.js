import ArrayProxy from '@ember/array/proxy';
import { guidFor } from '@ember/object/internals';

export default ArrayProxy.extend({
  createClass() {
    let contents = Object.keys([]).map(key => {
      let obj = {
        name: key,
        value: '',
      };
      guidFor(obj);
      return obj;
    });
    this.setObjects(
      // ARG TODO not sure I need the sorting
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
});
