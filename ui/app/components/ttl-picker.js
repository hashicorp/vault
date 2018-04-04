import Ember from 'ember';

const { computed, get, set } = Ember;

export default Ember.Component.extend({
  'data-test-component': 'ttl-picker',
  classNames: 'field',
  setDefaultValue: true,
  onChange: () => {},
  labelText: 'TTL',
  labelClass: '',
  time: 30,
  unit: 'm',
  initialValue: null,
  unitOptions: [
    { label: 'seconds', value: 's' },
    { label: 'minutes', value: 'm' },
    { label: 'hours', value: 'h' },
    { label: 'days', value: 'd' },
  ],

  ouputSeconds: false,

  convertToSeconds(time, unit) {
    const toSeconds = {
      s: 1,
      m: 60,
      h: 3600,
    };

    return time * toSeconds[unit];
  },

  TTL: computed('time', 'unit', function() {
    let { time, unit, outputSeconds } = this.getProperties('time', 'unit', 'outputSeconds');
    //convert to hours
    if (unit === 'd') {
      time = time * 24;
      unit = 'h';
    }
    const timeString = time + unit;
    return outputSeconds ? this.convertToSeconds(time, unit) : timeString;
  }),

  didRender() {
    this._super(...arguments);
    if (get(this, 'setDefaultValue') === false) {
      return;
    }
    get(this, 'onChange')(get(this, 'TTL'));
  },

  init() {
    this._super(...arguments);
    if (!get(this, 'onChange')) {
      throw new Ember.Error('`onChange` handler is a required attr in `' + this.toString() + '`.');
    }
    if (get(this, 'initialValue')) {
      this.parseAndSetTime();
    }
  },

  parseAndSetTime() {
    const value = get(this, 'initialValue');
    const seconds = Ember.typeOf(value) === 'number' ? value : Duration.parse(value).seconds();

    this.set('time', seconds);
    this.set('unit', 's');
  },

  actions: {
    changedValue(key, event) {
      const { type, value, checked } = event.target;
      const val = type === 'checkbox' ? checked : value;
      set(this, key, val);
      get(this, 'onChange')(get(this, 'TTL'));
    },
  },
});
