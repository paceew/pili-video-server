package dbops

import (
	"fmt"
	"testing"
)

func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	// clearTables()
	m.Run()
	// clearTables()
}

//test like funtion
// func TestLikeVideo(t *testing.T) {
// 	t.Run("fist get like count;\n", testLikeVideoCount)
// 	t.Run("fist get islike;\n", testIsLike)
// 	t.Run("fist like video;\n", testLikeVideo)
// 	t.Run("secend get like count;\n", testLikeVideoCount)
// 	t.Run("secend get islike;\n", testIsLike)
// 	t.Run("secend like video;\n", testLikeVideo)
// 	t.Run("third get like count;\n", testLikeVideoCount)
// 	t.Run("third get islike;\n", testIsLike)
// }

func testLikeVideoCount(t *testing.T) {
	count, err := LikeCount("567")
	if err != nil {
		t.Errorf("Error of like count: %v", err)
	}
	fmt.Printf("video: 567 like count :%v", count)
}

func testIsLike(t *testing.T) {
	yes, err := IsLike("567", "test")
	if err != nil {
		t.Errorf("Error of islike: %v", err)
	}
	fmt.Printf("video islike:%v", yes)
}

func testLikeVideo(t *testing.T) {
	err := LikeVideo("567", "test")
	if err != nil {
		t.Errorf("Error of like video: %v", err)
	}
	fmt.Printf("test like video 567\n")
}

//test user work flow
// func TestUserWorkFlow(t *testing.T) {
// 	t.Run("Add", testAddUser)
// 	t.Run("Get", testGetUser)
// 	t.Run("Modify", testModifyPwd)
// 	t.Run("Credential", testGetUserCredential)
// }

func testAddUser(t *testing.T) {
	err := AddUser("pace", "123")
	if err != nil {
		t.Errorf("Error of add user: %v", err)
	}
}

func testGetUser(t *testing.T) {
	user, err := GetUser("pace")
	if err != nil {
		t.Errorf("Get user error : %v", err)
	}
	fmt.Printf("user_id: %d", user.Id)
	fmt.Printf("user_name: %s", user.Username)
	fmt.Printf("user_pwd: %s", user.Pwd)
}

func testModifyPwd(t *testing.T) {
	err := ModifyUserPwd("pace", "345")
	if err != nil {
		t.Errorf("modify user pwd error: %v", err)
	}
}

func testGetUserCredential(t *testing.T) {
	pwd, err := GetUserCredential("pace")
	if err != nil {
		t.Errorf("Get user credantial error: %v", err)
	}
	fmt.Printf("user pace pwd: %s", pwd)
}

//test session work flow
// func TestSessionWorkFlow(t *testing.T) {
// 	t.Run("inser session", testInserSession)
// 	t.Run("read session", testRetrieveSession)
// 	t.Run("read all session", testRetrieveAllSession)
// 	t.Run("delete session", testDeletSession)
// }

// func testInserSession(t *testing.T) {
// 	err := InserSession("1", 312, "pace")
// 	if err != nil {
// 		t.Errorf("Inser session error: %v", err)
// 	}
// }

// func testRetrieveSession(t *testing.T) {
// 	session, err := RetrieveSession("1")
// 	if err != nil {
// 		t.Errorf("Retrieve session error: %v", err)
// 	}
// 	fmt.Printf("session: %v", session)
// }

// func testRetrieveAllSession(t *testing.T) {
// 	m, err := RetrieveAllSession()
// 	if err != nil {
// 		t.Errorf("retrieve all session error: %v", err)
// 	}

// 	m.Range(func(k, v interface{}) bool {
// 		fmt.Printf("session_id : %d,session: %v", k, v)
// 		return true
// 	})

// }

// func testDeletSession(t *testing.T) {
// 	err := DeleteSession("1")
// 	if err != nil {
// 		t.Errorf("delete session error: %v", err)
// 	}
// }

var (
	vid string
)

// test video work flow
func TestVideoWorkFlow(t *testing.T) {
	// t.Run("add new video", testAddNewVideo)
	// t.Run("get video", testGetVideoInfo)
	// t.Run("list video", testListVideoInfo)
	// t.Run("list video by mod hot", testListVideoInfoMod)
	// t.Run("delete video", testDeleteVideoInfo)
	t.Run("search video", testVdieoSearch)

}

func testVdieoSearch(t *testing.T) {
	_, err := VideoSearch("柯", 0, 10)
	if err != nil {
		t.Errorf("list videoinfo error:%v\n", err)
	}
}

// func testAddNewVideo(t *testing.T) {
// 	_, err := AddNewVideo(7, "video_1")
// 	if err != nil {
// 		t.Errorf("add new video error: %v\n", err)
// 	}
// 	//	vid = videoinfo.Id
// }

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(vid)
	if err != nil {
		t.Errorf("get videoinfo error: %v\n", err)
	}
}

func testListVideoInfo(t *testing.T) {
	_, err := ListVideoInfo("test", 0, 100)
	if err != nil {
		t.Errorf("list videoinfo error:%v\n", err)
	}
}

func testListVideoInfoMod(t *testing.T) {
	_, err := ListVideoInfoMod("other", 0, 100, "hot")
	if err != nil {
		t.Errorf("list videoinfo error:%v\n", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(vid)
	if err != nil {
		t.Errorf("delete videoinfo error: %v\n", err)
	}
}

//test comment work flow
// func TestCommentWorkFlow(t *testing.T) {
// 	t.Run("add comment", testAddNewComment)
// 	t.Run("list comment", testListComment)
// 	//	t.Run("delete comment", testDeleteComment)
// }

func testAddNewComment(t *testing.T) {
	err := AddNewComment(7, "5616", "good")
	if err != nil {
		t.Errorf("add new comment error:%v \n", err)
	}
}

func testListComment(t *testing.T) {
	_, err := ListComments("5616", 0, 100)
	if err != nil {
		t.Errorf("list comment error: %v\n", err)
	}
}

func testDeleteComment(t *testing.T) {
	//	err := DeleteComment()
}
