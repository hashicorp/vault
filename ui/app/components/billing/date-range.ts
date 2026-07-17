/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { dateFormat } from 'core/helpers/date-format';

import type { Month } from 'vault/vault/billing/overview';

interface Args {
  months: Month[];
  onDateChange: (selectedMonth: Month | null | undefined) => void;
  selectedDateOption: Month | null | undefined;
}

export default class BillingDateRange extends Component<Args> {
  get selectedDate() {
    return this.args.selectedDateOption;
  }

  get dateDropdownOptions() {
    const options = [];

    for (const option of this.args.months) {
      const formattedDate = dateFormat([option.month, 'MMMM yyyy'], {});
      options.push({ label: formattedDate, value: option.month });
    }

    return options;
  }

  @action
  updateSelectedDropdownOption(dropdownOption: string) {
    const selectedDateOption: Month | undefined = this.args.months.find(
      (option) => option.month === dropdownOption
    );
    this.args.onDateChange(selectedDateOption);
  }
}
