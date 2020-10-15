package core

type ISecretsEngine interface {
	SaveItem(path *string, item *[]byte) *error
	GetItem(path *string) (*[]byte, *error)
}
