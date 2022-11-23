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

import Component from '@glimmer/component';
import { typeOf } from '@ember/utils';
import Duration from '@icholy/duration';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { convertFromSeconds, goSafeConvertFromSeconds, secondsMap } from './ttl-form';
export default class TtlPicker2Component extends Component {
  @tracked enableTTL = false;
  @tracked time = ''; // if defaultValue is NOT set, then do not display a defaultValue.
  @tracked unit = 's';
  @tracked recalculateSeconds = false;

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

    let seconds = 0;
    let time = 0;
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
        seconds = Duration.parse(value).seconds();
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

    if (changeOnInit) {
      // Mock what TtlForm would return
      this.handleChange({
        seconds,
        timeString: time + unit,
        goSafeTimeString: goSafeConvertFromSeconds(seconds, unit),
      });
    }
    this.enableTTL = setEnable;
  }

  @action
  handleChange(ttlObj) {
    const ttl = {
      ...ttlObj,
      enabled: this.hideToggle || this.enableTTL,
    };
    this.args.onChange(ttl);
  }

  get helperText() {
    return this.enableTTL || this.hideToggle ? this.helperTextEnabled : this.helperTextDisabled;
  }

  @action
  toggleEnabled() {
    this.enableTTL = !this.enableTTL;
    this.handleChange();
  }
}
