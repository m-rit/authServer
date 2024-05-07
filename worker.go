package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Resource represents a shared resource that can be read from or written to.
type Resource struct {
	data string
	mu   sync.RWMutex // Mutex for read-write synchronization
}

// NewResource creates a new instance of Resource.
func NewResource(data string) *Resource {
	return &Resource{data: data}
}

// Read reads data from the resource within a specified timeout.
func (r *Resource) Read(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err() // Return error if context is canceled
	default:
		r.mu.RLock() // Acquire a read lock
		defer r.mu.RUnlock()
		return r.data, nil
	}
}

// Write writes data to the resource within a specified timeout.
func (r *Resource) Write(ctx context.Context, newData string) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // Return error if context is canceled
	default:
		r.mu.Lock() // Acquire a write lock
		defer r.mu.Unlock()
		r.data = newData
		return nil
	}
}

// Worker represents a worker that performs read or write operations on the resource.
type Worker struct {
	ID       int
	Resource *Resource
}

// NewWorker creates a new instance of Worker.
func NewWorker(id int, resource *Resource) *Worker {
	return &Worker{ID: id, Resource: resource}
}

// ReadFromResource reads data from the resource and prints it.
func (w *Worker) ReadFromResource(ctx context.Context) {
	data, err := w.Resource.Read(ctx)
	if err != nil {
		fmt.Printf("Worker %d: Read operation failed: %v\n", w.ID, err)
		return
	}
	fmt.Printf("Worker %d reading from resource: %s\n", w.ID, data)
}

// WriteToResource writes data to the resource.
func (w *Worker) WriteToResource(ctx context.Context, newData string) {
	err := w.Resource.Write(ctx, newData)
	if err != nil {
		fmt.Printf("Worker %d: Write operation failed: %v\n", w.ID, err)
		return
	}
	fmt.Printf("Worker %d writing to resource: %s\n", w.ID, newData)
}

// RunSimulation runs the authentication server simulation with the given number of workers and timeout duration.
func RunSimulation(numWorkers int, timeout time.Duration) {
	// Create a shared resource
	resource := NewResource("initial data")

	// Create a pool of workers
	workers := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = NewWorker(i+1, resource)
	}

	// Set timeout for read and write operations
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Simulate concurrent read and write operations with timeout
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(worker *Worker) {
			defer wg.Done()

			// Perform read operation
			worker.ReadFromResource(ctx)

			// Introduce some delay to simulate real-world scenarios
			time.Sleep(time.Second)

			// Perform write operation
			newData := fmt.Sprintf("new data written by Worker %d", worker.ID)
			worker.WriteToResource(ctx, newData)
		}(workers[i])
	}

	// Wait for all workers to finish
	wg.Wait()

	// Final state of the resource
	data, err := resource.Read(context.Background())
	if err != nil {
		fmt.Printf("Error reading final state of the resource: %v\n", err)
		return
	}
	fmt.Println("Final state of the resource:", data)
}

func main(){
	RunSimulation(3, time.Duration(100))
}
