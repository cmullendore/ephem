package mysql

const getSecret string = `
	SELECT secret
	FROM secrets
    WHERE path = ?;
`

const saveSecret string = `
	INSERT INTO secrets (path, secret, expires, readCount, timestamp)
	VALUES(?, ?, (SELECT DATE_ADD(?, INTERVAL ? SECOND)), ?, ?)
`

const incrementReadCount string = `
	UPDATE secrets
	SET readCount = readCount + 1
	WHERE path = ?
`

const deleteSecret string = `
	DELETE
	FROM secrets
	WHERE path = ?
`

const cleanupSecrets string = `
	DELETE
	FROM secrets
	WHERE (expires < NOW() OR readCount >= ?);
`

const createSecretsTable string = "CREATE TABLE IF NOT EXISTS `ephem`.`secrets` (`path` CHAR(128) NOT NULL, `secret` MEDIUMBLOB NULL, `expires` DATETIME NULL, `readCount` INT NULL, `timestamp` DATETIME NULL, PRIMARY KEY (`path`));"
