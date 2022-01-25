import ApplicationSerializer from '../application';

/* 
SAMPLE PAYLOAD BEFORE/AFTER:

payload.data.by_namespace = [
  {
    namespace_id: '5SWT8',
    namespace_path: 'namespacelonglonglong4/',
    _comment1: 'client counts are nested within own object', 
    counts: {
      entity_clients: 171,
      non_entity_clients: 20,
      clients: 191,
    },
    mounts: [
      {
        path: 'auth/method/uMGBU',
        counts: {
          clients: 35,
          entity_clients: 20,
          non_entity_clients: 15,
        },
      },
    ],
  },
];

transformedPayload.by_namespace = [
  {
    label: 'namespacelonglonglong4/',
    _comment2: 'remove nested object', 
    entity_clients: 171,
    non_entity_clients: 20,
    clients: 191,
    mounts: [
      {
        label: 'auth/method/uMGBU',
        entity_clients: 20,
        non_entity_clients: 15,
        clients: 35,
      },
    ],
  },
]
*/

export default ApplicationSerializer.extend({
  flattenDataset(payload) {
    let topTen = payload.slice(0, 10);

    return topTen.map((ns) => {
      // 'namespace_path' is an empty string for root
      if (ns['namespace_id'] === 'root') ns['namespace_path'] = 'root';
      let label = ns['namespace_path'] || ns['namespace_id']; // TODO CMB will namespace_path ever be empty?
      let flattenedNs = {};
      // we don't want client counts nested within the 'counts' object for stacked charts
      Object.keys(ns['counts']).forEach((key) => (flattenedNs[key] = ns['counts'][key]));

      // if mounts attribution unavailable, mounts will be undefined
      flattenedNs.mounts = ns.mounts?.map((mount) => {
        let flattenedMount = {};
        flattenedMount.label = mount['path'];
        Object.keys(mount['counts']).forEach((key) => (flattenedMount[key] = mount['counts'][key]));
        return flattenedMount;
      });
      return {
        label,
        ...flattenedNs,
      };
    });
  },

  // TODO CMB remove and use abstracted function above
  // prior to 1.10, client count key names are "distinct_entities" and "non_entity_tokens" so mapping below wouldn't work
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

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // needs to accept both /monthly and /activity payloads
    let transformedPayload = {
      ...payload,
      // TODO CMB should these be nested under "data" to go to model correctly?)
      by_namespace: this.flattenDataset(payload.data.by_namespace),
    };
    delete payload.data.by_namespace;
    return this._super(store, primaryModelClass, transformedPayload, id, requestType);
  },
});
