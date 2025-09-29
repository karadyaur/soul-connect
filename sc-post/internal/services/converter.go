package services

import (
	db "soul-connect/sc-post/internal/db/sqlc"
	"soul-connect/sc-post/internal/models"
	"soul-connect/sc-post/internal/utils"
)

func labelFromDB(l db.Label) models.Label {
	return models.Label{
		ID:   utils.UUIDToString(l.ID),
		Name: l.Name,
	}
}

func labelsFromDB(items []db.Label) []models.Label {
	labels := make([]models.Label, 0, len(items))
	for _, item := range items {
		labels = append(labels, labelFromDB(item))
	}
	return labels
}

func postFromDB(p db.Post, labels []db.Label) models.Post {
	likeCount := int32(0)
	if p.LikesCount.Valid {
		likeCount = p.LikesCount.Int32
	}

	return models.Post{
		ID:          utils.UUIDToString(p.ID),
		UserID:      utils.UUIDToString(p.UserID),
		Title:       p.Title,
		Description: utils.StringFromText(p.Description),
		LikesCount:  likeCount,
		CreatedAt:   utils.TimestampToTime(p.CreatedAt),
		UpdatedAt:   utils.TimestampToTime(p.UpdatedAt),
		Labels:      labelsFromDB(labels),
	}
}

func commentFromDB(c db.Comment) models.Comment {
	likeCount := int32(0)
	if c.LikesCount.Valid {
		likeCount = c.LikesCount.Int32
	}
	return models.Comment{
		ID:         utils.UUIDToString(c.ID),
		PostID:     utils.UUIDToString(c.PostID),
		UserID:     utils.UUIDToString(c.UserID),
		Content:    c.Content,
		LikesCount: likeCount,
		CreatedAt:  utils.TimestampToTime(c.CreatedAt),
		UpdatedAt:  utils.TimestampToTime(c.UpdatedAt),
	}
}

func commentsFromDB(items []db.Comment) []models.Comment {
	comments := make([]models.Comment, 0, len(items))
	for _, item := range items {
		comments = append(comments, commentFromDB(item))
	}
	return comments
}

func postSummaryFromRow(row db.GetPostsWithCommentsAndLikesRow) models.PostSummary {
	likeCount := int32(0)
	if row.PostLikes.Valid {
		likeCount = row.PostLikes.Int32
	}
	post := models.Post{
		ID:          utils.UUIDToString(row.PostID),
		UserID:      utils.UUIDToString(row.PostUserID),
		Title:       row.PostTitle,
		Description: utils.StringFromText(row.PostDescription),
		LikesCount:  likeCount,
	}
	return models.PostSummary{
		Post:          post,
		TotalComments: row.TotalComments,
		TotalLikes:    row.TotalLikes,
	}
}
