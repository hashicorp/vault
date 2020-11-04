import activity from './factories/metrics/activity';

export default function() {
  this.namespace = 'v1';
  // this.get('/api/users', function(db, request) {
  //   return {
  //     users: db.users,
  //   };
  // });
  this.get('sys/internal/counters/activity', function(db, request) {
    let data = {};
    console.log('db metric activity:', db['metrics/activities'].first());
    if (db['metrics/activities'].length > 0) {
      data = db['metrics/activities'].first();
    }
    return {
      data,
      request_id: data.request_id,
    };
  });
  this.get('sys/internal/counters/config', function(db) {
    console.log('db config');
    try {
      console.log(db['metrics/configs'].all());
    } catch (e) {
      console.log(e);
    }
    return {
      request_id: '00001',
      data: db['metrics/configs'].first(),
      // {
      //   defaultReportMonths: 6,
      //   enabled: 'default-enabled',
      //   queriesAvailable: false,
      //   retentionMonths: 12,
      // },
    };
  });
  this.passthrough();
}

/*
auth: null
​​
by_namespace: Array [ {…}, {…} ]
​​
end_time: "2020-09-30T23:59:59Z"
​​
id: "81035ff0-3217-2222-02b1-1521f98856c1"
​​
lease_duration: 0
​​
lease_id: ""
​​
renewable: false
​​
request_id: "81035ff0-3217-2222-02b1-1521f98856c1"
​​
start_time: "2020-03-01T00:00:00Z"
​​
total: Object { distinct_entities: Getter & Setter, non_entity_tokens: Getter & Setter, clients: Getter & Setter }
​​
warnings: null
​​
wrap_info: null
*/
