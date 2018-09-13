import { create } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';

const flash = create(flashMessage);

export default async function withFlash(promise, assertion) {
  await flash.waitForFlash();
  if (assertion) {
    assertion();
  }
  await flash.clickAll();
  await promise;
}
