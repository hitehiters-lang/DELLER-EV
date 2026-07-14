package blog_gen

import (
	"fmt"
	"time"
)

func Start() {
	fmt.Println("[BLOG_GEN] started background scheduler")

	if !alreadyGeneratedToday() {
		fmt.Println("[BLOG_GEN] no post for today, generating on startup...")
		go func() {
			if err := generatePost(); err != nil {
				fmt.Println("[BLOG_GEN] startup generation failed:", err)
			} else {
				fmt.Println("[BLOG_GEN] startup generation completed successfully")
			}
		}()
	}

	go run()
}

func run() {
	ticker := time.NewTicker(3 * time.Hour)
	defer ticker.Stop()

	for {
		check()
		<-ticker.C
	}
}

func check() {
	if !canGenerate() {
		return
	}

	fmt.Println("[BLOG_GEN] starting generation...")

	if err := generatePost(); err != nil {
		fmt.Println("[BLOG_GEN] generation failed:", err)
		return
	}

	fmt.Println("[BLOG_GEN] generation completed successfully")
}
