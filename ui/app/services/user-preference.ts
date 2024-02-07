import { action } from '@ember/object';
import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';

enum KvSecretDisplay {
  Unset = 'unset',
  Json = 'json',
  KeyValue = 'keyvalue',
}

export default class UserPreferenceService extends Service {
  @tracked kvDisplaySetting = KvSecretDisplay.Unset;
  calculateInitialKvJson = (isAdvanced: boolean) => {
    if (this.kvDisplaySetting === KvSecretDisplay.Unset) {
      // if user preference is unset, show json if advanced
      return isAdvanced;
    }
    // otherwise use user preference
    return this.kvDisplaySetting === KvSecretDisplay.Json;
  };
  @action setKvDisplayPreference(jsonToggle: boolean) {
    this.kvDisplaySetting = jsonToggle ? KvSecretDisplay.Json : KvSecretDisplay.KeyValue;
  }

  @action
  reset() {
    this.kvDisplaySetting = KvSecretDisplay.Unset;
  }
}
