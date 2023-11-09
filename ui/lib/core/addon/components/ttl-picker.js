/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module TtlPicker
 * TtlPicker components are used to enable and select duration values such as TTL.
 * This component renders a toggle by default, and passes all relevant attributes
 * to TtlForm. Please see that component for additional arguments
 * - allows TTL to be enabled or disabled
 * - recalculates the time when the unit is changed by the user (eg 60s -> 1m)
 *
 * @example
 * <TtlPicker @onChange={{this.handleChange}} @initialEnabled={{@model.myAttribute}} @initialValue={{@model.myAttribute}}/>
 *
 * @param onChange {Function} - This function will be passed a TTL object, which includes enabled{bool}, seconds{number}, timeString{string}, goSafeTimeString{string}.
 * @param initialEnabled=false {Boolean} - Set this value if you want the toggle on when component is mounted
 * @param label="Time to live (TTL)" {String} - Label is the main label that lives next to the toggle. Yielded values will replace the label
 * @param labelDisabled=Label to display when TTL is toggled off
 * @param helperTextEnabled="" {String} - This helper text is shown under the label when the toggle is switched on
 * @param helperTextDisabled="" {String} - This helper text is shown under the label when the toggle is switched off
 * @param initialValue=null {string} - InitialValue is the duration value which will be shown when the component is loaded. If it can't be parsed, will default to 0.
 * @param changeOnInit=false {boolean} - if true, calls the onChange hook when component is initialized
 * @param hideToggle=false {Boolean} - set this value if you'd like to hide the toggle and just leverage the input field
 */

import Component from '@glimmer/component';
import { typeOf } from '@ember/utils';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import Ember from 'ember';
import { restartableTask, timeout } from 'ember-concurrency';
import {
  convertFromSeconds,
  convertToSeconds,
  durationToSeconds,
  goSafeConvertFromSeconds,
  largestUnitFromSeconds,
} from 'core/utils/duration-utils';
export default class TtlPickerComponent extends Component {
  @tracked enableTTL = false;
  @tracked recalculateSeconds = false;
  @tracked time = ''; // if defaultValue is NOT set, then do not display a defaultValue.
  @tracked unit = 's';
  @tracked errorMessage = '';

  /* Used internally */
  recalculationTimeout = 5000;
  elementId = 'ttl-' + guidFor(this);

  get label() {
    if (this.args.label && this.args.labelDisabled) {
      return this.enableTTL ? this.args.label : this.args.labelDisabled;
    }
    return this.args.label || 'Time to live (TTL)';
  }
  get helperText() {
    return this.enableTTL || this.args.hideToggle
      ? this.args.helperTextEnabled
      : this.args.helperTextDisabled;
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
      const parseDuration = durationToSeconds(initialValue);
      // if parsing fails leave it empty
      if (parseDuration === null) return;
      seconds = parseDuration;
    }

    const unit = largestUnitFromSeconds(seconds);
    this.time = convertFromSeconds(seconds, unit);
    this.unit = unit;

    if (this.args.changeOnInit) {
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

  keepSecondsRecalculate(newUnit) {
    const newTime = convertFromSeconds(this.seconds, newUnit);
    if (Number.isInteger(newTime)) {
      // Only recalculate if time is whole number
      this.time = newTime;
    }
    this.unit = newUnit;
  }

  handleChange() {
    const { time, unit, seconds, enableTTL } = this;
    const ttl = {
      enabled: this.args.hideToggle || enableTTL,
      seconds,
      timeString: time + unit,
      goSafeTimeString: goSafeConvertFromSeconds(seconds, unit),
    };
    this.args.onChange(ttl);
  }

  @action
  toggleEnabled() {
    this.enableTTL = !this.enableTTL;
    this.handleChange();
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
