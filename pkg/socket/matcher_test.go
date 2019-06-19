package socket

//import (
//	"testing"
//
//	"github.com/nn/hobee/models"
//)

//const users = 50000
//
//func fillUsers() []*models.WSUser{
//	us := []*models.WSUser{}
//
//	for i := int64(0); i < users; i++ {
//		us = append(us, &models.WSUser{
//			User: models.User{
//				ID: i,
//				Group: 1,
//				RoomID: "",
//				Paired: make(chan bool),
//			},
//			Socket: nil,
//		})
//	}
//
//	return us
//}
//
//func BenchmarkNew(b *testing.B) {
//	matcher.users = fillUsers()
//
//	ch := make (chan [2]*models.WSUser)
//
//	for i := 0; i < b.N; i++ {
//		go func() {
//			for j := 0; j < users / 2; i++ {
//				<- ch
//			}
//		}()
//
//		matcher.getMatchingUsersUpgraded(ch)
//    }
//}

//func TestMatcher(t *testing.T) {
	//t.Run("testAddition", testAddition)
	//t.Run("testAddition2", testAddition2)
	//t.Run("testAddition3", testAddition3)
//}

//func testAddition(t *testing.T) {
//	expected := map[int][]*Unit{
//		0: []*Unit{
//			&Unit{ID: int64(1), Group: 1},
//			&Unit{ID: int64(2), Group: 1},
//		},
//	}
//
//	us := []*Unit{
//		&Unit{ID: int64(1), Group: 1},
//		&Unit{ID: int64(2), Group: 1},
//	}
//
//	matcher := NewMatcher()
//
//	for _, u := range us {
//		matcher.AddUnit(u)
//	}
//
//	matched := make(chan [2]*Unit)
//
//	go matcher.GetMatchingUnits(matched)
//
//	for i := 0; i < 1; i++ {
//		mu := <-matched
//
//		if mu[0].ID != expected[i][0].ID {
//			t.Error(fmt.Sprintf("Unexpected first unit at index %d: got %d, expected %d\n",
//				i, mu[0].ID, expected[i][0].ID))
//		}
//
//		if mu[1].ID != expected[i][1].ID {
//			t.Error(fmt.Printf("Unexpected second unit at index %d: got %d, expected %d\n",
//				i, mu[1].ID, expected[i][1].ID))
//		}
//	}
//
//	unitsLeft := matcher.Units()
//
//	if len(unitsLeft) != 0 {
//		t.Error("Expecting no units in the pool, have:", len(unitsLeft))
//	}
//}
//
//func testAddition2(t *testing.T) {
//	expected := map[int][]*Unit{
//		0: []*Unit{
//			&Unit{ID: int64(1), Group: 1},
//			&Unit{ID: int64(3), Group: 1},
//		},
//		1: []*Unit{
//			&Unit{ID: int64(4), Group: 3},
//			&Unit{ID: int64(5), Group: 3},
//		},
//	}
//
//	us := []*Unit{
//		&Unit{ID: int64(1), Group: 1},
//		&Unit{ID: int64(2), Group: 2},
//		&Unit{ID: int64(3), Group: 1},
//		&Unit{ID: int64(4), Group: 3},
//		&Unit{ID: int64(5), Group: 3},
//	}
//
//	matcher := NewMatcher()
//
//	for _, u := range us {
//		matcher.AddUnit(u)
//	}
//
//	matched := make(chan [2]*Unit)
//
//	go matcher.GetMatchingUnits(matched)
//
//	for i := 0; i < 2; i++ {
//		mu := <-matched
//
//		if mu[0].ID != expected[i][0].ID {
//			t.Error(fmt.Sprintf("Unexpected first unit at index %d: got %d, expected %d\n",
//				i, mu[0].ID, expected[i][0].ID))
//		}
//
//		if mu[1].ID != expected[i][1].ID {
//			t.Error(fmt.Printf("Unexpected second unit at index %d: got %d, expected %d\n",
//				i, mu[1].ID, expected[i][1].ID))
//		}
//	}
//
//	unitsLeft := matcher.Units()
//	if len(unitsLeft) != 1 {
//		t.Error("Expecting 1 unit in the pool, have:", len(unitsLeft))
//	}
//}
//
//func testAddition3(t *testing.T) {
//	expected := map[int][]*Unit{
//		0: []*Unit{
//			&Unit{ID: int64(1), Group: 3},
//			&Unit{ID: int64(2), Group: 3},
//		},
//		1: []*Unit{
//			&Unit{ID: int64(3), Group: 3},
//			&Unit{ID: int64(4), Group: 3},
//		},
//		2: []*Unit{
//			&Unit{ID: int64(6), Group: 3},
//			&Unit{ID: int64(7), Group: 3},
//		},
//	}
//
//	us := []*Unit{
//		&Unit{ID: int64(1), Group: 3},
//		&Unit{ID: int64(2), Group: 3},
//		&Unit{ID: int64(3), Group: 3},
//		&Unit{ID: int64(4), Group: 3},
//		&Unit{ID: int64(5), Group: 3},
//		&Unit{ID: int64(6), Group: 4},
//		&Unit{ID: int64(7), Group: 4},
//	}
//
//	matcher := NewMatcher()
//
//	for _, u := range us {
//		matcher.AddUnit(u)
//	}
//
//	matched := make(chan [2]*Unit)
//
//	go matcher.GetMatchingUnits(matched)
//
//	for i := 0; i < 3; i++ {
//		mu := <-matched
//
//		if mu[0].ID != expected[i][0].ID {
//			t.Error(fmt.Sprintf("Unexpected first unit at index %d: got %d, expected %d\n",
//				i, mu[0].ID, expected[i][0].ID))
//		}
//
//		if mu[1].ID != expected[i][1].ID {
//			t.Error(fmt.Printf("Unexpected second unit at index %d: got %d, expected %d\n",
//				i, mu[1].ID, expected[i][1].ID))
//		}
//	}
//
//	unitsLeft := matcher.Units()
//	if len(unitsLeft) != 1 {
//		t.Error("Expecting 1 unit in the pool, have:", len(unitsLeft))
//	}
//}
