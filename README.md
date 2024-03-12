CRUD FUNCTIONS:

ERD: https://drive.google.com/file/d/1KQ48G-8WMyVjEezBc1ACC49GUlzQmYqm/view?usp=drive_link

## Using the Multi-Strategy Logger

### How to use Logger

1. Initialise NewLoggerRepository and store inside Persistence

func NewPersistence() (\*Persistence, error) {

    // Initialise others

    // Initialise Logger
    logger := logger.NewLoggerRepositories([]string{"Honeycomb", "Zap"})

    return &Persistence{
        // Other databases and engines
        Logger:             logger,
    }

2. Then, call it from the Persistence as r.p.Logger
func (r orderRepo) GetOrder(id uint64) (\*entity.Order, error) {

    span := r.p.Logger.Start(r.c, "infrastructure/implementations/GetOrder", map[string]interface{}{"id": id})
    defer span.End()

    var order *entity.Order
    err := r.p.ProductDb.Debug().Unscoped().Preload("OrderedItems").Where("order_id = ?", id).Take(&order).Error
    if err != nil {
        r.p.Logger.Error(err.Error(), map[string]interface{}{"data": order})
        return nil, err
    }

    r.p.Logger.Info("get order", map[string]interface{}{"data": order})

    if errors.Is(err, gorm.ErrRecordNotFound) {
        r.p.Logger.Error(err.Error(), map[string]interface{}{"data": order})
        return nil, errors.New("order not found")
    }

    return order, nil

}

### Explain Trace

The overall representation of a request’s journey through a system. A trace is composed of one or more spans.

### Explain Span

A single unit of work within a trace. Each span represents a specific function or operation and contains metadata such as start time, end time, and attributes.

### Using the Start Function

The start function in Logger is mainly used for Honeycomb tracing for now. As seen from the function, it takes in the current \*gin.Context and information on the trace, and an array of optional functions. The Option type is a struct of LoggerRepo functions.

`func (l *LoggerRepo) Start(c *gin.Context, info string, options ...Option) trace.Span`

### Context

An example of calling this function is:

`span := or.Persistence.Logger.Start(c, "handler/SaveOrder", or.Persistence.Logger.SetContextWithSpanFunc())`

The optional function that was called is SetContextWithSpanFunc(), which sets the span into the current \*gin.Context. This is to ensure that the next function called will be recognised as its child. In functions that do not have any other child functions, e.g functions in the implementation layer, this optional function does not need to be called.

### Ending the span

In every function, if the Start is called, an End is needed.

`defer span.End()`
