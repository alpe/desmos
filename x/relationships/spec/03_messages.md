<!--
order: 3
-->

# Msg Service

## Msg/CreateRelationship
A new relationship can be created with the `MsgCreateRelationship`, which allows to specify the subspace inside which the relationship should live and the counterparty address.

+++ https://github.com/desmos-labs/desmos/blob/v3.0.0/proto/desmos/relationships/v1/msg_server.proto#L27-L39

It's expected to fail if a relationships between the same user and counterparty already exists inside the same subspace. 

## Msg/DeleteRelationship
An existing relationship can be deleted with the `MsgDeleteRelationship`. 

+++ https://github.com/desmos-labs/desmos/blob/v3.0.0/proto/desmos/relationships/v1/msg_server.proto#L47-L56

It's expected to fail if a relationships between the signer and counterparty does not exist inside the specified subspace.

## Msg/BlockUser
A new user block can be created with the `MsgBlockUser`, which allows to specify the subspace inside which the block should be valid, the address of the blocked user and an optional reason for the block.

+++ https://github.com/desmos-labs/desmos/blob/v3.0.0/proto/desmos/relationships/v1/msg_server.proto#L64-L74

It's expected to fail if a user block between the same user and blocker already exists inside the same subspace.

## Msg/UnblockUser
An existing user block can be deleted with the `MsgUnblockUser`.

+++ https://github.com/desmos-labs/desmos/blob/v3.0.0/proto/desmos/relationships/v1/msg_server.proto#L81-L89

It's expected to fail if the user block does not exist inside the specified subspace.
