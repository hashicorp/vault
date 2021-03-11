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
      <BoxRadioSet @title="Single Box Radio (should always be wrapped in BoxRadioSet)">
        <BoxRadio
          @type={{key}}
          @glyph={{key}}
          @displayName={{displayName}}
          @onRadioChange={{onRadioChange}}
          @disabled={{disabled}}
          @groupName="example-set-1"
        />
      </BoxRadioSet>
      <BoxRadioSet @title="Multiple Box Radios (should always be wrapped in BoxRadioSet)">
        <BoxRadio
          @key={{concat key "-1"}}
          @glyph={{key}}
          @displayName={{displayName}}
          @onRadioChange={{onRadioChange}}
          @disabled={{disabled}}
          @groupName="example-set"
        />
        <BoxRadio
          @key={{concat key "-2"}}
          @glyph={{key}}
          @displayName={{displayName}}
          @onRadioChange={{onRadioChange}}
          @groupName="example-set"
        />
        <BoxRadio
          @key={{concat key "-3"}}
          @glyph={{key}}
          @displayName={{displayName}}
          @onRadioChange={{onRadioChange}}
          @groupName="example-set"
        />
      </BoxRadioSet>
    `,
      context: {
        displayName: text('displayName', 'HashiCorp'),
        key: select('glyph', GLYPHS, 'hashicorp'),
        disabled: boolean('disabled', false),
        onRadioChange: e => {
          console.log('Radio changed! Should set the value and pass back into component as @groupValue', e);
        },
      },
    }),
    { notes }
  );
