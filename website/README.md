# Vault Website

[![Netlify Status](https://img.shields.io/netlify/f7fa8963-0022-4a0e-9ccf-f5385355906b?style=flat-square)](https://app.netlify.com/sites/vault-docs-platform/deploys)

This subdirectory contains the entire source for the [Vault Website](https://vaultproject.io/). This is a [NextJS](https://nextjs.org/) project, which builds a static site from these source files.

## Contributions Welcome!

If you find a typo or you feel like you can improve the HTML, CSS, or JavaScript, we welcome contributions. Feel free to open issues or pull requests like any normal GitHub project, and we'll merge it in 🚀

## Running the Site Locally

The website can be run locally through node.js or Docker. If you choose to run through Docker, everything will be a little bit slower due to the additional overhead, so for frequent contributors it may be worth it to use node. Also if you are a vim user, it's also worth noting that vim's swapfile usage can cause issues for the live reload functionality. In order to avoid these issues, make sure you have run `:set backupcopy=yes` within vim.

### With Docker

Running the site locally is simple. Provided you have Docker installed, clone this repo, run `make`, and then visit `http://localhost:3000`.

The docker image is pre-built with all the website dependencies installed, which is what makes it so quick and simple, but also means if you need to change dependencies and test the changes within Docker, you'll need a new image. If this is something you need to do, you can run `make build-image` to generate a local Docker image with updated dependencies, then `make website-local` to use that image and preview.

### With Node

If your local development environment has a supported version (v10.0.0+) of [node installed](https://nodejs.org/en/) you can run:

- `npm install`
- `npm start`

and then visit `http://localhost:3000`.

If you pull down new code from github, you should run `npm install` again. Otherwise, there's no need to re-run `npm install` each time the site is run, you can just run `npm start` to get it going.

## Editing Content

Documentation content is written in [Markdown](https://www.markdownguide.org/cheat-sheet/) and you'll find all files listed under the `/pages` directory.

To create a new page with Markdown, create a file ending in `.mdx` in the `pages/` directory. The path in the pages directory will be the URL route. For example, `pages/hello/world.mdx` will be served from the `/hello/world` URL.

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

> ⚠️Since `api` is a reserved directory within NextJS, all `/api/**` pages are listed under the `/pages/api-docs` path.

### Editing Sidebars

The structure of the sidebars are controlled by files in the [`/data` directory](data).

- Edit [this file](data/docs-navigation.js) to change the **docs** sidebar
- Edit [this file](data/api-navigation.js) to change the **api docs** sidebar

To nest sidebar items, you'll want to add a new `category` key/value accompanied by the appropriate embedded `content` values.

- `category` values will be **directory names** within the `pages` directory
- `content` values will be **file names** within their appropriately nested directory.


### Deployment

This website is hosted on Netlify and configured to automatically deploy anytime you push code to the `stable-website` branch. Any time a pull request is submitted that changes files within the `website` folder, a deployment preview will appear in the github checks which can be used to validate the way docs changes will look live. Deployments from `stable-website` will look and behave the same way as deployment previews.

## Known Issues

### Creating New Pages

There is currently a small bug with new page creation - if you create a new page and link it up via subnav data while the server is running, it will report an error saying the page was not found. This can be resolved by restarting the server.

### Editing Existing Content

There is currently an issue with hot-reload when certain editors, such as GoLand
and Vim, are used to edit content that causes the edited page to fail loading.
This is due to "safe write" behavior in such editors which conflicts with
NodeJS' file watching system.

If you encounter an error similar to the one below, restarting the server will resolve the issue.

```text
[ error ] ./pages/docs/commands/operator/migrate.mdx
Error: Cannot find module '/website/node_modules/babel-plugin-transform-define/lib/index.js' from '/website'
    at Array.map (<anonymous>)
    at cachedFunction.next (<anonymous>)
```
