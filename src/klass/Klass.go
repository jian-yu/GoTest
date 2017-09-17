package klass

import (
	"errors"
	"studentinfomanagement/src/lib"
	"sync"

	"github.com/astaxie/beego/orm"
)

//ClassOperator 学生数据管理类
type ClassOperator struct {
	myOrm orm.Ormer
	mux   sync.Mutex
}

//NewClassOperator ClassOperator构造函数
func NewClassOperator(orm orm.Ormer) *ClassOperator {
	nco := ClassOperator{myOrm: orm}
	return &nco
}

//Insert 插入一个班级到Class表中
func (co *ClassOperator) Insert(klass *lib.Class) (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	id, err := co.myOrm.Insert(klass)
	if err != nil {
		return nil, err
	}
	var classAfterInsert lib.Class
	err = co.myOrm.QueryTable("class").Filter("id", id).One(&classAfterInsert)
	if err != nil {
		return nil, err
	}
	return classAfterInsert, nil
}

//Delete  根据id从Class表删除一个Class
func (co *ClassOperator) Delete(id int) (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	num, err := co.myOrm.Delete(&lib.Class{Id: id})
	if err == nil && num > 0 {
		return num, nil
	} else if err == nil && num == 0 {
		return nil, errors.New("找不到该班级")
	} else if err == orm.ErrMissPK {
		return nil, errors.New("找不到该ID")
	} else if err == orm.ErrNoRows {
		return nil, errors.New("找不到该记录")
	}
	return nil, errors.New("未知错误")
}

//Update  根据ID 更改一个Class
func (co *ClassOperator) Update(id int, klass lib.Class) (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	err := co.myOrm.Read(&lib.Class{Id: id})
	if err == nil {
		_, err = co.myOrm.Update(&klass)
		if err != nil {
			return nil, err
		}
		var classAfterUpdate lib.Class
		co.myOrm.QueryTable("class").Filter("id", id).One(&classAfterUpdate)
		return classAfterUpdate, nil
	}
	if err == orm.ErrMissPK {
		return nil, errors.New("找不到该ID")
	}
	if err == orm.ErrNoRows {
		return nil, errors.New("找不到记录")
	}
	return nil, errors.New("未知错误")
}

//QueryBy 根据ID查询一个Class
func (co *ClassOperator) QueryBy(id int) (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	var klass lib.Class
	err := co.myOrm.QueryTable("class").Filter("id", id).One(&klass)
	if err == nil {
		return klass, nil
	}
	if err == orm.ErrMissPK {
		return nil, errors.New("找不到该ID")
	}
	if err == orm.ErrNoRows {
		return nil, errors.New("找不到该记录")
	}
	return nil, errors.New("未知错误")
}

//QueryAll 返回所有班级
func (co *ClassOperator) QueryAll() (interface{}, error) {
	co.mux.Lock()
	defer co.mux.Unlock()
	var klass []lib.Class
	_, err := co.myOrm.QueryTable("class").All(&klass)
	if err == nil {
		return klass, nil
	}
	if err == orm.ErrMissPK {
		return nil, errors.New("找不到该ID")
	}
	if err == orm.ErrNoRows {
		return nil, errors.New("找不到该记录")
	}
	return nil, errors.New("未知错误")
}
