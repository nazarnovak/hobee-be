package matcher

import (
	"fmt"
	"sort"
	"sync"

	"hobee-be/models"
)

var (
	matcher = &Matcher{}
	mutex   = &sync.Mutex{}
)

type Matcher struct {
	users []*models.WSUser
}

func New() *Matcher{
	return &Matcher{}
}


func AddUser(u *models.WSUser) {
	fmt.Println("Adding ID:", u.ID)
	mutex.Lock()
	matcher.users = append(matcher.users, u)
	mutex.Unlock()
}

func Users() []*models.WSUser {
	return matcher.users
}

func (m *Matcher) getMatchingUsers(ch chan<- [2]*models.WSUser) {
	if len(m.users) < 2 {
		return
	}

	type poolUser struct {
		poolIndex int
		user *models.WSUser
	}

	groupsUsers := map[int64][]poolUser{}

	// First we group users by groups
	for poolIndex, u := range m.users {
		pu := poolUser{
			poolIndex: poolIndex,
			user: u,
		}
		groupsUsers[u.Group] = append(groupsUsers[u.Group], pu)
	}

	toDelete := []int{}

	for _, gu := range groupsUsers {
		pairs := len(gu) / 2

		if pairs == 0 {
			continue
		}

		for i := 0; i < pairs; i++ {
			toDelete = append(toDelete, []int{gu[i].poolIndex, gu[i+1].poolIndex}...)
			ch <- [2]*models.WSUser{gu[i].user, gu[i+1].user}
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(toDelete)))

	mutex.Lock()
	for _, i := range toDelete {
		m.users = append(m.users[:i], m.users[i+1:]...)
	}
	mutex.Unlock()
}

func (m *Matcher) Run(ch chan<- [2]*models.WSUser) {
	for {
		m.getMatchingUsers(ch)
	}
}
