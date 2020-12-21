export default function() {
  this.namespace = 'v1';

  this.get('sys/internal/counters/activity', function(db) {
    let data = {};
    const firstRecord = db['metrics/activities'].first();
    if (firstRecord) {
      data = firstRecord;
    }
    return {
      data,
      request_id: '0001',
    };
  });

  this.get('sys/internal/counters/config', function(db) {
    return {
      request_id: '00001',
      data: db['metrics/configs'].first(),
    };
  });

  this.get('/sys/internal/vault-config', (db, request) => {
    return {
      data: {
        feature_flags: request.queryParams.empty ? null : ['MANAGED_NAMESPACE'],
      },
    };
  });

  this.passthrough();
}
