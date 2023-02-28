# A Generic CLI
##Concurrent Processing of Long-Running Tasks

The `main.go` program demonstrates how to perform concurrent processing of long-running tasks using goroutines and channels in Go. The program reads a list of `Thing` objects from a CSV file, and then uses goroutines to create each `Thing` object in parallel. The progress of each `Thing` creation is displayed using an ASCII progress bar.
Progress bars are implemented by using Vladimir Bauer - vbauerster [mpb](https://github.com/vbauerster/mpb) package.

## Code Description

### Thing struct

```go
type Thing struct {
    Name        string
    Description string
    Value       string
}
```

The Thing struct has three fields:
Name, Description and Value.

createThing function
```go
func createThing(thing Thing, progress chan<- int, cancel <-chan struct{})
```

The createThing function simulates a long-running task by sleeping for a random amount of time between 500 and 1500 milliseconds and sending progress updates on the progress channel. The function listens for a cancellation signal on the cancel channel and exits early if it receives the signal.
To use this example together with a module with real long-running tasks, replace the createThing function with a function that implements this method signature.

readThingsFromFile function
```go
func readThingsFromFile(fileName string) ([]Thing, error)
```
The readThingsFromFile function reads a CSV file and returns a slice of Thing objects. The function uses the encoding/csv package to read the file.

launchThings function
```go
func launchThings(things []Thing, progressChannels []chan int, cancelChannels []chan struct{}, semaphore chan struct{}, p *mpb.Progress)
```
The launchThings function creates a progress bar for each Thing object and launches a goroutine to create each Thing object in parallel. The function limits the number of concurrent goroutines using a semaphore to prevent overwhelming the system. The function updates the progress bars as each Thing object is created.

main function
```go
func main()
```
The main function reads a list of Thing objects from a CSV file, creates a progress bar for each Thing object, and launches a goroutine to create each Thing object in parallel. The function limits the number of concurrent goroutines using a semaphore to prevent overwhelming the system. The function waits for all Thing objects to be created before exiting.
The example has a hard-coded limit of 3 concurrent goroutines. To change the limit, modify the value of the semaphore variable:
```go
semaphore := make(chan struct{}, 3)
```
Output

When you run the program, it reads the input.csv file in the same directory and creates a Thing object for each row in the file. For each Thing object, the program displays a progress bar indicating the current progress of the createThing function.

Here is an example of the program output:
<pre>
Creating 1/10 | Thing1-Name ████████████████████████╟  Done
Creating 2/10 | Thing2-Name ████████████████████████╟  Done
Creating 3/10 | Thing3-Name ████████████████████████╟  Done
Creating 4/10 | Thing4-Name ██████████░░░░░░░░░░░░░░╟  40 %
Creating 5/10 | Thing5-Name ████████░░░░░░░░░░░░░░░░╟  30 %
Creating 6/10 | Thing6-Name ██████░░░░░░░░░░░░░░░░░░╟  20 %
Creating 7/10 | Thing7-Name █░░░░░░░░░░░░░░░░░░░░░░░╟   0 %
</pre>
