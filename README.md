# VALORANT Assistant

A cli tool for valorant that lets you:

- Check your stores
- Check your MMR (Rank)
- Check your wallet (VP, RP, Kingdom Credits, Free Agents)

## Auth

- Supports multi-factor authentication

Auth credentials (username, password) are store in the users home directory in `.valocli`, and the auth token(s), entitlement token and user id are cached in the same directory for an hour (riot has expiry for the auth token set to an hour)

## TODO

- Support more endpoints
- (Maybe) Add local endpoints support using the lockfile
- Better caching for responses (so that we dont hit a rate limit lol)

### Contributing

feel free to open a PR if there's a feature you want to add

also feel free to make my code cleaner/better (still improving at Go :D)
