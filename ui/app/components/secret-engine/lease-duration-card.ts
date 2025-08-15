/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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

interface Args {
  label: string;
  labelDisabled: boolean;
  helperTextEnabled: string;
  helperTextDisabled: string;
  initialEnabled?: boolean;
  initialValue: string;
  hideToggle?: boolean;
  changeOnInit?: boolean;
  onChange: (ttl: {
    enabled: boolean;
    seconds: number;
    timeString: string;
    goSafeTimeString: string;
  }) => void;
}

export default class LeaseDurationCard extends Component<Args> {
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

  constructor(owner: unknown, args: Args) {
    super(owner, args);
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
    this.time = convertFromSeconds(seconds, unit).toString();
    this.unit = unit;

    if (this.args.changeOnInit) {
      this.handleChange();
    }
  }

  get seconds() {
    return convertToSeconds(parseInt(this.time), this.unit);
  }
  get unitOptions() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }

  keepSecondsRecalculate(newUnit: any) {
    const newTime = convertFromSeconds(this.seconds, newUnit);
    if (Number.isInteger(newTime)) {
      // Only recalculate if time is whole number
      this.time = newTime.toString();
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
  *updateTime(newTime: any) {
    this.errorMessage = '';
    const parsedTime = parseInt(newTime, 10);
    if (!newTime) {
      this.errorMessage = 'This field is required';
      return;
    } else if (Number.isNaN(parsedTime)) {
      this.errorMessage = 'Value must be a number';
      return;
    }
    this.time = parsedTime.toString();
    this.handleChange();
    if (Ember.testing) {
      return;
    }
    this.recalculateSeconds = true;
    yield timeout(this.recalculationTimeout);
    this.recalculateSeconds = false;
  }

  @action
  updateUnit(newUnit: any) {
    if (this.recalculateSeconds) {
      this.unit = newUnit;
    } else {
      this.keepSecondsRecalculate(newUnit);
    }
    this.handleChange();
  }
}
