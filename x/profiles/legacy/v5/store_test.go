package v5_test

import (
	"encoding/hex"
	"testing"
	"time"

	v5 "github.com/desmos-labs/desmos/v3/x/profiles/legacy/v5"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	"github.com/desmos-labs/desmos/v3/testutil"
	profilestypes "github.com/desmos-labs/desmos/v3/x/profiles/types"

	"github.com/cosmos/cosmos-sdk/store"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/desmos/v3/app"
	"github.com/desmos-labs/desmos/v3/x/relationships/types"
)

func buildContext(
	keys map[string]*sdk.KVStoreKey, tKeys map[string]*sdk.TransientStoreKey, memKeys map[string]*sdk.MemoryStoreKey,
) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	}
	for _, tKey := range tKeys {
		cms.MountStoreWithDB(tKey, sdk.StoreTypeTransient, db)
	}
	for _, memKey := range memKeys {
		cms.MountStoreWithDB(memKey, sdk.StoreTypeMemory, nil)
	}

	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	return sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())
}

func TestMigrateStore(t *testing.T) {
	cdc, legacyAmino := app.MakeCodecs()

	// Build all the necessary keys
	keys := sdk.NewKVStoreKeys(authtypes.StoreKey, types.StoreKey)
	tKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	account := testutil.GetChainLinkAccount("cosmos", "cosmos")
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "application links owners are added properly",
			store: func(ctx sdk.Context) {
				kvStore := ctx.KVStore(keys[types.StoreKey])

				// Store an application link
				linkKey := profilestypes.UserApplicationLinkKey("cosmos10nsdxxdvy9qka3zv0lzw8z9cnu6kanld8jh773", "twitter", "twitteruser")
				kvStore.Set(
					linkKey,
					cdc.MustMarshal(&profilestypes.ApplicationLink{
						User:  "cosmos10nsdxxdvy9qka3zv0lzw8z9cnu6kanld8jh773",
						Data:  profilestypes.NewData("twitter", "twitteruser"),
						State: profilestypes.ApplicationLinkStateInitialized,
						OracleRequest: profilestypes.NewOracleRequest(
							0,
							1,
							profilestypes.NewOracleRequestCallData("twitter", "calldata"),
							"client_id",
						),
						Result:       nil,
						CreationTime: time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
					}),
				)

				// Store an application link client id
				kvStore.Set(profilestypes.ApplicationLinkClientIDKey("client_id"), linkKey)
			},
			check: func(ctx sdk.Context) {
				kvStore := ctx.KVStore(keys[types.StoreKey])

				// Check the application link owner
				key := profilestypes.ApplicationLinkOwnerKey(
					"twitter",
					"twitteruser",
					"cosmos10nsdxxdvy9qka3zv0lzw8z9cnu6kanld8jh773",
				)
				require.Equal(t, []byte{0x01}, kvStore.Get(key))
			},
		},
		{
			name: "invalid chain links are deleted",
			store: func(ctx sdk.Context) {
				kvStore := ctx.KVStore(keys[types.StoreKey])

				// Store the chain link
				signatureValue := []byte("custom value")
				signature := profilestypes.SingleSignatureData{
					Mode:      signing.SignMode_SIGN_MODE_TEXTUAL,
					Signature: signatureValue,
				}
				signatureAny := testutil.NewAny(&signature)

				chainLink := profilestypes.NewChainLink(
					"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
					profilestypes.NewBech32Address(account.Bech32Address().GetValue(), "cosmos"),
					profilestypes.Proof{
						PubKey:    account.PubKeyAny(),
						Signature: signatureAny,
						PlainText: hex.EncodeToString(signatureValue),
					},
					profilestypes.ChainConfig{Name: "cosmos"},
					time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
				)

				kvStore.Set(
					profilestypes.ChainLinksStoreKey(
						"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
						"cosmos",
						"cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
					),
					cdc.MustMarshal(&chainLink),
				)
			},
			check: func(ctx sdk.Context) {
				kvStore := ctx.KVStore(keys[types.StoreKey])

				// Make sure the chain link is deleted and the owner key is not added
				require.False(t, kvStore.Has(profilestypes.ChainLinkOwnerKey(
					"cosmos",
					account.Bech32Address().GetValue(),
					"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				)))
				require.False(t, kvStore.Has(profilestypes.ChainLinksStoreKey(
					"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
					"cosmos",
					account.Bech32Address().GetValue(),
				)))
			},
		},
		{
			name: "valid chain link owners are added properly",
			store: func(ctx sdk.Context) {
				kvStore := ctx.KVStore(keys[types.StoreKey])

				// Store the chain link
				chainLink := account.GetBech32ChainLink(
					"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
					time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
				)

				kvStore.Set(
					profilestypes.ChainLinksStoreKey(
						"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
						"cosmos",
						"cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
					),
					cdc.MustMarshal(&chainLink),
				)
			},
			check: func(ctx sdk.Context) {
				kvStore := ctx.KVStore(keys[types.StoreKey])

				key := profilestypes.ChainLinkOwnerKey(
					"cosmos",
					account.Bech32Address().GetValue(),
					"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				)
				require.Equal(t, []byte{0x01}, kvStore.Get(key))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := buildContext(keys, tKeys, memKeys)
			if tc.store != nil {
				tc.store(ctx)
			}

			err := v5.MigrateStore(ctx, keys[types.StoreKey], cdc, legacyAmino)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}