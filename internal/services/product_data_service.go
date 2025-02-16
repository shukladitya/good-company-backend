package services

import (
	"context"
	"errors"
	"theGoodCompany/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentService struct {
    collection *mongo.Collection
}

func NewDocumentService(collection *mongo.Collection) *DocumentService {
    return &DocumentService{
        collection: collection,
    }
}

var (
    ErrInvalidID = errors.New("invalid ID format")
    ErrNotFound  = errors.New("document not found")
)

func (s *DocumentService) GetDocumentByID(ctx context.Context, id string) (*models.Document, error) {
    // Convert string ID to ObjectID
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, ErrInvalidID
    }

    // Create filter
    filter := bson.M{"_id": objID}

    // Find document
    var result models.Document
    err = s.collection.FindOne(ctx, filter).Decode(&result)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, ErrNotFound
        }
        return nil, err
    }

    return &result, nil
}