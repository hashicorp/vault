import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text, boolean, select } from '@storybook/addon-knobs';
import notes from './box-radio.md';

const GLYPHS = {
  KMIP: 'kmip',
  Transform: 'transform',
  AWS: 'aws',
  Azure: 'azure',
  Cert: 'cert',
  GCP: 'gcp',
  Github: 'github',
  JWT: 'jwt',
  HashiCorp: 'hashicorp',
  LDAP: 'ldap',
  OKTA: 'okta',
  Radius: 'radius',
  Userpass: 'userpass',
  Secrets: 'kv',
  Consul: 'consul',
};

storiesOf('BoxRadio', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `BoxRadio`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Box Radio</h5>
      <BoxRadio
        @type={{type}}
        @glyph={{type}}
        @displayName={{displayName}}
        @onRadioChange={{onRadioChange}}
        @disabled={{disabled}}
      />
    `,
      context: {
        displayName: text('displayName', 'HashiCorp'),
        type: select('glyph', GLYPHS, 'hashicorp'),
        disabled: boolean('disabled', false),
        onRadioChange: e => {
          console.log('Radio changed!', e);
        },
      },
    }),
    { notes }
  );
