#log – Centralized slog configuration

This package provides a single helper function to configure the default
[`slog`](https://pkg.go.dev/log/slog) logger for the entire application.

It uses [`github.com/lmittmann/tint`](https://github.com/lmittmann/tint) as the
handler to get a clean, colored output in the terminal.

---

## Installation

```bash
go get github.com/lmittmann/tint
```

# Example

package main

import (
"log/slog"

    "your-module-path/tlog"

)

func main() {
// Configure global logger
tlog.SetLogHandler()

    // Now use slog anywhere in your app
    slog.Debug("debug message", "foo", "bar")
    slog.Info("app started", "version", "1.0.0")
    slog.Error("something went wrong", "err", "example error")

}
