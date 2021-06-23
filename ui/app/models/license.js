import Model, { attr } from '@ember-data/model';

/* sample response
{
  "data": {
    "autoloading_used": true,
    "autoloaded": {
      "expiration_time": "2017-11-14T16:34:36.546753-05:00",
      "license_id": "some-id",
      "start_time": "2017-11-14T16:04:36.546753-05:00"
      "features": [
        "UI",
        "HSM",
        "Performance Replication",
        "DR Replication"
      ],
    },
    "stored": {
      "expiration_time": "2017-11-14T16:34:36.546753-05:00",
      "license_id": "some-id",
      "start_time": "2017-11-14T16:04:36.546753-05:00"
      "features": [
        "UI",
        "HSM",
        "Performance Replication",
        "DR Replication"
      ],
    }
  },
  "warnings": [
    "time left on license is 29m33s"
  ]
}
*/

export default Model.extend({
  expirationTime: attr('string'),
  features: attr('array'),
  licenseId: attr('string'),
  startTime: attr('string'),
  performanceStandbyCount: attr('number'),
  autoloaded: attr('boolean'),
});
