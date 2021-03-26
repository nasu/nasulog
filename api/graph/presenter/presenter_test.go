package presenter

import (
	"reflect"
	"sort"
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/nasu/nasulog/api/domain/tag"
)

func TestUnion(t *testing.T) {
	tests := []struct {
		a    []string
		b    []string
		want []string
	}{
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			[]string{"a", "b"},
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			[]string{"a", "b", "c"},
			[]string{"a"},
			[]string{"a", "b", "c"},
		},
		{
			[]string{"a", "c"},
			[]string{"b"},
			[]string{"a", "b", "c"},
		},
		{
			[]string{},
			[]string{"a"},
			[]string{"a"},
		},
		{
			[]string{"a"},
			[]string{},
			[]string{"a"},
		},
		{
			[]string{},
			[]string{},
			[]string{},
		},
	}
	for _, tt := range tests {
		got := unionTags(tt.a, tt.b)
		sort.Strings(got)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("mismatch. a=%v, b =%v, got=%v, want=%v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		a    []string
		b    []string
		want []string
	}{
		{
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
			[]string{},
		},
		{
			[]string{"a", "b"},
			[]string{"a", "b", "c"},
			[]string{},
		},
		{
			[]string{"a", "b", "c"},
			[]string{"a"},
			[]string{"b", "c"},
		},
		{
			[]string{"a", "c"},
			[]string{"b"},
			[]string{"a", "c"},
		},
		{
			[]string{},
			[]string{"a"},
			[]string{},
		},
		{
			[]string{"a"},
			[]string{},
			[]string{"a"},
		},
		{
			[]string{},
			[]string{},
			[]string{},
		},
	}
	for _, tt := range tests {
		got := diffTags(tt.a, tt.b)
		sort.Strings(got)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("mismatch. a=%v, b =%v, got=%v, want=%v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMakeTagEntities(t *testing.T) {
	tests := []struct {
		newNames  []string
		oldNames  []string
		articleID string
		want      []*tag.Tag
	}{
		{
			[]string{"a", "b", "c", "d"},
			[]string{},
			"article-3",
			[]*tag.Tag{
				{"a", []string{"article-3"}},
				{"b", []string{"article-1", "article-3"}},
				{"c", []string{"article-1", "article-2", "article-3"}},
				{"d", []string{"article-3"}},
			},
		},
		{
			[]string{"a"},
			[]string{"a"},
			"article-3",
			[]*tag.Tag{},
		},
		{
			[]string{"b", "c"},
			[]string{"c"},
			"article-2",
			[]*tag.Tag{
				{"b", []string{"article-1", "article-2"}},
			},
		},
		{
			[]string{"c"},
			[]string{"b", "c"},
			"article-1",
			[]*tag.Tag{
				{"b", []string{}},
			},
		},
	}

	repo := &TestTagRepository{}
	for _, tt := range tests {
		got, err := makeTagEntities(repo, tt.newNames, tt.oldNames, tt.articleID)
		if err != nil {
			t.Error(err)
			continue
		}
		sort.SliceStable(got, func(i, j int) bool {
			return got[i].Name < got[j].Name
		})
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("mismatch. %s", pretty.Compare(got, tt.want))
		}
	}
}

type TestTagRepository struct{}

func (r *TestTagRepository) SelectByNames(names []string) ([]*tag.Tag, error) {
	tags := []*tag.Tag{
		{"a", []string{}},
		{"b", []string{"article-1"}},
		{"c", []string{"article-1", "article-2"}},
	}
	tagsMap := make(map[string]*tag.Tag)
	for _, t := range tags {
		tagsMap[t.Name] = t
	}
	res := make([]*tag.Tag, 0)
	for _, n := range names {
		if t, ok := tagsMap[n]; ok {
			res = append(res, t)
		}
	}
	return res, nil
}
