# Go Built-in Functions

Go provides a set of built-in functions that are predeclared and available for use without importing any packages. These functions perform various operations such as managing slices, maps, and channels, handling complex numbers, and error management.

## List of Built-in Functions

### Slice and Array Operations

-   **`append`** - Appends elements to the end of a slice and returns the updated slice

    ```go
    slice = append(slice, elem1, elem2, ...)
    ```

-   **`cap`** - Returns the capacity of a slice, array, or channel

    ```go
    capacity := cap(slice)
    ```

-   **`copy`** - Copies elements from a source slice to a destination slice and returns the number of elements copied

    ```go
    n := copy(dst, src)
    ```

-   **`len`** - Returns the length of a string, array, slice, map, or channel
    ```go
    length := len(slice)
    ```

### Map Operations

-   **`clear`** - Removes all elements from a map or slice (Go 1.21+)

    ```go
    clear(myMap)
    ```

-   **`delete`** - Removes an element from a map
    ```go
    delete(myMap, key)
    ```

### Memory Allocation

-   **`make`** - Allocates and initializes an object of type slice, map, or channel

    ```go
    slice := make([]int, length, capacity)
    myMap := make(map[string]int)
    ch := make(chan int)
    ```

-   **`new`** - Allocates memory for a new value and returns a pointer to it
    ```go
    ptr := new(Type)
    ```

### Channel Operations

-   **`close`** - Closes a channel, indicating that no more values will be sent on it
    ```go
    close(ch)
    ```

### Complex Number Operations

-   **`complex`** - Constructs a complex number from real and imaginary parts

    ```go
    c := complex(real, imag)
    ```

-   **`real`** - Returns the real part of a complex number

    ```go
    r := real(complexNum)
    ```

-   **`imag`** - Returns the imaginary part of a complex number
    ```go
    i := imag(complexNum)
    ```

### Comparison Operations

-   **`max`** - Returns the maximum value among the arguments (Go 1.21+)

    ```go
    maximum := max(a, b, c)
    ```

-   **`min`** - Returns the minimum value among the arguments (Go 1.21+)
    ```go
    minimum := min(a, b, c)
    ```

### Error Handling and Control Flow

-   **`panic`** - Stops ordinary flow of control and begins panicking

    ```go
    panic("something went wrong")
    ```

-   **`recover`** - Regains control of a panicking goroutine (must be called in a deferred function)
    ```go
    defer func() {
      if r := recover(); r != nil {
        // handle panic
      }
    }()
    ```

### Output Functions

-   **`print`** - Prints its arguments in an implementation-specific way (mainly for debugging)

    ```go
    print("debug message")
    ```

-   **`println`** - Like print but adds spaces between arguments and a newline at the end
    ```go
    println("debug message")
    ```

## Reference

For detailed information on each function, refer to the [Go Programming Language Specification](https://go.dev/ref/spec).

**Note:** The `print` and `println` functions are primarily for debugging and bootstrapping. For production code, use the `fmt` package instead.
