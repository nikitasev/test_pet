package service

type NewUserEventLogger interface {
	Log(userId int64) error
}
