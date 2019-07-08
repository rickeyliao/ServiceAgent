package db

type NbsDbInter interface {
	Load() NbsDbInter
	Insert(key string, value string) error
	Delete(key string)
	Find(key string) (string, error)
	Update(key string, value string)
	Save()
	Print()
}
