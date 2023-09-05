// Command flies provides a webserver which produces detailed logs of each
// request it receives.
//
// The following environment variables are available to configure flies:
//
//    FLIES_FORMAT                Can be "json", "wire", or "template". Default is "pretty".
//    FLIES_TEMPLATE_FILE         Path to template used for "template" format.
//    FLIES_RAW                   Show the raw TCP output only.
//    FLIES_PORT                  Port that the server should listen on.
//    FLIES_RESPONSE_STATUS_CODE  HTTP status code to respond with.
//    FLIES_RESPONSE_STATUS       HTTP status to respond with.
//    FLIES_RESPONSE_BODY_CONTENT Content of response body.
//
// To start the server:
//
//    go run .
//
package main
