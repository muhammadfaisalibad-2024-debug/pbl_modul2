package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"api-students/app/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrProjectLimit    = errors.New("jumlah proyek sudah mencapai batas")
	ErrProjectNotFound = errors.New("proyek tidak ditemukan")
)

const MaxProjects = 20

type PortfolioRepository interface {
	EnsureIndexes(ctx context.Context) error
	FindByStudent(ctx context.Context, studentID int) (model.Portfolio, error)
	Replace(ctx context.Context, p model.Portfolio) (model.Portfolio, error)
	AddProject(ctx context.Context, studentID int, pr model.Project) (model.Portfolio, error)
	RemoveProject(ctx context.Context, studentID int, projectID string) error
	TopTechnologies(ctx context.Context, limit int) ([]model.TechCount, error)
}

type portfolioDocument struct {
	ID           bson.ObjectID         `bson:"_id,omitempty"`
	StudentID    int                   `bson:"student_id"`
	Headline     string                `bson:"headline"`
	Bio          string                `bson:"bio"`
	Projects     []projectDocument     `bson:"projects"`
	Certificates []certificateDocument `bson:"certificates"`
	CreatedAt    time.Time             `bson:"created_at"`
	UpdatedAt    time.Time             `bson:"updated_at"`
}

type projectDocument struct {
	ID           bson.ObjectID `bson:"_id"`
	Title        string        `bson:"title"`
	Technologies []string      `bson:"technologies"`
	URL          string        `bson:"url,omitempty"`
	Year         int           `bson:"year"`
}

type certificateDocument struct {
	ID     bson.ObjectID `bson:"_id"`
	Name   string        `bson:"name"`
	Issuer string        `bson:"issuer"`
	Date   time.Time     `bson:"date"`
}

func (d portfolioDocument) toModel() model.Portfolio {
	projects := make([]model.Project, 0, len(d.Projects))
	for _, p := range d.Projects {
		techs := p.Technologies
		if techs == nil {
			techs = []string{}
		}
		projects = append(projects, model.Project{
			ID:           p.ID.Hex(),
			Title:        p.Title,
			Technologies: techs,
			URL:          p.URL,
			Year:         p.Year,
		})
	}

	certs := make([]model.Certificate, 0, len(d.Certificates))
	for _, c := range d.Certificates {
		certs = append(certs, model.Certificate{
			ID:     c.ID.Hex(),
			Name:   c.Name,
			Issuer: c.Issuer,
			Date:   c.Date,
		})
	}

	return model.Portfolio{
		ID:           d.ID.Hex(),
		StudentID:    d.StudentID,
		Headline:     d.Headline,
		Bio:          d.Bio,
		Projects:     projects,
		Certificates: certs,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type portfolioMongoRepository struct {
	coll *mongo.Collection
}

func NewPortfolioRepository(db *mongo.Database) PortfolioRepository {
	return &portfolioMongoRepository{coll: db.Collection("student_portfolios")}
}

func (r *portfolioMongoRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("student_unique"),
		},
		{
			Keys:    bson.D{{Key: "projects.technologies", Value: 1}},
			Options: options.Index().SetName("project_tech"),
		},
	})
	if err != nil {
		return fmt.Errorf("membuat index student_portfolios: %w", err)
	}
	return nil
}

func (r *portfolioMongoRepository) FindByStudent(ctx context.Context, studentID int) (model.Portfolio, error) {
	var doc portfolioDocument
	err := r.coll.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Portfolio{}, ErrNotFound
		}
		return model.Portfolio{}, fmt.Errorf("mengambil portofolio: %w", err)
	}
	return doc.toModel(), nil
}

