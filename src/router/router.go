package router

import (
    //"fmt"
    "strconv"
    "regexp"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "github.com/gorilla/mux"
    "pgv3"
)

type PostUserData struct {
    Name string `json:"name"`
}

type PostRelationshipData struct {
    State string `json:"state"`
}

// Do we support the given URL?
func urlSupported(url string) bool {
    routes:= []string{
        `^/users$`,
        `^/users/{user_id:[0-9]+}/relationships$`,
        `^/users/{user_id:[0-9]+}/relationships/{other_user_id:[0-9]+}$`,
    }
    for _, s := range routes {
        re := regexp.MustCompile(s)
        if re.Match([]byte(url)) {
            return true
        }
    }
    return false
}

func response(w http.ResponseWriter, code int, data []byte) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(code)
    w.Write(data)
    w.Write([]byte("\r\n"))
}

func responseErrorInternal(w http.ResponseWriter) {
    res := []byte(`{"code": 500, "message": "Internal exceptions"}`)
    response(w, http.StatusInternalServerError, res)
}

func responseErrorBadRequest(w http.ResponseWriter) {
    res := []byte(`{"code": 400, "message": "Bad Requests"}`)
    response(w, http.StatusBadRequest, res)
}

func responseErrorUnprocessableEntity(w http.ResponseWriter) {
    res := []byte(`{"code": 422, "message": "Unprocessable Entity"}`)
    response(w, 422, res)
}

func HandleBlackHoleRoute(w http.ResponseWriter, r *http.Request) {
    //fmt.Println("r.URL.Host: ", r.URL.Host)
    //fmt.Println("r.URL.Scheme: ", r.URL.Scheme)
    //fmt.Println("r.URL.Opaque: ", r.URL.Opaque)
    //fmt.Println("r.URL.Path: ", r.URL.Path)
    //fmt.Println("r.URL.RawPath: ", r.URL.RawPath)
    //fmt.Println("r.URL.RawQuery: ", r.URL.RawQuery)
    //fmt.Println("r.URL.Fragment: ", r.URL.Fragment)
    
    if urlSupported(r.URL.Path) {
        res := []byte(`{"code": 405, "message": "Method Not Allowed"}`)
        response(w, http.StatusMethodNotAllowed, res)
    } else {
        res := []byte(`{"code": 404, "message": "Not Found"}`)
        response(w, http.StatusNotFound, res)
    }
}

func RetrieveUsers(w http.ResponseWriter, r *http.Request) {
    db := pgv3.ConnectDatabase()
    users, err := pgv3.GetUsers(db)
    if err != nil {
        responseErrorInternal(w)
        return
    }

    for i := 0; i < len(users); i++ { // mark Type as 'user'
        users[i].Type = "user"
    }

    res, err := json.MarshalIndent(users, "", "  ")
    if err != nil {
        responseErrorInternal(w)
        return
    }

    response(w, http.StatusOK, res);
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseErrorInternal(w)
        return
    }
    var data PostUserData
    err = json.Unmarshal(body, &data)
    if err != nil {
        responseErrorUnprocessableEntity(w)
        return
    }
    //fmt.Println("Posted user name: ", data.Name)

    user := &pgv3.User{
        Name: data.Name,
        Type: "user",
    }
    db := pgv3.ConnectDatabase()
    err = pgv3.CreateUser(db, user)
    if err != nil {
        responseErrorInternal(w)
        return
    }

    res, err := json.MarshalIndent(user, "", "  ")
    if err != nil {
        responseErrorInternal(w)
        return
    }

    response(w, http.StatusCreated, res);
}

func RetrieveRelationship(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.ParseInt(params["user_id"], 10, 64)
    if err != nil {
        responseErrorUnprocessableEntity(w)
        return
    }
    db := pgv3.ConnectDatabase()
    rels, err := pgv3.GetUserRelationships(db, id)
    if err != nil {
        responseErrorInternal(w)
        return
    }

    for i := 0; i < len(rels); i++ {
        rels[i].Type = "relationship"
    }

    res, err := json.MarshalIndent(rels, "", "  ")
    if err != nil {
        responseErrorInternal(w)
        return
    }

    response(w, http.StatusOK, res);
}

func UpdateRelationship(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, err := strconv.ParseInt(params["user_id"], 10, 64)
    if err != nil {
        responseErrorUnprocessableEntity(w)
        return
    }
    otherId, err := strconv.ParseInt(params["other_user_id"], 10, 64)
    if err != nil {
        responseErrorUnprocessableEntity(w)
        return
    }
    if id == otherId {
        responseErrorBadRequest(w)
        return
    }
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseErrorUnprocessableEntity(w)
        return
    }
    var data PostRelationshipData
    err = json.Unmarshal(body, &data)
    if err != nil {
        responseErrorUnprocessableEntity(w)
        return
    }
    rel := &pgv3.Relationship{
        Id: id,
        OtherId: otherId,
        State: data.State,
        Type: "relationship",
    }
    //fmt.Println("Posted relainship state: ", id, otherId, data.State)

    db := pgv3.ConnectDatabase()
    err = pgv3.UpdateUserRelationship(db, rel)
    if err != nil {
        responseErrorInternal(w)
        return
    }

    res, err := json.MarshalIndent(rel, "", "  ")
    if err != nil {
        responseErrorInternal(w)
        return
    }

    response(w, http.StatusOK, res);
}

