package minisocial

import (
	"errors" // For handling errors
	"std"    // The standard Gno package
	"time"   // For handling time
)

var posts []*Post

// CreatePost creates a new post
func CreatePost(text string) error {
	// If the body of the post is empty, return an error
	if text == "" {
		return errors.New("empty post text")
	}

	// Append the new post to the list
	posts = append(posts, &Post{
		text:      text,                          // Set the input text
		author:    std.PreviousRealm().Address(), // The author of the address is the previous realm, the realm that called this one
		createdAt: time.Now(),                    // Capture the time of the transaction, in this case the block timestamp
	})

	return nil
}
