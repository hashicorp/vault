import ApplicationSerializer from '../application';
import { format } from 'date-fns';

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
          month: format(new Date(m.month_year), 'M/yy'),
          distinct_entities: m['new_clients']['counts']['entity_clients'],
          non_entity_tokens: m['new_clients']['counts']['non_entity_clients'],
          total: m['new_clients']['counts']['clients'],
        };
      });
    } else {
      return payload.map((m) => {
        return {
          month: format(new Date(m.month_year), 'M/yy'),
          distinct_entities: m['counts']['entity_clients'],
          non_entity_tokens: m['counts']['non_entity_clients'],
          total: m['counts']['clients'],
        };
      });
    }
  },

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    let byNamespace = this.flattenByNamespace(payload.data.by_namespace);
    let monthlyTotals = this.flattenByMonths(payload.data.months);
    let monthlyNew = this.flattenByMonths(payload.data.months, true);

    delete payload.by_namespace;
    let transformedPayload = {
      ...payload,
      flattenedByNamespace: byNamespace,
      monthlyTotalClients: monthlyTotals,
      monthlyNewClients: monthlyNew,
    };
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
