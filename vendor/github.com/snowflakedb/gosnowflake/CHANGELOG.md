## Version 1.1.18
- Ignore session gone error when closing session.
- Set the latest Go version to 1.12

## Version 1.1.17
- TIMESTAMP_LTZ is not converted properly
- Add SERVICE_NAME support to Golang

## Version 1.1.16
- Fix custom json parser (@mhseiden)
- Modify Support link on Gosnowflake github

## Version 1.1.15
- Perform check on error for pingpong response (@ChTimTsubasa)
- Add new 1.11.x for testing (@ChTimTsubasa)
- Handle extra snowflake parameter 'type'(@ChTimTsubasa)

## Version 1.1.14
- Disable tests for golang 1.8 (@ChTimTsubasa)
- Follow lint's new include path (@ChTimTsubasa)
- Do not sleep through a context timeout(@mhseiden)

## Version 1.1.13
- User configurable Retry timeout for downloading (@mhseiden)
- Implement retry-uuid in the url (@ChTimTsubasa)
- Remove unnecessary go routine and fix context cancel/timeout handling (@ChTimTsubasa)

## Version 1.1.12
- Allow users to customize their glog through different vendoring. (@ChTimTsubasa)
- Doc improvment for region parameter description (@smtakeda)

## Version 1.1.11

- (Private Preview) Added key pair authentication. (@ChTimTsubasa)
- Changed glog timestamp to UTC (@ChTimTsubasa)
- (Experimental) Added `MaxChunkDownloadWorkers` and `CustomJSONDecoderEnabled` to tune the result set download performance. (@mhseiden)

## Version 1.1.10

- Fixed heartbeat timings. It used to start a heartbeat per query. Now it starts per connection and closes in `Close` method. #181 
- Removed busy wait from 1) the main thread that waits for download chunk worker to finish downloads, 2) the heartbeat goroutine/thread to trigger the heartbeat to the server.

## Version 1.1.9

- Fixed proxy for OCSP (@brendoncarroll)
- Disabled megacheck for Go 1.8
- Changed the UUID dependency (@kenshaw)

## Version 1.1.8

- Removed username restriction for oAuth

## Version 1.1.7

- Added `client_session_keep_alive` option to have a heartbeat in the background every hour to keep the connection alive. Fixed #160
- Corrected doc about OCSP.
- Added OS session info to the session.

## Version 1.1.6

- Fixed memory leak in the large result set. The chunk of memory is freed as soon as the cursor moved forward.
- Removed glide dependency in favor of dep #149 (@tjj5036)
- Fixed username and password URL escape issue #151
- Added Go 1.10 test.

## Version 1.1.5

- Added externalbrowser authenticator support PR #141, #142 (@tjj5036)

## Version 1.1.4

- Raise HTTP 403 errors immediately after the authentication failure instead of retry until the timeout. Issue #138 (@dominicbarnes)
- Fixed vararg error message.

## Version 1.1.3

- Removed hardcoded `public` schema name in case not specified.
- Fixed `requestId` value

## Version 1.1.2

- `nil` should set to the target value instead of the pointer to the target

## Version 1.1.1

- Fixed HTTP 403 errors when getting result sets from AWS S3. The change in the server release 2.23.0 will enforce a signature of key for result set.

## Version 1.1.0

- Fixed #125. Dropped proxy parameters. HTTP_PROXY, HTTPS_PROXY and NO_PROXY should be used.
- Improved logging based on security code review. No sensitive information is logged.
- Added no connection pool example
- Fixed #110. Raise error if the specified db, schema or warehouse doesn't exist. role was already supported.
- Added go 1.9 config in TravisCI
- Added session parameter support in DSN.

## Vesrion 1.0.0

- Added [dep](https://github.com/golang/dep) manifest (@CrimsonVoid)
- Bumped up the version to 1.0.0
