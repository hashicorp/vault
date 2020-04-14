import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './replication-primary-card.md';

storiesOf('Replication/ReplicationPrimaryCard/', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `ReplicationPrimaryCard`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Replication Primary Card</h5>
      <ReplicationSecondaries/>
    `,
      context: {},
    }),
    { notes }
  );
