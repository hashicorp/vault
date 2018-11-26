import DateBase from './-date-base';
import { isBefore } from 'date-fns';

export default DateBase.extend({
  compute: function([date1, date2]) {
    this._super(...arguments);

    return isBefore(date1, date2);
  },
});
