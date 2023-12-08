{{!
  Copyright (c) HashiCorp, Inc.
  SPDX-License-Identifier: BUSL-1.1
~}}

<Messages::TabPageHeader @authenticated={{@authenticated}} @pageTitle="Create message" @breadcrumbs={{this.breadcrumbs}} />

<form id="message-create-edit-form" {{on "submit" (perform this.save)}} data-test-message-create-edit-form>
  <div class="box is-sideless is-fullwidth is-marginless has-top-padding-s">
    <Hds::Text::Body @tag="p" class="has-bottom-margin-l">
      Create a custom message for all users when they access a Vault system via the UI.
    </Hds::Text::Body>

    <MessageError @errorMessage={{this.errorBanner}} class="has-top-margin-s" />

    {{#each @messages.formFields as |attr|}}
      {{#if (eq attr.name "authenticated")}}
        <Messages::RadioSetForm @model={{@messages}} @attr={{attr}} />
      {{else if (eq attr.name "type")}}
        <Messages::RadioSetForm @model={{@messages}} @attr={{attr}} />
      {{else if (eq attr.name "linkTitle")}}
        <div class="field has-bottom-margin-m">
          <label for="link" class="is-label">Link
            <span class="has-text-grey-400 has-font-weight-normal">(optional)</span></label>
          <div class="control is-flex-between has-gap-m">
            <Input
              @type="text"
              @value={{@messages.linkTitle}}
              placeholder="Display text (e.g. Learn more)"
              id="link-title"
              class="input"
              {{on "input" (pipe (pick "target.value") (fn (mut @messages.linkTitle)))}}
              data-test-link="title"
            />
            <Input
              @type="text"
              @value={{@messages.linkHref}}
              placeholder="Paste URL (e.g. www.learnmore.com)"
              id="link-href"
              class="input"
              {{on "input" (pipe (pick "target.value") (fn (mut @messages.linkHref)))}}
              data-test-link="href"
            />
          </div>
        </div>
      {{else if (eq attr.name "startTime")}}
        <label for="message-start" class="has-text-weight-bold is-size-8">
          Message starts
        </label>
        <Hds::Text::Body @tag="p" @size="200" class="has-bottom-margin-s">
          Defaults to 12:00 a.m. the following day (local timezone). When toggled off, message is inactive.
        </Hds::Text::Body>
        <Input
          @type="datetime-local"
          class="input has-top-margin-xs is-auto-width is-block"
          data-test-date="startTime"
          name="startTime"
          {{on "input" this.updateDateTime}}
        />
      {{else if (eq attr.name "endTime")}}
        <Hds::Form::Radio::Group @layout="vertical" @name="endTime" class="has-top-margin-l" as |G|>
          <G.Legend>Message expires</G.Legend>
          <G.Radio::Field @id="never" checked @value="" {{on "change" (fn (mut @messages.endTime) "")}} as |F|>
            <F.Label>Never</F.Label>
            <F.HelperText>This message will never expire unless manually deleted by the operator.</F.HelperText>
          </G.Radio::Field>
          <G.Radio::Field
            @id="specific-date"
            @value=""
            {{on "change" (fn (mut @messages.endTime) @messages.endTime)}}
            as |F|
          >
            <F.Label>Specific date</F.Label>
            <F.HelperText>
              This message will expire at midnight (local timezone) at the specific date.
              <div class="control">
                <Input
                  @type="datetime-local"
                  class="input has-top-margin-xs is-auto-width"
                  data-test-link="title"
                  name="endTime"
                  {{on "input" this.updateDateTime}}
                />
              </div>
            </F.HelperText>
          </G.Radio::Field>
        </Hds::Form::Radio::Group>
      {{else}}
        <FormField
          class="has-bottom-margin-m"
          data-test-field={{true}}
          @attr={{attr}}
          @model={{@messages}}
          @modelValidations={{this.modelValidations}}
        />
      {{/if}}
    {{/each}}

    <Hds::ButtonSet class="has-top-margin-s has-bottom-margin-m has-top-margin-xl">
      {{! TODO: VAULT-21533 preview modal }}
      <Hds::Button @text="Preview" @color="tertiary" @icon="eye" />

      <Hds::Button @text="Create message" type="submit" />

      <Hds::Button @text="Cancel" @color="secondary" @route="messages.index" @query={{hash authenticated=true}} />
    </Hds::ButtonSet>
  </div>
</form>