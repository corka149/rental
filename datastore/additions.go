package datastore

func (u User) isDefaultUser() bool {
	return u.ID == 0
}
