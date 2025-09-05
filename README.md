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
err := healthy.New(
    healthy.TCP("http://dynamodb:8000"),
    healthy.HTTP("http://dependency:8080/health").Expect(http.StatusOK),
).Wait()
```

## Checks
Healthy includes checks for TCP, HTTP and files. Additional checks can be added by implementing `healthy.Check`. Alternatively `healthy.NewCheck` returns a new check using a supplied `CheckFunc` and information:
```
c := healthy.NewCheck(func(ctx context.Context) error {
    // implement check
}, "type", "custom", "target", "something")
```

## Execution
While checks can be executed directly by calling `check.Healthy`, they are intended to be executed as a group with multiple attempts using `healthy.New(checks...).Wait()`. `Wait` accepts a number of execution options relating to context, timeout and delay. Checks are executed until either context cancellation or timeout, specified via `WithContext` or `WithTimeout` respectively.

`WithCallback` accepts a callback function that is invoked for each check execution and includes the check info, attempt and result. The following example logs the result using `slog`:
```
healthy.New(healthy.TCP("http://dynamodb:8000")).
    Wait(healthy.WithCallback(func(ctx context.Context, r healthy.Result) {
        level := slog.LevelInfo
        args := []any{"attempt", r.Attempt}

        if r.Err != nil {
            level = slog.LevelError
            args = append(args, "err", r.Err.Error())
        }

        for k, v := range r.Info {
            args = append(args, k, v)
        }

        slog.Log(context.Background(), level, "check attempted", args...)
    }))
```