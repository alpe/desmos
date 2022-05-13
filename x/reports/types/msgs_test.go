package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/desmos/v3/x/reports/types"
)

var msgCreateReport = types.NewMsgCreateReport(
	1,
	1,
	"This post is spam",
	"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
	types.NewPostData(1),
)

func TestMsgCreateReport_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgCreateReport.Route())
}

func TestMsgCreateReport_Type(t *testing.T) {
	require.Equal(t, types.ActionCreateReport, msgCreateReport.Type())
}

func TestMsgCreateReport_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgCreateReport
		shouldErr bool
	}{
		{
			name: "invalid subspace id returns error",
			msg: types.NewMsgCreateReport(
				0,
				1,
				"This post is spam",
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				types.NewPostData(1),
			),
			shouldErr: true,
		},
		{
			name: "invalid reason id returns error",
			msg: types.NewMsgCreateReport(
				1,
				0,
				"This post is spam",
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				types.NewPostData(1),
			),
			shouldErr: true,
		},
		{
			name: "invalid reporter returns error",
			msg: types.NewMsgCreateReport(
				1,
				1,
				"This post is spam",
				"",
				types.NewPostData(1),
			),
			shouldErr: true,
		},
		{
			name: "invalid report data returns error",
			msg: types.NewMsgCreateReport(
				1,
				1,
				"This post is spam",
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
				types.NewPostData(0),
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgCreateReport,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgCreateReport_GetSignBytes(t *testing.T) {
	expected := `{"type":"desmos/MsgCreateReport","value":{"data":{"type":"desmos/PostData","value":{"post_id":"1"}},"message":"This post is spam","reason_id":1,"reporter":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","subspace_id":"1"}}`
	require.Equal(t, expected, string(msgCreateReport.GetSignBytes()))
}

func TestMsgCreateReport_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgCreateReport.Reporter)
	require.Equal(t, []sdk.AccAddress{addr}, msgCreateReport.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgDeleteReport = types.NewMsgDeleteReport(
	1,
	1,
	"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
)

func TestMsgDeleteReport_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgDeleteReport.Route())
}

func TestMsgDeleteReport_Type(t *testing.T) {
	require.Equal(t, types.ActionDeleteReport, msgDeleteReport.Type())
}

func TestMsgDeleteReport_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgDeleteReport
		shouldErr bool
	}{
		{
			name: "invalid subspace id returns error",
			msg: types.NewMsgDeleteReport(
				0,
				1,
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid report id returns error",
			msg: types.NewMsgDeleteReport(
				1,
				0,
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid signer returns error",
			msg: types.NewMsgDeleteReport(
				1,
				1,
				"",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgDeleteReport,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgDeleteReport_GetSignBytes(t *testing.T) {
	expected := `{"type":"desmos/MsgDeleteReport","value":{"report_id":"1","signer":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","subspace_id":"1"}}`
	require.Equal(t, expected, string(msgDeleteReport.GetSignBytes()))
}

func TestMsgDeleteReport_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgDeleteReport.Signer)
	require.Equal(t, []sdk.AccAddress{addr}, msgDeleteReport.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgSupportStandardReason = types.NewMsgSupportStandardReason(
	1,
	1,
	"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
)

func TestMsgSupportStandardReason_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgSupportStandardReason.Route())
}

func TestMsgSupportStandardReason_Type(t *testing.T) {
	require.Equal(t, types.ActionSupportReason, msgSupportStandardReason.Type())
}

func TestMsgSupportStandardReason_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgSupportStandardReason
		shouldErr bool
	}{
		{
			name: "invalid subspace id returns error",
			msg: types.NewMsgSupportStandardReason(
				0,
				1,
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid reason id returns error",
			msg: types.NewMsgSupportStandardReason(
				1,
				0,
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid reporter returns error",
			msg: types.NewMsgSupportStandardReason(
				1,
				1,
				"",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgSupportStandardReason,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgSupportStandardReason_GetSignBytes(t *testing.T) {
	expected := `{"type":"desmos/MsgSupportStandardReason","value":{"signer":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","standard_reason_id":1,"subspace_id":"1"}}`
	require.Equal(t, expected, string(msgSupportStandardReason.GetSignBytes()))
}

func TestMsgSupportStandardReason_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgSupportStandardReason.Signer)
	require.Equal(t, []sdk.AccAddress{addr}, msgSupportStandardReason.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgAddReason = types.NewMsgAddReason(
	1,
	"Spam",
	"This post is spam or the user is a spammer",
	"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
)

func TestMsgAddReason_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgAddReason.Route())
}

func TestMsgAddReason_Type(t *testing.T) {
	require.Equal(t, types.ActionAddReason, msgAddReason.Type())
}

func TestMsgAddReason_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgAddReason
		shouldErr bool
	}{
		{
			name: "invalid subspace id returns error",
			msg: types.NewMsgAddReason(
				0,
				"Spam",
				"This post is spam or the user is a spammer",
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid title returns error",
			msg: types.NewMsgAddReason(
				1,
				"",
				"This post is spam or the user is a spammer",
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid signer returns error",
			msg: types.NewMsgAddReason(
				1,
				"Spam",
				"This post is spam or the user is a spammer",
				"",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgAddReason,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgAddReason_GetSignBytes(t *testing.T) {
	expected := `{"type":"desmos/MsgAddReason","value":{"description":"This post is spam or the user is a spammer","signer":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","subspace_id":"1","title":"Spam"}}`
	require.Equal(t, expected, string(msgAddReason.GetSignBytes()))
}

func TestMsgAddReason_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgAddReason.Signer)
	require.Equal(t, []sdk.AccAddress{addr}, msgAddReason.GetSigners())
}

// --------------------------------------------------------------------------------------------------------------------

var msgRemoveReason = types.NewMsgRemoveReason(
	1,
	1,
	"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
)

func TestMsgRemoveReason_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgRemoveReason.Route())
}

func TestMsgRemoveReason_Type(t *testing.T) {
	require.Equal(t, types.ActionRemoveReason, msgRemoveReason.Type())
}

func TestMsgRemoveReason_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		msg       *types.MsgRemoveReason
		shouldErr bool
	}{
		{
			name: "invalid subspace id returns error",
			msg: types.NewMsgRemoveReason(
				0,
				1,
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid reason id returns error",
			msg: types.NewMsgRemoveReason(
				1,
				0,
				"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			),
			shouldErr: true,
		},
		{
			name: "invalid signer returns error",
			msg: types.NewMsgRemoveReason(
				1,
				1,
				"",
			),
			shouldErr: true,
		},
		{
			name:      "valid message returns no error",
			msg:       msgRemoveReason,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgRemoveReason_GetSignBytes(t *testing.T) {
	expected := `{"type":"desmos/MsgRemoveReason","value":{"reason_id":1,"signer":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","subspace_id":"1"}}`
	require.Equal(t, expected, string(msgRemoveReason.GetSignBytes()))
}

func TestMsgRemoveReason_GetSigners(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32(msgRemoveReason.Signer)
	require.Equal(t, []sdk.AccAddress{addr}, msgRemoveReason.GetSigners())
}