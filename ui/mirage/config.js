export default function() {
  this.namespace = 'v1';

  this.get('sys/internal/counters/activity', function(db) {
    console.log('getting sys/internal/counters/activity');
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
    console.log('getting sys/internal/counters/config');
    return {
      request_id: '00001',
      data: db['metrics/configs'].first(),
    };
  });

  this.get('/sys/internal/ui/feature-flags', db => {
    const featuresResponse = db.features.first();
    return {
      data: {
        feature_flags: featuresResponse ? featuresResponse.feature_flags : null,
      },
    };
  });

  this.passthrough();
}
