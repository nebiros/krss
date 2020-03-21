package entity

type Feed struct {
	FeedID int    `db:"feed.feed_id"`
	UserID int    `db:"feed.user_id"`
	Title  string `db:"feed.title"`
	URL    string `db:"feed.url"`
}

type Feeds []Feed

type CreateFeed struct {
	UserID int    `db:"feed.user_id"`
	Title  string `db:"feed.title"`
	URL    string `db:"feed.url"`
}
