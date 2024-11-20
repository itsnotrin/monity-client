package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func main() {
	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(403).JSON(fiber.Map{"message": "You've reached the rate limit. Please be patient!", "extra_details": fiber.Map{"IP:": c.IP(), "timeLeft": c.GetRespHeader("Retry-After")}})
		},
	}))

	app.Get("/mem", func(c *fiber.Ctx) error {
		v, _ := mem.VirtualMemory()
		swap, _ := mem.SwapMemory()
		var usedPercentage = toFixed(v.UsedPercent, 1)
		var usedRam = v.Used / 1000 / 1000 // Calculates it approximately to gigabytes from bytes.
		return c.JSON(fiber.Map{
			"used_percent": usedPercentage,
			"ram_usage_mb": usedRam,
			"total_mb":     v.Total / 1000 / 1000,
			"free_mb":      v.Free / 1000 / 1000,
			"cached_mb":    v.Cached / 1000 / 1000,
			"swap_used_mb": swap.Used / 1000 / 1000,
		})
	})

	app.Get("/disk", func(c *fiber.Ctx) error {
		partitions, _ := disk.Partitions(false)
		var diskInfo []fiber.Map

		for _, partition := range partitions {
			usage, _ := disk.Usage(partition.Mountpoint)
			diskInfo = append(diskInfo, fiber.Map{
				"device":       partition.Device,
				"mountpoint":   partition.Mountpoint,
				"used_percent": toFixed(usage.UsedPercent, 1),
				"total_gb":     usage.Total / 1024 / 1024 / 1024,
				"used_gb":      usage.Used / 1024 / 1024 / 1024,
				"available_gb": usage.Free / 1024 / 1024 / 1024,
				"fs_type":      partition.Fstype,
				"inodes_used":  usage.InodesUsed,
				"inodes_free":  usage.InodesFree,
			})
		}
		return c.JSON(fiber.Map{"disks": diskInfo})
	})

	app.Get("/cpu/percent", func(c *fiber.Ctx) error {
		v, _ := cpu.Percent(0, false)
		var usedPercentage = v[0]
		var cleanUsedPercentage = toFixed(v[0], 1)

		switch c.Get("detailed") {
		case "0":
			return c.JSON(fiber.Map{"usage_percentage": cleanUsedPercentage})
		case "1":
			return c.JSON(fiber.Map{"detailed_usage_percentage": usedPercentage})
		case "2":
			return c.JSON(fiber.Map{"detailed_usage_percentage": usedPercentage, "usage_percentage": cleanUsedPercentage})
		default:
			return c.JSON(fiber.Map{"detailed_usage_percentage": usedPercentage, "usage_percentage": cleanUsedPercentage})
		}
	})

	app.Get("/cpu/info", func(c *fiber.Ctx) error {
		info, _ := cpu.Info()
		return c.JSON(fiber.Map{
			"cores":      runtime.NumCPU(),
			"model":      info[0].ModelName,
			"mhz":        info[0].Mhz,
			"vendor":     info[0].VendorID,
			"cache_size": info[0].CacheSize,
			"microcode":  info[0].Microcode,
		})
	})

	app.Get("/system", func(c *fiber.Ctx) error {
		info, _ := host.Info()
		uptime := time.Duration(info.Uptime) * time.Second

		return c.JSON(fiber.Map{
			"hostname":       info.Hostname,
			"os":             info.OS,
			"platform":       info.Platform,
			"kernel":         info.KernelVersion,
			"uptime":         uptime.String(),
			"host_id":        info.HostID,
			"procs":          info.Procs,
			"os_version":     info.PlatformVersion,
			"host_arch":      runtime.GOARCH,
			"go_version":     runtime.Version(),
			"num_goroutines": runtime.NumGoroutine(),
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	app.Get("/load", func(c *fiber.Ctx) error {
		avgStat, _ := cpu.Times(false)
		return c.JSON(fiber.Map{
			"user":       toFixed(avgStat[0].User, 2),
			"system":     toFixed(avgStat[0].System, 2),
			"idle":       toFixed(avgStat[0].Idle, 2),
			"nice":       toFixed(avgStat[0].Nice, 2),
			"iowait":     toFixed(avgStat[0].Iowait, 2),
			"irq":        toFixed(avgStat[0].Irq, 2),
			"softirq":    toFixed(avgStat[0].Softirq, 2),
			"steal":      toFixed(avgStat[0].Steal, 2),
			"guest":      toFixed(avgStat[0].Guest, 2),
			"guest_nice": toFixed(avgStat[0].GuestNice, 2),
		})
	})

	app.Get("/metrics", func(c *fiber.Ctx) error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		return c.JSON(fiber.Map{
			"alloc":         m.Alloc,
			"total_alloc":   m.TotalAlloc,
			"sys":           m.Sys,
			"num_gc":        m.NumGC,
			"goroutines":    runtime.NumGoroutine(),
			"gc_pause_ns":   m.PauseTotalNs,
			"heap_objects":  m.HeapObjects,
			"heap_alloc":    m.HeapAlloc,
			"heap_sys":      m.HeapSys,
			"heap_idle":     m.HeapIdle,
			"heap_released": m.HeapReleased,
			"stack_in_use":  m.StackInuse,
			"stack_sys":     m.StackSys,
			"mspan_in_use":  m.MSpanInuse,
			"mspan_sys":     m.MSpanSys,
			"mcache_in_use": m.MCacheInuse,
			"mcache_sys":    m.MCacheSys,
		})
	})

	err := app.Listen(":3000")
	if err != nil {
		fmt.Println(err)
		return
	}
}
