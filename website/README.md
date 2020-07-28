# Vault Website

[![Netlify
Status](https://img.shields.io/netlify/d72a27d2-aba4-46fd-bf70-d7da0b19b578?style=flat-square)](https://app.netlify.com/sites/vault-www/deploys)

This subdirectory contains the entire source for the [Vault
Website](https://vaultproject.io/). This is a [NextJS](https://nextjs.org/)
project, which builds a static site from these source files.

## Contributions Welcome!

If you find a typo or you feel like you can improve the HTML, CSS, or JavaScript,
we welcome contributions. Feel free to open issues or pull requests like any
normal GitHub project, and we'll merge it in 🚀

## Running the Site Locally

The website can be run locally through node.js or Docker. If you choose to run
through Docker, everything will be a little bit slower due to the additional
overhead, so for frequent contributors it may be worth it to use node. Also if
you are a vim user, it's also worth noting that vim's swapfile usage can cause
issues for the live reload functionality. In order to avoid these issues, make
sure you have run `:set backupcopy=yes` within vim.

### With Docker

Running the site locally is simple. Provided you have Docker installed, clone
this repo, run `make`, and then visit `http://localhost:3000`.

The docker image is pre-built with all the website dependencies installed, which
is what makes it so quick and simple, but also means if you need to change
dependencies and test the changes within Docker, you'll need a new image. If this
is something you need to do, you can run `make build-image` to generate a local
Docker image with updated dependencies, then `make website-local` to use that
image and preview.

### With Node

If your local development environment has a supported version (v10.0.0+) of [node
installed](https://nodejs.org/en/) you can run:

- `npm install`
- `npm start`

and then visit `http://localhost:3000`.

If you pull down new code from github, you should run `npm install` again.
Otherwise, there's no need to re-run `npm install` each time the site is run, you
can just run `npm start` to get it going.

## Editing Content

Documentation content is written in
[Markdown](https://www.markdownguide.org/cheat-sheet/) and you'll find all files
listed under the `/pages` directory.

To create a new page with Markdown, create a file ending in `.mdx` in the
`pages/` directory. The path in the pages directory will be the URL route. For
example, `pages/hello/world.mdx` will be served from the `/hello/world` URL.

This file can be standard Markdown and also supports [YAML
frontmatter](https://middlemanapp.com/basics/frontmatter/). YAML frontmatter is
optional, there are defaults for all keys.

```yaml
---
title: 'My Title'
description: "A thorough, yet succinct description of the page's contents"
---

```

The significant keys in the YAML frontmatter are:

- `title` `(string)` - This is the title of the page that will be set in the HTML
  title.
- `description` `(string)` - This is a description of the page that will be set
  in the HTML description.

> ⚠️Since `api` is a reserved directory within NextJS, all `/api/**` pages are
> listed under the `/pages/api-docs` path.

### Code Highlighting

Code is highlighted using [prism](https://prismjs.com/). Feel free to check out [all the supported languages](https://prismjs.com/#supported-languages) that can be used for code blocks. All code blocks should be tagged with a language as such:

````md
```language
// code to be highlighted
```
````

If you have a code block that displays a command intended to be run from the terminal, it can be tagged with `shell-session`. This is distinct from `shell` which should represent a shell script. The following example shows a correctly formatted terminal command snippet:

````md
```shell-session
$ cowsay "hello world"
```
````

### Editing Sidebars

The structure of the sidebars are controlled by files in the [`/data`
directory](data).

- Edit [this file](data/docs-navigation.js) to change the **docs** sidebar
- Edit [this file](data/api-navigation.js) to change the **api docs** sidebar

To nest sidebar items, you'll want to add a new `category` key/value accompanied
by the appropriate embedded `content` values.

- `category` values will be **directory names** within the `pages` directory
- `content` values will be **file names** within their appropriately nested
  directory.

### Changing the Release Version

To change the version of Vault displayed for download on the website, head over
to `data/version.js` and change the number there. It's important to note that the
version number must match a version that has been released and is live on
`releases.hashicorp.com` -- if it does not, the website will be unable to fetch
links to the binaries and will not compile. So this version number should be
changed _only after a release_.

The `data/version.js` also contains a global variable, `CHANGELOG_URL`, that
should be updated to point to the latest changelog URL for the particular
release version. The URL should be based off the `master` blob such that
it always reflects the most up-to-date changes.

### Displaying a Prerelease

If there is a prerelease of any type that should be displayed on the downloads
page, this can be done by editing `pages/downloads/index.jsx`. By default, the
download component might look something like this:

```jsx
<ProductDownloader product="Vault" version={VERSION} downloads={downloadData} />
```

To add a prerelease, an extra `prerelease` property can be added to the component
as such:

```jsx
<ProductDownloader
  product="Vault"
  version={VERSION}
  downloads={downloadData}
  prerelease={{
    type: 'release candidate', // the type of prerelease: beta, release candidate, etc.
    name: 'v1.0.0', // the name displayed in text on the website
    version: '1.0.0-rc1', // the actual version tag that was pushed to releases.hashicorp.com
  }}
/>
```

This configuration would display something like the following text on the
website, emphasis added to the configurable parameters:

```
A {{ release candidate }} for Vault {{ v1.0.0 }} is available! The release can be <a href='https://releases.hashicorp.com/vault/{{ 1.0.0-rc1 }}'>downloaded here</a>.
```

You may customize the parameters in any way you'd like. To remove a prerelease
from the website, simply delete the `prerelease` paremeter from the above
component.

### Markdown Enhancements

There are several custom markdown plugins that are available by default that
enhance standard markdown to fit our use cases. This set of plugins introduces a
couple instances of custom syntax, and a couple specific pitfalls that are not
present by default with markdown, detailed below:

- If you see the symbols `~>`, `->`, `=>`, or `!>`, these represent [custom
  alerts](https://github.com/hashicorp/remark-plugins/tree/master/plugins/paragraph-custom-alerts#paragraph-custom-alerts).
  These render as colored boxes to draw the user's attention to some type of
  aside.
- If you see `@include '/some/path.mdx'`, this is a [markdown
  include](https://github.com/hashicorp/remark-plugins/tree/master/plugins/include-markdown#include-markdown-plugin).
  It's worth noting as well that all includes resolve from
  `website/pages/partials` by default.

  > **Note:** Changes to partials will not trigger a hot reload in development

- If you see `# Headline ((#slug))`, this is an example of an [anchor link
  alias](https://github.com/hashicorp/remark-plugins/tree/je.anchor-link-adjustments/plugins/anchor-links#anchor-link-aliases).
  It adds an extra permalink to a headline for compatibility and is removed from
  the output.
- Due to [automatically generated
  permalinks](https://github.com/hashicorp/remark-plugins/tree/je.anchor-link-adjustments/plugins/anchor-links#anchor-links),
  any text changes to _headlines_ or _list items that begin with inline code_ can
  and will break existing permalinks. Be very cautious when changing either of
  these two text items.

  Headlines are fairly self-explanitory, but here's an example of how list items
  that begin with inline code look.

  ```markdown
  - this is a normal list item
  - `this` is a list item that begins with inline code
  ```

  Its worth noting that _only the inline code at the beginning of the list item_
  will cause problems if changed. So if you changed the above markup to...

  ```markdown
  - lsdhfhksdjf
  - `this` jsdhfkdsjhkdsfjh
  ```

  ...while it perhaps would not be an improved user experience, no links would
  break because of it. The best approach is to **avoid changing headlines and
  inline code at the start of a list item**. If you must change one of these
  items, make sure to tag someone from the digital marketing development team on
  your pull request, they will help to ensure as much compatibility as possible.

There are also some custom components available for use within markdown files, see
the links below for more information on usage:

- [Enterprise Alert](components/enterprise-alert/README.md)
- [Tabs](components/tabs/README.md)

### Redirects

This website structures URLs based on the filesystem layout. This means that if a
file is moved, removed, or a folder is re-organized, links will break. If a path
change is necessary, it can be mitigated using redirects.

To add a redirect, head over to the `_redirects` file - the format is fairly
simple. On the left is the current path, and on the right is the path that should
be redirected to. It's important to note that if there are links to a `.html`
version of a page, that must also be explicitly redirected. For example:

```
/foo       /bar   301!
/foo.html  /bar   301!
```

This redirect rule will send all incoming links to `/foo` and `/foo.html` to
`/bar`. For more details on the redirects file format, [check out the docs on
netlify](https://docs.netlify.com/routing/redirects/rewrites-proxies). Note that
it is critical that `301!` is added to every one-to-one redirect - if it is left
off the redirect may not work.

There are a couple important caveats with redirects. First, redirects are applied
at the hosting layer, and therefore will not work by default in local dev mode.
To test in local dev mode, you can use [`netlify dev`](https://www.netlify.com/products/dev/), or just push a commit and check
using the deploy preview.

Second, redirects do not apply to client-side navigation. By default, all links
in the navigation and docs sidebar will navigate purely on the client side, which
makes navigation through the docs significantly faster, especially for those with
low-end devices and/or weak internet connections. In the future, we plan to
convert all internal links within docs pages to behave this way as well. This
means that if there is a link on this website to a given piece of content that
has changed locations in some way, we need to also _directly change existing
links to the content_. This way, if a user clicks a link that navigates on the
client side, or if they hit the url directly and the page renders from the server
side, either one will work perfectly.

Let's look at an example. Say you have a page called `/docs/foo` which needs to
be moved to `/docs/nested/foo`. Additionally, this is a page that has been around
for a while and we know there are links into `/docs/foo.html` left over from our
previous website structure. First, we move the page, then adjust the docs
sidenav, in `data/docs-navigation.js`. Find the category the page is in, and move
it into the appropriate subcategory. Next, we add to `_redirects` as such:

```
/foo       /nested/foo  301!
/foo.html  /nested/foo  301!
```

Finally, we run a global search for internal links to `/foo`, and make sure to
adjust them to be `/nested/foo` - this is to ensure that client-side navigation
still works correctly. _Adding a redirect alone is not enough_.

One more example - let's say that content is being moved to an external website.
A common example is guides moving to `learn.hashicorp.com`. In this case, we take
all the same steps, except that we need to make a different type of change to the
`docs-navigation` file. If previously the structure looked like:

```js
{
  category: 'docs',
  content: [
    'foo'
  ]
}
```

If we no longer want the link to be in the side nav, we can simply remove it. If
we do still want the link in the side nav, but pointing to an external
destnation, we need to slightly change the structure as such:

```js
{
  category: 'docs',
  content: [
    { title: 'Foo Title', href: 'https://learn.hashicorp.com/vault/foo' }
  ]
}
```

As the majority of items in the side nav are internal links, the structure makes
it as easy as possible to represent these links. This alternate syntax is the
most concise manner than an external link can be represented. External links can
be used anywhere within the docs sidenav.

It's also worth noting that it is possible to do glob-based redirects, for
example matching `/docs/*`, and you may see this pattern in the `_redirects`
file. This type of redirect is much higher risk and the behavior is a bit more
nuanced, so if you need to add a glob redirect, please reach out to the website
maintainers and ask about it first.

## Deployment

This website is hosted on Netlify and configured to automatically deploy anytime
you push code to the `stable-website` branch. Any time a pull request is
submitted that changes files within the `website` folder, a deployment preview
will appear in the github checks which can be used to validate the way docs
changes will look live. Deployments from `stable-website` will look and behave
the same way as deployment previews.

## Checking for Broken Links

There is a local script that can be used to check for broken links on the
_current product website_ - you can start it by running `npm run linkcheck`.
There will be a version of this script added as a github check in the near
future!

## Browser Support

We support the following browsers targeting roughly the versions specified.

| ![Chrome](https://raw.githubusercontent.com/alrra/browser-logos/master/src/chrome/chrome_24x24.png) | ![Firefox](https://raw.githubusercontent.com/alrra/browser-logos/master/src/firefox/firefox_24x24.png) | ![Opera](https://raw.githubusercontent.com/alrra/browser-logos/master/src/opera/opera_24x24.png) | ![Safari](https://raw.githubusercontent.com/alrra/browser-logos/master/src/safari/safari_24x24.png) | ![Internet Explorer](https://raw.githubusercontent.com/alrra/browser-logos/master/src/edge/edge_24x24.png) |
| --------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| **Latest**                                                                                          | **Latest**                                                                                             | **Latest**                                                                                       | **Latest**                                                                                          | **11+**                                                                                                    |

## Known Issues

### Creating New Pages

There is currently a small bug with new page creation - if you create a new page
and link it up via subnav data while the server is running, it will report an
error saying the page was not found. This can be resolved by restarting the
server.

### Editing Existing Content

There is currently an issue with hot-reload when certain editors, such as GoLand
and Vim, are used to edit content that causes the edited page to fail loading.
This is due to "safe write" behavior in such editors which conflicts with NodeJS'
file watching system.

If you encounter an error similar to the one below, restarting the server will
resolve the issue.

```text
[ error ] ./pages/docs/commands/operator/migrate.mdx
Error: Cannot find module '/website/node_modules/babel-plugin-transform-define/lib/index.js' from '/website'
    at Array.map (<anonymous>)
    at cachedFunction.next (<anonymous>)
```

Alternately, if you disable safe write by searching goland's settings for "safe
write" or by running `:set backupcopy=yes` in vim, this should solve the issue
permanently.
