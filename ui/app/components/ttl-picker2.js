/**
 * @module TtlPicker2
 * TtlPicker2 components are used to...
 *
 * @example
 * ```js
 * <TtlPicker2 @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {function} [onChange] - requiredParam is...
 * @param {boolean} [defaultEnabled] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@ember/component';
import { computed } from '@ember/object';

const TEMP_timeUpdatedRecently = false;

const convertToSeconds = (time, unit) => {
  const toSeconds = {
    s: 1,
    m: 60,
    h: 3600,
    d: 86400,
  };

  console.log('CONVERT TO SECONDS', time, unit);
  return time * toSeconds[unit];
};

const convertFromSeconds = (seconds, unit) => {
  const fromSeconds = {
    s: 1,
    m: 60,
    h: 3600,
    d: 86400,
  };
  const cfromSeconds = seconds / fromSeconds[unit];
  console.log(seconds, unit);
  console.log({ cfromSeconds });
  return cfromSeconds;
};

export default Component.extend({
  enableTTL: true,
  time: 30,
  unit: 'm',
  unitOptions: computed(function() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }),

  recalculateTime(newUnit) {
    // get converted value from current seconds value
    const newTime = convertFromSeconds(this.seconds, newUnit);
    this.setProperties({
      time: newTime,
      unit: newUnit,
    });
  },

  seconds: computed('time', 'unit', function() {
    return convertToSeconds(this.time, this.unit);
  }),
  label: 'Time to live (TTL)',
  helperTextDisabled: 'Allow tokens to be used indefinitely',
  helperTextEnabled: 'Disable the use of the token after',
  helperText: computed('enableTTL', 'helperTextUnset', 'helperTextSet', function() {
    return this.enableTTL ? this.helperTextEnabled : this.helperTextDisabled;
  }),
  errorMessage: null,
  actions: {
    updateUnit(newUnit) {
      console.log(newUnit);
      if (TEMP_timeUpdatedRecently) {
        this.set('unit', newUnit);
      } else {
        this.recalculateTime(newUnit);
      }
    },
    updateTime(newTime) {
      this.set('time', newTime);
    },
  },
});
