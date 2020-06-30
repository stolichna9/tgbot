package main

type User struct {
	//	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Username   string `bson:"username" json:"username"`
	TelegramID int    `bson:"telegramID" json:"telegramID"`
	Balance    int    `bson:"balance" json:"balance"`
}
