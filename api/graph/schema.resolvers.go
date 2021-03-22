package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nasu/nasulog/api/graph/generated"
	"github.com/nasu/nasulog/api/graph/model"
	"github.com/nasu/nasulog/api/graph/presenter"
)

func (r *mutationResolver) CreateArticle(ctx context.Context, input model.NewArticle) (*model.Article, error) {
	return presenter.CreateArticle(ctx, r.DB, &input)
}

func (r *mutationResolver) UpdateArticle(ctx context.Context, input model.UpdateArticle) (*model.Article, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteArticle(ctx context.Context, id string) (*bool, error) {
	return presenter.DeleteArticle(ctx, r.DB, id)
}

func (r *mutationResolver) DeleteTag(ctx context.Context, name string) (*bool, error) {
	return presenter.DeleteTag(ctx, r.DB, name)
}

func (r *queryResolver) Articles(ctx context.Context, cond *model.ArticleCondition) ([]*model.Article, error) {
	return presenter.GetArticles(ctx, r.DB, cond)
}

func (r *queryResolver) Article(ctx context.Context, id string) (*model.Article, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Tags(ctx context.Context) ([]*model.Tag, error) {
	return presenter.GetTags(ctx, r.DB)
}

func (r *queryResolver) Tag(ctx context.Context, name string) (*model.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
