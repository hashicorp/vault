# Overview

**duo_api_golang** - Go language bindings for the Duo APIs (both auth and admin).

## Duo Auth API

The Auth API is a low-level, RESTful API for adding strong two-factor authentication to your website or application.

This module's API client implementation is *complete*; corresponding methods are exported for all available endpoints.

For more information see the [Auth API guide](https://duo.com/docs/authapi).

## Duo Admin API

The Admin API provides programmatic access to the administrative functionality of Duo Security's two-factor authentication platform.

This module's API client implementation is *incomplete*; methods for fetching most entity types are exported, but methods that modify entities have (mostly) not yet been implemented. PRs welcome!

For more information see the [Admin API guide](https://duo.com/docs/adminapi).
