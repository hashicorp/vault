import ESASessionService from 'ember-simple-auth/services/session';

export default class VaultSessionService extends ESASessionService {
  handleAuthentication() {
    super.handleAuthentication('vault.cluster.dashboard');
  }
}
