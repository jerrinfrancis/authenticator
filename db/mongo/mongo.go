package mongo

import (
	"context"
	"encoding/json"
	"github.com/jerrinfrancis/authenticator/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type repo struct {
	cl *mongo.Client
}

type userDb struct {
	r repo
}

type User struct {
	ID          *primitive.ObjectID `bson:"_id"`
	FirstName   string
	LastName    string
	PhoneNumber uint64
	Email       string
	Passwd      string
	Salt        string
	Secret      string
}

func (ud userDb) Insert(u db.User) (*db.User, error) {
	col := ud.r.cl.Database("test").Collection("user")
	b, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	du := new(User)
	err = json.Unmarshal(b, du)
	id := primitive.NewObjectID()
	du.ID = &id
	u.ID = id.Hex()
	_, err = col.InsertOne(context.TODO(), du)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (ud userDb) Find(email string) (*db.User, error) {
	col := ud.r.cl.Database("test").Collection("user")
	filter := bson.D{{Key: "email", Value: email}}
	du := new(db.User)
	u := new(User)
	if err := col.FindOne(context.TODO(), filter).Decode(u); err != nil {
		return nil, err
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(b, du)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return du, nil
}

var client *mongo.Client

func Client() *mongo.Client {
	if client != nil {
		return client
	}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		cancel()
		log.Fatalln("Unable to connect to DB", err)
	}
	defer cancel()
	return client

}

func (r repo) User() db.UserDB {
	return userDb{r: r}
}

func New(cl *mongo.Client) db.DB {
	return repo{cl: cl}
}
