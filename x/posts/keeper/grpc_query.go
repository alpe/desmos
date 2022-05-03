package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/desmos-labs/desmos/v3/x/posts/types"
)

var _ types.QueryServer = &Keeper{}

// Posts implements the QueryPosts gRPC method
func (k Keeper) Posts(ctx context.Context, request *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
	if request.SubspaceId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid subspace id")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	postsSubspace := prefix.NewStore(store, types.SubspacePostsPrefix(request.SubspaceId))

	var posts []types.Post
	pageRes, err := query.Paginate(postsSubspace, request.Pagination, func(key []byte, value []byte) error {
		var post types.Post
		if err := k.cdc.Unmarshal(value, &post); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		posts = append(posts, post)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPostsResponse{
		Posts:      posts,
		Pagination: pageRes,
	}, nil
}

// Post implements the QueryPost gRPC method
func (k Keeper) Post(ctx context.Context, request *types.QueryPostRequest) (*types.QueryPostResponse, error) {
	if request.SubspaceId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid subspace id")
	}
	if request.PostId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid post id")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	post, found := k.GetPost(sdkCtx, request.SubspaceId, request.PostId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "post with id %d not found", request.PostId)
	}

	return &types.QueryPostResponse{
		Post: post,
	}, nil
}

// PostAttachments implements the QueryPostAttachments gRPC method
func (k Keeper) PostAttachments(ctx context.Context, request *types.QueryPostAttachmentsRequest) (*types.QueryPostAttachmentsResponse, error) {
	if request.SubspaceId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid subspace id")
	}
	if request.PostId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid post id")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	attachmentsStore := prefix.NewStore(store, types.PostAttachmentsPrefix(request.SubspaceId, request.PostId))

	var attachments []types.Attachment
	pageRes, err := query.Paginate(attachmentsStore, request.Pagination, func(key []byte, value []byte) error {
		var attachment types.Attachment
		if err := k.cdc.Unmarshal(value, &attachment); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		attachments = append(attachments, attachment)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPostAttachmentsResponse{
		Attachments: attachments,
		Pagination:  pageRes,
	}, nil
}

// PollAnswers implements the QueryPollAnswers gRPC method
func (k Keeper) PollAnswers(ctx context.Context, request *types.QueryPollAnswersRequest) (*types.QueryPollAnswersResponse, error) {
	if request.SubspaceId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid subspace id")
	}
	if request.PostId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid post id")
	}
	if request.PollId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid poll id")
	}

	var user sdk.AccAddress
	if request.User != "" {
		userAddr, err := sdk.AccAddressFromBech32(request.User)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid user address: %s", err)
		}
		user = userAddr
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)

	answersPrefix := types.PollAnswersPrefix(request.SubspaceId, request.PostId, request.PollId)
	if user != nil {
		answersPrefix = types.PollAnswerStoreKey(request.SubspaceId, request.PostId, request.PollId, user)
	}
	answersStore := prefix.NewStore(store, answersPrefix)

	var answers []types.UserAnswer
	pageRes, err := query.Paginate(answersStore, request.Pagination, func(key []byte, value []byte) error {
		var answer types.UserAnswer
		if err := k.cdc.Unmarshal(value, &answer); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		answers = append(answers, answer)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPollAnswersResponse{
		Answers:    answers,
		Pagination: pageRes,
	}, nil
}

// PollTallyResults implements the QueryPollTallyResults gRPC method
func (k Keeper) PollTallyResults(ctx context.Context, request *types.QueryPollTallyResultsRequest) (*types.QueryPollTallyResultsResponse, error) {
	if request.SubspaceId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid subspace id")
	}
	if request.PostId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid post id")
	}
	if request.PollId == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid poll id")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	results, found := k.GetPollTallyResults(sdkCtx, request.SubspaceId, request.PostId, request.PollId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "poll with id %d not found", request.PollId)
	}

	return &types.QueryPollTallyResultsResponse{
		Results: results,
	}, nil
}

// Params implements the QueryParams gRPC method
func (k Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params := k.GetParams(sdkCtx)
	return &types.QueryParamsResponse{Params: params}, nil
}
