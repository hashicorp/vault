import ApplicationSerializer from '../application';
import { format } from 'date-fns';

// TODO CMB: add before and after example of serializer

export default ApplicationSerializer.extend({
  // abstracting like this may not work as the order does matter in how the charts display
  // whatever is first will be the bottom/left bar, second will be top/right bar
  flattenDataset(dataset, nestedObjectKey = 'counts') {
    // nestedObjectKey needs to be passed in as a string, defaults to counts
    // because most of the data from the API is nested under the key 'counts'
    return dataset.map((d) => {
      let flattenedObject = {};
      Object.keys(d[nestedObjectKey]).forEach((k) => {
        flattenedObject[k] = d[nestedObjectKey][k];
      });
      return {
        label: d['namespace_path'] === '' ? 'root' : d['namespace_path'],
        ...flattenedObject,
      };
    });
  },

  // used for top 10 attribution charts
  flattenByNamespace(payload) {
    // keys in the object created here must match the legend keys in dashboard.js ('distinct_entities')
    let topTen = payload.slice(0, 10);
    return topTen.map((ns) => {
      let namespaceMounts = ns.mounts.map((m) => {
        // debugger
        return {
          label: m['path'],
          distinct_entities: m['counts']['entity_clients'],
          non_entity_tokens: m['counts']['non_entity_clients'],
          total: m['counts']['clients'],
        };
      });
      return {
        label: ns['namespace_path'] === '' ? 'root' : ns['namespace_path'],
        distinct_entities: ns['counts']['distinct_entities'],
        non_entity_tokens: ns['counts']['non_entity_tokens'],
        total: ns['counts']['clients'],
        mounts: namespaceMounts,
      };
    });
  },

  // for vault usage - vertical bar chart
  flattenByMonths(payload, isNew = false) {
    if (isNew) {
      return payload.map((m) => {
        return {
          month: format(new Date(m.timestamp), 'M/yy'), // format month as '1/22'
          distinct_entities: m['new_clients']['counts']['entity_clients'],
          non_entity_tokens: m['new_clients']['counts']['non_entity_clients'],
          total: m['new_clients']['counts']['clients'],
        };
      });
    } else {
      return payload.map((m) => {
        return {
          month: format(new Date(m.timestamp), 'M/yy'),
          distinct_entities: m['counts']['entity_clients'],
          non_entity_tokens: m['counts']['non_entity_clients'],
          total: m['counts']['clients'],
        };
      });
    }
  },

  formatTimestamp(payload) {
    return payload.map((m) => {
      let month = format(new Date(m.timestamp), 'M/yy');
      return {
        month,
        ...m,
      };
    });
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // conditionals to see if .months, if new_totals exist?
    // needs to accept both /monthly and /activity payloads
    // replace & delete old keys on payload
    let transformedPayload = {
      ...payload,
      // TODO CMB should these be nested under "data"?
      months: this.formatTimestamp(payload.data.months),
      by_namespace: this.flattenByNamespace(payload.data.by_namespace),
      by_month_total_clients: this.flattenByMonths(payload.data.months),
      by_month_new_clients: this.flattenByMonths(payload.data.months, { isNew: true }),
    };
    delete payload.data.by_namespace;
    delete payload.data.months;
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
