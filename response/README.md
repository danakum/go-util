### Response

Response helper package of [core](https://github.com/danakum/go-util) library

#### Usage

Render an error

```go
import "github.com/danakum/go-util/response"

response.HandleError(ctx, err, w)

```gitlab.mytaxi.lk/pickme/go-util
Above will check the error type and if
- Error is a [DomainError]() api response render helpers will generate following output

```json
    
```


