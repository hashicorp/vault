import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './message.md';

storiesOf('Confirm/Message/', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `Message`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Message</h5>
      <p>
        <code>Message</code> components should never render on their own. See the <code>Confirm</code> component for an example of what a <code>Message</code> looks like.
      </p>
    `,
      context: {},
    }),
    { notes }
  );
