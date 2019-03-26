/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './home-link.md';

storiesOf('HomeLink/', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    'HomeLink',
    () => ({
      template: hbs`
        <h5 class="title is-5">HomeLink</h5>
        <HomeLink />
        <br />
        <h5 class="title is-5">HomeLink with LogoEdition</h5>
        <HomeLink>
          <LogoEdition />
        </HomeLink>
    `,
    }),
    { notes }
  );
