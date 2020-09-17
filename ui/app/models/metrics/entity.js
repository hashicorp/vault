import Model, { attr } from '@ember-data/model';

/* sample response

{
  "request_id": "75cbaa46-e741-3eba-2be2-325b1ba8f03f",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "counters": {
      "entities": {
        "total": 1
      }
    }
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

*/

export default Model.extend({
  entities: attr('object'),
});
