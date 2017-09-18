package score

import (
	"errors"
	"fmt"
	"studentinfomanagement/src/lib"
	"sync"

	"github.com/astaxie/beego/orm"
)

//ScoreOparetor 成绩表操作类
type ScoreOparetor struct {
	myOrm orm.Ormer
	mux   sync.Mutex
}

//NewScoreOparetor ScoreOparetor构造函数
func NewScoreOparetor(orm orm.Ormer) *ScoreOparetor {
	so := ScoreOparetor{myOrm: orm}
	return &so
}

//Insert 插入一个成绩到成绩表
func (so *ScoreOparetor) Insert(sco *lib.Score) (interface{}, error) {
	so.mux.Lock()
	defer so.mux.Unlock()
	id, err := so.myOrm.Insert(sco)
	if err != nil {
		return nil, err
	}
	var soAfterInsert lib.Score
	err = so.myOrm.QueryTable("score").Filter("id", id).One(&soAfterInsert)
	if err != nil {
		return nil, err
	}
	return soAfterInsert, nil
}

//Delete 删除一个成绩从成绩表
func (so *ScoreOparetor) Delete(id int, courseID int) (interface{}, error) {
	so.mux.Lock()
	defer so.mux.Unlock()
	var sco lib.Score
	err := so.myOrm.QueryTable("score").Filter("student_id", id).Filter("course_id", courseID).One(&sco)
	if err == nil {
		num, _ := so.myOrm.Delete(&sco)
		return num, nil
	}
	if err == orm.ErrMissPK {
		return nil, errors.New("找不到主键")
	}
	if err == orm.ErrNoRows {
		return nil, errors.New("找不到该成绩")
	}
	return nil, errors.New("未知错误")
}

//Update 更新一个成绩从成绩表
func (so *ScoreOparetor) Update(id int, sc lib.Score) (interface{}, error) {
	so.mux.Lock()
	defer so.mux.Unlock()
	var sco lib.Score
	err := so.myOrm.QueryTable("score").Filter("student_id", id).Filter("course_id", sc.Course_id).One(&sco)
	if err == nil {
		sc.Id = sco.Id
		_, err = so.myOrm.Update(&sc)
		if err == nil {
			var score lib.Score
			so.myOrm.QueryTable("score").Filter("id", sc.Id).One(&score)
			return score, nil
		}
	}
	if err == orm.ErrNoRows {
		return nil, errors.New("找不到记录")
	}
	if err == orm.ErrMissPK {
		return nil, errors.New("找不到主键")
	}
	return nil, errors.New("未知错误")
}

//QueryAll 全部成绩从成绩表
func (so *ScoreOparetor) QueryAll() (interface{}, error) {
	so.mux.Lock()
	defer so.mux.Unlock()
	var soAfterQuery []lib.Score
	_, err := so.myOrm.QueryTable("score").All(&soAfterQuery)
	if err == nil {
		return soAfterQuery, nil
	} else if err == orm.ErrNoRows {
		return nil, errors.New("找不到记录")
	}
	return nil, errors.New("未知错误")
}

//QueryAll 查询某个同学所有科目的成绩
func (so *ScoreOparetor) QueryBy(SID int) (interface{}, error) {
	so.mux.Lock()
	defer so.mux.Unlock()
	var scores []lib.Score
	_, err := so.myOrm.QueryTable("score").Filter("student_id", SID).All(&scores)
	length := len(scores)
	if err == nil && length > 0 {
		return scores, nil
	} else if err == nil && length == 0 {
		return nil, errors.New("找不到成绩")
	} else {
		return nil, err
	}
}

//QueryClass 查找班级成绩
func (so *ScoreOparetor) QueryClass(classID int) (interface{}, error) {
	so.mux.Lock()
	defer so.mux.Unlock()
	var stuArr []lib.Student
	_, err := so.myOrm.QueryTable("student").Filter("class_id", classID).All(&stuArr)
	fmt.Println(stuArr)
	if err == nil && len(stuArr) > 0 {
		scoreArr := make([]lib.Score, 0)
		for _, stu := range stuArr {
			var score []lib.Score
			_, err := so.myOrm.QueryTable("score").Filter("student_id", stu.Id).All(&score)
			fmt.Println(score)
			if err == nil && len(score) > 0 {
				for _, sco := range score {
					fmt.Println(sco)
					scoreArr = append(scoreArr, sco)
				}
				return scoreArr, nil
			} else if err == nil && len(score) == 0 {
				return nil, errors.New("该班级无成员")
			}
		}

	} else if err == nil && len(stuArr) == 0 {
		return nil, errors.New("该班级无成员")
	} else if err == orm.ErrNoRows {
		return nil, errors.New("找不到记录")
	}
	return nil, errors.New("未知错误")
}
