/**
 * @module TtlPicker2
 * TtlPicker2 components are used to enable and select duration values such as TTL. This component renders a toggle by default and:
 * - allows TTL to be enabled or disabled
 * - recalculates the time when the unit is changed by the user (eg 60s -> 1m)
 *
 * @example
 * ```js
 * <TtlPicker2 @onChange={{handleChange}} @time={{defaultTime}} @unit={{defaultUnit}}/>
 * ```
 * @param onChange {Function} - This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}, goSafeTimeString{string}.
 * @param label="Time to live (TTL)" {String} - Label is the main label that lives next to the toggle.
 * @param helperTextDisabled="Allow tokens to be used indefinitely" {String} - This helper text is shown under the label when the toggle is switched off
 * @param helperTextEnabled="Disable the use of the token after" {String} - This helper text is shown under the label when the toggle is switched on
 * @param description="Longer description about this value, what it does, and why it is useful. Shows up in tooltip next to helpertext"
 * @param time='' {Number} - The time (in the default units) which will be adjustable by the user of the form
 * @param unit="s" {String} - This is the unit key which will show by default on the form. Can be one of `s` (seconds), `m` (minutes), `h` (hours), `d` (days)
 * @param recalculationTimeout=5000 {Number} - This is the time, in milliseconds, that `recalculateSeconds` will be be true after time is updated
 * @param initialValue=null {String} - This is the value set initially (particularly from a string like '30h')
 * @param initialEnabled=null {Boolean} - Set this value if you want the toggle on when component is mounted
 * @param changeOnInit=false {Boolean} - set this value if you'd like the passed onChange function to be called on component initialization
 * @param hideToggle=false {Boolean} - set this value if you'd like to hide the toggle and just leverage the input field
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { typeOf } from '@ember/utils';
import Duration from '@icholy/duration';
import { restartableTask, timeout } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

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
const goSafeConvertFromSeconds = (seconds, unit) => {
  // Go only accepts s, m, or h units
  const u = unit === 'd' ? 'h' : unit;
  return convertFromSeconds(seconds, u) + u;
};

export default class TtlPicker2Component extends Component {
  @tracked enableTTL = false;
  @tracked time = ''; // if defaultValue is NOT set, then do not display a defaultValue.
  @tracked unit = 's';

  get label() {
    return this.args.label || 'Time to live (TTL)';
  }
  get helperTextDisabled() {
    return this.args.helperTextDisabled || 'Allow tokens to be used indefinitely';
  }
  get helperTextEnabled() {
    return this.args.helperTextEnabled || 'Disable the use of the token after';
  }
  // initialValue: null,
  // changeOnInit: false,
  // hideToggle: false,

  constructor() {
    super(...arguments);
    const value = this.args.initialValue;
    const enable = this.args.initialEnabled;
    const changeOnInit = this.args.changeOnInit;
    // if initial value is unset use params passed in as defaults
    // and if no defaultValue is passed in display no time
    if (!value && value !== 0) {
      return;
    }

    let time = 30;
    let unit = 's';
    let setEnable = this.args.hideToggle || this.args.enableTTL;
    if (!!enable || typeOf(enable) === 'boolean') {
      // This allows non-boolean values passed in to be evaluated for truthiness
      setEnable = !!enable;
    }

    if (typeOf(value) === 'number') {
      // if the passed value is a number, assume unit is seconds
      // then check if the value can be converted into a larger unit
      if (value % secondsMap.d === 0) {
        unit = 'd';
      } else if (value % secondsMap.h === 0) {
        unit = 'h';
      } else if (value % secondsMap.m === 0) {
        unit = 'm';
      }
      time = convertFromSeconds(value, unit);
    } else {
      try {
        const seconds = Duration.parse(value).seconds();
        time = seconds;
        // get largest unit with no remainder
        if (seconds % secondsMap.d === 0) {
          unit = 'd';
        } else if (seconds % secondsMap.h === 0) {
          unit = 'h';
        } else if (seconds % secondsMap.m === 0) {
          unit = 'm';
        }

        if (unit !== 's') {
          time = convertFromSeconds(seconds, unit);
        }
      } catch (e) {
        // if parsing fails leave it empty
        return;
      }
    }

    this.time = time;
    this.unit = unit;
    this.enableTTL = setEnable;

    if (changeOnInit) {
      this.handleChange();
    }
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
    const { time, unit, enableTTL, seconds } = this;
    const ttl = {
      enabled: this.hideToggle || enableTTL,
      seconds,
      timeString: time + unit,
      goSafeTimeString: goSafeConvertFromSeconds(seconds, unit),
    };
    this.args.onChange(ttl);
  }

  get helperText() {
    return this.enableTTL || this.hideToggle ? this.helperTextEnabled : this.helperTextDisabled;
  }

  recalculateSeconds = false;
  keepSecondsRecalculate(newUnit) {
    const newTime = convertFromSeconds(this.seconds, newUnit);
    this.time = newTime;
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

  get seconds() {
    return convertToSeconds(this.time, this.unit);
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
  @action
  toggleEnabled() {
    this.enableTTL = !this.enableTTL;
    this.handleChange();
  }
}
