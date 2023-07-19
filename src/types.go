package main

type ABCIInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Response struct {
			Data             string `json:"data"`
			Version          string `json:"version"`
			LastBlockHeight  string `json:"last_block_height"`
			LastBlockAppHash string `json:"last_block_app_hash"`
		} `json:"response"`
	} `json:"result"`
}

type UpgradeInfo struct {
	Network       string `json:"network"`
	RPC           string `json:"rpc"`
	UpgradeHeight uint64 `json:"upgrade_height"`
	CheckSeconds  uint64 `json:"check_seconds"`
	Webhook       string `json:"webhook"`
}
