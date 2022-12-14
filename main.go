package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		// How the Limit works:
		// weightOfPreviousWindow = previous window's amount request * (whenNewWindow / Expiration)
		// rate = weightOfPreviousWindow + current window's amount request.
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(403).JSON(fiber.Map{"message": "You've reached the rate limit. Please be patient!", "extra_details": fiber.Map{"IP:": c.IP(), "timeleft": c.GetRespHeader("Retry-After")}})
		},
	}))

	app.Get("/mem", func(c *fiber.Ctx) error {
		v, _ := mem.VirtualMemory()
		return c.JSON(fiber.Map{"used_percent": v.UsedPercent})
	})

	app.Listen(":3000")
}
