package services

import (
	"ChatService/internal/models"
	"ChatService/pkg/database"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var chatCollection = database.Connection("chat_service", "chats")
var reviewsCollection = database.Connection("chat_service", "reviews")

func UserExists(id string) bool {
	url := fmt.Sprint("http://localhost:50051/api/user-exists/", id)

	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Пользователь существует")
		return true
	} else if resp.StatusCode == http.StatusNotFound {
		log.Println("Пользователь не найден")
		return false
	} else {
		log.Println("Неожиданный статус:", resp.StatusCode)
		return false
	}
}

func GetUserChats(userID string) ([]models.Chat, error) {
	var chats []models.Chat
	filter := bson.M{"participants": userID}

	cursor, err := chatCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &chats); err != nil {
		return nil, err
	}
	return chats, nil
}

func FindChatBetweenUsers(user1, user2 string) (*models.Chat, error) {
	var chat models.Chat

	filter := bson.M{
		"participants": bson.M{
			"$all": []string{user1, user2},
		},
	}

	err := chatCollection.FindOne(context.TODO(), filter).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &chat, nil
}

func FindOrCreateChat(user1, user2 string) (primitive.ObjectID, error) {
	chat, err := FindChatBetweenUsers(user1, user2)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if chat != nil {
		return chat.ID, nil
	}

	newChat := models.Chat{
		ID:           primitive.NewObjectID(),
		Participants: []string{user1, user2},
		Messages:     []models.Message{},
	}

	res, err := chatCollection.InsertOne(context.TODO(), newChat)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func SendMessage(chatID, sender, content string) error {
	objID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("неверный формат chatID: %v", err)
	}

	message := models.Message{
		Sender:    sender,
		Content:   content,
		Timestamp: time.Now(),
	}
	update := bson.M{
		"$push": bson.M{"messages": message},
		"$set":  bson.M{"lastMessage": message},
	}

	_, err = chatCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	return err
}

func GetUsersHistoryChat(userID, otherUserID string) (models.Chat, error) {
	var chat models.Chat
	filter := bson.M{
		"participants": bson.M{
			"$all": []string{userID, otherUserID},
		},
	}
	err := chatCollection.FindOne(context.TODO(), filter).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Chat{}, nil
		}
		return models.Chat{}, err
	}
	return chat, nil

}

func SendReview(productID, userName, content string, grade int) error {
	review := models.Review{
		Name:      userName,
		Content:   content,
		Grade:     grade,
		Timestamp: time.Now(),
	}

	update := bson.M{
		"$push": bson.M{"reviews": review},
	}

	_, err := reviewsCollection.UpdateOne(context.TODO(), bson.M{"productID": productID}, update)

	return err
}

func FindProductReviews(productID string) (*models.ProductReviews, error) {
	var reviews models.ProductReviews

	filter := bson.M{
		"productID": productID,
	}

	err := reviewsCollection.FindOne(context.TODO(), filter).Decode(&reviews)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &reviews, nil
}

func FindOrCreateProductReviews(productID string) error {
	reviews, err := FindProductReviews(productID)
	if err != nil {
		return err
	}
	if reviews != nil {
		return nil
	}

	newReviews := models.ProductReviews{
		ProductID: productID,
		Reviews:   []models.Review{},
	}

	_, err = reviewsCollection.InsertOne(context.TODO(), newReviews)
	if err != nil {
		return err
	}
	return nil
}

func GetAverageRating(products []string) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"productID": bson.M{"$in": products}}}},
		{{Key: "$unwind", Value: "$reviews"}},
		{{Key: "$group", Value: bson.M{
			"_id":      nil,
			"avgGrade": bson.M{"$avg": "$reviews.grade"},
		}}},
	}
	cursor, err := reviewsCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		AvgGrade float64 `bson:"avgGrade"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.AvgGrade, nil
	}

	return 0, nil
}
