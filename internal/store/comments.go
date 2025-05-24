package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID              int64  `json:"id"`
	PostID          int64  `json:"post_id"`
	UserID          *int64 `json:"user_id,omitempty"`
	Content         string `json:"content"`
	IsAIGenerated   bool   `json:"is_ai_generated"`
	ParentCommentID *int64 `json:"parent_comment_id,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content, is_ai_generated, parent_comment_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
		`

	err := s.db.QueryRowContext(ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
		comment.IsAIGenerated,
		comment.ParentCommentID,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT id, post_id, user_id, content, is_ai_generated, parent_comment_id, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.IsAIGenerated,
			&comment.ParentCommentID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}