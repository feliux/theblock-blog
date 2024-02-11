package database

import (
	"context"
	"fmt"
	"log"

	"github.com/feliux/theblock-blog/models"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoclient *mongo.Client
	database    string
)

func Connect(ctx context.Context) error {
	user := ctx.Value(models.Key("user")).(string)
	password := ctx.Value(models.Key("password")).(string)
	host := ctx.Value(models.Key("host")).(string)
	database = ctx.Value(models.Key("database")).(string)
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, password, host)
	var clientOptions = options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Fail connecting to database '%s' with ERROR: %s", database, err.Error())
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Error pining database: " + err.Error())
		return err
	}
	log.Printf("Connection to database '%s' is ready.", database)
	mongoclient = client
	return nil
}

func CheckDatabase() bool {
	err := mongoclient.Ping(context.TODO(), nil)
	return err == nil
}

func CheckUser(email string) (models.User, bool, string) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("users")

	condition := bson.M{"email": email}
	var results models.User
	err := col.FindOne(ctx, condition).Decode(&results)
	id := results.Id.Hex()
	if err != nil {
		return results, false, id
	}
	return results, true, id
}

func Insert(user models.User) (string, bool, error) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("users")
	// Encrypt password
	user.Password, _ = encryptPassword(user.Password)

	results, err := col.InsertOne(ctx, user)
	if err != nil {
		return "", false, err
	}
	objId, _ := results.InsertedID.(primitive.ObjectID)
	return objId.String(), true, nil
}

func encryptPassword(password string) (string, error) {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}

func TryLogin(email, password string) (models.User, bool) {
	user, isOk, _ := CheckUser(email)
	if !isOk {
		return user, false
	}
	passwordBytes := []byte(password)
	passwordFromDatabase := []byte(user.Password) // encrypted password
	err := bcrypt.CompareHashAndPassword(passwordFromDatabase, passwordBytes)
	if err != nil {
		return user, false
	}
	return user, true
}

func SearchProfile(id string) (models.User, error) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("users")
	var profile models.User
	objId, _ := primitive.ObjectIDFromHex(id)
	condition := bson.M{
		"_id": objId,
	}
	err := col.FindOne(ctx, condition).Decode(&profile)
	profile.Password = ""
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func ModifyRegister(user models.User, id string) (bool, error) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("users")
	register := make(map[string]interface{})
	if len(user.Nombre) > 0 {
		register["nombre"] = user.Nombre
	}
	if len(user.Apellidos) > 0 {
		register["apellidos"] = user.Apellidos
	}
	register["fechaNacimiento"] = user.FechaNacimiento
	if len(user.Avatar) > 0 {
		register["avatar"] = user.Avatar
	}
	if len(user.Banner) > 0 {
		register["banner"] = user.Banner
	}
	if len(user.Biografia) > 0 {
		register["biografia"] = user.Biografia
	}
	if len(user.Ubicacion) > 0 {
		register["ubicacion"] = user.Ubicacion
	}
	if len(user.SitioWeb) > 0 {
		register["sitioWeb"] = user.SitioWeb
	}
	updtString := bson.M{
		"$set": register,
	}

	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": objId}}

	_, err := col.UpdateOne(ctx, filter, updtString)
	if err != nil {
		return false, err
	}

	return true, nil
}

func InsertTweet(tweet models.SaveTweet) (string, bool, error) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("tweet")
	register := bson.M{
		"userid":  tweet.UserId,
		"mensaje": tweet.Message,
		"fecha":   tweet.Date,
	}
	result, err := col.InsertOne(ctx, register)
	if err != nil {
		return "", false, err
	}

	objId, _ := result.InsertedID.(primitive.ObjectID)
	return objId.String(), true, nil
}

