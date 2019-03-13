/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';

storiesOf('AuthForm/', module)
  .addParameters({ options: { showPanel: false } })
  .add(`AuthForm`, () => ({
    template: hbs`
        <h5 class="title is-5">Auth Form</h5>
        <AuthForm/>
        <h5 class="title is-5">Auth Form custom</h5>
        <AuthForm
          @wrappedToken={{wrappedToken}}
          @cluster={{model}}
          @namespace={{namespaceQueryParam}}
          @redirectTo={{redirectTo}}
          @selectedAuth={{authMethod}}
        />
    `,
    context: {},
  }));
