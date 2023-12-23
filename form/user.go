package form

// UserUpdateForm user update form
type UserUpdateForm struct {
	ID    int64 //非表单赋值
	Level int   `form:"level" json:"level"`
}
