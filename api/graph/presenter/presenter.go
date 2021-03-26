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

func GetArticle(ctx context.Context, db *dynamodb.DB, id string) (*model.Article, error) {
	repo := article.NewRepositoryWithContextAndDB(ctx, db)
	entity, err := repo.SelectByID(id)
	if err != nil {
		return nil, err
	}
	return toArticleModel(entity), nil
}

func GetTag(ctx context.Context, db *dynamodb.DB, name string) (*model.Tag, error) {
	repo := tag.NewRepositoryWithContextAndDB(ctx, db)
	entity, err := repo.SelectByName(name)
	if err != nil {
		return nil, err
	}
	return toTagModel(entity), nil
}

func GetArticles(ctx context.Context, db *dynamodb.DB, cond *model.ArticleCondition) ([]*model.Article, error) {
	repo := article.NewRepositoryWithContextAndDB(ctx, db)
	entities, err := repo.SelectAll()
	if err != nil {
		return nil, err
	}
	models := make([]*model.Article, len(entities), len(entities))
	for i, entity := range entities {
		models[i] = toArticleModel(entity)
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
		models[i] = toTagModel(entity)
	}
	return models, nil
}

func PostArticle(ctx context.Context, db *dynamodb.DB, input *model.PostArticle) (*model.Article, error) {
	repoArticle := article.NewRepositoryWithContextAndDB(ctx, db)
	repoTag := tag.NewRepositoryWithContextAndDB(ctx, db)

	var oldTags []string
	var art *article.Article
	if input.ID == nil {
		art = &article.Article{
			ID:      uuid.NewString(),
			Title:   input.Title,
			Content: input.Content,
			Tags:    input.Tags,
		}
	} else {
		var err error
		art, err = repoArticle.SelectByID(*input.ID)
		if err != nil {
			return nil, err
		}
		oldTags = art.Tags
		art.Title = input.Title
		art.Content = input.Content
		art.Tags = input.Tags
	}
	tags, err := makeTagEntities(repoTag, input.Tags, oldTags, art.ID)
	upsertTags := make([]*tag.Tag, 0)
	deleteTags := make([]*tag.Tag, 0)
	for _, t := range tags {
		if len(t.Articles) == 0 {
			deleteTags = append(deleteTags, t)
		} else {
			upsertTags = append(upsertTags, t)
		}
	}

	//TODO: transaction
	art, err = repoArticle.Upsert(art)
	if err != nil {
		return nil, err
	}
	if err := repoTag.UpsertMulti(upsertTags); err != nil {
		return nil, err
	}
	if err := repoTag.DeleteMulti(deleteTags); err != nil {
		return nil, err
	}

	return toArticleModel(art), nil
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

func toArticleModel(a *article.Article) *model.Article {
	return &model.Article{
		ID:        a.ID,
		Title:     a.Title,
		Content:   a.Content,
		Tags:      a.Tags,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}
}

func toTagModel(t *tag.Tag) *model.Tag {
	return &model.Tag{
		Name:     t.Name,
		Articles: t.Articles,
	}
}

func unionTags(a, b []string) []string {
	m := make(map[string]bool)
	for _, t := range a {
		m[t] = true
	}
	for _, t := range b {
		m[t] = true
	}
	n := make([]string, len(m), len(m))
	var i int
	for t, _ := range m {
		n[i] = t
		i++
	}
	return n
}

func diffTags(a, b []string) []string {
	am := make(map[string]bool)
	for _, t := range a {
		am[t] = true
	}
	bm := make(map[string]bool)
	for _, t := range b {
		bm[t] = true
	}
	n := make([]string, 0)
	for aa, _ := range am {
		if _, ok := bm[aa]; !ok {
			n = append(n, aa)
		}
	}
	return n
}

//TODO: write test
func makeTagEntities(repo tag.IRepository, newNames, oldNames []string, articleID string) ([]*tag.Tag, error) {
	tags, err := repo.SelectByNames(unionTags(newNames, oldNames))
	if err != nil {
		return nil, err
	}
	tagMap := make(map[string]*tag.Tag)
	for _, t := range tags {
		tagMap[t.Name] = t
	}

	changedTags := make([]*tag.Tag, 0)
	for _, n := range diffTags(newNames, oldNames) {
		if _, ok := tagMap[n]; ok {
			tagMap[n].Articles = append(tagMap[n].Articles, articleID)
		} else {
			tagMap[n] = &tag.Tag{
				Name:     n,
				Articles: []string{articleID},
			}
		}
		changedTags = append(changedTags, tagMap[n])
	}
	for _, n := range diffTags(oldNames, newNames) {
		for i, a := range tagMap[n].Articles {
			if a != articleID {
				continue
			}
			tagMap[n].Articles = append(tagMap[n].Articles[:i], tagMap[n].Articles[i+1:]...)
			changedTags = append(changedTags, tagMap[n])
			break
		}
	}
	return changedTags, nil
}
