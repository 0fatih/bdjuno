package types

import (
	"fmt"
	"os"

	"github.com/cometbft/cometbft/libs/log"

	"cosmossdk.io/simapp/params"
	"github.com/forbole/juno/v5/node/remote"

	wasmapp "github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v5/node/local"

	nodeconfig "github.com/forbole/juno/v5/node/config"

	banksource "github.com/forbole/bdjuno/v5/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v5/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v5/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v5/modules/distribution/source"
	remotedistrsource "github.com/forbole/bdjuno/v5/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v5/modules/gov/source"
	localgovsource "github.com/forbole/bdjuno/v5/modules/gov/source/local"
	remotegovsource "github.com/forbole/bdjuno/v5/modules/gov/source/remote"
	mintsource "github.com/forbole/bdjuno/v5/modules/mint/source"
	localmintsource "github.com/forbole/bdjuno/v5/modules/mint/source/local"
	remotemintsource "github.com/forbole/bdjuno/v5/modules/mint/source/remote"
	slashingsource "github.com/forbole/bdjuno/v5/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v5/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v5/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v5/modules/staking/source"
	localstakingsource "github.com/forbole/bdjuno/v5/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v5/modules/staking/source/remote"

	wasmsource "github.com/forbole/bdjuno/v5/modules/wasm/source"
	localwasmsource "github.com/forbole/bdjuno/v5/modules/wasm/source/local"
	remotewasmsource "github.com/forbole/bdjuno/v5/modules/wasm/source/remote"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MintSource     mintsource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
	WasmSource     wasmsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := wasmapp.NewWasmApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, wasm.EnableAllProposals, nil, nil)

	// app := simapp.NewSimApp(
	// 	log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, nil, nil,
	// )

	sources := &Sources{
		BankSource: localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		// DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govtypesv1.QueryServer(app.GovKeeper), nil),
		MintSource:     localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
		WasmSource:     localwasmsource.NewSource(source, wasmkeeper.Querier(&app.WasmKeeper)),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(app, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(app, "memKeys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn), govtypesv1beta1.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		WasmSource:     remotewasmsource.NewSource(source, wasmtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
