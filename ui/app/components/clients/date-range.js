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
import { format, parseISO } from 'date-fns';

/**
 * @module ClientsDateRange
 * ClientsDateRange components are used to display the current date range and provide a modal interface for editing the date range.
 *
 * @example
 * ```js
 * <Clients::DateRange @startTime={{this.startTimeISO}} @endTime={{@endTimeISO}} @onChange={{this.handleDateChange}} />
 * ```
 * @param {function} onChange - callback when a new range is saved.
 * @param {string} [startTime] - ISO string timestamp of the current start date
 * @param {string} [endTime] - ISO string timestamp of the current end date
 */

export default class ClientsDateRangeComponent extends Component {
  @service version;

  @tracked showEditModal = false;
  @tracked startDate; // format yyyy-MM
  @tracked endDate; // format yyyy-MM
  currentMonth = format(timestamp.now(), 'yyyy-MM');

  constructor() {
    super(...arguments);
    if (this.args.startTime) {
      this.startDate = format(parseISO(this.args.startTime.toString()), 'yyyy-MM');
    }
    if (this.args.endTime) {
      this.endDate = format(parseISO(this.args.endTime.toString()), 'yyyy-MM');
    }
  }

  formattedDate = (isoTimestamp) => {
    return parseAPITimestamp(isoTimestamp, 'MMMM yyyy');
  };

  get useDefaultDates() {
    return !this.startDate && !this.endDate;
  }

  get validationError() {
    if (this.useDefaultDates) {
      // this means we want to reset, which is fine
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

  @action resetDates() {
    this.startDate = '';
    this.endDate = '';
  }

  @action updateDate(evt) {
    const { name, value } = evt.target;
    if (name === 'end') {
      this.endDate = value;
    } else {
      this.startDate = value;
    }
  }

  @action handleSave() {
    if (this.validationError) return;
    const returnObj = {};
    if (this.startDate) {
      const [year, month] = this.startDate.split('-');
      returnObj.start_time = formatDateObject({ monthIdx: parseInt(month) - 1, year: parseInt(year) }, false);
    }
    if (this.endDate) {
      const [year, month] = this.endDate.split('-');
      returnObj.end_time = formatDateObject({ monthIdx: parseInt(month) - 1, year: parseInt(year) }, true);
    }
    this.args.onChange(returnObj);
    this.showEditModal = false;
  }
}
