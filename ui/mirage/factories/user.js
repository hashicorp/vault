import Mirage from 'ember-cli-mirage';

export default Mirage.Factory.extend({
  name(i) {
    return `Person ${i}`;
  },
  age: 28,
  admin: false,
});
