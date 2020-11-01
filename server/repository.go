package main

import "database/sql"

type Repository interface {
	RunSchemaUpgrade(curVersion, targetVersion int)
}

func CheckNullString(s NullString) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: string(s),
		Valid:  true,
	}
}