func ReadTweets(id string, page int64) ([]*models.RetrieveTweets, bool) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("tweet")

	var results []*models.RetrieveTweets

	condition := bson.M{
		"userid": id,
	}

	opts := options.Find()
	opts.SetLimit(20)
	opts.SetSort(bson.D{{Key: "fecha", Value: -1}})
	opts.SetSkip((page - 1) * 20)

	cursor, err := col.Find(ctx, condition, opts)
	if err != nil {
		return results, false
	}

	for cursor.Next(ctx) {
		var register models.RetrieveTweets
		err := cursor.Decode(&register)
		if err != nil {
			return results, false
		}
		results = append(results, &register)
	}
	return results, true
}

func DeleteTweet(tweetId string, userId string) error {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("tweet")
	objId, _ := primitive.ObjectIDFromHex(tweetId)
	condition := bson.M{
		"_id":    objId,
		"userid": userId,
	}
	_, err := col.DeleteOne(ctx, condition)
	return err
}

func InsertRelation(relation models.Relationship) (bool, error) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("relation")

	_, err := col.InsertOne(ctx, relation)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteRelation(relation models.Relationship) (bool, error) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("relation")

	_, err := col.DeleteOne(ctx, relation)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetRelation(relation models.Relationship) bool {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("relation")

	condition := bson.M{
		"userid":         relation.UserId,
		"userrelationid": relation.UserRelationId,
	}

	var result models.Relationship
	err := col.FindOne(ctx, condition).Decode(&result)
	if err != nil {
		return false
	}
	return true
}

func GetAllUsers(id string, page int64, search string, userType string) ([]*models.User, bool) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("users")
	var allUsers []*models.User
	opts := options.Find()
	opts.SetLimit(20)
	opts.SetSkip((page - 1) * 20)
	query := bson.M{
		"nombre": bson.M{"$regex": `(?i)` + search},
	}
	cursor, err := col.Find(ctx, query, opts)
	if err != nil {
		log.Printf("Could not find users data with ERROR: %s" + err.Error())
		return allUsers, false
	}

	var toInclude bool

	for cursor.Next(ctx) {
		var user models.User

		err := cursor.Decode(&user)
		if err != nil {
			log.Printf("Fail decoding user data from GetAllUsers function with ERROR: %s" + err.Error())
			return allUsers, false
		}

		var relation models.Relationship
		relation.UserId = id
		relation.UserRelationId = user.Id.Hex()

		toInclude = false

		existUser := GetRelation(relation)
		if userType == "new" && !existUser {
			toInclude = true
		}
		if userType == "follow" && existUser {
			toInclude = true
		}
		if relation.UserRelationId == id { // If is the same person (user from jwt and user from database)
			toInclude = false
		}
		if toInclude {
			user.Password = ""
			allUsers = append(allUsers, &user)
		}
	}
	err = cursor.Err()
	if err != nil {
		log.Printf("Fail decoding user data from GetAllUsers function... cursor.Err(): %s" + err.Error())
		return allUsers, false
	}
	cursor.Close(ctx)
	return allUsers, true
}

func GetFollowerTweets(id string, pag int) ([]models.ResponseFollowerTweets, bool) {
	ctx := context.TODO()
	db := mongoclient.Database(database)
	col := db.Collection("relation")
	skip := (pag - 1) * 20
	conditions := make([]bson.M, 0)
	conditions = append(conditions, bson.M{"$match": bson.M{"userid": id}})
	conditions = append(conditions, bson.M{
		"$lookup": bson.M{
			"from":         "tweet",
			"localField":   "userrelationid",
			"foreignField": "userid",
			"as":           "tweet",
		}})
	conditions = append(conditions, bson.M{"$unwind": "$tweet"})
	conditions = append(conditions, bson.M{"$sort": bson.M{"tweet.fecha": -1}})
	conditions = append(conditions, bson.M{"$skip": skip})
	conditions = append(conditions, bson.M{"$limit": 20})

	var results []models.ResponseFollowerTweets

	cursor, err := col.Aggregate(ctx, conditions)
	if err != nil {
		log.Printf("Failed at GetFollowerTweets aggregating conditions into cursor with ERROR: %s" + err.Error())
		return results, false
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		log.Printf("Failed at GetFollowerTweets retrieving all results from cursor... cursor.Err(): %s" + err.Error())
		return results, false
	}
	return results, true
}
