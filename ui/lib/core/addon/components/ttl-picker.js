/**
 * @module TtlPicker
 * TtlPicker components are used to enable and select duration values such as TTL.
 * This component renders a toggle by default, and passes all relevant attributes
 * to TtlForm. Please see that component for additional arguments
 * - allows TTL to be enabled or disabled
 * - recalculates the time when the unit is changed by the user (eg 60s -> 1m)
 *
 * @example
 * ```js
 * <TtlPicker @onChange={{this.handleChange}} @initialEnabled={{@model.myAttribute}} @initialValue={{@model.myAttribute}}/>
 * ```
 * @param onChange {Function} - This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}, goSafeTimeString{string}.
 * @param initialEnabled=false {Boolean} - Set this value if you want the toggle on when component is mounted
 * @param label="Time to live (TTL)" {String} - Label is the main label that lives next to the toggle.
 * @param helperTextDisabled="Allow tokens to be used indefinitely" {String} - This helper text is shown under the label when the toggle is switched off
 * @param helperTextEnabled="Disable the use of the token after" {String} - This helper text is shown under the label when the toggle is switched on
 * @param hideToggle=false {Boolean} - DEPRECATED set this value if you'd like to hide the toggle and just leverage the input field
 */

import Component from '@glimmer/component';
import { typeOf } from '@ember/utils';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import Duration from '@icholy/duration';
import { convertFromSeconds, goSafeConvertFromSeconds, largestUnitFromSeconds } from './ttl-form';

const DEFAULT_TTL = { seconds: 0, timeString: '0s', goSafeTimeString: '0s' };
export default class TtlPickerComponent extends Component {
  @tracked enableTTL = false;
  @tracked recalculateSeconds = false;
  @tracked ttl = DEFAULT_TTL; // internal tracking

  get label() {
    return this.args.label || 'Time to live (TTL)';
  }
  get helperTextDisabled() {
    return this.args.helperTextDisabled || 'Allow tokens to be used indefinitely';
  }
  get helperTextEnabled() {
    return this.args.helperTextEnabled || 'Disable the use of the token after';
  }

  constructor() {
    super(...arguments);
    const enable = this.args.initialEnabled;

    let setEnable = !!this.args.hideToggle;
    if (!!enable || typeOf(enable) === 'boolean') {
      // This allows non-boolean values passed in to be evaluated for truthiness
      setEnable = !!enable;
    }

    this.enableTTL = setEnable;
    this.initializeTtl();
  }

  initializeTtl() {
    const initialValue = this.args.initialValue;
    let seconds = 0;
    if (typeof initialValue === 'number') {
      // if the passed value is a number, assume unit is seconds
      seconds = initialValue;
    } else {
      try {
        seconds = Duration.parse(initialValue).seconds();
      } catch (e) {
        // if parsing fails leave it empty
        return;
      }
    }
    const unit = largestUnitFromSeconds(seconds);
    const time = convertFromSeconds(seconds, unit);
    this.ttl = {
      seconds,
      timeString: time + unit,
      goSafeTimeString: goSafeConvertFromSeconds(seconds, unit),
    };
    if (this.args.changeOnInit) {
      this.handleChange();
    }
  }

  @action
  handleChange(ttlObj) {
    if (ttlObj) {
      // Update local TTL object if triggered from child TtlForm
      this.ttl = ttlObj;
    }
    const ttl = {
      ...this.ttl,
      enabled: this.args.hideToggle || this.enableTTL,
    };
    this.args.onChange(ttl);
  }

  get helperText() {
    return this.enableTTL || this.args.hideToggle ? this.helperTextEnabled : this.helperTextDisabled;
  }

  @action
  toggleEnabled() {
    this.enableTTL = !this.enableTTL;
    this.handleChange();
  }
}
