package rxlog

import (
	"context"
	"time"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
)

//SaveRxSystemComment save system comment for Rx
func SaveRxSystemComment(ctx context.Context, rxID, comment string) {
	rxComment := new(data.RxComment)
	rxComment.RxID = rxID
	rxComment.AgentID = "SYSTEM"
	rxComment.AgentName = "Phil System"

	rxComment.Important = false
	rxComment.Type = "SYSTEM"
	rxComment.Message = comment

	SaveRxComment(ctx, rxComment)
}

//SaveRxSystemUserComment save system comment for Rx, but with Type = "USER"
func SaveRxSystemUserComment(ctx context.Context, rxID, comment string) {
	rxComment := new(data.RxComment)
	rxComment.RxID = rxID
	rxComment.AgentID = "SYSTEM"
	rxComment.AgentName = "Phil System"

	rxComment.Important = false
	rxComment.Type = "USER"
	rxComment.Message = comment

	SaveRxComment(ctx, rxComment)
}

// SaveRxSystemCommentAsImportant save rx system comment with option of marking it important
func SaveRxSystemCommentAsImportant(ctx context.Context, rxID, comment string, important bool) {
	rxComment := new(data.RxComment)
	rxComment.RxID = rxID
	rxComment.AgentID = "SYSTEM"
	rxComment.AgentName = "Phil System"

	rxComment.Important = important
	rxComment.Type = "SYSTEM"
	rxComment.Message = comment

	SaveRxComment(ctx, rxComment)
}

// SaveRxCommentForUser saves rx comment for user and prescription with given ids
func SaveRxCommentForUser(ctx context.Context, userID, rxID, comment string) {
	rxComment := new(data.RxComment)
	rxComment.RxID = rxID
	rxComment.AgentID = userID
	rxComment.AgentName = "User"

	rxComment.Important = false
	rxComment.Type = "USER"
	rxComment.Message = comment

	SaveRxComment(ctx, rxComment)
}

// SaveRxCommentForUserWithName saves rx comment for user and prescription with given ids, and includes their username
func SaveRxCommentForUserWithName(ctx context.Context, userID, userName, rxID, comment string) {
	rxComment := new(data.RxComment)
	rxComment.RxID = rxID
	rxComment.AgentID = userID
	rxComment.AgentName = userName

	rxComment.Important = false
	rxComment.Type = "USER"
	rxComment.Message = comment

	SaveRxComment(ctx, rxComment)
}

//SaveRxComment saves Rx comment
func SaveRxComment(ctx context.Context, comment *data.RxComment) (*data.RxComment, error) {
	session := data.Session()
	defer session.Close()

	if comment.StringID() == "" {
		comment.InitData()
	} else {
		t := time.Now().UTC()
		comment.CreatedDate = &t
		comment.UpdatedDate = &t
	}
	return comment, session.Save(comment)
}

//FindRxComments finds all the rx comments for the given rx id
func FindRxComments(ctx context.Context, rxID string) ([]*data.RxComment, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"rxId": rxID}, new(data.RxComment))
	if err != nil {
		return nil, err
	}
	return results.([]*data.RxComment), nil
}

//FindRxCommentsSinceDate finds latest rx comments for the given rx id and current fill
func FindRxCommentsSinceDate(ctx context.Context, rxID string, sinceDate *time.Time) ([]*data.RxComment, error) {
	if sinceDate == nil {
		return FindRxComments(ctx, rxID)
	}

	session := data.Session()
	defer session.Close()

	q := gmgo.Q{
		"rxId": rxID,
		"createdDate": gmgo.Q{
			"$gte": sinceDate,
		},
	}

	results, err := session.FindAll(q, new(data.RxComment))
	if err != nil {
		return nil, err
	}
	return results.([]*data.RxComment), nil
}
