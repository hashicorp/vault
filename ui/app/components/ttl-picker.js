import { typeOf } from '@ember/utils';
import EmberError from '@ember/error';
import Component from '@ember/component';
import { set, computed } from '@ember/object';
import Duration from 'Duration.js';

const ERROR_MESSAGE = 'TTLs must be specified in whole number increments, please enter a whole number.';

export default Component.extend({
  'data-test-component': 'ttl-picker',
  classNames: 'field',
  setDefaultValue: true,
  onChange: () => {},
  labelText: 'TTL',
  labelClass: '',
  time: 30,
  unit: 'm',
  initialValue: null,
  errorMessage: null,
  unitOptions: computed(function() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }),

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

  didInsertElement() {
    this._super(...arguments);
    if (this.setDefaultValue === false) {
      return;
    }
    this.onChange(this.TTL);
  },

  init() {
    this._super(...arguments);
    if (!this.onChange) {
      throw new EmberError('`onChange` handler is a required attr in `' + this.toString() + '`.');
    }
    if (this.initialValue != undefined) {
      this.parseAndSetTime();
    }
  },

  parseAndSetTime() {
    let value = this.initialValue;
    let seconds = typeOf(value) === 'number' ? value : 30;
    try {
      seconds = Duration.parse(value).seconds();
    } catch (e) {
      // if parsing fails leave as default 30
    }

    this.set('time', seconds);
    this.set('unit', 's');
  },

  actions: {
    changedValue(key, event) {
      let { type, value, checked } = event.target;
      let val = type === 'checkbox' ? checked : value;
      if (val && key === 'time') {
        val = parseInt(val, 10);
        if (Number.isNaN(val)) {
          this.set('errorMessage', ERROR_MESSAGE);
          return;
        }
      }
      this.set('errorMessage', null);
      set(this, key, val);
      this.onChange(this.TTL);
    },
  },
});
