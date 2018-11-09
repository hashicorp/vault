# Vault Website

This subdirectory contains the entire source for the [Vault Website][vault].
This is a [Middleman][middleman] project, which builds a static site from these
source files.

## Updating Navigation

There are a couple different places on the website that present navigation interfaces with differing levels of detail.

On the homepage, docs index page, and api docs index page, there are grids of major categories [that look like this](https://cl.ly/73df9722848d/Screen%20Shot%202018-11-09%20at%2011.40.56%20AM.png). These major category grids can be updated through [`data/docs_basic_categories.yml`](data/docs_basic_categories.yml) and [`data/api_basic_categories.yml`](data/api_basic_categories.yml).

On the docs and api index pages, there are more detailed breakdowns of top-level documentation pages within each category [that look like this](https://cl.ly/b05cf42402eb/Screen%20Shot%202018-11-09%20at%2011.43.25%20AM.png). These more detailed category listings can be updated through [`data/docs_detailed_categories.yml`](data/docs_detailed_categories.yml) and [`data/api_detailed_categories.yml`](data/api_detailed_categories.yml).

Finally, within a given docs page, there is a sidebar which displays a fully nested version of all docs pages. This sidebar navigation can be updated through via middleman's layouts, found at [`source/layouts/docs.erb`](source/layouts/docs.erb) and [`source/layouts/api.erb`](source/layouts/api.erb). You will see within these files that it is no longer necessary to type out full nested html list item and link tags, you can simply add the documentation page's slug, defined as `sidebar_current` within the frontmatter of any docs markdown file. The sidebar nav component will go find the page by slug and render out its human-readable title and a link for you. This component does not allow broken links or nesting mistakes, so if you make a typo on the slug or put a page in the wrong category, the build will fail.

## Contributions Welcome!

If you find a typo or you feel like you can improve the HTML, CSS, or
JavaScript, we welcome contributions. Feel free to open issues or pull requests
like any normal GitHub project, and we'll merge it in.

## Running the Site Locally

Running the site locally is simple. Clone this repo and run `make website`. If it is your first time running the site, the build will take a little longer as it needs to download a docker image and a bunch of dependencies, so maybe go grab a coffee. On subsequent runs, it will be much faster as dependencies are cached.

Then open up `http://localhost:4567`. Note that some URLs you may need to append
".html" to make them work (in the navigation).

> **Note:** We are currently working through some issues with the `make website` command introduced by architecture changes to the site. You will likely experience slow performance while we get these issues patched up. In order to bypass these issues, you can run middleman directly on your machine by running `gem install middleman`, then `bundle && bundle exec middleman`. Assuming you have a reasonably recent version of ruby installed, this will run the site locally.

[middleman]: https://www.middlemanapp.com
[vault]: https://www.vaultproject.io
