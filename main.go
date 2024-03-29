package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
			return c.Status(403).JSON(fiber.Map{"message": "You've reached the rate limit. Please be patient!", "extra_details": fiber.Map{"IP:": c.IP(), "timeLeft": c.GetRespHeader("Retry-After")}})
		},
	}))

	app.Get("/mem", func(c *fiber.Ctx) error {
		v, _ := mem.VirtualMemory() // Use mem from the import.
		var usedPercentage = toFixed(v.UsedPercent, 1)
		var usedRam = v.Used / 1000 / 1000 // Calculates it approximately to gigabytes from bytes.
		return c.JSON(fiber.Map{"used_percent": usedPercentage, "ram_usage_mb": usedRam})
	})

	app.Get("/cpu/percent", func(c *fiber.Ctx) error {
		v, _ := cpu.Percent(0, false) // CPU Usage
		var usedPercentage = v[0]
		var cleanUsedPercentage = toFixed(v[0], 1)

		//TODO: Refactor as Switch Case
		if c.Get("detailed") == "0" {
			return c.JSON(fiber.Map{"usage_percentage": cleanUsedPercentage})
		}
		if c.Get("detailed") == "1" {
			return c.JSON(fiber.Map{"detailed_usage_percentage": usedPercentage})
		}
		if c.Get("detailed") == "2" {
			return c.JSON(fiber.Map{"detailed_usage_percentage": usedPercentage, "usage_percentage": cleanUsedPercentage})
		}
		// If no option is passed, just send both regardless.
		return c.JSON(fiber.Map{"detailed_usage_percentage": usedPercentage, "usage_percentage": cleanUsedPercentage})

	})

	err := app.Listen(":3000")
	if err != nil {
		fmt.Println(err)
		return
	}
}
