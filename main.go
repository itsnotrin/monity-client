package main

import (
	"time"
	"math"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/shirou/gopsutil/v3/mem"
)

func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}



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
		v, _ := mem.VirtualMemory() // Use mem from the import.
        var used_percentage = toFixed(v.UsedPercent, 1)
        var usedRam = v.Used/1000/1000 // Calculates it approximately to gigabytes from bytes.
        return c.JSON(fiber.Map{"used_percent": used_percentage, "ram_usage_mb": usedRam})
	})

	app.Listen(":3000")
}
