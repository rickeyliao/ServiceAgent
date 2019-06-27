package db

type NbsDbInter interface {
	Insert(key string, value string) error
	Delete(key string)
	Find(key string) (string, error)
	Update(key string, value string)
	Save()
}
