package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"os"
	"strconv"
	"sort"
	"math"
)

type ircline struct {
	time int
	message string
}

//For sorting
type irclines []ircline
func (p irclines) Len() int           { return len(p) }
func (p irclines) Less(i, j int) bool { return p[i].time < p[j].time }
func (p irclines) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//Main struct
type UserInfo struct {
	Lines     irclines
	MeanDelay float64
	StdDelay  float64

	dirty bool
}

//Our aliases
var people = map[string][]string{
	"jeff": []string{"zeebo", "zeebo_", "zeeboo"},
	"jake": []string{"jbaikge"},
	"jamie": []string{"blooze", "jamie_okc", "jqbokco", "jburdette", "jamieb"},
	"mark": []string{"ldapmonk", "ldapmonk1", "pixelmonk", "ldapmonky"},
	"trish": []string{"trishy1", "trishytr1", "trishytri", "trishy"},
	"chad": []string{"Chad2", "Chad1", "Chad"},
	"trafbot": []string{"TRAFBOT_", "TRAFBOT", "OKCOBOT", "OKCOBOT_"},
}

func Analyze(file string) (map[string]*UserInfo, os.Error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	userInfo := make(map[string]*UserInfo)

	chunks := strings.Split(string(data), "\n", -1)
	for _, line := range chunks {
		info := strings.Split(line, " ", 3)

		time, err := strconv.Atoi(info[0])
		if err != nil {
			return userInfo, err
		}

		newLine := ircline{
			time: time,
			message: info[2],
		}
		
		username := strings.ToLower(info[1])
		user, exists := userInfo[username]
		if !exists {
			userInfo[username] = NewUserInfo()
			user = userInfo[username]
		}

		user.AddLine(newLine)
	}

	combinedInfo := make(map[string]*UserInfo)
	for name, aliases := range people {
		combinedInfo[name] = NewUserInfo()

		for _, alias := range aliases {
			if user, exists := userInfo[strings.ToLower(alias)]; exists {
				combinedInfo[name].Combine(user)
			}
		}
	}

	return combinedInfo, nil
}

func NewUserInfo() *UserInfo {
	return &UserInfo{
		Lines: make(irclines, 0),
		dirty: true,
	}
}

func (u *UserInfo) AddLine(newLine ircline) {
	u.Lines = append(u.Lines, newLine)
	u.dirty = true
}

func (u *UserInfo) Combine(u2 *UserInfo) {
	u.Lines = append(u.Lines, u2.Lines...)
	u.dirty = true
}

func (u *UserInfo) String() string {
	if u.dirty {
		u.Calculate()
	}
	return fmt.Sprintf("%d lines [%f:%f]", len(u.Lines), u.MeanDelay, u.StdDelay)
}

func isum(a []int) (s int) {
	for _, v := range a { s += v }
	return
}

func fsum(a []float64) (s float64) {
	for _, v := range a { s += v }
	return
}

func filter(a []int, f func(int) bool) []int {
	ret := make([]int, 0)
	for _, v := range a {
		if f(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func sq(x float64) float64 { return x * x }

func (u *UserInfo) Calculate() {
	sort.Sort(u.Lines)
	delays, first := make([]int, 0), u.Lines[0]
	for _, second := range u.Lines[1:] {
		delays, first = append(delays, second.time - first.time), second
	}

	delays = filter(delays, func(x int) bool { return x < 60*60 })
	u.MeanDelay = float64(isum(delays)) / float64(len(delays))

	diffs := make([]float64, len(delays))
	for i := range delays {
		diffs[i] = sq(float64(delays[i]) - u.MeanDelay)
	}

	u.StdDelay = math.Sqrt(fsum(diffs))
}