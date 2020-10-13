package shashankMongo

import (
	"fmt"
	"context"
	"os"
	"log"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectToDataBase struct {
	CustomApplyURI string 
	DatabaseName string 
	CollectionName string 
}

type ProfileConfig struct{
	Zone int64 `bson: "zone" json: "zone"`
	MessagePlan int64 `bson: "messageplan" json: "messageplan"`
	Tracking bool `bson: "tracking" json: "tracking"`
	ZoneID []string `bson: "zoneid" json: "zoneid"`
}

type BusinessAccount struct{
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserName string `bson: "username" json: "username"`
	BusinessName string `bson: "businessname" json: "businessname"`
	Password string `bson: "password" json: "password"`
	City string `bson: "city" json: "city"`
	BusinessPlan string `bson: "businessplan" json: "businessplan"`
	ProfileConfig ProfileConfig `bson: "profileConfig" json: "profileConfig"`
	UserID string
}

type ZoneInfo struct {
	Name string `bson: "name" json: "name"`
	BusinessUID string `bson: "businessUid" json: "businessUid"`
}

var resultID string
var profileConfig *ProfileConfig
var businessAccount *BusinessAccount
var zones *[]ZoneInfo

func initializeClient(applyURI string) (*mongo.Client,context.Context){
	c,err:= mongo.NewClient(options.Client().ApplyURI(applyURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx:= context.Background()
	err = c.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return c,ctx
}

func InsertOne(connectionInfo *ConnectToDataBase,customInsertStruct map[string]interface{}) string {
	
	client,ctx:= initializeClient(connectionInfo.CustomApplyURI)
	databaseName := client.Database(connectionInfo.DatabaseName)
	collectionName := databaseName.Collection(connectionInfo.CollectionName)

	result, insertErr := collectionName.InsertOne(ctx, customInsertStruct)
	if insertErr != nil {
	fmt.Println("InsertOne ERROR:", insertErr)
	os.Exit(1) // safely exit script on error
	} else {
	fmt.Println("InsertOne() API result:", result)

	newID := result.InsertedID
	fmt.Println("InsertOne() newID:", newID)
	resultID = newID.(primitive.ObjectID).Hex()
	}
	return resultID

}

func UpdateOneByID(connectionInfo *ConnectToDataBase,docID string,insertKey string, insertValue string) int64 {

	client,ctx:= initializeClient(connectionInfo.CustomApplyURI)
	databaseName := client.Database(connectionInfo.DatabaseName)
	collectionName := databaseName.Collection(connectionInfo.CollectionName)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{insertKey: insertValue}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Fatal(err)
		}

	return res.ModifiedCount
}

func FetchProfileConfiguration(connectionInfo *ConnectToDataBase,filterValue string) *ProfileConfig{

	client,ctx:= initializeClient(connectionInfo.CustomApplyURI)
	databaseName := client.Database(connectionInfo.DatabaseName)
	collectionName := databaseName.Collection(connectionInfo.CollectionName)
	
	filter := bson.M{"plan": filterValue}
    err:= collectionName.FindOne(ctx, filter).Decode(&profileConfig)
	if err != nil {
		log.Println(err)
	}
    return profileConfig
}

func UpdateProfileConfiguration(connectionInfo *ConnectToDataBase,docID string,config *ProfileConfig) int64 {

	client,ctx:= initializeClient(connectionInfo.CustomApplyURI)
	databaseName := client.Database(connectionInfo.DatabaseName)
	collectionName := databaseName.Collection(connectionInfo.CollectionName)

	id, _ := primitive.ObjectIDFromHex(docID)
	update := bson.M{"$set": bson.M{"profileConfig": config}}
		filter := bson.M{"_id": id}
		res,err := collectionName.UpdateOne(ctx,filter, update)
		if err!=nil{
			log.Fatal(err)
		}

	fmt.Println("profile created")
	return res.ModifiedCount

}

func FetchProfile(connectionInfo *ConnectToDataBase,docID string) *BusinessAccount{

	client,ctx:= initializeClient(connectionInfo.CustomApplyURI)
	databaseName := client.Database(connectionInfo.DatabaseName)
	collectionName := databaseName.Collection(connectionInfo.CollectionName)
	
	id, _ := primitive.ObjectIDFromHex(docID)
	filter := bson.M{"_id": id}
    err:= collectionName.FindOne(ctx, filter).Decode(&businessAccount)
	if err != nil {
		log.Println(err)
	}
	businessAccount.UserID=docID
    return businessAccount
}

//FetchLogin is exported
func FetchLogin(connectionInfo *ConnectToDataBase,username string, password string) (*BusinessAccount, error){
	
	client,ctx:= initializeClient(connectionInfo.CustomApplyURI)
	databaseName := client.Database(connectionInfo.DatabaseName)
	collectionName := databaseName.Collection(connectionInfo.CollectionName)
	
	filter := bson.M{"username": username,"password": password}
    err:= collectionName.FindOne(ctx, filter).Decode(&businessAccount)
	if err != nil {
		log.Println(err)
	}
	resultID = businessAccount.ID.Hex()
	businessAccount.UserID = resultID
	return businessAccount,err
}

func getZone(connectionInfo *ConnectToDataBase,docID string) {

	client,ctx:= initializeClient(connectionInfo.CustomApplyURI)
	databaseName := client.Database(connectionInfo.DatabaseName)
	collectionName := databaseName.Collection(connectionInfo.CollectionName)

	cursor, err := collectionName.Find(ctx, bson.M{"businessUid":docID})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &zones); err != nil {
		log.Fatal(err)
	}
	fmt.Println(zones)
}