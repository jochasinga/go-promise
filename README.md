# go-promise

Light-weight channel-compatible [Promise](https://promisesaplus.com/) implementation.

## Usage
Use a `Promise` to easily handle an async routine, like an HTTP request

```go

p := promise.New(func() (interface{}, error) {
        return http.Get("www.google.com")
})
_ = p.Then(func(res interface{}) {
        fmt.Println(res)
}, func(err error) {
        fmt.Println(err)
})

```

Convert channels to Promise and vice versa

```go

rc := make(chan interface{})

go func() {
        <-time.After(100)
        rc <- "hello"
        close(rc)
}()

p1 := promise.From(rc)
p1.Then(func(result interface{}) {
        fmt.Print(result)  // "hello"
})

p2 := promise.New(func() (interface{}, error) {
        <-time.After(100)
        return "hello", nil
})
rc, errc := p2.To()
select {
case <-errc:
case result := <-rc:
        fmt.Println(result)  // "hello"
}

```

Compose `*Promise.Resolve` and `*Promise.Reject` chains

```go

p := New(func() (interface{}, error) {
        return "hello", nil
})

_ = p.Then(func(result interface{}) {
        fmt.Println(result == "hello")        // true
        word := result.(string) + " world"
        p.Resolve(word)
}, func(err error) {
        fmt.Println(err == nil)               // true
}).Then(func(result interface{}) {
        fmt.Println(result == "hello world")  // true
        p.Reject("too bad")
}).Then(nil, func(err error) {
        fmt.Println(err.Error() == "too bad") // true
})

```
