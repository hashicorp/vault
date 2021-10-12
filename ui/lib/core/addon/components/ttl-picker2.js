/**
 * @module TtlPicker2
 * TtlPicker2 components are used to enable and select 'time to live' values. Use this TtlPicker2 instead of TtlPicker if you:
 * - Want the TTL to be enabled or disabled
 * - Want to have the time recalculated by default when the unit changes (eg 60s -> 1m)
 *
 * @example
 * ```js
 * <TtlPicker2 @onChange={{handleChange}} @time={{defaultTime}} @unit={{defaultUnit}}/>
 * ```
 * @param onChange {Function} - This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}.
 * @param label="Time to live (TTL)" {String} - Label is the main label that lives next to the toggle.
 * @param helperTextDisabled="Allow tokens to be used indefinitely" {String} - This helper text is shown under the label when the toggle is switched off
 * @param helperTextEnabled="Disable the use of the token after" {String} - This helper text is shown under the label when the toggle is switched on
 * @param description="Longer description about this value, what it does, and why it is useful. Shows up in tooltip next to helpertext"
 * @param time=30 {Number} - The time (in the default units) which will be adjustable by the user of the form
 * @param unit="s" {String} - This is the unit key which will show by default on the form. Can be one of `s` (seconds), `m` (minutes), `h` (hours), `d` (days)
 * @param recalculationTimeout=5000 {Number} - This is the time, in milliseconds, that `recalculateSeconds` will be be true after time is updated
 * @param initialValue=null {String} - This is the value set initially (particularly from a string like '30h')
 * @param initialEnabled=null {Boolean} - Set this value if you want the toggle on when component is mounted
 * @param changeOnInit=false {Boolean} - set this value if you'd like the passed onChange function to be called on component initialization
 */

import { computed } from '@ember/object';
import { typeOf } from '@ember/utils';
import Duration from '@icholy/duration';
import TtlForm from './ttl-form';
import layout from '../templates/components/ttl-picker2';

const secondsMap = {
  s: 1,
  m: 60,
  h: 3600,
  d: 86400,
};
const validUnits = ['s', 'm', 'h', 'd'];
const convertFromSeconds = (seconds, unit) => {
  return seconds / secondsMap[unit];
};

export default TtlForm.extend({
  layout,
  enableTTL: false,
  label: 'Time to live (TTL)',
  helperTextDisabled: 'Allow tokens to be used indefinitely',
  helperTextEnabled: 'Disable the use of the token after',
  description: '',
  time: 30,
  unit: 's',
  initialValue: null,
  changeOnInit: false,

  init() {
    this._super(...arguments);
    const value = this.initialValue;
    const enable = this.initialEnabled;
    const changeOnInit = this.changeOnInit;
    // if initial value is unset use params passed in as defaults
    if (!value && value !== 0) {
      return;
    }

    let time = 30;
    let unit = 's';
    let setEnable = this.enableTTL;
    if (!!enable || typeOf(enable) === 'boolean') {
      // This allows non-boolean values passed in to be evaluated for truthiness
      setEnable = !!enable;
    }

    if (typeOf(value) === 'number') {
      // if the passed value is a number, assume unit is seconds
      time = value;
    } else {
      try {
        const seconds = Duration.parse(value).seconds();
        const lastDigit = value.toString().substring(value.length - 1);
        if (validUnits.indexOf(lastDigit) >= 0 && lastDigit !== 's') {
          time = convertFromSeconds(seconds, lastDigit);
          unit = lastDigit;
        } else {
          time = seconds;
        }
      } catch (e) {
        // if parsing fails leave as default 30s
      }
    }

    this.setProperties({
      time,
      unit,
      enableTTL: setEnable,
    });

    if (changeOnInit) {
      this.handleChange();
    }
  },

  unitOptions: computed(function() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }),
  handleChange() {
    let { time, unit, enableTTL, seconds } = this;
    const ttl = {
      enabled: enableTTL,
      seconds,
      timeString: time + unit,
    };
    this.onChange(ttl);
  },

  helperText: computed(
    'enableTTL',
    'helperTextDisabled',
    'helperTextEnabled',
    'helperTextSet',
    'helperTextUnset',
    function() {
      return this.enableTTL ? this.helperTextEnabled : this.helperTextDisabled;
    }
  ),

  recalculateSeconds: false,
  actions: {
    toggleEnabled() {
      this.toggleProperty('enableTTL');
      this.handleChange();
    },
  },
});
