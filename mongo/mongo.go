package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Note представляет заметку пользователя
type Note struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  int64              `bson:"user_id"`
	Text    string             `bson:"text"`
	Created time.Time          `bson:"created_at"`
}

var client *mongo.Client
var notesCollection *mongo.Collection

// InitMongoDB инициализирует подключение к MongoDB
func InitMongoDB(uri string, dbName string, collectionName string) error {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Подключено к MongoDB")

	notesCollection = client.Database(dbName).Collection(collectionName)
	return nil
}

// CreateNote создает новую заметку для пользователя
func CreateNote(userID int64, text string) error {
	note := Note{
		UserID:  userID,
		Text:    text,
		Created: time.Now(),
	}

	_, err := notesCollection.InsertOne(context.TODO(), note)
	return err
}

// GetNotes получает все заметки пользователя
func GetNotes(userID int64) ([]Note, error) {
	var notes []Note
	filter := bson.M{"user_id": userID}

	cursor, err := notesCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var note Note
		if err := cursor.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// DeleteNoteByID удаляет заметку по её ID
func DeleteNoteByID(noteID string) error {
	id, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return err
	}

	_, err = notesCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
