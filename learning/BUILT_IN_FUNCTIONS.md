# Go Built-in Functions

## Overview

Go exposes a small set of predeclared (built-in) functions. They are available without imports and cover common language operations.

## Collections (arrays, slices, strings, maps, channels)

### Length and capacity

- `len(x)` returns the length of a string, array, slice, map, or channel.
- `cap(x)` returns the capacity of a slice, array, or channel.

```go
length := len(slice)
capacity := cap(slice)
```

### Append and copy

- `append` adds elements to a slice and returns the new slice.
- `copy` copies from a source slice to a destination slice and returns the number of elements copied.

```go
slice = append(slice, elem1, elem2)
n := copy(dst, src)
```

### Clear (Go 1.21+)

- `clear` removes all elements from a map or zeroes a slice.

```go
clear(myMap)
clear(mySlice)
```

## Map operations

### Delete

- `delete(map, key)` removes a key/value pair from a map.

```go
delete(myMap, key)
```

## Allocation

### Make vs. new

- `make` allocates and initializes slices, maps, and channels.
- `new` allocates zeroed storage for any type and returns a pointer.

```go
slice := make([]int, 0, 10)
myMap := make(map[string]int)
ch := make(chan int)

ptr := new(MyType)
```

## Channels

### Close

- `close(ch)` closes a channel, signaling that no more values will be sent.

```go
close(ch)
```

## Complex numbers

- `complex(real, imag)` constructs a complex number.
- `real(z)` and `imag(z)` extract parts.

```go
z := complex(2, 3)
r := real(z)
i := imag(z)
```

## Comparison helpers (Go 1.21+)

- `min(a, b, ...)` and `max(a, b, ...)` return the smallest/largest argument.

```go
minimum := min(a, b, c)
maximum := max(a, b, c)
```

## Panic and recovery

- `panic` stops normal execution.
- `recover` regains control of a panicking goroutine (only from a deferred function).

```go
defer func() {
    if r := recover(); r != nil {
        // handle panic
    }
}()

panic("something went wrong")
```

## Debug output (use sparingly)

- `print` and `println` are for debugging and bootstrapping only.
- Prefer `fmt` for production code.

```go
print("debug")
println("debug")
```

## Notes and pitfalls

- `recover` only works when called in a deferred function.
- `clear` resets elements to zero values (slices) or removes all entries (maps).
- `make` initializes runtime-managed types; `new` returns a pointer to a zero value.

## Reference

See the Go specification: https://go.dev/ref/spec
