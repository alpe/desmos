package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/desmos-labs/desmos/v3/x/reports/types"
	subspaceskeeper "github.com/desmos-labs/desmos/v3/x/subspaces/keeper"
	subspacestypes "github.com/desmos-labs/desmos/v3/x/subspaces/types"
)

// RegisterInvariants registers all posts invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper Keeper) {
	ir.RegisterRoute(types.ModuleName, "valid-subspaces",
		ValidSubspacesInvariant(keeper))
	ir.RegisterRoute(types.ModuleName, "valid-reasons",
		ValidReasonsInvariant(keeper))
	ir.RegisterRoute(types.ModuleName, "valid-reports",
		ValidReportsInvariant(keeper))
}

// --------------------------------------------------------------------------------------------------------------------

// ValidSubspacesInvariant checks that all the subspaces have a valid reason id and report id to them
func ValidSubspacesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (message string, broken bool) {
		var invalidSubspaces []subspacestypes.Subspace
		k.sk.IterateSubspaces(ctx, func(index int64, subspace subspacestypes.Subspace) (stop bool) {
			invalid := false

			// Make sure the next reason id exists
			if !k.HasNextReasonID(ctx, subspace.ID) {
				invalid = true
			}

			if !k.HasNextReportID(ctx, subspace.ID) {
				invalid = true
			}

			if invalid {
				invalidSubspaces = append(invalidSubspaces, subspace)
			}

			return false
		})

		return sdk.FormatInvariant(types.ModuleName, "invalid subspaces",
			fmt.Sprintf("the following subspaces are invalid:\n%s", subspaceskeeper.FormatOutputSubspaces(invalidSubspaces)),
		), invalidSubspaces != nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ValidReasonsInvariant checks that all the reasons are valid
func ValidReasonsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidReasons []types.Reason
		k.IterateReasons(ctx, func(index int64, reason types.Reason) (stop bool) {
			invalid := false

			// Make sure the subspace exists
			if !k.HasSubspace(ctx, reason.SubspaceID) {
				invalid = true
			}

			nextReasonID, err := k.GetNextReasonID(ctx, reason.SubspaceID)
			if err != nil {
				invalid = true
			}

			// Make sure the reason id is always less than the next one
			if reason.ID >= nextReasonID {
				invalid = true
			}

			// Validate the reason
			err = reason.Validate()
			if err != nil {
				invalid = true
			}

			if invalid {
				invalidReasons = append(invalidReasons, reason)
			}

			return false

		})

		return sdk.FormatInvariant(types.ModuleName, "invalid reasons",
			fmt.Sprintf("the following reasons are invalid:\n%s", formatOutputReasons(invalidReasons)),
		), invalidReasons != nil
	}
}

// formatOutputReasons concatenates the given reasons information into a string
func formatOutputReasons(reasons []types.Reason) (output string) {
	for _, reason := range reasons {
		output += fmt.Sprintf("SuspaceID: %d, ReasonID: %d\n", reason.SubspaceID, reason.ID)
	}
	return output
}

// --------------------------------------------------------------------------------------------------------------------

// ValidReportsInvariant checks that all the reports are valid
// TODO: Test this once the x/posts gets merged (we need the posts keeper to save posts)
func ValidReportsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var invalidReports []types.Report
		k.IterateReports(ctx, func(index int64, report types.Report) (stop bool) {
			invalid := false

			// Make sure the subspace exists
			if !k.HasSubspace(ctx, report.SubspaceID) {
				invalid = true
			}

			nextReportID, err := k.GetNextReportID(ctx, report.SubspaceID)
			if err != nil {
				invalid = true
			}

			// Make sure the report id is always less than the next one
			if report.ID >= nextReportID {
				invalid = true
			}

			if data, ok := report.Data.GetCachedValue().(*types.PostData); ok {
				// Make sure the reported post exists
				if !k.HasPost(ctx, report.SubspaceID, data.PostID) {
					invalid = true
				}
			}

			// Validate the report
			err = report.Validate()
			if err != nil {
				invalid = true
			}

			if invalid {
				invalidReports = append(invalidReports, report)
			}

			return false

		})

		return sdk.FormatInvariant(types.ModuleName, "invalid reports",
			fmt.Sprintf("the following reasons are invalid:\n%s", formatOutputReports(invalidReports)),
		), invalidReports != nil
	}
}

// formatOutputReports concatenates the given reports information into a string
func formatOutputReports(reports []types.Report) (output string) {
	for _, report := range reports {
		output += fmt.Sprintf("SuspaceID: %d, ReportID: %d\n", report.SubspaceID, report.ID)
	}
	return output
}