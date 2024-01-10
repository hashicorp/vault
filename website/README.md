# Vault Website

This subdirectory contains the content for the [Vault Website](https://vaultproject.io/).

<!--
  This readme file contains several blocks of generated text, to make it easier to share common information
  across documentation website readmes. To generate these blocks from their source, run `npm run generate:readme`

  Any edits to the readme are welcome outside the clearly noted boundaries of the blocks. Alternately, a
  block itself can be safely "forked" from the central implementation simply by removing the "BEGIN" and
  "END" comments from around it.
-->

## Table of Contents

- [Contributions](#contributions-welcome)
- [Running the Site Locally](#running-the-site-locally)
- [Editing Markdown Content](#editing-markdown-content)
- [Editing Navigation Sidebars](#editing-navigation-sidebars)
- [Changing the Release Version](#changing-the-release-version)
- [Redirects](#redirects)
- [Browser Support](#browser-support)
- [Deployment](#deployment)

<!-- BEGIN: contributions -->
<!-- Generated text, do not edit directly -->

## Contributions Welcome!

If you find a typo or you feel like you can improve the HTML, CSS, or JavaScript, we welcome contributions. Feel free to open issues or pull requests like any normal GitHub project, and we'll merge it in 🚀

<!-- END: contributions -->

<!-- BEGIN: local-development -->
<!-- Generated text, do not edit directly -->

## Running the Site Locally

The website can be run locally through node.js or [Docker](https://www.docker.com/get-started). If you choose to run through Docker, everything will be a little bit slower due to the additional overhead, so for frequent contributors it may be worth it to use node.

> **Note:** If you are using a text editor that uses a "safe write" save style such as **vim** or **goland**, this can cause issues with the live reload in development. If you turn off safe write, this should solve the problem. In vim, this can be done by running `:set backupcopy=yes`. In goland, search the settings for "safe write" and turn that setting off.

### With Docker

Running the site locally is simple. Provided you have Docker installed, clone this repo, run `make`, and then visit `http://localhost:3000`.

The docker image is pre-built with all the website dependencies installed, which is what makes it so quick and simple, but also means if you need to change dependencies and test the changes within Docker, you'll need a new image. If this is something you need to do, you can run `make build-image` to generate a local Docker image with updated dependencies, then `make website-local` to use that image and preview.

### With Node

If your local development environment has a supported version (v10.0.0+) of [node installed](https://nodejs.org/en/) you can run:

- `npm install`
- `npm start`

...and then visit `http://localhost:3000`.

If you pull down new code from github, you should run `npm install` again. Otherwise, there's no need to re-run `npm install` each time the site is run, you can just run `npm start` to get it going.

<!-- END: local-development -->

<!-- BEGIN: editing-markdown -->
<!-- Generated text, do not edit directly -->

## Editing Markdown Content

Documentation content is written in [Markdown](https://www.markdownguide.org/cheat-sheet/) and you'll find all files listed under the `/content` directory.

To create a new page with Markdown, create a file ending in `.mdx` in a `content/<subdirectory>`. The path in the content directory will be the URL route. For example, `content/docs/hello.mdx` will be served from the `/docs/hello` URL.

> **Important**: Files and directories will only be rendered and published to the website if they are [included in sidebar data](#editing-navigation-sidebars). Any file not included in sidebar data will not be rendered or published.

This file can be standard Markdown and also supports [YAML frontmatter](https://middlemanapp.com/basics/frontmatter/). YAML frontmatter is optional, there are defaults for all keys.

```yaml
---
title: 'My Title'
description: "A thorough, yet succinct description of the page's contents"
---
```

The significant keys in the YAML frontmatter are:

- `title` `(string)` - This is the title of the page that will be set in the HTML title.
- `description` `(string)` - This is a description of the page that will be set in the HTML description.

> ⚠️ If there is a need for a `/api/*` url on this website, the url will be changed to `/api-docs/*`, as the `api` folder is reserved by next.js.

### Validating Content

Content changes are automatically validated against a set of rules as part of the pull request process. If you want to run these checks locally to validate your content before comitting your changes, you can run the following command:

```
npm run content-check
```

If the validation fails, actionable error messages will be displayed to help you address detected issues.

### Creating New Pages

There is currently a small bug with new page creation - if you create a new page and link it up via subnav data while the server is running, it will report an error saying the page was not found. This can be resolved by restarting the server.

### Markdown Enhancements

There are several custom markdown plugins that are available by default that enhance [standard markdown](https://commonmark.org/) to fit our use cases. This set of plugins introduces a couple instances of custom syntax, and a couple specific pitfalls that are not present by default with markdown, detailed below:

- > **Warning**: We are deprecating the current [paragraph alerts](https://github.com/hashicorp/remark-plugins/tree/master/plugins/paragraph-custom-alerts#paragraph-custom-alerts), in favor of the newer [MDX Inline Alert](#inline-alerts) components. The legacy paragraph alerts are represented by the symbols `~>`, `->`, `=>`, or `!>`.
- If you see `@include '/some/path.mdx'`, this is a [markdown include](https://github.com/hashicorp/remark-plugins/tree/master/plugins/include-markdown#include-markdown-plugin). It's worth noting as well that all includes resolve from `website/content/partials` by default, and that changes to partials will not live-reload the website.
- If you see `# Headline ((#slug))`, this is an example of an [anchor link alias](https://github.com/hashicorp/remark-plugins/tree/je.anchor-link-adjustments/plugins/anchor-links#anchor-link-aliases). It adds an extra permalink to a headline for compatibility and is removed from the output.
- Due to [automatically generated permalinks](https://github.com/hashicorp/remark-plugins/tree/je.anchor-link-adjustments/plugins/anchor-links#anchor-links), any text changes to _headlines_ or _list items that begin with inline code_ can and will break existing permalinks. Be very cautious when changing either of these two text items.

  Headlines are fairly self-explanatory, but here's an example of how to list items that begin with inline code look.

  ```markdown
  - this is a normal list item
  - `this` is a list item that begins with inline code
  ```

  Its worth noting that _only the inline code at the beginning of the list item_ will cause problems if changed. So if you changed the above markup to...

  ```markdown
  - lsdhfhksdjf
  - `this` jsdhfkdsjhkdsfjh
  ```

  ...while it perhaps would not be an improved user experience, no links would break because of it. The best approach is to **avoid changing headlines and inline code at the start of a list item**. If you must change one of these items, make sure to tag someone from the digital marketing development team on your pull request, they will help to ensure as much compatibility as possible.

### Custom Components

A number of custom [mdx components](https://mdxjs.com/) are available for use within any `.mdx` file. Each one is documented below:

#### Inline Alerts

There are custom MDX components available to author alert data. [See the full documentation here](https://developer.hashicorp.com/swingset/components/mdxinlinealert). They render as colored boxes to draw the user's attention to some type of aside.

```mdx
## Alert types

### Tip

<Tip>
  To provide general information to the user regarding the current context or
  relevant actions.
</Tip>

### Highlight

<Highlight>
  To provide general or promotional information to the user prominently.
</Highlight>

### Note

<Note>
  To help users avoid an issue. Provide guidance and actions if possible.
</Note>

### Warning

<Warning>
  To indicate critical issues that need immediate action or help users
  understand something critical.
</Warning>

## Title override prop

<Note title="Hashiconf 2027">To provide general information.</Note>
```

#### Tabs

The `Tabs` component creates tabbed content of any type, but is often used for code examples given in different languages. Here's an example of how it looks from the Vagrant documentation website:

![Tabs Component](https://p176.p0.n0.cdn.getcloudapp.com/items/WnubALZ4/Screen%20Recording%202020-06-11%20at%2006.03%20PM.gif?v=1de81ea720a8cc8ade83ca64fb0b9edd)

> Please refer to the [Swingset](https://react-components.vercel.app/?component=Tabs) documentation for the latest examples and API reference.

It can be used as such within a markdown file:

````mdx
Normal **markdown** content.

<Tabs>
<Tab heading="CLI command">
            <!-- Intentionally skipped line.. -->
```shell-session
$ command ...
```
            <!-- Intentionally skipped line.. -->
</Tab>
<Tab heading="API call using cURL">

```shell-session
$ curl ...
```

</Tab>
</Tabs>

Continued normal markdown content
````

The intentionally skipped line is a limitation of the mdx parser which is being actively worked on. All tabs must have a heading, and there is no limit to the number of tabs, though it is recommended to go for a maximum of three or four.

#### Enterprise Alert

This component provides a standard way to call out functionality as being present only in the enterprise version of the software. It can be presented in two contexts, inline or standalone. Here's an example of standalone usage from the Consul docs website:

![Enterprise Alert Component - Standalone](https://p176.p0.n0.cdn.getcloudapp.com/items/WnubALp8/Screen%20Shot%202020-06-11%20at%206.06.03%20PM.png?v=d1505b90bdcbde6ed664831a885ea5fb)

The standalone component can be used as such in markdown files:

```mdx
# Page Headline

<EnterpriseAlert />

Continued markdown content...
```

It can also receive custom text contents if you need to change the messaging but wish to retain the style. This will replace the text `This feature is available in all versions of Consul Enterprise.` with whatever you add. For example:

```mdx
# Page Headline

<EnterpriseAlert>
  My custom text here, and <a href="#">a link</a>!
</EnterpriseAlert>

Continued markdown content...
```

It's important to note that once you are adding custom content, it must be html and can not be markdown, as demonstrated above with the link.

Now let's look at inline usage, here's an example:

![Enterprise Alert Component - Inline](https://p176.p0.n0.cdn.getcloudapp.com/items/L1upYLEJ/Screen%20Shot%202020-06-11%20at%206.07.50%20PM.png?v=013ba439263de8292befbc851d31dd78)

And here's how it could be used in your markdown document:

```mdx
### Some Enterprise Feature <EnterpriseAlert inline />

Continued markdown content...
```

It's also worth noting that this component will automatically adjust to the correct product colors depending on the context.

#### Other Components

Other custom components can be made available on a per-site basis, the above are the standards. If you have questions about custom components that are not documented here, or have a request for a new custom component, please reach out to @hashicorp/digital-marketing.

### Syntax Highlighting

When using fenced code blocks, the recommendation is to tag the code block with a language so that it can be syntax highlighted. For example:

````
```
// BAD: Code block with no language tag
```

```javascript
// GOOD: Code block with a language tag
```
````

Check out the [supported languages list](https://prismjs.com/#supported-languages) for the syntax highlighter we use if you want to double check the language name.

It is also worth noting specifically that if you are using a code block that is an example of a terminal command, the correct language tag is `shell-session`. For example:

🚫**BAD**: Using `shell`, `sh`, `bash`, or `plaintext` to represent a terminal command

````
```shell
$ terraform apply
```
````

✅**GOOD**: Using `shell-session` to represent a terminal command

````
```shell-session
$ terraform apply
```
````

<!-- END: editing-markdown -->

<!-- BEGIN: editing-docs-sidebars -->
<!-- Generated text, do not edit directly -->

## Editing Navigation Sidebars

The structure of the sidebars are controlled by files in the [`/data` directory](data). For example, [data/docs-nav-data.json](data/docs-nav-data.json) controls the **docs** sidebar. Within the `data` folder, any file with `-nav-data` after it controls the navigation for the given section.

The sidebar uses a simple recursive data structure to represent _files_ and _directories_. The sidebar is meant to reflect the structure of the docs within the filesystem while also allowing custom ordering. Let's look at an example. First, here's our example folder structure:

```text
.
├── docs
│   └── directory
│       ├── index.mdx
│       ├── file.mdx
│       ├── another-file.mdx
│       └── nested-directory
│           ├── index.mdx
│           └── nested-file.mdx
```

Here's how this folder structure could be represented as a sidebar navigation, in this example it would be the file `website/data/docs-nav-data.json`:

```json
[
  {
    "title": "Directory",
    "routes": [
      {
        "title": "Overview",
        "path": "directory"
      },
      {
        "title": "File",
        "path": "directory/file"
      },
      {
        "title": "Another File",
        "path": "directory/another-file"
      },
      {
        "title": "Nested Directory",
        "routes": [
          {
            "title": "Overview",
            "path": "directory/nested-directory"
          },
          {
            "title": "Nested File",
            "path": "directory/nested-directory/nested-file"
          }
        ]
      }
    ]
  }
]
```

A couple more important notes:

- Within this data structure, ordering is flexible, but hierarchy is not. The structure of the sidebar must correspond to the structure of the content directory. So while you could put `file` and `another-file` in any order in the sidebar, or even leave one or both of them out, you could not decide to un-nest the `nested-directory` object without also un-nesting it in the filesystem.
- The `title` property on each node in the `nav-data` tree is the human-readable name in the navigation.
- The `path` property on each leaf node in the `nav-data` tree is the URL path where the `.mdx` document will be rendered, and the
- Note that "index" files must be explicitly added. These will be automatically resolved, so the `path` value should be, as above, `directory` rather than `directory/index`. A common convention is to set the `title` of an "index" node to be `"Overview"`.

Below we will discuss a couple of more unusual but still helpful patterns.

### Index-less Categories

Sometimes you may want to include a category but not have a need for an index page for the category. This can be accomplished, but as with other branch and leaf nodes, a human-readable `title` needs to be set manually. Here's an example of how an index-less category might look:

```text
.
├── docs
│   └── indexless-category
│       └── file.mdx
```

```json
// website/data/docs-nav-data.json
[
  {
    "title": "Indexless Category",
    "routes": [
      {
        "title": "File",
        "path": "indexless-category/file"
      }
    ]
  }
]
```

### Custom or External Links

Sometimes you may have a need to include a link that is not directly to a file within the docs hierarchy. This can also be supported using a different pattern. For example:

```json
[
  {
    "name": "Directory",
    "routes": [
      {
        "title": "File",
        "path": "directory/file"
      },
      {
        "title": "Another File",
        "path": "directory/another-file"
      },
      {
        "title": "Tao of HashiCorp",
        "href": "https://www.hashicorp.com/tao-of-hashicorp"
      }
    ]
  }
]
```

If the link provided in the `href` property is external, it will display a small icon indicating this. If it's internal, it will appear the same way as any other direct file link.

<!-- END: editing-docs-sidebars -->

<!-- BEGIN: releases -->
<!-- Generated text, do not edit directly -->

## Changing the Release Version

To change the version displayed for download on the website, head over to `data/version.js` and change the number there. It's important to note that the version number must match a version that has been released and is live on `releases.hashicorp.com` -- if it does not, the website will be unable to fetch links to the binaries and will not compile. So this version number should be changed _only after a release_.

### Displaying a Prerelease

If there is a prerelease of any type that should be displayed on the downloads page, this can be done by editing `pages/downloads/index.jsx`. By default, the download component might look something like this:

```jsx
<ProductDownloader
  product="<Product>"
  version={VERSION}
  downloads={downloadData}
  community="/resources"
/>
```

To add a prerelease, an extra `prerelease` property can be added to the component as such:

```jsx
<ProductDownloader
  product="<Product>"
  version={VERSION}
  downloads={downloadData}
  community="/resources"
  prerelease={{
    type: 'release candidate', // the type of prerelease: beta, release candidate, etc.
    name: 'v1.0.0', // the name displayed in text on the website
    version: '1.0.0-rc1', // the actual version tag that was pushed to releases.hashicorp.com
  }}
/>
```

This configuration would display something like the following text on the website, emphasis added to the configurable parameters:

```
A {{ release candidate }} for <Product> {{ v1.0.0 }} is available! The release can be <a href='https://releases.hashicorp.com/<product>/{{ 1.0.0-rc1 }}'>downloaded here</a>.
```

You may customize the parameters in any way you'd like. To remove a prerelease from the website, simply delete the `prerelease` parameter from the above component.

<!-- END: releases -->

<!-- BEGIN: redirects -->
<!-- Generated text, do not edit directly -->

## Redirects

This website structures URLs based on the filesystem layout. This means that if a file is moved, removed, or a folder is re-organized, links will break. If a path change is necessary, it can be mitigated using redirects. It's important to note that redirects should only be used to cover for external links -- if you are moving a path which internal links point to, the internal links should also be adjusted to point to the correct page, rather than relying on a redirect.

To add a redirect, head over to the `redirects.js` file - the format is fairly simple - there's a `source` and a `destination` - fill them both in, indicate that it's a permanent redirect or not using the `permanent` key, and that's it. Let's look at an example:

```
{
  source: '/foo',
  destination: '/bar',
  permanent: true
}
```

This redirect rule will send all incoming links to `/foo` to `/bar`. For more details on the redirects file format, [check out the docs on vercel](https://vercel.com/docs/configuration#project/redirects). All redirects will work both locally and in production exactly the same way, so feel free to test and verify your redirects locally. In the past testing redirects has required a preview deployment -- this is no longer the case. Please note however that if you add a redirect while the local server is running, you will need to restart it in order to see the effects of the redirect.

There is still one caveat though: redirects do not apply to client-side navigation. By default, all links in the navigation and docs sidebar will navigate purely on the client side, which makes navigation through the docs significantly faster, especially for those with low-end devices and/or weak internet connections. In the future, we plan to convert all internal links within docs pages to behave this way as well. This means that if there is a link on this website to a given piece of content that has changed locations in some way, we need to also _directly change existing links to the content_. This way, if a user clicks a link that navigates on the client side, or if they hit the url directly and the page renders from the server side, either one will work perfectly.

Let's look at an example. Say you have a page called `/docs/foo` which needs to be moved to `/docs/nested/foo`. Additionally, this is a page that has been around for a while and we know there are links into `/docs/foo.html` left over from our previous website structure. First, we move the page, then adjust the docs sidenav, in `data/docs-navigation.js`. Find the category the page is in, and move it into the appropriate subcategory. Next, we add to `_redirects` as such. The `.html` version is covered automatically.

```js
{ source: '/foo', destination: '/nested/foo', permanent: true }
```

Next, we run a global search for internal links to `/foo`, and make sure to adjust them to be `/nested/foo` - this is to ensure that client-side navigation still works correctly. _Adding a redirect alone is not enough_.

One more example - let's say that content is being moved to an external website. A common example is guides moving to `learn.hashicorp.com`. In this case, we take all the same steps, except that we need to make a different type of change to the `docs-navigation` file. If previously the structure looked like:

```js
{
  category: 'docs',
  content: [
    'foo'
  ]
}
```

If we no longer want the link to be in the side nav, we can simply remove it. If we do still want the link in the side nav, but pointing to an external destination, we need to slightly change the structure as such:

```js
{
  category: 'docs',
  content: [
    { title: 'Foo Title', href: 'https://learn.hashicorp.com/<product>/foo' }
  ]
}
```

As the majority of items in the side nav are internal links, the structure makes it as easy as possible to represent these links. This alternate syntax is the most concise manner than an external link can be represented. External links can be used anywhere within the docs sidenav.

It's also worth noting that it is possible to do glob-based redirects, for example matching `/docs/*`, and you may see this pattern in the redirects file. This type of redirect is much higher risk and the behavior is a bit more nuanced, so if you need to add a glob redirect, please reach out to the website maintainers and ask about it first.

<!-- END: redirects -->

<!-- BEGIN: browser-support -->
<!-- Generated text, do not edit directly -->

## Browser Support

We support the following browsers targeting roughly the versions specified.

| ![Chrome](https://raw.githubusercontent.com/alrra/browser-logos/main/src/chrome/chrome.svg) | ![Edge](https://raw.githubusercontent.com/alrra/browser-logos/main/src/edge/edge.svg) | ![Opera](https://raw.githubusercontent.com/alrra/browser-logos/main/src/opera/opera.svg) | ![Firefox](https://raw.githubusercontent.com/alrra/browser-logos/main/src/firefox/firefox.svg) | ![Safari](https://raw.githubusercontent.com/alrra/browser-logos/main/src/safari/safari.svg) |
| ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| **Latest**                                                                                  | **Latest**                                                                            | **Latest**                                                                               | **Latest**                                                                                     | **Latest**                                                                                  |

<!-- END: browser-support -->

<!-- BEGIN: deployment -->
<!-- Generated text, do not edit directly -->

## Deployment

This website is hosted on Vercel and configured to automatically deploy anytime you push code to the `stable-website` branch. Any time a pull request is submitted that changes files within the `website` folder, a deployment preview will appear in the github checks which can be used to validate the way docs changes will look live. Deployments from `stable-website` will look and behave the same way as deployment previews.

<!-- END: deployment -->
