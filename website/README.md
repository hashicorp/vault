# Vault Website

[![Netlify Status](https://img.shields.io/netlify/d72a27d2-aba4-46fd-bf70-d7da0b19b578?style=flat-square)](https://app.netlify.com/sites/vault-www/deploys)

This subdirectory contains the entire source for the [Vault
Website](https://vaultproject.io/). This is a [NextJS](https://nextjs.org/)
project, which builds a static site from these source files.

## Contributions Welcome!

If you find a typo or you feel like you can improve the HTML, CSS, or
JavaScript, we welcome contributions. Feel free to open issues or pull requests
like any normal GitHub project, and we'll merge it in 🚀

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
dependencies and test the changes within Docker, you'll need a new image. If
this is something you need to do, you can run `make build-image` to generate a
local Docker image with updated dependencies, then `make website-local` to use
that image and preview.

### With Node

If your local development environment has a supported version (v10.0.0+) of
[node installed](https://nodejs.org/en/) you can run:

- `npm install`
- `npm start`

and then visit `http://localhost:3000`.

If you pull down new code from github, you should run `npm install` again.
Otherwise, there's no need to re-run `npm install` each time the site is run,
you can just run `npm start` to get it going.

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

- `title` `(string)` - This is the title of the page that will be set in the
  HTML title.
- `description` `(string)` - This is a description of the page that will be set
  in the HTML description.

> ⚠️Since `api` is a reserved directory within NextJS, all `/api/**` pages are
> listed under the `/pages/api-docs` path.

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

To change the version of Vault displayed for download on the website, head over to `data/version.js` and change the number there. It's important to note that the version number must match a version that has been released and is live on `releases.hashicorp.com` -- if it does not, the website will be unable to fetch links to the binaries and will not compile. So this version number should be changed _only after a release_.

### Displaying a Prerelease

If there is a prerelease of any type that should be displayed on the downloads page, this can be done by editing `pages/downloads/index.jsx`. By default, the download component might look something like this:

```jsx
<ProductDownloader product="Vault" version={VERSION} downloads={downloadData} />
```

To add a prerelease, an extra `prerelease` property can be added to the component as such:

```jsx
<ProductDownloader
  product="Vault"
  version={VERSION}
  downloads={downloadData}
  prerelease={{
    type: 'release candidate', // the type of prerelease: beta, release candidate, etc.
    name: 'v1.0.0', // the name displayed in text on the website
    version: '1.0.0-rc1' // the actual version tag that was pushed to releases.hashicorp.com
  }}
/>
```

This configuration would display something like the following text on the website, emphasis added to the configurable parameters:

```
A {{ release candidate }} for Vault {{ v1.0.0 }} is available! The release can be <a href='https://releases.hashicorp.com/vault/{{ 1.0.0-rc1 }}'>downloaded here</a>.
```

You may customize the parameters in any way you'd like. To remove a prerelease from the website, simply delete the `prerelease` paremeter from the above component.

## Deployment

This website is hosted on Netlify and configured to automatically deploy anytime
you push code to the `stable-website` branch. Any time a pull request is
submitted that changes files within the `website` folder, a deployment preview
will appear in the github checks which can be used to validate the way docs
changes will look live. Deployments from `stable-website` will look and behave
the same way as deployment previews.

## Checking for Broken Links

There is a local script that can be used to check for broken links on the _current product website_ - you can start it by running `npm run linkcheck`. There will be a version of this script added as a github check in the near future!

## Known Issues

### Creating New Pages

There is currently a small bug with new page creation - if you create a new page
and link it up via subnav data while the server is running, it will report an
error saying the page was not found. This can be resolved by restarting the
server.

### Editing Existing Content

There is currently an issue with hot-reload when certain editors, such as GoLand
and Vim, are used to edit content that causes the edited page to fail loading.
This is due to "safe write" behavior in such editors which conflicts with
NodeJS' file watching system.

If you encounter an error similar to the one below, restarting the server will
resolve the issue.

```text
[ error ] ./pages/docs/commands/operator/migrate.mdx
Error: Cannot find module '/website/node_modules/babel-plugin-transform-define/lib/index.js' from '/website'
    at Array.map (<anonymous>)
    at cachedFunction.next (<anonymous>)
```
