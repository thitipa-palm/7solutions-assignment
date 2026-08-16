package service

import (
	"context"
	"log"
	"time"
)

func RunUserCountMonitor(
	ctx context.Context,
	userService *UserService,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("user count monitor stopped")
			return

		case <-ticker.C:
			count, err := userService.Count(ctx)
			if err != nil {
				log.Printf("count users: %v", err)
				continue
			}

			log.Printf("total users: %d", count)
		}
	}
}
