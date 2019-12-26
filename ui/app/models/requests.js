import DS from 'ember-data';
const { attr } = DS;

/* sample response

{
  "request_id": "75cbaa46-e741-3eba-2be2-325b1ba8f03f",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "counters": [
      {
        "start_time": "2019-05-01T00:00:00Z",
        "total": 50
      },
      {
        "start_time": "2019-04-01T00:00:00Z",
        "total": 45
      }
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

*/

export default DS.Model.extend({
  counters: attr('array'),
});
