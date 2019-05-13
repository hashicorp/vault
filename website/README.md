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

When running the site locally, you can choose between running it directly on your machine, or running it through Docker. Docker has the advantage of requiring only Docker to be installed - no other dependencies are needed on your machine. However, Docker's overhead makes the site's compilation perform much slower than running it directly on your machine. If you are a frequent contributor, are bothered by the performance in Docker, or have no issues with installing ruby and node / already have them installed, it might be an advantage to try running the site directly on your machine. Instructions for both approaches are included below.

### Running the Site with Docker

First, make sure that [docker](docker) is installed. It can be installed in many ways, [the desktop app](docker-desktop) is the simplest. To run the site, clone this repo down, `cd` into the `website` directory, and run `make website`. If it is your first time running the site, the build will take a little longer as it needs to download a docker image and a bunch of dependencies, so maybe go grab a coffee. On subsequent runs, it will be faster as dependencies are cached.

### Running the Site Directly

This site requires a recent version of ruby as well as nodejs to be installed in order to run. There are [many ways to install ruby](https://www.ruby-lang.org/en/documentation/installation/), we recommend [rbenv](rbenv), which has very clear installation instructions in its readme, linked here, and installing ruby version `2.4.3`. Once ruby has been installed, you will need to install `bundler` as well, using `gem install bundler`. Node is quite easy to install [via universal binary](node) or [homebrew](homebrew) if you are a mac user.

Once ruby and node have been installed, within this directory, you can run `sh bootstrap.sh` to install all the dependencies needed to run the site, then run `middleman` to start the dev server.

### Browsing the Site Locally

Once you have the local dev server running, head to `http://localhost:4567` in your browser. Note that for some URLs, you may need to append
".html" to make them work (in the navigation).

[middleman]: https://www.middlemanapp.com
[vault]: https://www.vaultproject.io
[docker]: https://www.docker.com/
[docker-desktop]: https://www.docker.com/products/docker-desktop
[rbenv]: https://github.com/rbenv/rbenv#installation
[node]: https://nodejs.org/en/
[homebrew]: https://brew.sh/
