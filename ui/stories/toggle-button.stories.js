/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';

storiesOf('ToggleButton', module)
  .add('showOptions', () => ({
    template: hbs`
    <ToggleButton
      @toggleAttr="showOptions"
      @toggleTarget={{this}}
      />
    `,
  }))
  .add('with label', () => ({
    template: hbs`
    <ToggleButton
      @openLabel="Expand Options"
      @closedLabel="Hide Options"
      @toggleTarget={{this}}
      @toggleAttr="use_pgp"
      />
    `,
  }));
