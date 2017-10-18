package main

import (
        "net/http"
)

var sess = map[string]bool{}

func isSession(r *http.Request) bool {
        cookie, err := r.Cookie("session-id")
        if err != nil {
                return false
        }
        _, ok := sess[cookie.Value]
        return ok
}

func addSession(w http.ResponseWriter) {
        u, err := genUUID()
        if err != nil {
                return
        }
        sess[u] = true
        http.SetCookie(w, &http.Cookie{
                Name: "session-id",
                Value: u,
        })
}

func endSession(w http.ResponseWriter, r *http.Request) {
        if isSession(r) {
                cookie, err := r.Cookie("session-id")
                if err == nil {
                        delete(sess, cookie.Value)
                }
        }
        http.Redirect(w, r, "/", 307)
}
