import { match } from '@ember/object/computed';
import DS from 'ember-data';
import KeyMixin from 'vault/mixins/key-mixin';
const { attr } = DS;

/* sample response
{
  "id": "auth/token/create/25c75065466dfc5f920525feafe47502c4c9915c",
  "issue_time": "2017-04-30T10:18:11.228946471-04:00",
  "expire_time": "2017-04-30T11:18:11.228946708-04:00",
  "last_renewal": null,
  "renewable": true,
  "ttl": 3558
}

*/

export default DS.Model.extend(KeyMixin, {
  issueTime: attr('string'),
  expireTime: attr('string'),
  lastRenewal: attr('string'),
  renewable: attr('boolean'),
  ttl: attr('number'),
  isAuthLease: match('id', /^auth/),
});
