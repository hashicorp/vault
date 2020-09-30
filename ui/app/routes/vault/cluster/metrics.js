import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';
import { hash } from 'rsvp';
import { endOfMonth, isValid } from 'date-fns';

const getActivityParams = ({ start, end }) => {
  // Expects MM-YYYY format
  // TODO: minStart, maxEnd
  let params = {};
  if (start) {
    let datePieces = start.split('-');
    if (datePieces.length > 1) {
      let startDate = new Date(Date.UTC(datePieces[1], datePieces[0] - 1, 1));
      if (isValid(startDate)) {
        // TODO: Replace with formatRFC3339
        params.start_time = Math.round(startDate.getTime() / 1000);
      }
    }
  }
  if (end) {
    let datePieces = end.split('-');
    if (datePieces.length > 1) {
      let endDate = new Date(Date.UTC(datePieces[1], datePieces[0], 1));
      if (isValid(endDate)) {
        // TODO: Replace with formatRFC3339
        params.end_time = Math.round(endOfMonth(endDate).getTime() / 1000);
      }
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
      activity,
      config,
    });
  },
});
