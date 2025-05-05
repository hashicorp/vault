/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { formatDateObject } from 'core/utils/client-count-utils';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import timestamp from 'core/utils/timestamp';
import { format } from 'date-fns';
import type VersionService from 'vault/services/version';
import type { HTMLElementEvent } from 'forms';

interface OnChangeParams {
  start_time: number | undefined;
  end_time: number | undefined;
}
interface Args {
  onChange: (callback: OnChangeParams) => void;
  startTime: string;
  endTime: string;
  billingStartTime: string;
  retentionMonths: number;
}
/**
 * @module ClientsDateRange
 * ClientsDateRange components are used to display the current date range and provide a modal interface for editing the date range.
 *
 * @example
 *
 * <Clients::DateRange @startTime="2018-01-01T14:15:30Z" @endTime="2019-01-31T14:15:30Z" @onChange={{this.handleDateChange}} />
 *
 * @param {function} onChange - callback when a new range is saved.
 * @param {string} [startTime] - ISO string timestamp of the current start date
 * @param {string} [endTime] - ISO string timestamp of the current end date
 * @param {int} [retentionMonths=48] - number of months for historical billing
 * @param {string} [billingStartTime] - ISO string timestamp of billing start date
 */

export default class ClientsDateRangeComponent extends Component<Args> {
  @service declare readonly version: VersionService;

  @tracked showEditModal = false;
  @tracked startDate = ''; // format yyyy-MM
  @tracked endDate = ''; // format yyyy-MM
  @tracked selectedStart = this.args.billingStartTime;
  currentMonth = format(timestamp.now(), 'yyyy-MM');

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.setTrackedFromArgs();
  }

  setTrackedFromArgs() {
    if (this.args.startTime) {
      this.startDate = parseAPITimestamp(this.args.startTime, 'yyyy-MM') as string;
      this.selectedStart = this.formattedDate(this.args.startTime) as string;
    }
    if (this.args.endTime) {
      this.endDate = parseAPITimestamp(this.args.endTime, 'yyyy-MM') as string;
    }
  }

  formattedDate = (isoTimestamp: string) => {
    return parseAPITimestamp(isoTimestamp, 'MMMM yyyy');
  };

  get historicalBillingPeriods() {
    // we want whole billing periods
    const count = Math.floor(this.args.retentionMonths / 12);
    const periods: string[] = [];

    for (let i = 1; i <= count; i++) {
      const startDate = new Date(this.args.billingStartTime);
      const utcMonth = startDate.getUTCMonth();
      const utcYear = startDate.getUTCFullYear() - i;

      startDate.setUTCFullYear(utcYear);
      startDate.setUTCMonth(utcMonth);

      periods.push(startDate.toISOString());
    }

    return periods;
  }

  get useDefaultDates() {
    return !this.startDate && !this.endDate;
  }

  get validationError() {
    if (this.useDefaultDates && this.version.isEnterprise) {
      // this means we want to reset, which is fine for ent only
      return null;
    }
    if (!this.startDate || !this.endDate) {
      return 'You must supply both start and end dates.';
    }
    if (this.startDate > this.endDate) {
      return 'Start date must be before end date.';
    }
    return null;
  }

  @action onClose() {
    // since the component never gets torn down, we have to manually re-set this on close
    this.setTrackedFromArgs();
    this.showEditModal = false;
  }

  @action resetDates() {
    this.startDate = '';
    this.endDate = '';
  }

  @action updateDate(evt: HTMLElementEvent<HTMLInputElement>) {
    const { name, value } = evt.target;
    if (name === 'end') {
      this.endDate = value;
    } else {
      this.startDate = value;
    }
  }

  // used for CE date picker
  @action handleSave() {
    if (this.validationError) return;
    const params: OnChangeParams = {
      start_time: undefined,
      end_time: undefined,
    };
    if (this.startDate) {
      const [year, month] = this.startDate.split('-');
      if (year && month) {
        params.start_time = formatDateObject({ monthIdx: parseInt(month) - 1, year: parseInt(year) }, false);
      }
    }
    if (this.endDate) {
      const [year, month] = this.endDate.split('-');
      if (year && month) {
        params.end_time = formatDateObject({ monthIdx: parseInt(month) - 1, year: parseInt(year) }, true);
      }
    }

    this.args.onChange(params);
    this.onClose();
  }

  @action
  updateEnterpriseDateRange(start: string) {
    const params: OnChangeParams = {
      start_time: undefined,
      end_time: undefined,
    };

    const [year, month] = start.split('-');
    if (year && month) {
      // pass true for isEnd even for start because we want to go off last day of month here, otherwise we risk
      // setting it to a start_time that is for the previous billing period
      params.start_time = formatDateObject({ monthIdx: parseInt(month) - 1, year: parseInt(year) }, true);
    }

    this.args.onChange(params);
  }
}
