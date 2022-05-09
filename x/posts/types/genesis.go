package types

import "fmt"

// NewSubspaceDataEntry returns a new SubspaceDataEntry instance
func NewSubspaceDataEntry(subspaceID uint64, initialPostID uint64) SubspaceDataEntry {
	return SubspaceDataEntry{
		SubspaceID:    subspaceID,
		InitialPostID: initialPostID,
	}
}

// Validate returns an error if something is wrong within the entry data
func (e SubspaceDataEntry) Validate() error {
	if e.SubspaceID == 0 {
		return fmt.Errorf("invalid subspace id: %d", e.SubspaceID)
	}

	if e.InitialPostID == 0 {
		return fmt.Errorf("invalid initial post id: %d", e.InitialPostID)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewGenesisPost returns a new GenesisPost instance
func NewGenesisPost(initialAttachmentID uint32, post Post) GenesisPost {
	return GenesisPost{
		Post:                post,
		InitialAttachmentID: initialAttachmentID,
	}
}

// Validate returns an error if something is wrong within the entry data
func (p GenesisPost) Validate() error {
	if p.InitialAttachmentID == 0 {
		return fmt.Errorf("invalid initial attachment id: %d", p.InitialAttachmentID)
	}

	return p.Post.Validate()
}

// --------------------------------------------------------------------------------------------------------------------

// NewGenesisState returns a new GenesisState instance
func NewGenesisState(
	subspacesData []SubspaceDataEntry,
	posts []GenesisPost,
	attachments []Attachment,
	userAnswers []UserAnswer,
	params Params,
) *GenesisState {
	return &GenesisState{
		SubspacesData: subspacesData,
		GenesisPosts:  posts,
		Attachments:   attachments,
		UserAnswers:   userAnswers,
		Params:        params,
	}
}

// DefaultGenesisState returns a default GenesisState
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil, nil, nil, nil, DefaultParams())
}

// getInitialPostID returns the initial post id for the given subspace, 0 if not found
func (e *GenesisState) getInitialPostID(subspaceID uint64) uint64 {
	for _, entry := range e.SubspacesData {
		if entry.SubspaceID == subspaceID {
			return entry.InitialPostID
		}
	}
	return 0
}

// getInitialAttachmentID returns the initial attachment id for the given post
func (e *GenesisState) getInitialAttachmentID(subspaceID uint64, postID uint64) uint32 {
	for _, post := range e.GenesisPosts {
		if post.SubspaceID == subspaceID && post.ID == postID {
			return post.InitialAttachmentID
		}
	}
	return 0
}

// --------------------------------------------------------------------------------------------------------------------

// ValidateGenesis validates the given genesis state and returns an error if something is invalid
func ValidateGenesis(data *GenesisState) error {
	for _, entry := range data.SubspacesData {
		if containsDuplicatedSubspaceDataEntry(data.SubspacesData, entry) {
			return fmt.Errorf("duplicated subspace data entry: %d", entry.SubspaceID)
		}

		err := entry.Validate()
		if err != nil {
			return err
		}
	}

	for _, post := range data.GenesisPosts {
		if containsDuplicatedPost(data.GenesisPosts, post) {
			return fmt.Errorf("duplicated post: subspace id %d, post id %d", post.Post.SubspaceID, post.Post.ID)
		}

		initialPostID := data.getInitialPostID(post.SubspaceID)
		if post.ID >= initialPostID {
			return fmt.Errorf("post id must be lower than initial post id: subspace id %d", post.SubspaceID)
		}

		err := post.Validate()
		if err != nil {
			return err
		}
	}

	for _, attachment := range data.Attachments {
		if containsDuplicatedAttachment(data.Attachments, attachment) {
			return fmt.Errorf("duplicated attachmend: subspace id %d, post id %d, attachment id %d",
				attachment.SubspaceID, attachment.PostID, attachment.ID)
		}

		initialAttachmentID := data.getInitialAttachmentID(attachment.SubspaceID, attachment.PostID)
		if attachment.ID >= initialAttachmentID {
			return fmt.Errorf("attachment id must be lower than initial attachment it: subspace id %d, post id %d",
				attachment.SubspaceID, attachment.PostID)
		}

		err := attachment.Validate()
		if err != nil {
			return err
		}
	}

	for _, answer := range data.UserAnswers {
		if containsDuplicatedAnswer(data.UserAnswers, answer) {
			return fmt.Errorf("duplicated user answer: subspace id %d, post id %d, poll id %d, user: %s",
				answer.SubspaceID, answer.PostID, answer.PollID, answer.User)
		}

		err := answer.Validate()
		if err != nil {
			return err
		}
	}

	return data.Params.Validate()
}

// containsDuplicatedSubspaceDataEntry tells whether the given entries slice contains
// two or more entries for the same subspace
func containsDuplicatedSubspaceDataEntry(entries []SubspaceDataEntry, entry SubspaceDataEntry) bool {
	var count = 0
	for _, s := range entries {
		if s.SubspaceID == entry.SubspaceID {
			count++
		}
	}
	return count > 1
}

// containsDuplicatedPost tells whether the given posts slice contains two or more posts
// having the same id of the given one
func containsDuplicatedPost(posts []GenesisPost, post GenesisPost) bool {
	var count = 0
	for _, s := range posts {
		if s.Post.SubspaceID == post.SubspaceID && s.ID == post.ID {
			count++
		}
	}
	return count > 1
}

// containsDuplicatedAttachment tells whether the given attachments slice contains two or more attachments
// having the same id of the given one
func containsDuplicatedAttachment(attachments []Attachment, attachment Attachment) bool {
	var count = 0
	for _, s := range attachments {
		if s.SubspaceID == attachment.SubspaceID && s.PostID == attachment.PostID && s.ID == attachment.ID {
			count++
		}
	}
	return count > 1
}

// containsDuplicatedAnswer tells whether the given user answers slice contains two or more answers
// by the same user as the given one
func containsDuplicatedAnswer(answers []UserAnswer, answer UserAnswer) bool {
	var count = 0
	for _, s := range answers {
		if s.SubspaceID == answer.SubspaceID && s.PostID == answer.PostID && s.PollID == answer.PollID && answer.User.Equals(answer.User) {
			count++
		}
	}
	return count > 1
}
