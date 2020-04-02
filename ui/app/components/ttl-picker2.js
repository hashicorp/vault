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
 * @param {function} onChange - This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}.
 * @param {string} [label='Time to live (TTL)'] - Label is the main label that lives next to the toggle.
 * @param {string} [helperTextDisabled='Allow tokens to be used indefinitely'] - This helper text is shown under the label when the toggle is switched off
 * @param {string} [helperTextEnabled='Disable the use of the token after'] - This helper text is shown under the label when the toggle is switched on
 * @param {number} [time=30] - The time (in the default units) which will be adjustable by the user of the form
 * @param {string} [unit='s'] - This is the unit key which will show by default on the form. Can be one of `s` (seconds), `m` (minutes), `h` (hours), `d` (days)
 * @param {number} [recalculationTimeout=5000] - This is the time, in milliseconds, that `recalculateSeconds` will be be true after time is updated
 */

import Ember from 'ember';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { task, timeout } from 'ember-concurrency';

const secondsMap = {
  s: 1,
  m: 60,
  h: 3600,
  d: 86400,
};
const convertToSeconds = (time, unit) => {
  return time * secondsMap[unit];
};
const convertFromSeconds = (seconds, unit) => {
  return seconds / secondsMap[unit];
};

export default Component.extend({
  enableTTL: false,
  label: 'Time to live (TTL)',
  helperTextDisabled: 'Allow tokens to be used indefinitely',
  helperTextEnabled: 'Disable the use of the token after',
  time: 30,
  unit: 'm',
  recalculationTimeout: 5000,
  unitOptions: computed(function() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }),

  TTL: computed('enableTTL', 'seconds', function() {
    let { time, unit, enableTTL, seconds } = this.getProperties('time', 'unit', 'enableTTL', 'seconds');
    return {
      enabled: enableTTL,
      seconds,
      timeString: time + unit,
    };
  }),

  updateTime: task(function*(newTime) {
    this.set('errorMessage', '');
    let parsedTime;
    parsedTime = parseInt(newTime, 10);
    if (!newTime) {
      this.set('errorMessage', 'This field is required');
      return;
    } else if (Number.isNaN(parsedTime)) {
      this.set('errorMessage', 'Value must be a number');
      return;
    }
    this.set('time', parsedTime);
    this.onChange(this.TTL);
    if (Ember.testing) {
      return;
    }
    this.set('recalculateSeconds', true);
    yield timeout(this.recalculationTimeout);
    this.set('recalculateSeconds', false);
  }).restartable(),

  recalculateTime(newUnit) {
    const newTime = convertFromSeconds(this.seconds, newUnit);
    this.setProperties({
      time: newTime,
      unit: newUnit,
    });
  },

  seconds: computed('time', 'unit', function() {
    return convertToSeconds(this.time, this.unit);
  }),
  helperText: computed('enableTTL', 'helperTextUnset', 'helperTextSet', function() {
    return this.enableTTL ? this.helperTextEnabled : this.helperTextDisabled;
  }),
  errorMessage: null,
  recalculateSeconds: false,
  actions: {
    updateUnit(newUnit) {
      if (this.recalculateSeconds) {
        this.set('unit', newUnit);
      } else {
        this.recalculateTime(newUnit);
      }
      this.onChange(this.TTL);
    },
    toggleEnabled() {
      this.toggleProperty('enableTTL');
      this.onChange(this.TTL);
    },
  },
});
