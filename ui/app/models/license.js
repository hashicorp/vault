import DS from 'ember-data';
const { attr } = DS;

/* sample response
{
  "data": {
    "expiration_time": "2017-11-14T16:34:36.546753-05:00",
    "features": [
      "UI",
      "HSM",
      "Performance Replication",
      "DR Replication"
    ],
    "license_id": "temporary",
    "start_time": "2017-11-14T16:04:36.546753-05:00"
  },
  "warnings": [
    "time left on license is 29m33s"
  ]
}
*/

export default DS.Model.extend({
  expirationTime: attr('string'),
  features: attr('array'),
  licenseId: attr('string'),
  startTime: attr('string'),
  text: attr('string'),
});
