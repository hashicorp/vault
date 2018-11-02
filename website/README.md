# Vault Website

This subdirectory contains the entire source for the [Vault Website][vault].
This is a [Middleman][middleman] project, which builds a static site from these
source files.

## Contributions Welcome!

If you find a typo or you feel like you can improve the HTML, CSS, or
JavaScript, we welcome contributions. Feel free to open issues or pull requests
like any normal GitHub project, and we'll merge it in.

## Running the Site Locally

Running the site locally is simple. Clone this repo and run `make website`. If it is your first time running the site, the build will take a little longer as it needs to download a docker image and a bunch of dependencies, so maybe go grab a coffee. On subsequent runs, it will be much faster as dependencies are cached.

Then open up `http://localhost:4567`. Note that some URLs you may need to append
".html" to make them work (in the navigation).

[middleman]: https://www.middlemanapp.com
[vault]: https://www.vaultproject.io
