# go-promise

Light-weight channel-compatible [Promise](https://promisesaplus.com/) implementation.

## Switchable from-to channels

```go

        rc := make(chan interface{})
        p1 := promise.From(rc)
        p1.Then(func(result interface{}) {
                fmt.Print(result)
        })

        p2 := promise.New(func() (interface{}, error) {
		<-time.After(100)
		return "hello", nil
	})
        rc, errc := p2.To()
        select {
        case <-rc:
        case <-errc:
        }

```

## Usage

```go

        p := promise.New(func() (interface{}, error) {
                return http.Get("www.google.com")
        })
        p.Then(func(result interface{}) {
                if res, ok := result.(*http.Response); ok {
                        fmt.Print(res.StatusCode)
                }
        }, func(err error) {
                panic(err)
        })

```
