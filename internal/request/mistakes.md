# HTTP Request Parser: Common Pitfalls and Solutions

## 1. Parse loop could spin forever

### What you did
```go
func (r *Request) parse(data []byte) (int, error) {
    totalBytesParsed := 0
    for r.enum != requestStateDone {
        n, err := r.parseSingle(data[totalBytesParsed:])
        if err != nil {
            return 0, err
        }
        totalBytesParsed += n
    }
    return totalBytesParsed, nil
}
```

### Why it's a problem
If `parseSingle` returns `0, nil` (meaning "I need more data, cannot progress"), the loop keeps calling `parseSingle` on the same data slice forever. That's exactly when your tests "hung indefinitely".

### Correct idea
You need to break when there is no progress:

```go
if n == 0 {
    break
}
totalBytesParsed += n
```

So the loop waits for more bytes instead of spinning.

## 2. Confusing "need more data" with "we're done"

At one point you suggested:

```go
if done || n == 0 {
    r.enum = requestStateDone
}
```

### Why that's logically wrong
- `done == true` → "I finished parsing this section" → yes, you should go to Done.
- `n == 0` (with `err == nil`) from `parseSingle` / `Headers.Parse` → "I need more data" → you must **not** mark state as done.

Treating `n == 0` as "done" will:
- End parsing early.
- Ignore headers that haven't arrived yet.

The right place to react to `n == 0` is in `parse` (stop looping and wait for more data), not by flipping the state to done.

## 3. Not initializing Headers before using it

### What you had
```go
type Request struct {
    RequestLine RequestLine
    enum        int
    Headers     headers.Headers
}

func RequestFromReader(reader io.Reader) (*Request, error) {
    ...
    req := Request{enum: requestStateParsingLine}
    ...
}
```

### Why it broke
- `headers.Headers` is a `map[string]string`.
- A zero-value map is `nil`.
- `Headers.Parse` writes to the map:
  ```go
  h[string(keyToLower)] = string(fieldValue)
  ```
- Writing to a nil map panics: `assignment to entry in nil map`.

### Fix
Initialize the headers map when you create the request:

```go
req := Request{
    enum:    requestStateParsingLine,
    Headers: headers.NewHeaders(),
}
```

## 4. Early header parsing: operating on all remaining data, not just one line (fixed)

Your earlier `Headers.Parse` logic used to operate on the entire data slice after trimming, rather than a single header line, which led to values bleeding into the next line.

### The correct logic:

1. Find the first `\r\n`:
   ```go
   idx := bytes.Index(data, []byte("\r\n"))
   if idx == -1 {
       return 0, false, nil
   }
   line := data[:idx]
   ```

2. If line is empty → end of headers.

3. Parse only that line (`SplitN` on `:` and `TrimSpace` the value).

4. Return `idx + 2` as bytes consumed.

Your final `Parse` matches this, which avoids mixing pieces of the next header into the value.

## 5. State transitions in parseSingle

You eventually ended up with:

```go
func (r *Request) parseSingle(data []byte) (int, error) {
    switch r.enum {
    case requestStateParsingLine:
        reqLine, bytes, err := parseRequestLine(string(data))
        ...
        r.enum = requestStateParsingHeaders
        r.RequestLine = reqLine
        return bytes, nil

    case requestStateParsingHeaders:
        n, done, err := r.Headers.Parse(data)
        ...
        if done {
            r.enum = requestStateDone
        }
        return n, nil

    case requestStateDone:
        return 0, errors.New("error: trying to read data in a done state")
    }
}
```

The earlier conceptual mistake was not having a separate state for `requestStateParsingHeaders` and trying to go straight from "parsed line" to "done" without looping through header parsing. That would have prevented you from properly parsing multiple headers across multiple reads.

Once you added the header-parsing state and let `parseSingle` drive transitions, your state machine matched the intended design.

---

If you'd like, I can quiz you a bit on when `parseSingle` should return `(0, nil)` versus when it should change the state to help cement the understanding.