package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

// query and return server
func queryHostPort(query model.Query) (server string) {
	fmt.Println(query.Head)
	fmt.Println("example:", query.Example)
	fmt.Println("default value:", query.Default)
	fmt.Printf("enter your value[default if blank] ->")
	scanner.Scan()
	server = strings.TrimSpace(scanner.Text())

	if server == "" {
		server = query.Default
	}
	fmt.Println()
	return
}

// set timeout with query input like "enter the Install timeout for client:"
// return the timeOut which is time.Duration
func queryTimeout(query model.Query) (timeOut time.Duration) {
	timeString := ""

	fmt.Println(query.Head)
	fmt.Println("example:", query.Example)
	fmt.Println("default value:", query.Default)
	fmt.Printf("enter your value[default if blank] ->")
	scanner.Scan()
	timeString = strings.TrimSpace(scanner.Text())

	if timeString == "" {
		timeString = query.Default
	}

	timeOut, err := time.ParseDuration(timeString)
	if err != nil {
		fmt.Printf("error when parse %s", timeString)
		fmt.Println(err)
		timeOut = time.Second * 0
	}

	fmt.Println()
	return
}

// set Servertokens with query like ">>set the read tokens for server<<"
func queryServerToken(query model.Query) (tokens []string) {
	var tokenNum int
	var tempToken string

	fmt.Println(query.Head)
	fmt.Println("example:", query.Example)
	fmt.Println("default value:", query.Default)

	fmt.Printf("enter the number of tocken you want to add:")
	scanner.Scan()
	numStr := strings.TrimSpace(scanner.Text())

	if numStr == "" {
		tokens = []string{}
		return
	}

	_, err := fmt.Sscanf(numStr, "%d", &tokenNum)
	if err != nil || tokenNum < 0 {
		fmt.Println("invalid number, skipping token configuration")
		fmt.Println()
		return []string{}
	}
	fmt.Println("setting your value[generated randomly if blank]")
	for i := 0; i < tokenNum; i++ {
		fmt.Printf("enter the %d th token ->", i+1)
		scanner.Scan()
		tempToken = strings.TrimSpace(scanner.Text())

		if tempToken == "" {
			tempToken = generateRandomToken(32)
			fmt.Printf("generated random token: %s\n", tempToken)
		}

		tokens = append(tokens, tempToken)
	}
	fmt.Println()
	return
}

func queryClientToken(query model.Query) (token string) {
	fmt.Println(query.Head)
	fmt.Println("example:", query.Example)
	fmt.Println("default value:", query.Default)

	fmt.Printf("enter your token [press Enter to generate randomly] -> ")
	scanner.Scan()
	token = strings.TrimSpace(scanner.Text())

	if token == "" {
		token = generateRandomToken(32)
		fmt.Printf("generated random token: %s\n", token)
	}

	fmt.Println()
	return
}

// generateRandomToken generate random token with given length
func generateRandomToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Printf("warning: failed to generate secure random token: %v\n", err)
		fmt.Println("using fallback random generation")
		return generateFallbackToken()
	}
	return hex.EncodeToString(bytes)
}

// generateFallbackToken generate token by time stape
func generateFallbackToken() string {
	bytes := make([]byte, 32)
	for i := range bytes {
		bytes[i] = byte(time.Now().UnixNano() >> (i * 8))
	}
	return hex.EncodeToString(bytes)
}
