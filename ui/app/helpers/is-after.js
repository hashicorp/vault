import DateBase from './-date-base';
import { isAfter } from 'date-fns';

export default DateBase.extend({
  compute: function([date1, date2]) {
    this._super(...arguments);

    return isAfter(date1, date2);
  },
});
