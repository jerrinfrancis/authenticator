package user

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jerrinfrancis/authenticator/db"
	"github.com/jerrinfrancis/authenticator/db/mongo"
	"github.com/jerrinfrancis/authenticator/pkg/hash"
)

const secret = "(*$#^#HJASD@#%2dasflkadsf)"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := db.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		panic(err)
	}
	defer r.Body.Close()
	cl := mongo.Client()
	repo := mongo.New(cl)
	du, _ := repo.User().Find(u.Email)

	pltxt, _ := base64.StdEncoding.DecodeString(u.Passwd)
	if hash.Verify(du.Passwd, string(pltxt), du.Salt) {
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		fmt.Println("DUI > ID", du.ID, du.Secret)
		claims["sub"] = du.ID
		claims["test"] = "jerrin"
		token.Claims = claims
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			fmt.Println(err)
		}
		w.Write([]byte(tokenString))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := db.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		panic(err)
	}
	defer r.Body.Close()
	pltxt, err := base64.StdEncoding.DecodeString(u.Passwd)
	if err != nil {
		panic(err)
	}
	salt := hash.Salt(16)
	hp := hash.Hash(string(pltxt), salt)
	cl := mongo.Client()
	repo := mongo.New(cl)
	u.Passwd = hex.EncodeToString([]byte(hp))
	u.Salt = hex.EncodeToString([]byte(salt))
	ru, err := repo.User().Insert(u)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(ru.UserInfo)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
