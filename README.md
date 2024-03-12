CRUD FUNCTIONS:

ERD: https://drive.google.com/file/d/1KQ48G-8WMyVjEezBc1ACC49GUlzQmYqm/view?usp=drive_link

## Using the Multi-Strategy Logger

- How to use Logger

- - Explain Trace
- - Explain Span
- Using Start
- - Explain Options Functions
- - - Context
- - Explain End

### Using the Start Function

The start function in Logger is mainly used for Honeycomb tracing for now. As seen from the function, it takes in the current *gin.Context and information on the trace, and an array of optional functions. The Option type is a struct of LoggerRepo functions.
`func (l *LoggerRepo) Start(c \*gin.Context, info string, options ...Option) trace.Span`

### Context

An example of calling this function is:
`span := or.Persistence.Logger.Start(c, "handler/SaveOrder", or.Persistence.Logger.SetContextWithSpanFunc())`

The optional function that was called is SetContextWithSpanFunc(), which sets the span into the current \*gin.Context. This is to ensure that the next function called will be recognised as its child. In functions that do not have any other child functions, e.g functions in the implementation layer, this optional function does not need to be called.

### Ending the span

In every function, if the Start is called, an End is needed.
`defer span.End()`
