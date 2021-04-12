package xzbase

import (
	"context"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type Model struct {
	Id string `json:"id" bson:"id" form:"id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" bson:"deletedAt"`
}
type ZError struct {
	Code int64
	Msg  string
}

func (z ZError) Error() string {
	return z.Msg
}


var addIndexFuncs = make([]func(), 0)

var db *mongo.Database

func InitDB(uri,project,username,pwd string) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cp := options.Client()
	cp.Auth = &options.Credential{}
	cp.Auth.Username = username
	cp.Auth.Password = pwd
	//client, err := mongo.Connect(ctx, cp.ApplyURI(conf.Config.String("mongoUri")))
	client, err := mongo.NewClient(cp.ApplyURI(uri))
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	// 判断服务是不是可用
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logs.Error(err)
		panic(err)
	}
	// 2, 选择数据库my_db
	db = client.Database(project)

	for _, l := range addIndexFuncs {
		l()
	}

}
func AddIndex(table string, uniqueNames []string, names []string) {
	uniqueIndex := make([]mongo.IndexModel, 0)
	index := make([]mongo.IndexModel, 0)
	uniqueIndex = append(uniqueIndex, mongo.IndexModel{
		Keys:    bsonx.Doc{{"model.id", bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	})
	for _, n := range uniqueNames {
		uniqueIndex = append(uniqueIndex, mongo.IndexModel{
			Keys:    bsonx.Doc{{n, bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		})
	}
	index = append(uniqueIndex, mongo.IndexModel{
		Keys: bsonx.Doc{{"model.createdAt", bsonx.Int32(1)}},
	})
	index = append(uniqueIndex, mongo.IndexModel{
		Keys: bsonx.Doc{{"model.updatedAt", bsonx.Int32(1)}},
	})
	index = append(uniqueIndex, mongo.IndexModel{
		Keys: bsonx.Doc{{"model.deletedAt", bsonx.Int32(1)}},
	})

	for _, n := range names {
		index = append(uniqueIndex, mongo.IndexModel{
			Keys: bsonx.Doc{{n, bsonx.Int32(1)}},
		})
	}
	for _, n := range append(index, uniqueIndex...) {
		_, err := db.Collection(table).Indexes().CreateOne(context.Background(), n)
		if err != nil {
			logs.Error(err)
		}
	}

}
func DB() *mongo.Database {
	return db

}

