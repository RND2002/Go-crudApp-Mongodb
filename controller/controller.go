package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Rnd2002/mongodb/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://<Your-ussername>:<password>@cluster0.mehjx1r.mongodb.net/?retryWrites=true&w=majority"

const dbName = "Netflix"

const colName = "WatchList"

// most imp
var collection *mongo.Collection

func init() {
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongodb

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		panic(err)
	}
	fmt.Println("mongo db connection success")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("collection instance is ready")

}

//mongodb helpers
//insert one records

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted one movie with id", inserted.InsertedID)

}

func updateMovieData(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count", result.ModifiedCount)
}

func deleteMoviesByid(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)

	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Movie deleted", deleteCount)

}

func deleteAllRecord() int64 {
	filter := bson.D{{}}
	deleted, err := collection.DeleteMany(context.Background(), filter, nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All record deleted", deleted.DeletedCount)
	return deleted.DeletedCount
}

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	defer cur.Close(context.Background())
	return movies
}

func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")

	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "POST")

	var movie model.Netflix

	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func MarkWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "POST")

	params := mux.Vars(r)

	updateMovieData(params["id"])
}
func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteMoviesByid(params["id"])
	json.NewEncoder(w).Encode("Movie deleted successfully with id")
}
func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-control-Allow-Methods", "DELETE")

	count := deleteAllRecord()

	json.NewEncoder(w).Encode(count)

}
