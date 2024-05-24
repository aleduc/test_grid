### Principles:
* It is possible to use external libraries
* KISS
* Error logs, logger - almost omitted.
* Parametrization - ignored, in most cases I've just used some constants, instead of params.
* Async part - doesn't make sense to add there(mocks), there are no goroutine blocks( I/O operations).
but in real app I assume that we have a lot of chargers and call them byt HTTP/GRPC, consequently it should be added there. 
you can find it as **SetValueAsync**


### Requirements:
1. It is not clear, should I charge batteries based on *actual*(post factum) or forecast values.
I decided to update it based on forecast.
1. Cron runs in exact time, instead every half-hour, it matches better with API.

### CMD
**cmd/test_trigger/test_grid.go** - main run.

### How to run:
1. From project root(go must be version 1.21):

    **Example:** go run cmd/test_grid/main.go
2. **From project root:**

    **Example**: ./for_run
