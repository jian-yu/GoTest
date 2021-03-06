package lib

import (
	_ "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Class_id int   `json:"class_id"`
}

type StudentPassword struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type Course struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type CourseRequset struct {
	Body Course `json:"body"`
}

type Score struct {
	Id         int `json:"id"`
	Student_id int `json:"student_id"`
	Course_id  int `json:"course_id"`
	Score      int `json:"score"`
}
type Score1 struct {
	Id         int    `json:"id",orm:"column(id)"`
	Name       string `json:"name",orm:"-"`
	Student_id int    `json:"student_id",orm:"column(student_id)"`
	Course_id  int    `json:"course_id",orm:"column(course_id)"`
	Class_id   int    `json:"class_id",orm:"column(class_id)"`
	Score      int    `json:"score",orm:"column(score)"`
}

type ScoreRequest struct {
	Body Score `json:"body"`
}

type Class struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type ClassRequest struct {
	Body Class `json:"body"`
}

type Manager struct {
	Id       int    `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

type LoginRequest struct {
	Body Login `json:"body"`
}
type Login struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type StudentRequest struct {
	Body Student `json:"body"`
}

type ManagerRequest struct {
	Body Manager `json:"body"`
}
