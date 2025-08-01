import sinon from 'sinon';

export function setupAuth(hooks) {
  hooks.beforeEach(function () {
    const auth = this.owner.lookup('service:auth');
    sinon.stub(auth, 'currentToken').value('test-token');
  });
}
