## Introduction

Conventional approach towards errors in *golang* is quite limited.

The typical case implies an error being created at some point:
```go
return errors.New("now this is unfortunate")
```

Then being handled with a no-brainer:
```go
if err != nil {
  return nil, err
}
```

And, finally, handled by printing it to the log file:
```go
log.Errorf("Error: %s", err)
```

It does'n take long to find out that this is not often enough. There's little fun in solving the issue when everything a developer is able to observe is a line in the log that looks like on of those:
> Error: EOF

> Error: unexpected '>' at the beginning of value

> Error: wrong argument value

An *errorx* library makes an approach to create a toolset that would help remedy this issue with these considerations in mind:
* No extra care should be required for an error to have all the required debug information; it is the opposite that may constitute a special case
* There must be a way to distinguish one kind of error from another, as they may imply or require a different handling in user code
* Errors must be composable, and patterns like ```if err == io.EOF``` defetat that purpose, so they should be avoided
* Some context information may be added to the error along the way, and there must be a way to do so without altering the semantics of the error
* It must be easy to create an error, add some context to it, check for it
* A kind of error that requires a special treatment by the caller *is* a part of a public API; an excessive amount of such kinds is a code smell

As s result, the goal of the library is to provide a brief, expressive syntax for a conventional error handling and to discourage usage patterns that bring less value than harm.

Error-related, negative codepath is typically less well tested, though of, and may confuse the reader more than its positive counterpart. Therefore, an error system could do well without too much of a flexibility and unpredictability

# errorx

With *errorx*, the pattern above looks like this:

```go
if err != nil {
  return nil, errorx.IllegalState.New("unfortunate")
}
```
```go
if err != nil {
  return nil, errorx.Decorate(err, "this could be be much better")
}
```
```go
log.Errorf("Error: %+v", err)
```

An error message will look something like this:

```
Error: this could be be much better, cause: common.illegal_state: unfortunate
 at main.culprit()
	main.go:21
 at main.innocent()
	main.go:16
 at main.main()
	main.go:11
  ```

Now we have some context to our little problem, as well as a full stack trace of the original cause - which is, in effect, all that you really need, most of the time. ```errorx.Decorate``` is handy to add some info which a stack trace does not already hold: an id of the relevant entity, a portion of the failed request, etc. In all other cases, the good old ```if err != nil {return nil}``` still works for you.

And this, frankly, may be quite enough. With a set of standard error types provided with *errorx* and a syntax to create your own, the best way to deal with errors is in an opaque manner: create them, add information and log as tome point. Whenever this is sufficient, don't go any further. The simpler, the better.
