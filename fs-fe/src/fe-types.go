package main

type IGAuthCred struct {
	Token string `json:"access_token"`
	User  struct {
		Username string `json:"username"`
		Id       string `json:"id"`
	} `json:"user"`
}

type Threads struct {
	Pagination struct {
	} `json:"pagination"`
	Data []struct {
		ID   string `json:"id"`
		User struct {
			ID             string `json:"id"`
			FullName       string `json:"full_name"`
			ProfilePicture string `json:"profile_picture"`
			Username       string `json:"username"`
		} `json:"user"`
		Images struct {
			Thumbnail struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"thumbnail"`
			LowResolution struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"low_resolution"`
			StandardResolution struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"standard_resolution"`
		} `json:"images"`
		CreatedTime string `json:"created_time"`
		Caption     struct {
			ID          string `json:"id"`
			Text        string `json:"text"`
			CreatedTime string `json:"created_time"`
			From        struct {
				ID             string `json:"id"`
				FullName       string `json:"full_name"`
				ProfilePicture string `json:"profile_picture"`
				Username       string `json:"username"`
			} `json:"from"`
		} `json:"caption"`
		UserHasLiked bool `json:"user_has_liked"`
		Likes        struct {
			Count int `json:"count"`
		} `json:"likes"`
		Tags     []interface{} `json:"tags"`
		Filter   string        `json:"filter"`
		Comments struct {
			Count int `json:"count"`
		} `json:"comments"`
		Type         string        `json:"type"`
		Link         string        `json:"link"`
		Location     interface{}   `json:"location"`
		Attribution  interface{}   `json:"attribution"`
		UsersInPhoto []interface{} `json:"users_in_photo"`
	} `json:"data"`
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
}

type ThreadsContent struct {
	Threads map[string][]string
}

type SentimentContent struct {
	Sentiment []string
}
