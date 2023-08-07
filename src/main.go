package main

// An application which waits for an upgrade to process.
// Once the blocks increase past the upgrade, it sends a discord webhook notification. 'BLOCKS'

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	upgrades, err := readUpgradesFromFile("secret.json")
	if err != nil {
		fmt.Println("Error reading upgrades: ", err)
		return
	}

	for _, upgrade := range upgrades {
		u := upgrade

		// TODO: on a per upgrade basis
		var startTime time.Time

		go func() {
			for {
				cb, err := getCurrentBlock(u.RPC)
				if err != nil {
					fmt.Println("Error getting current block: ", err)
					sleep(u.CheckSeconds)
					continue
				}

				height, err := strconv.ParseUint(cb.Result.Response.LastBlockHeight, 10, 64)
				if err != nil {
					fmt.Println("Error parsing height: ", err)
					sleep(u.CheckSeconds)
					continue
				}

				if height-1 >= u.UpgradeHeight {
					if height > u.UpgradeHeight+10 {
						fmt.Printf("Too late to start for %v %d\n", u.Network, height)
						time.Sleep(5 * time.Minute)
						continue
					}

					if startTime.IsZero() {
						fmt.Println("upgrade height reached! Waiting")
						startTime = time.Now()

						NewDiscordTimeToUpgrade(u.Webhook, u.Network, height)
					}

					if height >= u.UpgradeHeight+1 {
						fmt.Println("BLOCKS")
						diff := time.Since(startTime)

						NewDiscordBlocks(u.Webhook, u.Network, diff, height)
						break
					}
				} else {
					fmt.Printf("Blocks until upgrade: %d (%v)\n", (u.UpgradeHeight - height), height)
				}

				sleep(u.CheckSeconds)
			}
		}()
	}

	// wait forever
	select {}
}

func readUpgradesFromFile(filename string) ([]UpgradeInfo, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return nil, err
	}

	file := filepath.Join(currentDir, filename)

	upgrades := make([]UpgradeInfo, 0)

	if data, err := os.ReadFile(file); err != nil {
		fmt.Println("Error reading upgrade.json: ", err)
		return upgrades, err
	} else {
		if err := json.Unmarshal(data, &upgrades); err != nil {
			fmt.Println("Error unmarshalling upgrade.json: ", err)
			return upgrades, err
		}
	}
	return upgrades, nil
}

func getCurrentBlock(url string) (ABCIInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return ABCIInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ABCIInfo{}, fmt.Errorf("HTTP request failed with status code: %d (%s)", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ABCIInfo{}, err
	}

	var data ABCIInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return ABCIInfo{}, err
	}

	return data, nil
}

func sleep(seconds uint64) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
