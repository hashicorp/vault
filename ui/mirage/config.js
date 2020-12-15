import Response from 'ember-cli-mirage/response';

export default function() {
  this.get(
    'http://localhost:4200/vault-config',
    function() {
      console.log('getting vault config');
      return new Response(
        201,
        {
          'Content-Type': 'application/json',
        },
        {
          managedNamespaceRoot: 'admin',
        }
      );
    },
    { timing: 3000 }
  );

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

  this.passthrough();
}
