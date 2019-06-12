package mggo

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
)

const (
	iterations = 1e4
	salt       = "abcdefghijklmnopqrstuvwxyz"
)

// User is user struct
type User struct {
	ID    int    `sql:",pk"`
	Name  string `sql:",notnull"`
	Login string `sql:",unique"`
	Right int    `sql:"default:4"`
}

// UserPassword is user password
type UserPassword struct {
	UserID   int
	Password string
}

// Identity is get user id by login and password
func (u *User) Identity(login, password string) int {
	var id int

	SQL().QueryOne(Scan(&id), `
        select "id" from "users" 
        join "user_passwords" on "user_id" = "id"
        where "login" = ? and "password" = ?
    `, login, GenerateFromPassword(password))
	return id
}

// GetCurrentUserInfo is get current user info from session
func (u *User) GetCurrentUserInfo(ctx *BaseContext) User {
	if ctx.CurrentUser.ID != 0 {
		return ctx.CurrentUser
	}
	id := SAP{}.SessionUserID(ctx)
	if id == 0 {
		return ctx.CurrentUser
	}
	cache := "User.GetCurrentUserInfo" + cacheUserPrefix + string(id)
	if value, ok := Cache.get(cache); ok {
		ctx.CurrentUser = value.(User)
		return ctx.CurrentUser
	}

	u.ID = id
	ctx.CurrentUser = u.Read(ctx)
	//1 day
	Cache.set("User.GetCurrentUserInfo", cache, ctx.CurrentUser, 60*60*24)
	return ctx.CurrentUser
}

// Read is read user by user id
func (u *User) Read(ctx *BaseContext) User {
	SQL().Select(u)
	return *u
}

// Update is insert or update user
func (u *User) Update(ctx *BaseContext) int {
	if u.ID == 0 {
		SQL().Insert(&u)
	} else {
		SQL().Update(&u)
	}
	return u.ID
}

// SetPassword is set password in user
func (u *User) SetPassword(id int, password string) bool {
	if id != 0 && password != "" {
		up := UserPassword{id, GenerateFromPassword(password)}
		SQL().Insert(&up)
		return true
	}
	return false
}

// GenerateFromPassword is generate password hash
func GenerateFromPassword(password string) string {
	hash := generateFromPassword([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)
	return fmt.Sprintf("%x", hash)
}

func generateFromPassword(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}
