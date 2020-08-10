package model

import "github.com/jinzhu/gorm"

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (a Article) Update(db *gorm.DB, values interface{}) error {
	if err := db.Model(&a).Updates(values).Where("id = ? and is_del = ?", a.ID, a.IsDel).Error; err != nil {
		return err
	}
	return nil
}

//func (a Article) Get(db *gorm.DB) (Article, error) {
//	var article Article
//	db = db.Where("id = ? and state = ? and is_del = ?", article.ID, article.State, 0)
//	err := db.First(&article).Error
//	if err != nil {
//		return article, err
//	}
//	return article, nil
//}

func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? AND state = ? AND is_del = ?", a.ID, a.State, 0)
	err := db.First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}

	return article, nil
}

func (a Article) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? and is_del = ?", a.Model.ID, 0).Delete(&a).Error; err != nil {
		return err
	}
	return nil
}

type ArticleRow struct {
	ArticleID     uint32
	TagID         uint32
	TagName       string
	ArticleTitle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
}

func (a Article) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int) ([]*ArticleRow, error) {
	fields := []string{"ar.id AS article_id", "ar.title AS article_title", "ar.desc AS article_desc", "ar.cover_image_url", "ar.content"}
	fields = append(fields, []string{"t.id AS tag_id", "t.name AS tag_name"}...)

	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	rows, err := db.Select(fields).Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		if err := rows.Scan(&r.ArticleID, &r.ArticleTitle, &r.ArticleDesc, &r.CoverImageUrl, &r.Content, &r.TagID, &r.TagName); err != nil {
			return nil, err
		}

		articles = append(articles, r)
	}

	return articles, nil
}

func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a ArticleTag) GetByID(db *gorm.DB) (ArticleTag, error) {
	var articleTag ArticleTag
	err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, 0).First(&articleTag).Error
	if err != nil {
		return articleTag, err
	}
	return articleTag, nil
}

func (a ArticleTag) ListByTID(db *gorm.DB) ([]*ArticleTag, error) {
	var articleTags []*ArticleTag
	if err := db.Where("article_id in (?) AND is_del = ?", articleTags, 0).Find(&articleTags).Error; err != nil {
		return nil, err
	}
	return articleTags, nil
}

func (a ArticleTag) ListByAIDs(db *gorm.DB, articleIDs []uint32) ([]*ArticleTag, error) {
	var articleTags []*ArticleTag
	err := db.Where("article_id in (?) and is_del = ?", articleIDs, 0).Find(&articleTags).Error
	if err != nil {
		return nil, err
	}
	return articleTags, nil
}

func (a ArticleTag) Create(db *gorm.DB) error {
	if err := db.Create(&a).Error; err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) UpdateOne(db *gorm.DB, values interface{}) error {
	err := db.Model(&a).Where("article_id = ? AND is_del = ?", a.ArticleID, 0).
		Limit(1).Updates(values).Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) Delete(db *gorm.DB) error {
	err := db.Where("id = ? AND is_del = ?", a.Model.ID, 0).Delete(&a).Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleTag) DeleteOne(db *gorm.DB) error {
	if err := db.Where("article_id = ? AND is_del = ?", a.ArticleID, 0).
		Delete(&a).Limit(1).Error; err != nil {
		return err
	}
	return nil
}
