package wasm_test

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/v3/app"
	profilestypes "github.com/desmos-labs/desmos/v3/x/profiles/types"
	"github.com/desmos-labs/desmos/v3/x/subspaces/wasm"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

func marshalMsg(msg proto.Marshaler) []byte {
	bzMsg, _ := msg.Marshal()
	return bzMsg
}

func TestMsgsParser_ParseCustomMsgs(t *testing.T) {
	cdc, _ := app.MakeCodecs()
	parser := wasm.NewWasmMsgParser(cdc)
	contractAddr, err := sdk.AccAddressFromBech32("cosmos14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s4hmalr")
	require.NoError(t, err)

	testCases := []struct {
		name      string
		msg       json.RawMessage
		shouldErr bool
		expMsgs   []sdk.Msg
	}{
		{
			name:      "Parse wrong module message returns error",
			msg:       json.RawMessage(marshalMsg(profilestypes.NewMsgDeleteProfile(""))),
			shouldErr: true,
			expMsgs:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msgs, err := parser.ParseCustomMsgs(contractAddr, tc.msg)
			if tc.shouldErr {
				require.Error(t, err)
			}
			require.Equal(t, tc.expMsgs, msgs)
		})
	}
}
