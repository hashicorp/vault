import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';
import { endOfMonth } from 'date-fns';
import { parseDateString } from 'vault/helpers/parse-date-string';

const getActivityParams = ({ start, end }) => {
  // Expects MM-YYYY format
  // TODO: minStart, maxEnd
  let params = {};
  if (start) {
    let startDate = parseDateString(start);
    if (startDate) {
      // TODO: Replace with formatRFC3339 when date-fns is updated
      params.start_time = Math.round(startDate.getTime() / 1000);
    }
  }
  if (end) {
    let endDate = parseDateString(end);
    if (endDate) {
      // TODO: Replace with formatRFC3339 when date-fns is updated
      params.end_time = Math.round(endOfMonth(endDate).getTime() / 1000);
    }
  }
  return params;
};

export default Route.extend(ClusterRoute, {
  queryParams: {
    start: {
      refreshModel: true,
    },
    end: {
      refreshModel: true,
    },
  },

  model(params) {
    let config = this.store.queryRecord('metrics/config', {}).catch(e => {
      // swallowing error so activity can show if no config permissions
      return {};
    });
    const activityParams = getActivityParams(params);
    let activity = this.store.queryRecord('metrics/activity', activityParams);

    return hash({
      queryStart: params.start,
      queryEnd: params.end,
      activity,
      config,
    });
  },
});
