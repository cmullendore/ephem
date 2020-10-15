package ephem

type IDatabase interface {
	SaveItem(path *string, item *string, lifetime int, readCount int) *error
	GetItem(path *string) (*string, *error)
	IncrementReadCount(path *string) *error
	DeleteItem(path *string) *error
	CleanupSecrets(maxReads int) *error
}
