/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { buildISOTimestamp, parseAPITimestamp } from 'core/utils/date-formatters';
import timestamp from 'core/utils/timestamp';
import { format } from 'date-fns';

import type VersionService from 'vault/services/version';
import type { HTMLElementEvent } from 'forms';

interface OnChangeParams {
  start_time: string;
  end_time: string;
}

interface Args {
  onChange: (callback: OnChangeParams) => void;
  setEditModalVisible: (visible: boolean) => void;
  showEditModal: boolean;
  startTimestamp: string;
  endTimestamp: string;
  billingStartTime: string;
  retentionMonths: number;
}
/**
 * @module ClientsDateRange
 * ClientsDateRange components are used to display the current date range and provide a modal interface for editing the date range.
 *
 * @example
 *
 * <Clients::DateRange @startTimestamp="2018-01-01T14:15:30Z" @endTimestamp="2019-01-31T14:15:30Z" @onChange={{this.handleDateChange}} />
 *
 * @param {function} onChange - callback when a new range is saved.
 * @param {function} setEditModalVisible - callback to tell parent header when modal is opened/closed
 * @param {boolean} showEditModal - boolean for when parent header triggers the modal open
 * @param {string} [startTimestamp] - ISO string timestamp of the start date for the displayed client count data
 * @param {string} [endTimestamp] - ISO string timestamp of the end date for the displayed client count data
 * @param {int} [retentionMonths=48] - number of months for historical billing
 * @param {string} [billingStartTime] - ISO string timestamp of billing start date
 */

export default class ClientsDateRangeComponent extends Component<Args> {
  @service declare readonly version: VersionService;

  @tracked modalStart = ''; // format yyyy-MM
  @tracked modalEnd = ''; // format yyyy-MM

  currentMonth = timestamp.now();
  previousMonth = format(
    new Date(this.currentMonth.getUTCFullYear(), this.currentMonth.getUTCMonth() - 1, 1),
    'yyyy-MM'
  );

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.setTrackedFromArgs();
  }

  get historicalBillingPeriods() {
    // we want whole billing periods
    const count = Math.floor(this.args.retentionMonths / 12);
    const periods: string[] = [];

    for (let i = 1; i <= count; i++) {
      const startDate = parseAPITimestamp(this.args.billingStartTime) as Date;
      const utcYear = startDate.getUTCFullYear() - i;
      startDate.setUTCFullYear(utcYear);
      periods.push(startDate.toISOString());
    }
    return periods;
  }

  get validationError() {
    if (!this.modalStart || !this.modalEnd) {
      return 'You must supply both start and end dates.';
    }
    if (this.modalStart > this.modalEnd) {
      return 'Start date must be before end date.';
    }
    if (this.modalStart > this.previousMonth || this.modalEnd > this.previousMonth) {
      return 'You cannot select the current month or beyond.';
    }
    return null;
  }

  @action onClose() {
    // since the component never gets torn down, we have to manually re-set this on close
    this.setTrackedFromArgs();
    this.args.setEditModalVisible(false);
  }

  @action updateDate(evt: HTMLElementEvent<HTMLInputElement>) {
    const { name, value } = evt.target;
    this[name as 'modalStart' | 'modalEnd'] = value;
  }

  // used for CE date picker
  @action handleSave() {
    if (this.validationError) return;
    const params: OnChangeParams = { start_time: '', end_time: '' };

    if (this.modalStart) {
      params.start_time = this.formatModalTimestamp(this.modalStart, false);
    }

    if (this.modalEnd) {
      params.end_time = this.formatModalTimestamp(this.modalEnd, true);
    }

    this.args.onChange(params);
    this.onClose();
  }

  @action
  updateEnterpriseDateRange(start: string) {
    // We do not send an end_time so the backend handles computing the expected billing period
    this.args.onChange({ start_time: start, end_time: '' });
  }

  // HELPERS
  formatModalTimestamp(modalValue: string, isEndDate: boolean) {
    const [yearString, month] = modalValue.split('-');
    const monthIdx = Number(month) - 1;
    const year = Number(yearString);
    return buildISOTimestamp({ monthIdx, year, isEndDate });
  }

  setTrackedFromArgs() {
    if (this.args.startTimestamp) {
      this.modalStart = parseAPITimestamp(this.args.startTimestamp, 'yyyy-MM') as string;
    }
    if (this.args.endTimestamp) {
      this.modalEnd = parseAPITimestamp(this.args.endTimestamp, 'yyyy-MM') as string;
    }
  }

  // TEMPLATE HELPERS
  formatDropdownDate = (isoTimestamp: string) => parseAPITimestamp(isoTimestamp, 'MMMM yyyy');

  isSelected = (dropdownTimestamp: string) => {
    // Compare against this.args.startTimestamp because it's from the URL query param
    // which is used to query the client count activity API.
    const selectedStart = this.formatDropdownDate(this.args.startTimestamp);
    return this.formatDropdownDate(dropdownTimestamp) === selectedStart;
  };
}
