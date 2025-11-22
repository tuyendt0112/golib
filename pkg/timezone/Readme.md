# config timezone all services

Default timezone is UTC
You can change timezone by set environment variable APP_TIMEZONE in your service
package main

import (
\_ "go.alireviews.dev/shared/pkg/tzinit" // import this package file main server or worker or scheduler
)

func main() {
// your code here
}
