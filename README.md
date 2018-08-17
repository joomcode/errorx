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

It doesn't take long to find out that this is not often enough. There's little fun in solving the issue when everything a developer is able to observe is a line in the log that looks like on of those:
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
  return nil, errorx.IllegalState.New("unfortunate")
```
```go
if err != nil {
  return nil, errorx.Decorate(err, "this could be so much better")
}
```
```go
log.Errorf("Error: %+v", err)
```

An error message will look something like this:

```
Error: this could be so much better, cause: common.illegal_state: unfortunate
 at main.culprit()
	main.go:21
 at main.innocent()
	main.go:16
 at main.main()
	main.go:11
  ```

Now we have some context to our little problem, as well as a full stack trace of the original cause - which is, in effect, all that you really need, most of the time. ```errorx.Decorate``` is handy to add some info which a stack trace does not already hold: an id of the relevant entity, a portion of the failed request, etc. In all other cases, the good old ```if err != nil {return nil}``` still works for you.

And this, frankly, may be quite enough. With a set of standard error types provided with *errorx* and a syntax to create your own (note that a name of the type is a good way to express its semantics), the best way to deal with errors is in an opaque manner: create them, add information and log as some point. Whenever this is sufficient, don't go any further. The simpler, the better.

## Error check

If an error requires special treatment, it may be done like this:
```go
if errorx.IsOfType(err, MyError) {
  // handle
}
```

Note that it is never a good idea to inspect a message of an error. Type check, on the other hand, is sometimes OK, especially if this technique is used inside of a package rather than forced upon API users.

An alternative is a mechanisms called **traits**:
```go
TimeoutElapsed       = CommonErrors.NewType("timeout", Timeout())
```

Here, ```TimeoutElapsed``` error type is created with a Timeout() trait, and errors may be checked against it:
```go
if errorx.HasTrait(err, errorx.Timeout()) {
  // handle
}
```

Note that here a check is made against a trait, not a type, so any type with the same trait would pass it. Type check is more restricted this way and creates tighter dependency if used outside of an originating package. It allows for some little flexibility, though: via a subtype feature a broader type check can be made.

## Wrap

The example above introduced ```errorx.Decorate()```, a syntax used to add message as an error is passed along. This mechanism is highly non-intrusive: any properties an original error possessed, a result of a  ```Decorate()``` will possess, too.

Sometimes, though, it is not the desired effect. A possibility to make a type check as a double edged one, and should be restricted as often as it is allowed. The bad way to do so would be to create a new error and to pass an ```Error()``` output as a message. Among other possible issues, this would either lose or duplicate the stack trace information.

A better alternative is:
```go
return MyError.Wrap(err, "fail")
```

With ```Wrap()```, an original error is fully retained for the log, but hidden from type checks by the caller.

See ```WrapMany()``` and ```DecorateMany()``` for more sophisticated cases.

## Stack traces

As an essential part of debug information, stack traces are included in all *errorx* errors by default.

When an error is passed along, the original stack trace is simply retained, as this typically takes place along the lines if the same frames that were originally captured. When an error is received from another goroutine, use this to add frames that would otherwise be missing:

```go
return EnhanceStackTrace(<-errorChan, "task failed")
```

Result would look like this:
```
Error: common.illegal_state: unfortunate
 at main.proxy()
	main.go:17
 at main.main()
	main.go:11
 ----------------------------------
 at main.culprit()
	main.go:26
 at main.innocent()
	main.go:21
  ```

On the other hand, some errors do not require a stack trace. Some may be used as a control flow mark, other are known to be benign. Stack trace could be omitted by not using the ```%+v``` formatting, but the better alternative is to modify the error type:

```go
ErrInvalidToken    = AuthErrors.NewType("invalid_token").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
```

This way, a receiver of an error always treats it the same way, and it is the producer who modifies the behaviour. Following, again, the principle of opacity.

Other relevant tools include ```EnsureStackTrace(err)``` to provide an error of unknown nature with a stack trace, if it lacks one.

## More

See godoc for other *errorx* features:
* Namespaces
* Type switches
* ```errorx.Ignore```
* Trait inheritance
* Dynamic properties
* Panic-related utils
* Type registry
* etc.
