package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// data is interface{} so in case of data retrieve from mongodb, Data may be bson.D
// bson.D is slice so it convert into Array, so we need convert into bson.M for json Object
func TryConvertBsonDToM(data interface{}) interface{} {
	switch v := data.(type) {
	case primitive.D:
		return v.Map()
	}

	return data
}
