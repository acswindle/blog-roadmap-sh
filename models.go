package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "articles.db")
	if err != nil {
		log.Fatalf("could not open database %v", err)
	}
	_, err = db.Exec(`

		CREATE TABLE IF NOT EXISTS Articles (
		id integer primary key autoincrement,
		title text not null,
		date datetime default current_timestamp,
		content text not null
		);
		`)
	if err != nil {
		log.Fatalf("could not create table articles %v", err)
	}
	return db
}

func (app *application) AddArticle(article Article) error {
	_, err := app.db.Exec(`
		insert into Articles (title,content) values 
		(?,?)
		;
		`, article.Title, article.Content)
	return err
}

func (app *application) GetArticles() (Articles, error) {
	rows, err := app.db.Query(`
		select id, title, date, content 
		from Articles
		order by date DESC;
		`)
	if err != nil {
		return nil, err
	}
	articles := []Article{}
	for rows.Next() {
		article := Article{}
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Date,
			&article.Content,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (app *application) GetArticle(id int) (Article, error) {
	rows, err := app.db.Query(`
		select id,title, date, content
		from Articles
		where id = ?
		;
		`, id)
	if err != nil {
		return Article{}, nil
	}
	if hasNext := rows.Next(); !hasNext {
		return Article{}, fmt.Errorf("no article with specified id %d", id)
	}
	article := Article{}
	err = rows.Scan(
		&article.ID,
		&article.Title,
		&article.Date,
		&article.Content,
	)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}
