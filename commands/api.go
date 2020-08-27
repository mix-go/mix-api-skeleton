package commands

import (
    "context"
    "fmt"
    "github.com/mix-go/dotenv"
    "github.com/mix-go/gin"
    "github.com/mix-go/mix-api-skeleton/globals"
    "github.com/mix-go/mix-api-skeleton/routes/api"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"
)

type APICommand struct {
}

func (t *APICommand) Main() {
    // server
    gin.SetMode(dotenv.Getenv("GIN_MODE").String(gin.ReleaseMode))
    router := gin.New(api.RouteDefinitions...)
    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // signal
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-ch
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        if err := srv.Shutdown(ctx); err != nil {
            globals.Logger().Errorf("Server shutdown error: %s", err)
        }
    }()

    // logger
    router.Use(gin.LogrusWithFormatter(globals.Logger(), func(params gin.LogFormatterParams) string {
        return fmt.Sprintf("%s|%s|%d|%s",
            params.Method,
            params.Path,
            params.StatusCode,
            params.ClientIP,
        )
    }))

    // run
    if err := srv.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "http: Server closed") {
        panic(err)
    }
}
