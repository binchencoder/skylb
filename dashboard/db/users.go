package db

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	ld "gopkg.in/ldap.v2"
	db "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	"binchencoder.com/letsgo/ldap"
	pb "binchencoder.com/skylb/dashboard/proto"
)

const (
	RoleAdmin = iota
	RoleDeveloper
	RoleOperation

	tableUsers = "users"
)

var ldapEndpoint = flag.String("ldap-endpoint", "", "The address of ldap endpoint. eg: ldap.eff.com:389")

// RoleSet represents a set of roles for a user.
type RoleSet map[int]struct{}

// User represents a row in database table "users".
type User struct {
	LoginName string `db:"login_name"`
	Roles     string `db:"role"`
	Disabled  bool   `db:"disabled"`
	Version   int32  `db:"version"`
}

// AddRole adds the given role into the user.
func (u *User) AddRole(role int32) {
	if u.Roles == "" {
		u.Roles = strconv.Itoa(int(role))
	} else {
		u.Roles += "," + strconv.Itoa(int(role))
	}
}

// GetRoles returns a RoleSet for the user.
func (u *User) GetRoles() RoleSet {
	rs := RoleSet{}
	for _, r := range strings.Split(u.Roles, ",") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		rvalue, err := strconv.Atoi(r)
		if err != nil {
			glog.Errorf("Failed to convert user role, %v", err)
		}
		rs[rvalue] = struct{}{}
	}
	return rs
}

// ToUserInfo converts the User struct to a pb.UserInfo struct.
func (u *User) ToUserInfo() *pb.UserInfo {
	info := pb.UserInfo{
		LoginName: u.LoginName,
		Disabled:  u.Disabled,
		Version:   u.Version,
	}

	rs := u.GetRoles()
	if _, ok := rs[RoleAdmin]; ok {
		info.IsAdmin = true
	}
	if _, ok := rs[RoleDeveloper]; ok {
		info.IsDev = true
	}
	if _, ok := rs[RoleOperation]; ok {
		info.IsOps = true
	}

	return &info
}

// FromUserInfo creates a new User struct from the given UserInfo struct.
func FromUserInfo(info *pb.UserInfo) *User {
	user := User{
		LoginName: info.LoginName,
		Disabled:  info.Disabled,
		Version:   info.Version,
	}

	if info.IsAdmin {
		user.AddRole(RoleAdmin)
	}
	if info.IsDev {
		user.AddRole(RoleDeveloper)
	}
	if info.IsOps {
		user.AddRole(RoleOperation)
	}

	return &user
}

// Authenticate logs user in the system with LDAP, finds the user in database
//  and merge in its roles.
func Authenticate(loginname, passwd string) (*User, error) {
	client, err := ld.Dial("tcp", *ldapEndpoint)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ud := ldap.NewUserDN(client)

	ok, err := ud.AuthUserByID(loginname, passwd)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid user name or password")
	}

	user := User{LoginName: loginname}

	if err = MergeUserRoles(context.Background(), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// MergeUserRoles finds the user in database and merges in user roles.
func MergeUserRoles(ctx context.Context, user *User) error {
	tx, err := dbClient.NewTx(ctx)
	if err != nil {
		glog.Errorf("Failed to begin transaction, %v", err)
		return err
	}

	u, err := loadUserByLoginName(tx, user.LoginName)
	if err != nil && err != db.ErrNoMoreRows {
		tx.Rollback()
		glog.Errorf("Failed to load user from database, %v", err)
		return err
	}
	tx.Commit()

	if u == nil {
		// If user does not exist in db, set developer role and insert.
		user.AddRole(RoleDeveloper)
		system := User{LoginName: "System"}
		if err := insertUser(ctx, &system, user); err != nil {
			return err
		}
	} else {
		user.Disabled = u.Disabled
		user.Version = u.Version
		user.Roles = u.Roles
	}

	return nil
}

// GetUsers returns a copy of all users.
func GetUsers(ctx context.Context) ([]*User, error) {
	tx, err := dbClient.NewTx(ctx)
	if err != nil {
		glog.Errorf("Failed to begin transaction, %v", err)
		return nil, err
	}

	us, err := loadUsers(tx)
	if err != nil {
		glog.Errorf("Failed to load users from database, %v", err)
		return nil, err
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return us, nil
}

// UpsertUser updates or insert the given user into database.
func UpsertUser(ctx context.Context, curUser, newUser *User, isNew bool) error {
	var err error
	if isNew {
		err = insertUser(ctx, curUser, newUser)
	} else {
		err = updateUser(ctx, curUser, newUser)
	}

	if err != nil {
		return err
	}

	return nil
}

func loadUsers(tx sqlbuilder.Tx) ([]*User, error) {
	us := []*User{}
	err := withTx(tx, func() error {
		if err := tx.Collection(tableUsers).Find().OrderBy("login_name").All(&us); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return us, nil
}

func loadUserByLoginName(tx sqlbuilder.Tx, loginname string) (*User, error) {
	us := User{}
	err := withTx(tx, func() error {
		if err := tx.Collection(tableUsers).Find("login_name", loginname).One(&us); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func insertUser(ctx context.Context, curUser, newUser *User) error {
	tx, err := dbClient.NewTx(ctx)
	if err != nil {
		glog.Errorf("Failed to begin transaction, %v", err)
		return err
	}

	err = withTx(tx, func() error {
		if _, err = tx.Collection(tableUsers).Insert(newUser); err != nil {
			return err
		}

		return insertLog(tx, curUser.LoginName, 0, fmt.Sprintf("Insert user \"%s\"", newUser.LoginName))
	})
	if err != nil {
		glog.Errorf("Failed to insert user to database, %v", err)
		return err
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func updateUser(ctx context.Context, curUser, targetUser *User) error {
	tx, err := dbClient.NewTx(ctx)
	if err != nil {
		glog.Errorf("Failed to begin transaction, %v", err)
		return err
	}

	err = withTx(tx, func() error {
		u := User{}
		res := tx.Collection(tableUsers).Find("login_name", targetUser.LoginName)
		if err = res.One(&u); err != nil {
			return err
		}

		u.Roles = targetUser.Roles
		u.Disabled = targetUser.Disabled
		u.Version++

		if err := res.Update(&u); err != nil {
			return err
		}

		return insertLog(tx, curUser.LoginName, 0, fmt.Sprintf("Update user \"%s\"", targetUser.LoginName))
	})
	if err != nil {
		tx.Rollback()
		glog.Errorf("Failed to update user to database, %v", err)
		return err
	}
	tx.Commit()

	return nil
}
