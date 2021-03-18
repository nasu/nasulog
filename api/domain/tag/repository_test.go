package tag

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nasu/nasulog/domain/article"
	"github.com/nasu/nasulog/infrastructure/dynamodb"
)

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

func TestInsertMultiAndDeleteByPK(t *testing.T) {
	ctx := context.TODO()
	db, err := dynamodb.GetDB(ctx)
	if err != nil {
		t.Errorf("failed to get db client. err=%v", err)
	}

	// insert
	tags := []string{
		"unit-test-" + uuid.NewString(),
		"unit-test-" + uuid.NewString()}
	if err := InsertMulti(ctx, db, tags); err != nil {
		t.Fatalf("failed to insert multi. err=%v", err)
	}

	// confirm
	if got, err := SelectAll(ctx, db); err != nil {
		t.Fatalf("failed to select all. err=%v", err)
	} else {
		tagMap := make(map[string]bool)
		for _, tag := range tags {
			tagMap[tag] = true
		}
		cnt := 0
		for _, tag := range got {
			if _, ok := tagMap[tag]; ok {
				cnt++
			}
		}
		if cnt != len(tagMap) {
			t.Errorf("tag mismatch. (got,want)=(%d,%d)", cnt, len(tagMap))
		}
	}

	// delete
	for _, tag := range tags {
		if err := DeleteByPK(ctx, db, tag); err != nil {
			t.Fatalf("failed to delete. err=%v", err)
		}
	}
	if got, err := SelectAll(ctx, db); err != nil {
		t.Fatalf("failed to select all. err=%v", err)
	} else {
		tagMap := make(map[string]bool)
		for _, tag := range tags {
			tagMap[tag] = true
		}
		cnt := 0
		for _, tag := range got {
			if _, ok := tagMap[tag]; ok {
				cnt++
			}
		}
		if cnt != 0 {
			t.Errorf("tag mismatch. (got,want)=(%d,%d)", cnt, len(tagMap))
		}

		article.SelectAll(ctx, db)
	}
}
