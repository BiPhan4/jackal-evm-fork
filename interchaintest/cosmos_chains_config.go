package interchaintest

import (
	"encoding/json"
	"fmt"

	testtypes "github.com/JackalLabs/storage-outpost/e2e/interchaintest/types"
	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

var genesisAllowICH = map[string]interface{}{

	"interchainaccounts": map[string]interface{}{

		"controller_genesis_state": map[string]interface{}{
			"active_channels":     []interface{}{},
			"interchain_accounts": []interface{}{},
			"params": map[string]interface{}{
				"controller_enabled": true,
			},
			"ports": []interface{}{},
		},

		"host_genesis_state": map[string]interface{}{
			"active_channels":     []interface{}{},
			"interchain_accounts": []interface{}{},
			"params": map[string]interface{}{
				"allow_messages": []interface{}{"*"},
				"host_enabled":   true,
			},
			"port": "icahost",
		},
	},
	"storage": map[string]interface{}{

		"active_providers_list": []interface{}{},
		"attest_forms":          []interface{}{},
		"collateral_list":       []interface{}{},
		"file_list":             []interface{}{},

		"params": map[string]interface{}{
			"attestFormSize":             "5",
			"attestMinToPass":            "3",
			"check_window":               "100",
			"chunk_size":                 "1024",
			"collateralPrice":            "10000000000",
			"deposit_account":            "jkl12g4qwenvpzqeakavx5adqkw203s629tf6k8vdg", // Let's deposit all storage payments to the Danny user
			"max_contract_age_in_blocks": "100",
			"misses_to_burn":             "3",
			"price_feed":                 "jklprice",
			"price_per_tb_per_month":     "8",
			"proof_window":               "50",
		},

		"payment_info_list": []interface{}{},
		"providers_list":    []interface{}{},
		"report_forms":      []interface{}{},
	},
}

var chainSpecs = []*interchaintest.ChainSpec{

	// -- CANINED --
	{
		ChainConfig: ibc.ChainConfig{
			Type:    "cosmos",
			Name:    "canined",
			ChainID: "puppy-1",
			Images: []ibc.DockerImage{
				{
					Repository: "biphan4/many-bindings", // FOR LOCAL IMAGE USE: Docker Image Name
					Version:    "0.0.7",                 // FOR LOCAL IMAGE USE: Docker Image Tag
				},
			},
			Bin:            "canined",
			Bech32Prefix:   "jkl",
			Denom:          "ujkl", // do we have to use ujkl or is jkl ok?
			GasPrices:      "0.00ujkl",
			GasAdjustment:  1.3,
			EncodingConfig: testtypes.JackaklEncoding(),
			TrustingPeriod: "508h",
			NoHostMount:    false,
			ModifyGenesis:  modifyGenesisAtPath(genesisAllowICH, "app_state"),
			// NOTE: I'm not sure this does anything?
			HostPortOverride: map[int]int{
				16657: 26657, // NOTE: tried switching these around and it still doesn't work
				1317:  1317,
				9090:  9090,
			},
		},
	},
}

func modifyGenesisAtPath(insertedBlock map[string]interface{}, key string) func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(chainConfig ibc.ChainConfig, genbz []byte) ([]byte, error) {
		g := make(map[string]interface{})
		if err := json.Unmarshal(genbz, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}

		//Get the section of the genesis file under the given key (e.g. "app_state")
		app_state, ok := g[key].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("genesis json does not have top level key: %s", key)
		}

		// Replace or add each entry from the insertedBlock into the appState

		for k, v := range insertedBlock {
			app_state[k] = v
		}

		g[key] = app_state
		out, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}
		return out, nil
	}
}
