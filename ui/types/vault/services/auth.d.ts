// temporary interface for auth service until it can be updated to ts
// add properties as needed

import Service from '@ember/service';

export interface AuthData {
  entity_id: string;
}

export default class AuthService extends Service {
  authData: AuthData;
}
