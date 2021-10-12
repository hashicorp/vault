/**
 * @module TtlForm
 * TtlForm components are used to enter a Time To Live (TTL) input.
 * This component does not include a label and is designed to take
 * a time and unit, and pass an object including seconds and
 * timestring when those two values are changed.
 *
 * @example
 * ```js
 * <TtlForm @onChange={{action handleChange}} @unit="m"/>
 * ```
 * @param {function} onChange - This function will be called when the user changes the value. An object will be passed in as a parameter with values seconds{number}, timeString{string}
 * @param {number} [time] - Time is the value that will be passed into the value input. Can be null/undefined to start if input is required.
 * @param {unit} [unit="s"] - This is the unit key which will show by default on the form. Can be one of `s` (seconds), `m` (minutes), `h` (hours), `d` (days)
 * @param {number} [recalculationTimeout=5000] - This is the time, in milliseconds, that `recalculateSeconds` will be be true after time is updated
 */

import Ember from 'ember';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { task, timeout } from 'ember-concurrency';
import layout from '../templates/components/ttl-form';

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
  layout,
  time: '',
  unit: 's',

  /* Used internally */
  recalculationTimeout: 5000,
  recalculateSeconds: false,
  errorMessage: null,
  unitOptions: computed(function() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }),
  handleChange() {
    let { time, unit, seconds } = this;
    const ttl = {
      seconds,
      timeString: time + unit,
    };
    this.onChange(ttl);
  },
  keepSecondsRecalculate(newUnit) {
    const newTime = convertFromSeconds(this.seconds, newUnit);
    this.setProperties({
      time: newTime,
      unit: newUnit,
    });
  },
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
    this.handleChange();
    if (Ember.testing) {
      return;
    }
    this.set('recalculateSeconds', true);
    yield timeout(this.recalculationTimeout);
    this.set('recalculateSeconds', false);
  }).restartable(),

  seconds: computed('time', 'unit', function() {
    return convertToSeconds(this.time, this.unit);
  }),

  actions: {
    updateUnit(newUnit) {
      if (this.recalculateSeconds) {
        this.set('unit', newUnit);
      } else {
        this.keepSecondsRecalculate(newUnit);
      }
      this.handleChange();
    },
  },
});
