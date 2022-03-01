# ADR 008: Fee modules

## Changelog

- March 1, 2022: Initial draft;

## Status
DRAFTED

## Abstract
This ADR defines the `x/fees` module which enables the possibility to set custom fees for each desmos custom message
through governance proposals.

## Context
In order to better prevent any kind of spam, from smart contracts implementation to subspaces creation and, in the near future,
posts, it is useful to have a system with which the community can tweak the costs of desmos operations.

## Decision
We will create a module named `fees` which give desmos community the possibility to add minimum fees to any Desmos custom
message through governance.

### Types
Minimum (or custom) fees will contains the following elements:  

* a string identifying the message type delivered with the transaction;
* the amount of fees associated with the give message.

#### Min Fees
```go
type MinFee struct {
	// The message type identifying the Desmos message
	messageType string
	
	// The amount of fees to be paid for the message
	amount sdk.Coins
}
```

## Consequences

> This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

### Backwards Compatibility

> All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

{positive consequences}

### Negative

{negative consequences}

### Neutral

{neutral consequences}

## Further Discussions

While an ADR is in the DRAFT or PROPOSED stage, this section should contain a summary of issues to be solved in future iterations (usually referencing comments from a pull-request discussion).
Later, this section can optionally list ideas or improvements the author or reviewers found during the analysis of this ADR.

## Test Cases [optional]

Test cases for an implementation are mandatory for ADRs that are affecting consensus changes. Other ADRs can choose to include links to test cases if applicable.

## References

- {reference link}