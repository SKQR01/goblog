package store

//Store is a abstract database store implementation.
type Store interface {
	User() UserRepository
	Post() PostRepository
}
