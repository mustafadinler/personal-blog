package service

import (
	"context"
	"log"
	"os"
	"time"

	repository "personal-blog/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostDto struct {
	Title      string
	Body       string
	CreateDate time.Time
	Id         string
	CategoryId int
}

type PostService interface {
	FindAll() ([]*PostDto, error)
	FindById(id string) (*PostDto, error)
	FindByPagination(page int64, pageSize int64) ([]*PostDto, error)
	FindByCategoryPagination(categoryid string, page int, pageSize int) []*PostDto
	Add(PostDto string) (bool, error)
}

func FindById(id string) PostDto {
	repo, err := createNewRepo()
	if err != nil {
		log.Fatal(err)
	}

	data, err := repo.FindById(id)
	if err != nil {
		log.Fatal(err)
	}
	return mapSingleDataToDto(*data)
}

func FindByCategoryPagination(categoryid string, page int64, pageSize int64) []PostDto {
	repo, err := createNewRepo()
	if err != nil {
		log.Fatal(err)
	}

	data, err := repo.FindByCategoryPagination(categoryid, page, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	return mapDataToDto(data)
}

func FindAll() []PostDto {
	repo, err := createNewRepo()
	if err != nil {
		log.Fatal(err)
	}

	data, err := repo.FindAll()
	if err != nil {
		log.Fatal(err)
	}
	return mapDataToDto(data)
}

func FindByPagination(page int64, pageSize int64) []PostDto {
	repo, err := createNewRepo()
	if err != nil {
		log.Fatal(err)
	}

	data, err := repo.FindByPagination(page, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	return mapDataToDto(data)
}

func Add(dto *PostDto) (bool, error) {
	repo, err := createNewRepo()
	if err != nil {
		log.Fatal(err)
	}

	data := mapSingleDtoToData(*dto)
	data.Id = primitive.NewObjectID()
	data.IsActive = true
	result, err := repo.Add(data)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

func mapSingleDataToDto(data repository.PostData) PostDto {
	dto := PostDto{Title: data.Title, Body: data.Body, CreateDate: data.CreateDate, Id: data.Id.Hex(), CategoryId: data.CategoryId}
	return dto
}

func mapDataToDto(data []*repository.PostData) []PostDto {
	var slice []PostDto
	for _, element := range data {
		dto := mapSingleDataToDto(*element)
		slice = append(slice, dto)
	}
	return slice
}

func mapSingleDtoToData(dto PostDto) repository.PostData {
	objectId, _ := primitive.ObjectIDFromHex(dto.Id)
	data := repository.PostData{Title: dto.Title, Body: dto.Body, CreateDate: dto.CreateDate, Id: objectId, CategoryId: dto.CategoryId}
	return data
}

func mapDtoToData(dtos []*PostDto) []repository.PostData {
	var slice []repository.PostData
	for _, element := range dtos {
		data := mapSingleDtoToData(*element)
		slice = append(slice, data)
	}
	return slice
}

func createNewRepo() (*repository.MongoRepository, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("db-uri"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	repository, err := repository.NewMongoRepository(client.Database(os.Getenv("db-name")))
	return repository, err
}
