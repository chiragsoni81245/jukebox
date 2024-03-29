package utils

import (
	"database/sql"
	"errors"
	"reflect"
	"github.com/gin-gonic/gin"
)

func GetDB(c *gin.Context) (*sql.DB, error) {
    _db, _ := c.Get("DB")
    db, ok := _db.(*sql.DB)
    if !ok {
        return nil, errors.New("Something went wrong in getting DB instance!")
    }
    return db, nil
}

// Helper function to append a value to the params slice if it is not zero
func AppendIfNonZero(params []interface{}, value interface{}) []interface{}{
    if value != reflect.Zero(reflect.TypeOf(value)).Interface() {
        return append(params, value)
    } else {
        return append(params, nil)
    }
}

// Helper function to append a value to the params slice if it is not empty
func AppendIfNotEmpty(params []interface{}, value string) []interface{}{
    if len(value) > 0 {
        return append(params, value)
    } else {
        return append(params, nil)
    }
}
