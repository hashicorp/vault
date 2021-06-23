import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';
import { getTime } from 'date-fns';
import { parseDateString } from 'vault/helpers/parse-date-string';

const getActivityParams = ({ start, end }) => {
  // Expects MM-yyyy format
  // TODO: minStart, maxEnd
  let params = {};
  if (start) {
    let startDate = parseDateString(start);
    if (startDate) {
      // TODO: Replace with formatRFC3339 when date-fns is updated
      // converts to milliseconds, divide by 1000 to get epoch
      params.start_time = getTime(startDate) / 1000;
    }
  }
  if (end) {
    let endDate = parseDateString(end);
    if (endDate) {
      // TODO: Replace with formatRFC3339 when date-fns is updated
      params.end_time = getTime(endDate) / 1000;
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
      console.debug(e);
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
