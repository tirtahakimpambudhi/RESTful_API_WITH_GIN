// Import necessary packages
package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"go_gin/internal/domain/model"
)


// ShutDown initiates a graceful shutdown with a timeout and a map of cleanup operations
func ShutDown(ctx context.Context, timeout time.Duration, ops map[string]model.Operation) <-chan struct{} {
	// Create a channel to signal the completion of shutdown
	wait := make(chan struct{})

	// Goroutine for handling OS signals
	go func() {
		// Create a channel for OS signals
		s := make(chan os.Signal, 1)

		// Notify the channel for specified signals (SIGINT, SIGTERM, SIGHUP)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		// Wait for a signal
		<-s

		// Log that shutdown process has started
		log.Println("shutting down")

		// Create a new context with timeout
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Use a WaitGroup to wait for all cleanup operations to finish
		var wg sync.WaitGroup

		// Iterate through the map of cleanup operations
		for key, op := range ops {
			// Increment the WaitGroup counter
			wg.Add(1)

			// Capture variables for the goroutine
			innerOp := op
			innerKey := key

			// Goroutine for executing cleanup operation
			go func() {
				// Decrement the WaitGroup counter when the goroutine completes
				defer wg.Done()

				// Log that cleanup is in progress
				log.Printf("cleaning up: %s", innerKey)

				// Execute the cleanup operation with the specified context
				if err := innerOp(ctx); err != nil {
					// Log if cleanup operation fails
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				// Log that cleanup operation was successful
				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		// Wait for all cleanup operations to finish
		wg.Wait()

		// Close the channel to signal the completion of shutdown
		close(wait)
	}()

	// Return the channel for signaling completion
	return wait
}