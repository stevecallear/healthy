# healthy
[![build](https://github.com/stevecallear/healthy/actions/workflows/build.yml/badge.svg)](https://github.com/stevecallear/healthy/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/stevecallear/healthy/graph/badge.svg?token=3JBUN06BOD)](https://codecov.io/gh/stevecallear/healthy)
[![Go Report Card](https://goreportcard.com/badge/github.com/stevecallear/healthy)](https://goreportcard.com/report/github.com/stevecallear/healthy)

`healthy` provides a simple mechanism to check the health of dependencies and wait for them to become available. While fundamentally an excuse to explore `synctest.Test` introduced in Go 1.25 the module can simplify dependency checking in non-production scenarios such as local docker compose setups.

## Getting Started
```
go get github.com/stevecallear/healthy@latest
```
```
err := healthy.Wait(healthy.HTTP("http://dependency:8080/healthy").Expect(http.StatusOK))
```

## Checks
Healthy includes checks for TCP, HTTP and files. Additional checks can be added by implementing `Check` or providing a `CheckFunc`. 

## Metadata
Check metadata can be provided by implementing `MetadataCheck` or wrapping a `CheckFunc` with `WithMetadata`.
```
c := healthy.WithMetadata(func(ctx context.Context) error {
    // implement check
}, "type", "custom", "target", "something")
```

## Execution
While checks can be executed directly by calling `check.Healthy`, they are intended to be executed as a group with multiple attempts using `healthy.New(checks...).Wait()`. `Wait` accepts a number of execution options relating to context, timeout and delay. Checks are executed until either context cancellation or timeout, specified via `WithContext` or `WithTimeout` respectively.

`WithCallback` accepts a callback function that is invoked for each check execution and error. The following example logs the result using `slog`:
```
healthy.Wait(
    healthy.TCP("host:8080"),
    healthy.WithCallback(func(ctx context.Context, err error) {
        args := []any{}
        for k, v := range healthy.GetContextMetadata(ctx) {
            args = append(args, k, v)
        }

        level := slog.LevelInfo
        if err != nil {
            level = slog.LevelError
            args = append(args, "err", err)
        }

        slog.Log(ctx, level, "health check result", args...)  
    }),
)
```

## Parallel Checks
The initial implementation of the module allowed multiple checks to be executed in parallel. While convenient, this created a confused API and limited configuration for what was effectively a wrapper over `errgroup.Group`.

If parallel check execution is desired, then it is trivial (if a bit more verbose) to execute them using `errgroup`:
```
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error {
    return healthy.Wait(
        healthy.TCP("host:8080"),
        healthy.WithContext(ctx),
        healthy.WithCallback(callback),
    )
})
g.Go(func() error {
    return healthy.Wait(
        healthy.HTTP("http://dependency:8080/health"),
        healthy.WithContext(ctx),
        healthy.WithCallback(callback),
    )
})
err = g.Wait()
```