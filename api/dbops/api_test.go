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
	clearTables()
	m.Run()
	clearTables()
}

//test user work flow
func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Modify", testModifyPwd)
	t.Run("Credential", testGetUserCredential)
}

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
func TestSessionWorkFlow(t *testing.T) {
	t.Run("inser session", testInserSession)
	t.Run("read session", testRetrieveSession)
	t.Run("read all session", testRetrieveAllSession)
	t.Run("delete session", testDeletSession)
}

func testInserSession(t *testing.T) {
	err := InserSession("1", 312, "pace")
	if err != nil {
		t.Errorf("Inser session error: %v", err)
	}
}

func testRetrieveSession(t *testing.T) {
	session, err := RetrieveSession("1")
	if err != nil {
		t.Errorf("Retrieve session error: %v", err)
	}
	fmt.Printf("session: %v", session)
}

func testRetrieveAllSession(t *testing.T) {
	m, err := RetrieveAllSession()
	if err != nil {
		t.Errorf("retrieve all session error: %v", err)
	}

	m.Range(func(k, v interface{}) bool {
		fmt.Printf("session_id : %d,session: %v", k, v)
		return true
	})

}

func testDeletSession(t *testing.T) {
	err := DeleteSession("1")
	if err != nil {
		t.Errorf("delete session error: %v", err)
	}
}

var (
	vid string
)

//test video work flow
// func TestVideoWorkFlow(t *testing.T) {
// 	t.Run("add new video", testAddNewVideo)
// 	t.Run("get video", testGetVideoInfo)
// 	t.Run("list video", testListVideoInfo)
// 	t.Run("delete video", testDeleteVideoInfo)

// }

func testAddNewVideo(t *testing.T) {
	err := AddNewVideo(7, "video_1")
	if err != nil {
		t.Errorf("add new video error: %v\n", err)
	}
	//	vid = videoinfo.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(vid)
	if err != nil {
		t.Errorf("get videoinfo error: %v\n", err)
	}
}

func testListVideoInfo(t *testing.T) {
	_, err := ListVideoInfo("pace-wang", 1, 100)
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
func TestCommentWorkFlow(t *testing.T) {
	t.Run("add comment", testAddNewComment)
	t.Run("list comment", testListComment)
	//	t.Run("delete comment", testDeleteComment)
}

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
