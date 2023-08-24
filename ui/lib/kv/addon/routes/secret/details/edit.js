import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';

@withConfirmLeave('model.newVersion')
export default class KvSecretDetailsEditRoute extends Route {
  @service store;

  model() {
    const parentModel = this.modelFor('secret.details');
    const { backend, path, secret, metadata } = parentModel;
    return hash({
      secret,
      metadata,
      backend,
      path,
      newVersion: this.store.createRecord('kv/data', {
        backend,
        path,
        secretData: secret?.secretData,
        // see serializer for logic behind setting casVersion
        casVersion: metadata?.currentVersion || secret?.version,
      }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
      ...breadcrumbsForSecret(resolvedModel.path),
      { label: 'edit' },
    ];
  }
}
