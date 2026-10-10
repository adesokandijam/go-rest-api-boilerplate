package data

import "database/sql"

type Models struct{ Items ItemModel }

func NewModels(db *sql.DB) Models { return Models{Items: ItemModel{DB: db}} }
