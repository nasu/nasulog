package article

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nasu/nasulog/infrastructure/dynamodb"
)

func TestInsertAndSelectByPK(t *testing.T) {
	ctx := context.TODO()
	db, err := dynamodb.GetDB(ctx)
	if err != nil {
		t.Fatalf("failed to get db client. err=%v", err)
	}

	// case: no hit
	entity, err := SelectByID(ctx, db, uuid.NewString())
	if err != nil {
		t.Fatalf("failed to get an entity. err=%v", err)
	}
	if entity != nil {
		t.Fatalf("should be no hit. entity=%+v", entity)
	}

	// case: hit
	entity, err = Insert(ctx, db, makeDummyArticle())
	if err != nil {
		t.Fatalf("failed to insert an entity. err=%v", err)
	}
	got, err := SelectByID(ctx, db, entity.ID)
	if err != nil {
		t.Fatalf("failed to get an entity. err=%v", err)
	}
	if got.ID != entity.ID {
		t.Errorf("wrong ID")
	}
}

func TestSelectAll(t *testing.T) {
	ctx := context.TODO()
	db, err := dynamodb.GetDB(ctx)
	if err != nil {
		t.Fatalf("failed to get db client. err=%v", err)
	}

	_, err = SelectAll(ctx, db)
	if err != nil {
		t.Fatalf("failed to select all. err=%v", err)
	}
}

func TestInsertAndDeleteByPK(t *testing.T) {
	ctx := context.TODO()
	db, err := dynamodb.GetDB(ctx)
	if err != nil {
		t.Errorf("failed to get db client. err=%v", err)
	}

	// insert
	entity, err := Insert(ctx, db, makeDummyArticle())
	if err != nil {
		t.Fatalf("failed to insert an entity. err=%v", err)
	}

	// confirmation
	got, err := SelectByID(ctx, db, entity.ID)
	if err != nil {
		t.Fatalf("failed to get an entity. err=%v", err)
	}
	if got.ID != entity.ID {
		t.Fatalf("wrong ID. (got,want)=(%s,%s)", got.ID, entity.ID)
	}

	// delete
	if DeleteByID(ctx, db, entity.ID) != nil {
		t.Fatalf("failed to delete an entity. err=%v", err)
	}
	got, err = SelectByID(ctx, db, entity.ID)
	if err != nil {
		t.Fatalf("failed to get an entity")
	}
	if got != nil {
		t.Fatalf("failed to delete")
	}
}

func makeDummyArticle() *Article {
	return &Article{
		Title:   "unit test",
		Content: "This article is for unit test.",
		Tags:    []string{uuid.NewString()},
	}
}
