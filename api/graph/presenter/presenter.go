package presenter

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nasu/nasulog/api/domain/article"
	"github.com/nasu/nasulog/api/domain/tag"
	"github.com/nasu/nasulog/api/graph/model"
	"github.com/nasu/nasulog/api/infrastructure/dynamodb"
)

func GetArticles(ctx context.Context, db *dynamodb.DB, cond *model.ArticleCondition) ([]*model.Article, error) {
	repo := article.NewRepositoryWithContextAndDB(ctx, db)
	entities, err := repo.SelectAll()
	if err != nil {
		return nil, err
	}
	models := make([]*model.Article, len(entities), len(entities))
	for i, entity := range entities {
		models[i] = &model.Article{
			ID:        entity.ID,
			Title:     entity.Title,
			Content:   entity.Content,
			Tags:      entity.Tags,
			CreatedAt: entity.CreatedAt.Format(time.RFC3339),
			UpdatedAt: entity.UpdatedAt.Format(time.RFC3339),
		}
	}
	return models, nil
}

func GetTags(ctx context.Context, db *dynamodb.DB) ([]*model.Tag, error) {
	repo := tag.NewRepositoryWithContextAndDB(ctx, db)
	entities, err := repo.SelectAll()
	if err != nil {
		return nil, err
	}
	models := make([]*model.Tag, len(entities), len(entities))
	for i, entity := range entities {
		models[i] = &model.Tag{
			Name:     entity.Name,
			Articles: entity.Articles,
		}
	}
	return models, nil
}

func CreateArticle(ctx context.Context, db *dynamodb.DB, input *model.NewArticle) (*model.Article, error) {
	art := &article.Article{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	}
	repoArticle := article.NewRepositoryWithContextAndDB(ctx, db)
	repoTag := tag.NewRepositoryWithContextAndDB(ctx, db)

	//TODO: transaction
	art.ID = uuid.NewString()
	entity, err := repoArticle.Insert(art)
	if err != nil {
		return nil, err
	}

	tags, err := repoTag.SelectByNames(art.Tags)
	if err != nil {
		return nil, err
	}
	tagsMap := make(map[string]*tag.Tag)
	for _, t := range tags {
		tagsMap[t.Name] = t
	}
	for _, it := range art.Tags {
		if et, ok := tagsMap[it]; ok {
			et.Articles = append(et.Articles, entity.ID)
		} else {
			tags = append(tags, &tag.Tag{
				Name:     it,
				Articles: []string{entity.ID},
			})
		}
	}

	if err := repoTag.InsertMulti(tags); err != nil {
		return nil, err
	}

	return &model.Article{
		ID:        art.ID,
		Title:     art.Title,
		Content:   art.Content,
		Tags:      art.Tags,
		CreatedAt: art.CreatedAt.Format(time.RFC3339),
		UpdatedAt: art.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func DeleteArticle(ctx context.Context, db *dynamodb.DB, id string) (*bool, error) {
	repoArticle := article.NewRepositoryWithContextAndDB(ctx, db)
	repoTag := tag.NewRepositoryWithContextAndDB(ctx, db)

	entity, err := repoArticle.SelectByID(id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		res := false
		return &res, nil
	}

	// Remove an article matched to the article from tag.articles
	tags, err := repoTag.SelectByNames(entity.Tags)
	if err != nil {
		return nil, err
	}
	tagsShouldDelete := make([]*tag.Tag, 0)
	tagsShouldUpdate := make([]*tag.Tag, 0)
	for _, t := range tags {
		for i, a := range t.Articles {
			if a == id {
				t.Articles = append(t.Articles[:i], t.Articles[i+1:]...)
				break
			}
		}
		if len(t.Articles) == 0 {
			tagsShouldDelete = append(tagsShouldDelete, t)
		} else {
			tagsShouldUpdate = append(tagsShouldUpdate, t)
		}
	}

	if err := repoArticle.DeleteByID(id); err != nil {
		return nil, err
	}
	if err := repoTag.UpsertMulti(tagsShouldUpdate); err != nil {
		return nil, err
	}
	if err := repoTag.DeleteMulti(tagsShouldDelete); err != nil {
		return nil, err
	}
	res := true
	return &res, nil
}

func DeleteTag(ctx context.Context, db *dynamodb.DB, name string) (*bool, error) {
	repoArticle := article.NewRepositoryWithContextAndDB(ctx, db)
	repoTag := tag.NewRepositoryWithContextAndDB(ctx, db)

	entity, err := repoTag.SelectByName(name)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		res := false
		return &res, nil
	}

	// Remove a tag matched to the tag from article.tags
	articles, err := repoArticle.SelectByIDs(entity.Articles)
	if err != nil {
		return nil, err
	}
	for _, a := range articles {
		for i, t := range a.Tags {
			if t == name {
				a.Tags = append(a.Tags[:i], a.Tags[i+1:]...)
				break
			}
		}
	}

	// transaction
	if err := repoTag.DeleteByName(name); err != nil {
		return nil, err
	}
	if err := repoArticle.UpsertMulti(articles); err != nil {
		return nil, err
	}
	res := true
	return &res, err
}
