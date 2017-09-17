package student

import (
	"errors"
	"studentinfomanagement/src/lib"
	"sync"

	"github.com/astaxie/beego/orm"
)

//StudentOperator 学生数据管理类
type StudentOperator struct {
	myOrm orm.Ormer
	mux   sync.Mutex
}

//NewStudentOperator 返回一个StudentOperator指针
func NewStudentOperator(orm orm.Ormer) (*StudentOperator, error) {
	sto := StudentOperator{myOrm: orm}
	return &sto, nil
}

//Insert 插入一个学生到学生表
func (sto *StudentOperator) Insert(stu *lib.Student) (interface{}, error) {
	sto.mux.Lock()
	defer sto.mux.Unlock()
	id, err := sto.myOrm.Insert(stu)
	if err != nil {
		return nil, err
	}
	var stuAfterInsert lib.Student
	err = sto.myOrm.QueryTable("student").Filter("id", id).One(&stuAfterInsert)
	if err != nil {
		return nil, err
	}
	return stuAfterInsert, nil
}

//Delete 删除一个学生从学生表
func (sto *StudentOperator) Delete(id int) (interface{}, error) {
	sto.mux.Lock()
	defer sto.mux.Unlock()
	num, err := sto.myOrm.Delete(&lib.Student{Id: id})
	if err == nil && num > 0 {
		return nil, nil
	} else if err == nil && num == 0 {
		return nil, errors.New("找不到该学生")
	} else if err == orm.ErrMissPK {
		return nil, errors.New("找不到主键")
	} else if err == orm.ErrNoRows {
		return nil, errors.New("找不到该记录")
	} else {
		return nil, errors.New("未知错误")
	}
}

//Update 更新一个学生从学生表
func (sto *StudentOperator) Update(id int, stu lib.Student) (interface{}, error) {
	sto.mux.Lock()
	defer sto.mux.Unlock()
	err := sto.myOrm.Read(&lib.Student{Id: id})
	if err == nil {
		_, err = sto.myOrm.Update(&stu)
		if err != nil {
			return nil, err
		}
		var stuAfterUpdate lib.Student
		sto.myOrm.QueryTable("student").Filter("id", id).One(&stuAfterUpdate)
		return stuAfterUpdate, nil
	}
	if err == orm.ErrNoRows {
		return nil, errors.New("找不到记录")
	}
	if err == orm.ErrMissPK {
		return nil, errors.New("找不到主键")
	}
	return nil, errors.New("未知错误")
}

//Query 查询一个学生从学生表
func (sto *StudentOperator) Query(id int) (interface{}, error) {
	sto.mux.Lock()
	defer sto.mux.Unlock()
	var stu []lib.Student
	err := sto.myOrm.QueryTable("student").Filter("id", id).One(&stu)
	if err == nil {
		return stu, nil
	}
	if err == orm.ErrMultiRows {
		return nil, errors.New("找到多条记录")
	}
	if err == orm.ErrNoRows {
		return nil, errors.New("找不到记录")
	}
	return nil, errors.New("未知错误")
}

//QueryStudentByClassID 查询同一个班级的学生
func (sto *StudentOperator) QueryStudentByClassID(classID int) (interface{}, error) {
	sto.mux.Lock()
	defer sto.mux.Unlock()
	var stuArr []lib.Student
	_, err := sto.myOrm.QueryTable("student").Filter("class_id", classID).All(&stuArr)
	if err == nil && len(stuArr) > 0 {
		return stuArr, nil
	} else if err == orm.ErrNoRows || len(stuArr) == 0 {
		return nil, errors.New("找不到记录")
	} else {
		return nil, errors.New("未知错误")
	}
}

//QueryStudentAll 查询所有学生
func (sto *StudentOperator) QueryStudentAll() (interface{}, error) {
	sto.mux.Lock()
	defer sto.mux.Unlock()
	var stuArr []lib.Student
	_, err := sto.myOrm.QueryTable("student").All(&stuArr)
	if err == nil && len(stuArr) > 0 {
		return stuArr, nil
	} else if err == orm.ErrNoRows || len(stuArr) == 0 {
		return nil, errors.New("找不到记录")
	} else {
		return nil, errors.New("未知错误")
	}
}
