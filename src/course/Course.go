package course

import (
	"errors"
	"fmt"
	"studentinfomanagement/src/lib"
	"sync"

	"github.com/astaxie/beego/orm"
)

//CourseOperator 课程管理类
type CourseOperator struct {
	myOrm orm.Ormer
	mux   sync.Mutex
}

//NewCourseOperator CourseOperator构造函数
func NewCourseOperator(orm orm.Ormer) *CourseOperator {
	co := CourseOperator{myOrm: orm}
	return &co
}

//Insert 插入一个课程到课程表
func (co *CourseOperator) Insert(cou *lib.Course) (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	var course lib.Course
	err := co.myOrm.QueryTable("course").Filter("name", cou.Name).One(&course)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println(course)
		_, err := co.myOrm.Insert(cou)
		if err == nil {
			return cou, nil
		}
	} else if err == nil && course.Id == 0 {
		_, err := co.myOrm.Insert(&cou)
		if err == nil {
			return cou, nil
		}

	} else if err == nil && course.Id != 0 {
		return nil, errors.New("该课程已存在")
	}
	return nil, err
}

//Delete 删除一个课程从课程表，并且删除所有对应的学生的成绩记录
func (co *CourseOperator) Delete(cou *lib.Course) (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	_, err := co.myOrm.Delete(&cou)
	if err == nil {
		var score []lib.Score
		_, err = co.myOrm.QueryTable("sorce").Filter("id", cou.Id).All(&score)
		if err == nil && len(score) > 0 {
			for v := range score {
				co.myOrm.Delete(&v)
			}
		} else if err == nil && len(score) == 0 {
			return nil, errors.New("找不到该记录")
		}
	} else if err == orm.ErrMissPK {
		return nil, errors.New("找不到主键")
	} else if err == orm.ErrNoRows {
		return nil, errors.New("找不到该记录")
	}
	return nil, errors.New("未知错误")
}

//QueryAll 查询所有课程
func (co *CourseOperator) QueryAll() (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	var courses []lib.Course
	_, err := co.myOrm.QueryTable("course").OrderBy("id").All(&courses)
	if err == nil && len(courses) > 0 {
		return courses, nil
	} else if err == nil && len(courses) == 0 {
		return nil, errors.New("无班级")
	}
	return nil, err
}
