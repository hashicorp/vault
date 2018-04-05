import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/section';

moduleForAcceptance('Acceptance | settings/auth/configure/section', {
  beforeEach() {
    return authLogin();
  },
});

test('it can save options', function(assert) {
  const path = `approle-${new Date().getTime()}`;
  const type = 'approle';
  const section = 'options';
  enablePage.visit().enableAuth(type, path);
  page.visit({ path, section });
  andThen(() => {
    page.fields().findByName('description').textarea('This is AppRole!');
    page.save();
  });
  andThen(() => {
    assert.equal(
      page.flash.latestMessage,
      `The configuration options were saved successfully.`,
      'success flash shows'
    );
  });
});
