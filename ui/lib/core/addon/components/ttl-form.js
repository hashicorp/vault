/**
 * @module TtlForm
 * TtlForm components are used to enter a Time To Live (TTL) input.
 * This component includes the label and is designed to take
 * a time and unit, and pass an object including seconds,
 * timestring, and go-safe timestring when either values are changed.
 * This picker recalculates the time when the unit is changed by the user (eg 60s -> 1m)
 * To allow the user to toggle this form, use TtlPicker
 *
 * @example
 * ```js
 * <TtlForm @onChange={{this.handleChange}} @initialValue="30m"/>
 * ```
 * @param onChange {Function} - This function will be called when the user changes the value. An object will be passed in as a parameter with values seconds{number}, timeString{string}
 * @param initialValue=null {string} - InitialValue is the duration value which will be shown when the component is loaded. If it can't be parsed, will default to 0.
 * @param changeOnInit=false {boolean} - if true, calls the onChange hook when component is initialized
 * @param label='' {string} - label for the TTL Picker inputs group
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { restartableTask, timeout } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import Duration from '@icholy/duration';
import { guidFor } from '@ember/object/internals';

export const secondsMap = {
  s: 1,
  m: 60,
  h: 3600,
  d: 86400,
};
export const convertToSeconds = (time, unit) => {
  return time * secondsMap[unit];
};
export const convertFromSeconds = (seconds, unit) => {
  return seconds / secondsMap[unit];
};
export const goSafeConvertFromSeconds = (seconds, unit) => {
  // Go only accepts s, m, or h units
  const u = unit === 'd' ? 'h' : unit;
  return convertFromSeconds(seconds, u) + u;
};
export const largestUnitFromSeconds = (seconds) => {
  let unit = 's';
  // get largest unit with no remainder
  if (seconds % secondsMap.d === 0) {
    unit = 'd';
  } else if (seconds % secondsMap.h === 0) {
    unit = 'h';
  } else if (seconds % secondsMap.m === 0) {
    unit = 'm';
  }
  return unit;
};

export default class TtlFormComponent extends Component {
  @tracked time = ''; // if defaultValue is NOT set, then do not display a defaultValue.
  @tracked unit = 's';
  @tracked recalculateSeconds = false;
  @tracked errorMessage = '';

  /* Used internally */
  recalculationTimeout = 5000;
  elementId = 'ttl-' + guidFor(this);

  constructor() {
    super(...arguments);
    const value = this.args.initialValue;
    const changeOnInit = this.args.changeOnInit;
    // if initial value is unset use params passed in as defaults
    // and if no defaultValue is passed in display no time
    if (!value && value !== 0) {
      return;
    }

    // let unit = 's';
    let seconds = 0;
    if (typeof value === 'number') {
      // if the passed value is a number, assume unit is seconds
      seconds = value;
    } else {
      try {
        seconds = Duration.parse(value).seconds();
      } catch (e) {
        // if parsing fails leave it empty
        return;
      }
    }
    const unit = largestUnitFromSeconds(seconds);
    const time = convertFromSeconds(seconds, unit);
    this.time = time;
    this.unit = unit;

    if (changeOnInit) {
      this.handleChange();
    }
  }

  get seconds() {
    return convertToSeconds(this.time, this.unit);
  }
  get unitOptions() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }

  handleChange() {
    const { time, unit, seconds } = this;
    const ttl = {
      seconds,
      timeString: time + unit,
      goSafeTimeString: goSafeConvertFromSeconds(seconds, unit),
    };
    this.args.onChange(ttl);
  }

  keepSecondsRecalculate(newUnit) {
    const newTime = convertFromSeconds(this.seconds, newUnit);
    if (Number.isInteger(newTime)) {
      // Only recalculate if time is whole number
      this.time = newTime;
    }
    this.unit = newUnit;
  }

  @restartableTask
  *updateTime(newTime) {
    this.errorMessage = '';
    const parsedTime = parseInt(newTime, 10);
    if (!newTime) {
      this.errorMessage = 'This field is required';
      return;
    } else if (Number.isNaN(parsedTime)) {
      this.errorMessage = 'Value must be a number';
      return;
    }
    this.time = parsedTime;
    this.handleChange();
    if (Ember.testing) {
      return;
    }
    this.recalculateSeconds = true;
    yield timeout(this.recalculationTimeout);
    this.recalculateSeconds = false;
  }

  @action
  updateUnit(newUnit) {
    if (this.recalculateSeconds) {
      this.unit = newUnit;
    } else {
      this.keepSecondsRecalculate(newUnit);
    }
    this.handleChange();
  }
}
