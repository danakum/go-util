## Traceable Context

**Traceable Context** is a wrapper around go context which will provide a way to share a traceable UUID between
different contexts


```go
    ctx := traceable_context.WithUUID(uuid.New())
```
