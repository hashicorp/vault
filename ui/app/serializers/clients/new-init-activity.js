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
    // keys in the object created here must match the legend keys in dashboard.js ('entity_clients')
    let topTen = payload.slice(0, 10);
    return topTen.map((ns) => {
      if (ns['namespace_path'] === '') ns['namespace_path'] = 'root';
      // this may need to change when we have real data
      // right now under months, namespaces have key value of "path" or "id", not "namespace_path"
      let label = ns['namespace_path'] || ns['id'];
      let namespaceMounts = ns.mounts.map((m) => {
        return {
          label: m['path'],
          entity_clients: m['counts']['entity_clients'],
          non_entity_clients: m['counts']['non_entity_clients'],
          total: m['counts']['clients'],
        };
      });
      return {
        label,
        entity_clients: ns['counts']['entity_clients'],
        non_entity_clients: ns['counts']['non_entity_clients'],
        total: ns['counts']['clients'],
        mounts: namespaceMounts,
      };
    });
  },

  // for vault usage - vertical bar chart
  flattenByMonths(payload, isNewClients = false) {
    if (isNewClients) {
      return payload.map((m) => {
        return {
          month: format(new Date(m.timestamp), 'M/yy'),
          entity_clients: m['new_clients']['counts']['entity_clients'],
          non_entity_clients: m['new_clients']['counts']['non_entity_clients'],
          total: m['new_clients']['counts']['clients'],
          namespaces: this.flattenByNamespace(m['new_clients']['namespaces']),
        };
      });
    } else {
      return payload.map((m) => {
        return {
          month: format(new Date(m.timestamp), 'M/yy'),
          entity_clients: m['counts']['entity_clients'],
          non_entity_clients: m['counts']['non_entity_clients'],
          total: m['counts']['clients'],
          namespaces: this.flattenByNamespace(m['namespaces']),
          new_clients: {
            entity_clients: m['new_clients']['counts']['entity_clients'],
            non_entity_clients: m['new_clients']['counts']['non_entity_clients'],
            total: m['new_clients']['counts']['clients'],
            namespaces: this.flattenByNamespace(m['new_clients']['namespaces']),
          },
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
      by_namespace: this.flattenByNamespace(payload.data.by_namespace),
      by_month: this.flattenByMonths(payload.data.months),
      by_month_new_clients: this.flattenByMonths(payload.data.months, { isNewClients: true }),
    };

    delete payload.data.by_namespace;
    delete payload.data.months;
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
