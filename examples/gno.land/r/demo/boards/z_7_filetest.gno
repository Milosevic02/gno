// PKGPATH: gno.land/r/demo/boards_test
package boards_test

// SEND: 200000000ugnot

import (
	"std"

	"gno.land/p/demo/testutils"
	"gno.land/r/demo/boards"
	"gno.land/r/demo/users"
)

func init() {
	caller := testutils.TestAddress("caller")
	std.TestSetRealm(std.NewUserRealm(caller))
	// register
	users.Register("", "gnouser", "my profile")

	// create board and post
	bid := boards.CreateBoard("test_board")
	boards.CreateThread(bid, "First Post (title)", "Body of the first post. (body)")
}

func main() {
	println(boards.Render("test_board"))
}

// Output:
// \[[post](/r/demo/boards$help&func=CreateThread&bid=1)]
//
// ----------------------------------------
// ## [First Post (title)](/r/demo/boards:test_board/1)
//
// Body of the first post. (body)
// \- [@gnouser](/r/demo/users:gnouser), [2009-02-13 11:31pm UTC](/r/demo/boards:test_board/1) \[[x](/r/demo/boards$help&func=DeletePost&bid=1&postid=1&threadid=1)] (0 replies) (0 reposts)
//
//
