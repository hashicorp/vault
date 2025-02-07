import Service from '@ember/service';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import config from '../config/environment';

const { host } = config;

export default class HostService extends Service {

  @tracked _host = '';

  constructor() {
    super();
    this._host = host
  }

  get host() {
    return this._host;
  }

  set host(value: string) {
    this._host = value;
  }
}
