package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"compro-backend/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("item not found")

// DB holds the DynamoDB client and table names.
type DB struct {
	Client        *dynamodb.Client
	ContentTable  string
	NewsTable     string
	AdminsTable   string
	NewsSlugIndex string
	AdminUserIndex string
}

// NewDB creates a DB using the DynamoDB client and reads table names from env vars.
func NewDB(client *dynamodb.Client) *DB {
	return &DB{
		Client:         client,
		ContentTable:   getEnv("DYNAMODB_CONTENT_TABLE", "compro-content"),
		NewsTable:      getEnv("DYNAMODB_NEWS_TABLE", "compro-news"),
		AdminsTable:    getEnv("DYNAMODB_ADMINS_TABLE", "compro-admins"),
		NewsSlugIndex:  getEnv("DYNAMODB_NEWS_SLUG_INDEX", "slug-index"),
		AdminUserIndex: getEnv("DYNAMODB_ADMIN_USERNAME_INDEX", "username-index"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ---------------------------------------------------------------------------
// Content operations
// ---------------------------------------------------------------------------

const contentID = "site-content"

func (d *DB) GetContent(ctx context.Context) (*models.SiteContent, error) {
	out, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.ContentTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: contentID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb GetItem content: %w", err)
	}
	if out.Item == nil {
		return nil, ErrNotFound
	}
	var content models.SiteContent
	if err := attributevalue.UnmarshalMap(out.Item, &content); err != nil {
		return nil, fmt.Errorf("unmarshal content: %w", err)
	}
	return &content, nil
}

func (d *DB) PutContent(ctx context.Context, content *models.SiteContent) error {
	content.ID = contentID
	item, err := attributevalue.MarshalMap(content)
	if err != nil {
		return fmt.Errorf("marshal content: %w", err)
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.ContentTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("dynamodb PutItem content: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// News operations
// ---------------------------------------------------------------------------

func (d *DB) ListNews(ctx context.Context) ([]models.NewsPost, error) {
	out, err := d.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(d.NewsTable),
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb Scan news: %w", err)
	}
	var posts []models.NewsPost
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &posts); err != nil {
		return nil, fmt.Errorf("unmarshal news list: %w", err)
	}
	// Sort descending by PublishedAt in-memory (Scan doesn't guarantee order)
	sortNewsDesc(posts)
	return posts, nil
}

func (d *DB) GetNewsByID(ctx context.Context, id string) (*models.NewsPost, error) {
	out, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.NewsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb GetItem news: %w", err)
	}
	if out.Item == nil {
		return nil, ErrNotFound
	}
	var post models.NewsPost
	if err := attributevalue.UnmarshalMap(out.Item, &post); err != nil {
		return nil, fmt.Errorf("unmarshal news: %w", err)
	}
	return &post, nil
}

func (d *DB) GetNewsBySlug(ctx context.Context, slug string) (*models.NewsPost, error) {
	out, err := d.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.NewsTable),
		IndexName:              aws.String(d.NewsSlugIndex),
		KeyConditionExpression: aws.String("slug = :s"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: slug},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb Query news by slug: %w", err)
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	var post models.NewsPost
	if err := attributevalue.UnmarshalMap(out.Items[0], &post); err != nil {
		return nil, fmt.Errorf("unmarshal news: %w", err)
	}
	return &post, nil
}

func (d *DB) CreateNews(ctx context.Context, post *models.NewsPost) error {
	post.ID = uuid.New().String()
	item, err := attributevalue.MarshalMap(post)
	if err != nil {
		return fmt.Errorf("marshal news: %w", err)
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.NewsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("dynamodb PutItem news: %w", err)
	}
	return nil
}

func (d *DB) UpdateNews(ctx context.Context, post *models.NewsPost) error {
	item, err := attributevalue.MarshalMap(post)
	if err != nil {
		return fmt.Errorf("marshal news: %w", err)
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.NewsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("dynamodb PutItem news update: %w", err)
	}
	return nil
}

func (d *DB) DeleteNews(ctx context.Context, id string) error {
	_, err := d.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.NewsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb DeleteItem news: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Admin operations
// ---------------------------------------------------------------------------

func (d *DB) GetAdminByUsername(ctx context.Context, username string) (*models.Admin, error) {
	out, err := d.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(d.AdminsTable),
		IndexName:              aws.String(d.AdminUserIndex),
		KeyConditionExpression: aws.String("username = :u"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":u": &types.AttributeValueMemberS{Value: username},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb Query admin by username: %w", err)
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	var admin models.Admin
	if err := attributevalue.UnmarshalMap(out.Items[0], &admin); err != nil {
		return nil, fmt.Errorf("unmarshal admin: %w", err)
	}
	return &admin, nil
}

func (d *DB) CountAdmins(ctx context.Context) (int64, error) {
	out, err := d.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:          aws.String(d.AdminsTable),
		Select:             types.SelectCount,
	})
	if err != nil {
		return 0, fmt.Errorf("dynamodb Scan admins count: %w", err)
	}
	return int64(out.Count), nil
}

func (d *DB) CreateAdmin(ctx context.Context, admin *models.Admin) error {
	admin.ID = uuid.New().String()
	item, err := attributevalue.MarshalMap(admin)
	if err != nil {
		return fmt.Errorf("marshal admin: %w", err)
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.AdminsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("dynamodb PutItem admin: %w", err)
	}
	return nil
}

func (d *DB) ListAdmins(ctx context.Context) ([]models.Admin, error) {
	out, err := d.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(d.AdminsTable),
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb Scan admins: %w", err)
	}
	var admins []models.Admin
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &admins); err != nil {
		return nil, fmt.Errorf("unmarshal admins: %w", err)
	}
	return admins, nil
}

func (d *DB) GetAdminByID(ctx context.Context, id string) (*models.Admin, error) {
	out, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.AdminsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb GetItem admin: %w", err)
	}
	if out.Item == nil {
		return nil, ErrNotFound
	}
	var admin models.Admin
	if err := attributevalue.UnmarshalMap(out.Item, &admin); err != nil {
		return nil, fmt.Errorf("unmarshal admin: %w", err)
	}
	return &admin, nil
}

func (d *DB) DeleteAdmin(ctx context.Context, id string) error {
	_, err := d.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.AdminsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb DeleteItem admin: %w", err)
	}
	return nil
}

// sortNewsDesc sorts news posts descending by PublishedAt (in-place).
func sortNewsDesc(posts []models.NewsPost) {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[j].PublishedAt > posts[i].PublishedAt {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
}
