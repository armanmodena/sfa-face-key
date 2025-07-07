package model

type User struct {
	Id              int    `json:"id" db:"id" bson:"id"`
	Username        string `json:"username" db:"username" bson:"username"`
	Nik             string `json:"nik" db:"nik" bson:"nik"`
	FullName        string `json:"full_name" db:"full_name" bson:"full_name"`
	Email           string `json:"email" db:"email" bson:"email"`
	Phone           string `json:"phone" db:"phone" bson:"phone"`
	IsActive        bool   `json:"is_active" db:"is_active" bson:"is_active"`
	GoFaceEmbedding string `json:"go_face_embedding" db:"go_face_embedding" bson:"go_face_embedding"`
	GoFaceImageUrl  string `json:"go_face_image_url" db:"go_face_image_url" bson:"go_face_image_url"`
}
