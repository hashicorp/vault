import { RestSerializer } from 'ember-cli-mirage';

export default RestSerializer.extend({
  embed: true,
  root: false,
});