func (r *portfolioMongoRepository) Replace(ctx context.Context, p model.Portfolio) (model.Portfolio, error) {
	t := now()

	// Proyek dan sertifikat selalu diberi ID baru oleh server
	projDocs := make([]projectDocument, 0, len(p.Projects))
	for _, pr := range p.Projects {
		techs := pr.Technologies
		if techs == nil {
			techs = []string{}
		}
		projDocs = append(projDocs, projectDocument{
			ID:           bson.NewObjectID(),
			Title:        pr.Title,
			Technologies: techs,
			URL:          pr.URL,
			Year:         pr.Year,
		})
	}

	certDocs := make([]certificateDocument, 0, len(p.Certificates))
	for _, c := range p.Certificates {
		certDocs = append(certDocs, certificateDocument{
			ID:     bson.NewObjectID(),
			Name:   c.Name,
			Issuer: c.Issuer,
			Date:   c.Date,
		})
	}

	filter := bson.M{"student_id": p.StudentID}
	setFields := bson.M{
		"headline":     p.Headline,
		"bio":          p.Bio,
		"projects":     projDocs,
		"certificates": certDocs,
		"updated_at":   t,
	}
	setOnInsert := bson.M{
		"_id":        bson.NewObjectID(),
		"student_id": p.StudentID,
		"created_at": t,
	}
	update := bson.M{
		"$set":         setFields,
		"$setOnInsert": setOnInsert,
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var doc portfolioDocument
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		err = r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
		if err == nil {
			return doc.toModel(), nil
		}
		if !mongo.IsDuplicateKeyError(err) {
			break
		}
	}

	if err != nil {
		return model.Portfolio{}, fmt.Errorf("menyimpan portofolio: %w", err)
	}
	return doc.toModel(), nil
}

func (r *portfolioMongoRepository) AddProject(ctx context.Context, studentID int, pr model.Project) (model.Portfolio, error) {
	techs := pr.Technologies
	if techs == nil {
		techs = []string{}
	}

	doc := projectDocument{
		ID:           bson.NewObjectID(),
		Title:        pr.Title,
		Technologies: techs,
		URL:          pr.URL,
		Year:         pr.Year,
	}

	limitField := "projects." + strconv.Itoa(MaxProjects-1)
	filter := bson.M{
		"student_id": studentID,
		limitField:   bson.M{"$exists": false},
	}
	update := bson.M{
		"$push": bson.M{"projects": doc},
		"$set":  bson.M{"updated_at": now()},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated portfolioDocument
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menambah proyek: %w", err)
	}

	count, countErr := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if countErr != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", countErr)
	}
	if count == 0 {
		return model.Portfolio{}, ErrNotFound
	}
	return model.Portfolio{}, ErrProjectLimit
}

func (r *portfolioMongoRepository) RemoveProject(ctx context.Context, studentID int, projectID string) error {
	oid, err := parseID(projectID)
	if err != nil {
		return err
	}

	filter := bson.M{"student_id": studentID, "projects._id": oid}
	update := bson.M{
		"$pull": bson.M{"projects": bson.M{"_id": oid}},
		"$set":  bson.M{"updated_at": now()},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("menghapus proyek: %w", err)
	}

	if res.MatchedCount == 0 {
		n, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
		if err != nil {
			return fmt.Errorf("memeriksa portofolio: %w", err)
		}
		if n == 0 {
			return ErrNotFound
		}
		return ErrProjectNotFound
	}
	return nil
}

func (r *portfolioMongoRepository) TopTechnologies(ctx context.Context, limit int) ([]model.TechCount, error) {
	if limit <= 0 {
		limit = 10
	}

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$projects"}},
		{{Key: "$unwind", Value: "$projects.technologies"}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$projects.technologies",
			"total": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "total", Value: -1},
			{Key: "_id", Value: 1}, // Deterministic tie-breaker
		}}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("menghitung statistik teknologi: %w", err)
	}
	defer cursor.Close(ctx)

	var items []struct {
		Technology string `bson:"_id"`
		Total      int    `bson:"total"`
	}
	if err := cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("membaca statistik teknologi: %w", err)
	}

	result := make([]model.TechCount, 0, len(items))
	for _, it := range items {
		result = append(result, model.TechCount{
			Technology: it.Technology,
			Total:      it.Total,
		})
	}
	return result, nil
}
