import { belongsTo, attr } from '@ember-data/model';
import timestamp from 'core/utils/timestamp';
import SecretModel from './secret';

export default class SecretV2VersionModel extends SecretModel {
  @attr('boolean') failedServerRead;
  @attr('number') version;
  @attr('string') path;
  @attr('string') deletionTime;
  @attr('string') createdTime;
  @attr('boolean') destroyed;
  @attr('number') currentVersion;
  @belongsTo('secret-v2') secret;

  pathAttr = 'path';

  get deleted() {
    const deletionTime = new Date(this.deletionTime);
    const now = timestamp.now();
    return deletionTime <= now;
  }
}
