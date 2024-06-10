# slogger

`slogger` is a Go package that wraps the `slog` package to add some convenient functionalities. It provides a flexible and easy-to-use interface for logging with additional options like different output modes, attribute customization, and more.

## Features

- **Output Modes**: Supports both text and JSON output modes.
- **Custom Attributes**: Allows adding custom attributes to logs.
- **Multiple Destinations**: Supports logging to different `io.Writer` destinations.
- **Log Level Configuration**: Easy configuration of log levels.
- **Sub-Loggers**: Allows creation of sub-loggers for different components.
- **DevNullLogger**: A logger that writes to nowhere, useful for testing.

## Installation

To install the package, run:

```sh
go get github.com/xorima/slogger
```

## Usage

### Creating a Logger

To create a new logger with default options:

```go
import (
    "github.com/xorima/slogger"
)

loggerOpts := slogger.NewLoggerOpts("ServiceName", "ApplicationName")
logger := slogger.NewLogger(loggerOpts)
```

### Customizing Logger

You can customize the logger by providing additional options:

```go
loggerOpts := slogger.NewLoggerOpts(
    "ServiceName",
    "ApplicationName",
    slogger.WithDestination(os.Stderr),
    slogger.WithJsonOutput(),
    slogger.WithAttr(slog.String("environment", "production")),
)

logger := slogger.NewLogger(loggerOpts, slogger.WithLevel("debug"))
```

### Sub-Loggers

Create sub-loggers for different components:

```go
componentLogger := slogger.SubLogger(logger, "ComponentName")
componentLogger.Info("This is a log message from the component")
```

### DevNullLogger

For testing purposes, you can use `DevNullLogger`:

```go
devNullLogger := slogger.NewDevNullLogger()
devNullLogger.Info("This message will not be logged anywhere")
```

### Error Attributes

To log errors with attributes:

```go
err := errors.New("something went wrong")
logger.Error("An error occurred", slogger.ErrorAttr(err))
```

## API Reference

### `NewLoggerOpts`

Creates a new `LoggerOpts` instance.

```go
func NewLoggerOpts(serviceName, applicationName string, opts ...func(o *LoggerOpts)) *LoggerOpts
```

### `WithDestination`

Sets the log destination.

```go
func WithDestination(destination io.Writer) func(o *LoggerOpts)
```

### `WithJsonOutput`

Sets the output mode to JSON.

```go
func WithJsonOutput() func(o *LoggerOpts)
```

### `WithAttr`

Adds a custom attribute to the logger.

```go
func WithAttr(attr slog.Attr) func(o *LoggerOpts)
```

### `NewLogger`

Creates a new `slog.Logger` with the specified options.

```go
func NewLogger(loggerOpts *LoggerOpts, handlerOpts ...func(o *slog.HandlerOptions)) *slog.Logger
```

### `WithSource`

Includes the source of the log message.

```go
func WithSource() func(o *slog.HandlerOptions)
```

### `WithLevel`

Sets the log level.

```go
func WithLevel(level string) func(o *slog.HandlerOptions)
```

### `WithReplaceAttr`

Replaces an attribute with a custom function.

```go
func WithReplaceAttr(fn func(groups []string, a slog.Attr) slog.Attr) func(o *slog.HandlerOptions)
```

### `SubLogger`

Creates a sub-logger for a specific component.

```go
func SubLogger(logger *slog.Logger, componentName string) *slog.Logger
```

### `DevNullLogger`

A logger that writes to nowhere.

```go
type DevNullLogger struct{}
func (d *DevNullLogger) Write(p []byte) (n int, err error)
func NewDevNullLogger() *slog.Logger
```

### `ErrorAttr`

Creates an error attribute.

```go
func ErrorAttr(err error) slog.Attr
```

## License

This project is licensed under the Apache 2 License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## Contact

For any questions or suggestions, feel free to reach out.
