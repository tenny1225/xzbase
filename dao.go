package xzbase

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type CountViewModel struct {
	Count int64 `json:"count" bson:"count"`
}

type ModelDao interface {
	Create(m interface{}) error
	UpdateMany(id string, m interface{}) error
	Update(id string, m map[string]interface{}) error
	Inc(id string, m map[string]interface{}) error
	Pull(id string, m map[string]interface{}) error
	Push(id string, m map[string]interface{}) error
	Save(m interface{}, filter map[string]interface{}) error
	Delete(id string) error
	DeleteByMany(map[string]interface{}) error
	ForceDelete(id string) error
	ForceDeleteByMany(map[string]interface{}) error
	Get(v interface{}, id string) error
	GetBy(m interface{}, key, value string) error
	GetManyBy(m interface{}, key, value string) error
	GetManyByMany(m interface{}, filter map[string]interface{}) error
	GetOneOrder(m interface{}, key string, order int) error
	All(m interface{}) error
	Begin() *mongo.Collection
}
type modelDao struct {
	Ctx       context.Context
	TableName string
}


func NewModelDao(tableName string) ModelDao {
	return &modelDao{Ctx: context.Background(), TableName: tableName}
}

func (d *modelDao) Create(m interface{}) error {
	_, e := d.Begin().InsertOne(d.Ctx, m)
	return e
}

func (d *modelDao) UpdateMany(id string, m interface{}) error {
	d.ForceDelete(id)
	return d.Create(m)
}

func (d *modelDao) Update(id string, m map[string]interface{}) error {
	m["model.updatedAt"] = time.Now()
	_, e := d.Begin().UpdateOne(d.Ctx, bson.M{"model.id": id}, bson.M{"$set": m})

	return e
}
func (d *modelDao) Inc(id string, m map[string]interface{}) error {
	_, e := d.Begin().UpdateOne(d.Ctx, bson.M{"model.id": id}, bson.M{"$inc": m})

	return e
}
func (d *modelDao) Pull(id string, m map[string]interface{}) error {
	var e error

	_, e = d.Begin().UpdateOne(d.Ctx, bson.M{"model.id": id}, bson.M{"$pull": m})

	return e
}
func (d *modelDao) Push(id string, m map[string]interface{}) error {
	_, e := d.Begin().UpdateOne(d.Ctx, bson.M{"model.id": id}, bson.M{"$push": m})

	return e
}
func (d *modelDao) Save(m interface{}, filter map[string]interface{}) error {
	f := bson.M{}

	for k, v := range filter {
		f[k] = v
	}
	_, e := d.Begin().DeleteMany(d.Ctx, f)
	if e != nil && e != mongo.ErrNoDocuments {
		return e
	}
	return d.Create(m)
}
func (d *modelDao) Delete(id string) error {
	_, e := d.Begin().UpdateOne(d.Ctx, bson.M{"model.id": id}, bson.M{"$set": bson.M{
		"model.deletedAt": time.Now(),
	}})
	return e
}
func (d *modelDao) DeleteByMany(m map[string]interface{}) error {
	_, e := d.Begin().DeleteMany(d.Ctx, m)
	return e
}
func (d *modelDao) ForceDelete(id string) error {
	_, e := d.Begin().DeleteOne(d.Ctx, bson.M{"model.id": id})
	return e
}
func (d *modelDao) ForceDeleteByMany(m map[string]interface{}) error {
	filter := bson.M{}
	for k, v := range m {
		filter[k] = v
	}

	_, e := d.Begin().DeleteMany(d.Ctx, filter)
	return e
}
func (d *modelDao) GetOneOrder(m interface{}, key string, order int) error {
	opts := options.FindOne()
	opts.SetSort(bson.D{{Key: key, Value: order}})
	return d.Begin().FindOne(d.Ctx, bson.M{}, opts).Decode(m)
}
func (d *modelDao) Get(m interface{}, id string) error {

	return d.Begin().FindOne(d.Ctx, bson.M{"model.id": id}).Decode(m)
}
func (d *modelDao) GetBy(m interface{}, key, value string) error {

	return d.Begin().FindOne(d.Ctx, bson.M{key: value}).Decode(m)
}
func (d *modelDao) GetManyBy(m interface{}, key, value string) error {

	c, e := d.Begin().Find(d.Ctx, bson.M{key: value})

	if e != nil {
		return e
	}
	return c.All(d.Ctx, m)
}
func (d *modelDao) GetManyByMany(m interface{}, filter map[string]interface{}) error {

	f := bson.M{"model.deletedAt": nil}

	for k, v := range filter {
		f[k] = v
	}

	c, e := d.Begin().Find(d.Ctx, f)
	if e != nil {
		return e
	}
	return c.All(d.Ctx, m)
}
func (d *modelDao) All(m interface{}) error {
	c, e := d.Begin().Find(d.Ctx, bson.M{"model.deletedAt": nil})
	if e != nil {
		return e
	}
	return c.All(d.Ctx, m)
}

func (m *modelDao) Begin() *mongo.Collection {

	return db.Collection(m.TableName)
}
